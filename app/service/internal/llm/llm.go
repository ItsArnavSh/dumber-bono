package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const groqChatCompletionsURL = "https://api.groq.com/openai/v1/chat/completions"

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// non-streaming response shape (kept for reference / non-stream callers)
type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// streamChunk mirrors the SSE "data: {...}" payload shape for
// Groq's OpenAI-compatible streaming chat completions endpoint.
type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// QueryLLM sends a non-streaming prompt to Groq and returns the full generated response text.
func QueryLLM(ctx context.Context, systemPrompt, userPrompt, model string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY environment variable is not set")
	}
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}

	var messages []chatMessage
	if systemPrompt != "" {
		messages = append(messages, chatMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, chatMessage{Role: "user", Content: userPrompt})

	reqBody := chatCompletionRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, groqChatCompletionsURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(respBody, &completion); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}
	if completion.Error != nil {
		return "", fmt.Errorf("groq API error (%s): %s", completion.Error.Type, completion.Error.Message)
	}
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from groq API")
	}
	return completion.Choices[0].Message.Content, nil
}

// StreamLLM sends a streaming prompt to Groq (SSE) and writes each token's
// text content into w as it arrives.
//
// This is intended to be run in its own goroutine, typically paired with
// an io.Pipe:
//
//	pr, pw := io.Pipe()
//	go func() {
//	    err := llm.StreamLLM(ctx, systemPrompt, userPrompt, model, pw)
//	    // CloseWithError propagates err (nil is fine, it just closes cleanly)
//	    pw.CloseWithError(err)
//	}()
//	// consume pr on the calling side, e.g. io.Copy(os.Stdout, pr)
//	// or bufio.NewScanner(pr), etc.
//
// StreamLLM does NOT close w itself — the caller owns the writer's
// lifecycle (this matters for io.PipeWriter, which needs CloseWithError
// called exactly once by the producer side to unblock any reader).
//
// If w implements http.Flusher (e.g. it's wrapping an http.ResponseWriter),
// StreamLLM will flush after every chunk write so consumers see data
// incrementally rather than buffered.
func StreamLLM(ctx context.Context, systemPrompt, userPrompt, model string, w io.Writer) error {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GROQ_API_KEY environment variable is not set")
	}
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}

	var messages []chatMessage
	if systemPrompt != "" {
		messages = append(messages, chatMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, chatMessage{Role: "user", Content: userPrompt})

	reqBody := chatCompletionRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, groqChatCompletionsURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// Non-200: body is a plain JSON error object, not SSE. Read it all
		// for a useful error message.
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("groq API returned status %d: %s", resp.StatusCode, string(errBody))
	}

	flusher, _ := w.(interface{ Flush() })

	scanner := bufio.NewScanner(resp.Body)
	// SSE lines can be long; bump the scanner's buffer to be safe.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue // blank line separating SSE events
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}
		if data == "" {
			continue
		}

		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf("unmarshal stream chunk: %w (raw: %s)", err, data)
		}
		if chunk.Error != nil {
			return fmt.Errorf("groq API error (%s): %s", chunk.Error.Type, chunk.Error.Message)
		}

		for _, choice := range chunk.Choices {
			if choice.Delta.Content == "" {
				continue
			}
			if _, werr := io.WriteString(w, choice.Delta.Content); werr != nil {
				return fmt.Errorf("write to output: %w", werr)
			}
			if flusher != nil {
				flusher.Flush()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read stream: %w", err)
	}

	return nil
}
