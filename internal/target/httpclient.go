package target

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/FAR747/imaprelay/internal/config"
	"golang.org/x/net/proxy"
)

func NewHTTPClient(proxyConfig *config.ProxyConfig) (*http.Client, error) {
	if proxyConfig == nil {
		return http.DefaultClient, nil
	}

	switch proxyConfig.Type {
	case "http":
		proxyURL := &url.URL{
			Scheme: "http",
			Host:   proxyConfig.Address,
		}

		return &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}, nil

	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", proxyConfig.Address, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("create socks5 dialer: %w", err)
		}

		return &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network string, address string) (net.Conn, error) {
					return dialer.Dial(network, address)
				},
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported proxy type %q", proxyConfig.Type)
	}
}
