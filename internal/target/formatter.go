package target

import (
	"github.com/FAR747/imaprelay/internal/imapclient"
	"strings"
)

const MaxPushLength = 1500
const MaxTitleLength = 128
const MaxFromLength = 128

const DefaultHeader = "Account: {account} ({username})\\nFrom: {from}\\nTitle: {title}\\nReceived: {received}\\n\\n"
const DefaultTimeFormat = "02.01.2006 15:04"

func FormatMessage(msg imapclient.Message) string {
	header := formatTags(DefaultHeader, msg)

	bodyLimit := MaxPushLength - len([]rune(header))
	if bodyLimit < 0 {
		bodyLimit = 0
	}
	body := truncateRunes(msg.Body, bodyLimit)

	return header + body
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
