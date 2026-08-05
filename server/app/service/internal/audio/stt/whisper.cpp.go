package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
		model = "whisper-large-v3" // Default fallback model if not specified
	}

	return &Whisper{
		ApiKey: os.Getenv("GROQ_API_KEY"),
		Model:  model,
	}
}

func (a *STT) transcribeWAV(ctx context.Context) {
	defer a.wg.Done()

	// Initialize Whisper client using the updated constructor
	whisperClient := newWhisper()

	for wav := range a.incoming {
		text, err := whisperClient.transcribe(ctx, wav)
		if err != nil {
			fmt.Println(err)
			a.done <- err // Send error to done channel
			return
		}
		a.StrChunk.WriteString(text)
	}
	// All chunks processed successfully
	a.done <- nil
}

func (w *Whisper) transcribe(ctx context.Context, wav []byte) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "speech.wav")
	if err != nil {
		return "", err
	}
	if _, err := part.Write(wav); err != nil {
		return "", err
	}

	// Required form field for Groq API model specification
	if err := writer.WriteField("model", w.Model); err != nil {
		return "", err
	}

	prompt := `F1 radio: DRS, ERS, Undercut, Overcut, Lift and Coast, Lock Up, Understeer, Oversteer, Delta, Box, Track Limits. Verstappen, Norris, Piastri, Leclerc, Hamilton, Russell, Sainz, Alonso, Stroll, Gasly, Ocon, Albon, Tsunoda, Lawson, Hadjar, Antonelli, Bearman, Bortoleto, Hülkenberg. Red Bull, McLaren, Ferrari, Mercedes, Aston Martin, Alpine, Williams, Haas, Racing Bulls, Sauber.`

	if err := writer.WriteField("prompt", prompt); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.groq.com/openai/v1/audio/transcriptions",
		body,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.ApiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return "", fmt.Errorf("groq api error (status %d): %v", resp.StatusCode, errResp)
	}

	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Text, nil
}
