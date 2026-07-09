package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type BotSender struct {
	BotToken string
	ChatID   string
	Client   *http.Client
}

func NewBotSender(botToken string, chatID string, client *http.Client) *BotSender {
	if client == nil {
		client = http.DefaultClient
	}

	return &BotSender{
		BotToken: botToken,
		ChatID:   chatID,
		Client:   client,
	}
}

func (s *BotSender) Send(ctx context.Context, text string) error {
	if s.BotToken == "" {
		return fmt.Errorf("telegram bot token is required")
	}

	if s.ChatID == "" {
		return fmt.Errorf("telegram chat ID is required")
	}

	payload := sendMessagePayload{
		ChatID:    s.ChatID,
		Text:      text,
		ParseMode: "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal telegram payload: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.BotToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create telegram request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ImapRelay")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram returned status %s", resp.Status)
	}

	return nil
}

type sendMessagePayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}
