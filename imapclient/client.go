package imapclient

import (
	"context"
	"fmt"
	"github.com/FAR747/imaprelay/internal/config"
	"github.com/emersion/go-imap/v2"
	"strings"
)

const maxFetchedBodyBytes = 64 * 1024

func FetchUnseen(ctx context.Context, account config.IMAPConfig, proxyConfig *config.ProxyConfig) ([]Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	client, err := connect(account, proxyConfig)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer client.Close()

	if err := client.Login(account.Username, account.Password).Wait(); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	defer func() {
		_ = client.Logout().Wait()
	}()

	if _, err := client.Select(account.Mailbox, nil).Wait(); err != nil {
		return nil, fmt.Errorf("select mailbox %q: %w", account.Mailbox, err)
	}

	searchData, err := client.UIDSearch(&imap.SearchCriteria{
		NotFlag: []imap.Flag{imap.FlagSeen},
	}, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("search unseen: %w", err)
	}

	uids := searchData.AllUIDs()
	if len(uids) == 0 {
		return []Message{}, nil
	}

	bodySection := &imap.FetchItemBodySection{
		Peek: true,
		Partial: &imap.SectionPartial{
			Offset: 0,
			Size:   maxFetchedBodyBytes,
		},
	}

	fetchOptions := &imap.FetchOptions{
		UID:          true,
		Envelope:     true,
		InternalDate: true,
		BodySection:  []*imap.FetchItemBodySection{bodySection},
	}

	fetchCmd := client.Fetch(imap.UIDSetNum(uids...), fetchOptions)
	defer fetchCmd.Close()

	messages := make([]Message, 0, len(uids))

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		data := fetchCmd.Next()
		if data == nil {
			break
		}

		buf, err := data.Collect()
		if err != nil {
			return nil, fmt.Errorf("collect message: %w", err)
		}

		msg := Message{
			UID:        UID(buf.UID),
			Account:    account.Name,
			Mailbox:    account.Mailbox,
			ReceivedAt: buf.InternalDate,
		}

		if buf.Envelope != nil {
			msg.From = formatAddresses(buf.Envelope.From)
			msg.Title = buf.Envelope.Subject
		}

		body := buf.FindBodySection(bodySection)
		msg.Body = extractPlainText(body)

		messages = append(messages, msg)
	}

	if err := fetchCmd.Close(); err != nil {
		return nil, fmt.Errorf("fetch messages: %w", err)
	}

	return messages, nil
}

func formatAddresses(addresses []imap.Address) string {
	if len(addresses) == 0 {
		return ""
	}

	parts := make([]string, 0, len(addresses))

	for _, address := range addresses {
		email := address.Addr()

		switch {
		case address.Name != "" && email != "":
			parts = append(parts, fmt.Sprintf("%s <%s>", address.Name, email))
		case email != "":
			parts = append(parts, email)
		case address.Name != "":
			parts = append(parts, address.Name)
		}
	}

	return strings.Join(parts, ", ")
}
