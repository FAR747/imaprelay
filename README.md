# ImapRelay

[![GitHub Release](https://img.shields.io/github/v/release/far747/imaprelay)](https://github.com/FAR747/imaprelay/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/far747/imaprelay)](go.mod)
[![GitHub License](https://img.shields.io/github/license/far747/imaprelay)](LICENSE)


ImapRelay is a small self-hosted daemon for forwarding unread IMAP emails to Discord, Telegram, etc.

It checks configured mailboxes, sends a short notification for each unread email, and marks the email as read after delivery.

Built as a simple Go binary with minimal setup.

## Features

- Multiple IMAP accounts and mailboxes
- Discord webhook and Telegram bot targets
- SOCKS5 and HTTP proxy support
- TLS and STARTTLS connections


## Installation

Download the latest build for your platform from the [Releases page](https://github.com/FAR747/imaprelay/releases/latest).

Download the [example configuration file](https://github.com/FAR747/imaprelay/blob/main/config.example.yaml), rename it to `config.yaml`, and place it in the same directory as the binary.

Edit the configuration and run ImapRelay.


## Configuration

ImapRelay reads its configuration from `config.yaml`.

Use [`config.example.yaml`](https://github.com/FAR747/imaprelay/blob/main/config.example.yaml) as a starting point.

The configuration contains three main sections:

- `targets` - Discord and Telegram destinations
- `imaps` - mailboxes to check
- `proxy` - optional HTTP or SOCKS5 proxy

Each IMAP entry may define its own `targets`. If `targets` is omitted, ImapRelay uses all targets marked with `default: true`.

Supported IMAP security modes:

- `tls`
- `starttls`
- `none`

Typical combinations are port `993` with `tls`, or port `143` with `starttls`.

The optional `proxy` section applies to both IMAP connections and notification requests.

Secrets may be provided through environment variables using `${VARIABLE_NAME}`. ImapRelay also loads an optional `.env` file located next to `config.yaml`. System environment variables take priority over values from `.env`.


## Building from source

ImapRelay requires Go `1.26` or newer.

Clone the repository and build the binary:

```bash
git clone https://github.com/FAR747/imaprelay.git
cd imaprelay
```

### Windows

```powershell
go build -ldflags "-X main.version=selfbuild" -o imaprelay.exe .
```

### Linux

```bash
go build -ldflags "-X main.version=selfbuild" -o imaprelay .
```


## Feedback

Bug reports and feature requests are welcome in [GitHub Issues](https://github.com/FAR747/imaprelay/issues).