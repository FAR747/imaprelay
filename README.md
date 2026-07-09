# ImapRelay

ImapRelay is a small self-hosted daemon for forwarding unread IMAP emails to Discord, Telegram, etc.

It checks configured mailboxes, sends a short notification for each unread email, and marks the email as read after delivery.

Built as a simple Go binary with minimal setup.