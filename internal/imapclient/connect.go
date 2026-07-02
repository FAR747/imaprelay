package imapclient

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/FAR747/imaprelay/internal/config"
	"github.com/emersion/go-imap/v2/imapclient"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
)

func connect(account config.IMAPConfig, proxyConfig *config.ProxyConfig) (*imapclient.Client, error) {
	address := fmt.Sprintf("%s:%d", account.Host, account.Port)

	conn, err := dial(address, proxyConfig)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	switch account.Security {
	case "tls":
		tlsConn := tls.Client(conn, &tls.Config{
			ServerName: account.Host,
		})

		if err := tlsConn.Handshake(); err != nil {
			_ = tlsConn.Close()
			return nil, fmt.Errorf("tls handshake: %w", err)
		}

		return imapclient.New(tlsConn, nil), nil

	case "starttls":
		client, err := imapclient.NewStartTLS(conn, nil)
		if err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("starttls: %w", err)
		}

		return client, nil

	case "none":
		return imapclient.New(conn, nil), nil

	default:
		_ = conn.Close()
		return nil, fmt.Errorf("unsupported security mode %q", account.Security)
	}
}

func dial(address string, proxyConfig *config.ProxyConfig) (net.Conn, error) {
	if proxyConfig == nil {
		return net.Dial("tcp", address)
	}

	switch proxyConfig.Type {
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", proxyConfig.Address, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("create socks5 dialer: %w", err)
		}

		return dialer.Dial("tcp", address)

	case "http":
		return dialHTTPProxy(proxyConfig.Address, address)

	default:
		return nil, fmt.Errorf("unsupported proxy type %q", proxyConfig.Type)
	}
}

func dialHTTPProxy(proxyAddress string, targetAddress string) (net.Conn, error) {
	conn, err := net.Dial("tcp", proxyAddress)
	if err != nil {
		return nil, fmt.Errorf("dial http proxy: %w", err)
	}

	request := fmt.Sprintf(
		"CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n",
		targetAddress,
		targetAddress,
	)

	if _, err := conn.Write([]byte(request)); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("send CONNECT request: %w", err)
	}

	reader := bufio.NewReader(conn)

	response, err := http.ReadResponse(reader, nil)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("read CONNECT response: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		_ = conn.Close()
		return nil, fmt.Errorf("http proxy CONNECT failed: %s", response.Status)
	}

	return conn, nil
}
