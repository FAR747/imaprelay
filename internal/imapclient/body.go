package imapclient

import (
	"bytes"
	"io"
	"strings"

	"github.com/emersion/go-message/mail"
)

func extractPlainText(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}

	reader, err := mail.CreateReader(bytes.NewReader(raw))
	if err != nil {
		return strings.TrimSpace(string(raw))
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		switch part.Header.(type) {
		case *mail.InlineHeader:
			data, err := io.ReadAll(part.Body)
			if err != nil {
				return ""
			}

			return strings.TrimSpace(string(data))
		}
	}

	return strings.TrimSpace(string(raw))
}
