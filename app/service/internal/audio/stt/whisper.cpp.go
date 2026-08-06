package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Whisper struct {
	ApiKey string
	Model  string
}

func newWhisper() *Whisper {
	model := os.Getenv("GROQ_STT_MODEL")
	if model == "" {
		model = "whisper-large-v3" // Default fallback model
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Println("[WHISPER WARN] GROQ_API_KEY environment variable is empty!")
	}

	return &Whisper{
		ApiKey: apiKey,
		Model:  model,
	}
}

func (a *STT) transcribeWAV(ctx context.Context) {
	defer a.wg.Done()

	log.Println("[STT WORKER] Starting transcribeWAV worker loop...")
	whisperClient := newWhisper()

	for wav := range a.incoming {
		log.Printf("[STT WORKER] Received WAV chunk (%d bytes). Sending to Groq...", len(wav))

		text, err := whisperClient.transcribe(ctx, wav)
		if err != nil {
			log.Printf("[STT WORKER ERROR] Groq API call failed: %v", err)
			// Log error but CONTINUE processing next incoming chunks instead of exiting loop
			continue
		}

		log.Printf("[STT WORKER OUTPUT] Transcribed text: %q", text)
		a.StrChunk.WriteString(text)
	}

	log.Println("[STT WORKER] 'a.incoming' channel closed. Worker loop finished successfully.")
	a.done <- nil
}

func (w *Whisper) transcribe(ctx context.Context, wav []byte) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "speech.wav")
	if err != nil {
		return "", fmt.Errorf("failed to create multipart form file: %w", err)
	}
	if _, err := part.Write(wav); err != nil {
		return "", fmt.Errorf("failed to write WAV bytes to form: %w", err)
	}

	if err := writer.WriteField("model", w.Model); err != nil {
		return "", fmt.Errorf("failed to write model field: %w", err)
	}

	// Domain context prompt to guide Whisper jargon recognition
	prompt := `F1 radio: DRS, ERS, Undercut, Overcut, Lift and Coast, Lock Up, Understeer, Oversteer, Delta, Box, Track Limits.`
	if err := writer.WriteField("prompt", prompt); err != nil {
		return "", fmt.Errorf("failed to write prompt field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.groq.com/openai/v1/audio/transcriptions",
		body,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.ApiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("network error during Groq API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return "", fmt.Errorf("groq API HTTP %d error: %v", resp.StatusCode, errResp)
	}

	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response JSON: %w", err)
	}

	return result.Text, nil
}
