package target

import (
	"fmt"
	"github.com/FAR747/imaprelay/internal/imapclient"
	"strings"
)

const MaxPushLength = 1500
const MaxTitleLength = 128
const MaxFromLength = 128

type FormatType string

const (
	FormatDiscord  FormatType = "discord"
	FormatTelegram FormatType = "telegram"
)

// Discord
const DefaultHeaderDiscord = "**NEW MESSAGE**\n> Account: `{account}` (`{username}`)\n> From: `{from}`\n> Title: `{title}`\n> Received: `{received}`\n\n"

// Telegram
const DefaultHeaderTelegram = "*NEW MESSAGE*\nAccount: `{account}` (`{username}`)\nFrom: `{from}`\nTitle: `{title}`\nReceived: `{received}`\n\n"

const DefaultTimeFormat = "02.01.2006 15:04"

func FormatMessage(msg imapclient.Message, formatTypes ...FormatType) (string, error) {
	if len(formatTypes) > 1 {
		return "", fmt.Errorf("only one format type is allowed")
	}

	formatType := FormatDiscord
	rawHeader := DefaultHeaderDiscord
	if len(formatTypes) == 1 {
		formatType = formatTypes[0]
	}

	switch formatType {
	case FormatDiscord:
		rawHeader = DefaultHeaderDiscord

	case FormatTelegram:
		rawHeader = DefaultHeaderTelegram

	default:
		return "", fmt.Errorf("unknown format type %q", formatType)
	}

	header := formatTags(rawHeader, msg)

	bodyLimit := MaxPushLength - len([]rune(header))
	if bodyLimit < 0 {
		bodyLimit = 0
	}
	body := truncateRunes(msg.Body, bodyLimit)

	return header + body, nil
}

func formatTags(text string, msg imapclient.Message) string {
	values := map[string]string{
		"{account}":  msg.Account,
		"{username}": msg.Username,
		"{from}":     truncateRunes(msg.From, MaxFromLength),
		"{title}":    truncateRunes(msg.Title, MaxTitleLength),
		"{received}": msg.ReceivedAt.Format(DefaultTimeFormat),
	}
	result := text
	for key, value := range values {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)

	if len(runes) <= max {
		return s
	}

	if max <= 3 {
		return string(runes[:max])
	}

	return string(runes[:max-3]) + "..."
}
