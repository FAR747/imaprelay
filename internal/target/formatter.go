package target

import (
	"fmt"
	"github.com/FAR747/imaprelay/internal/imapclient"
	"strings"
)

const MaxPushLength = 1500
const MaxAccountLength = 64
const MaxUsernameLength = 128
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
const DefaultHeaderTelegram = "<b>NEW MESSAGE</b>\nAccount: <code>{account}</code> (<code>{username}</code>)\nFrom: <code>{from}</code>\nTitle: <code>{title}</code>\nReceived: <code>{received}</code>\n\n"

const DefaultTimeFormat = "02.01.2006 15:04"

func FormatMessage(msg imapclient.Message, formatTypes ...FormatType) (string, error) {
	if len(formatTypes) > 1 {
		return "", fmt.Errorf("only one format type is allowed")
	}

	formatType := FormatDiscord
	if len(formatTypes) == 1 {
		formatType = formatTypes[0]
	}

	switch formatType {
	case FormatDiscord:
		return renderMessage(msg, DefaultHeaderDiscord, escapeDiscord), nil

	case FormatTelegram:
		return renderMessage(msg, DefaultHeaderTelegram, escapeHTML), nil

	default:
		return "", fmt.Errorf("unknown format type %q", formatType)
	}
}

func renderMessage(msg imapclient.Message, rawHeader string, escape func(string) string) string {
	header := formatTags(rawHeader, msg, escape)

	if len([]rune(header)) >= MaxPushLength {
		return truncateRunes(header, MaxPushLength)
	}

	bodyLimit := MaxPushLength - len([]rune(header))
	body := truncateRunes(escape(msg.Body), bodyLimit)

	return header + body
}

func formatTags(text string, msg imapclient.Message, escape func(string) string) string {
	values := map[string]string{
		"{account}":  escape(truncateRunes(msg.Account, MaxAccountLength)),
		"{username}": escape(truncateRunes(msg.Username, MaxUsernameLength)),
		"{from}":     escape(truncateRunes(msg.From, MaxFromLength)),
		"{title}":    escape(truncateRunes(msg.Title, MaxTitleLength)),
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

func escapeDiscord(s string) string {
	return strings.ReplaceAll(s, "`", "'")
}

func escapeHTML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return replacer.Replace(s)
}
