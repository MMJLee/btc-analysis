package client

import (
	"io"
	"net/http"
	"net/url"

	"github.com/mmjlee/btc-analysis/internal/util"
)

type APIClient struct {
	*http.Client
}

func GetNewAPIClient() APIClient {
	return APIClient{&http.Client{}}
}

func (a APIClient) NewRequest(method string, url url.URL, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return req, util.WrappedError{Err: err, Message: "Client-NewRequest-NewRequest"}
	}

	jwt, err := BuildJWT(method, url.Host, url.Path)
	if err != nil {
		return req, util.WrappedError{Err: err, Message: "Client-NewRequest-BuildJWT"}
	}

	bearer := "Bearer " + jwt
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
