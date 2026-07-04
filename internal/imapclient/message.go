package imapclient

import "time"

type UID uint32

type Message struct {
	UID      UID
	Account  string
	Username string
	Mailbox  string

	From  string
	Title string
	Body  string

	// IMAP internal date
	ReceivedAt time.Time
}
