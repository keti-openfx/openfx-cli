package rest

import (
	"net"
	"net/http"
	"time"
)

// MakeHTTPClient makes a HTTP client with good defaults for timeouts.
func MakeHTTPClient(timeout *time.Duration) http.Client {
	if timeout != nil {
		return http.Client{
			Timeout: *timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout: *timeout,
					// KeepAlive: 0,
				}).DialContext,
				// MaxIdleConns:          1,
				// DisableKeepAlives:     true,
				IdleConnTimeout:       120 * time.Millisecond,
				ExpectContinueTimeout: 1500 * time.Millisecond,
			},
		}
	}

	// This should be used for openfx-cli invoke etc.
	return http.Client{}
}
