package openai

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Client is an OpenAI API client for audio transcription.
type Client struct {
	APIKey  string
	BaseURL string
}

// Transcribe sends a WAV audio file to the OpenAI Whisper API and returns the transcribed text.
func (c *Client) Transcribe(ctx context.Context, wav io.Reader) (string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "audio.wav")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, wav); err != nil {
		return "", err
	}
	w.WriteField("model", "whisper-1")
	if err := w.Close(); err != nil {
		return "", err
	}
	baseURL := cmp.Or(c.BaseURL, "https://api.openai.com/v1")
	transcriptionsURL, err := url.JoinPath(baseURL, "/audio/transcriptions")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", transcriptionsURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read error body: %w", err)
		}
		return "", fmt.Errorf("openai: %s: %s", resp.Status, string(body))
	}
	var out struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Text, nil
}
