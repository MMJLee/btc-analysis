package client

import (
	"io"
	"log"
	"net/http"
	"net/url"
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
		log.Panicf("Error: Client-NewRequest-NewRequest: %v", err)
	}

	jwt, err := BuildJWT(method, url.Host, url.Path)
	if err != nil {
		log.Panicf("Error: Client-NewRequest-BuildJWT: %v", err)
	}

	bearer := "Bearer " + jwt
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
