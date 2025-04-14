package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type CoinbaseClient struct {
	*http.Client
}

func NewCoinbaseClient() CoinbaseClient {
	return CoinbaseClient{&http.Client{Timeout: time.Duration(5) * time.Second}}
}

func NewRequest(method string, url url.URL, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return req, fmt.Errorf("NewRequest-%w", err)
	}

	jwt, err := BuildJWT(method, url.Host, url.Path)
	if err != nil {
		return req, fmt.Errorf("NewRequest-%w", err)
	}

	bearer := "Bearer " + jwt
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
