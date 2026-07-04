package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhookSender struct {
	WebhookURL string
	Client     *http.Client
}

func NewWebhookSender(webhookURL string, client *http.Client) *WebhookSender {
	if client == nil {
		client = http.DefaultClient
	}

	return &WebhookSender{
		WebhookURL: webhookURL,
		Client:     client,
	}
}

func (s *WebhookSender) Send(ctx context.Context, text string) error {
	if s.WebhookURL == "" {
		return fmt.Errorf("discord webhook URL is required")
	}

	client := s.Client
	if client == nil {
		client = http.DefaultClient
	}

	payload := webhookPayload{
		Content: text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal discord webhook payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create discord webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ImapRelay")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned status %s", resp.Status)
	}

	return nil
}

type webhookPayload struct {
	Content string `json:"content"`
}
