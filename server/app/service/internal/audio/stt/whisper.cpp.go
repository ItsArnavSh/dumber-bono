package audio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

type Whisper struct {
	Binary string
	Model  string
}

func newWhisper(binary, model string) *Whisper {
	return &Whisper{
		Binary: binary,
		Model:  model,
	}
}

func (a *STT) transcribeWAV(ctx context.Context) {
	defer a.wg.Done()

	for wav := range a.incoming {
		text, err := transcribe(ctx, wav)
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

func transcribe(ctx context.Context, wav []byte) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "speech.wav")
	if err != nil {
		return "", err
	}
	if _, err := part.Write(wav); err != nil {
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
		"http://127.0.0.1:8088/inference",
		body,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Text, nil
}
