package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// StreamLLM sends a prompt to Groq with streaming enabled and pushes each
// raw text delta onto tokens as it arrives. It blocks until the stream ends
// or ctx is cancelled, then closes tokens.
func StreamLLM(ctx context.Context, systemPrompt, userPrompt, model string, tokens chan<- string) error {
	defer close(tokens)

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GROQ_API_KEY environment variable is not set")
	}

	if model == "" {
		model = "openai/gpt-oss-20b"
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("groq API returned status %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // allow longer SSE lines

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue // SSE keep-alive/blank separator lines
		}

		data, ok := strings.CutPrefix(line, "data: ")
		if !ok {
			continue
		}

		if data == "[DONE]" {
			return nil
		}

		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf("unmarshal stream chunk: %w (data=%s)", err, data)
		}

		if chunk.Error != nil {
			return fmt.Errorf("groq API error (%s): %s", chunk.Error.Type, chunk.Error.Message)
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		if delta := chunk.Choices[0].Delta.Content; delta != "" {
			select {
			case tokens <- delta:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read stream: %w", err)
	}

	return nil
}
