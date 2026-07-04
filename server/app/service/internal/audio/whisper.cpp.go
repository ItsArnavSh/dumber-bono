package audio

import (
	"bytes"
	"context"
	"encoding/json"
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

func Transcribe(ctx context.Context, wav []byte) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "speech.wav")
	if err != nil {
		return "", err
	}
	if _, err := part.Write(wav); err != nil {
		return "", err
	}

	prompt := `This is a Formula 1 race engineer conversation.

Common terminology:
Copy, DRS, ERS, KERS, DRS detection,
Push Lap, Cooldown Lap, Undercut, Overcut, Dirty Air, Clean Air,
Brake Bias, Differential, Engine Braking, Lift and Coast, Lock Up, Wheelspin,
Understeer, Oversteer, Sector One, Sector Two, Sector Three,
Soft, Medium, Hard, Intermediate, Wet, Tyres, Fuel, Puncture,
Pit Window, Box,  Fastest Lap, Delta Time, Track Limits.

Driver names:
Verstappen, Norris, Piastri, Leclerc, Hamilton, Russell, Sainz, Alonso,
Stroll, Gasly, Ocon, Albon, Tsunoda, Lawson, Hadjar, Antonelli, Bearman, Bortoleto, Hülkenberg.

Teams:
Red Bull Racing, McLaren, Ferrari, Mercedes, Aston Martin,
Alpine, Williams, Haas, Racing Bulls, Sauber.`

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
