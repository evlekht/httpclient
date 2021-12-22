package httpclient

import (
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	baseURL url.URL
	client  http.Client
}

type KV []struct {
	Key   string
	Value string
}

type unmarshallable interface {
	UnmarshalJSON([]byte) error
}

type marshallable interface {
	MarshalJSON() ([]byte, error)
}

func New(scheme, host, baseUrl string, timeout, idleTimeout time.Duration) *HTTPClient {
	return &HTTPClient{
		baseURL: url.URL{
			Scheme: scheme,
			Host:   host,
			Path:   baseUrl,
		},
		client: http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				IdleConnTimeout: idleTimeout,
			},
		},
	}
}
