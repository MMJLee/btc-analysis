package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"mjlee.dev/btc-analysis/util"
)

type APIClient struct {
	*http.Client
}
type Candlestick struct {
	Start  util.StringUInt64  `json:"Start"`
	Low    util.StringUInt32  `json:"Low"`
	High   util.StringUInt32  `json:"High"`
	Open   util.StringUInt32  `json:"Open"`
	Close  util.StringUInt32  `json:"Close"`
	Volume util.StringFloat32 `json:"Volume"`
}

func (a APIClient) NewRequest(method string, url url.URL, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		log.Fatal(err)
	}

	jwt, err := BuildJWT(method, url.Host, url.Path)
	if err != nil {
		log.Fatal("error building jwt: %v", err)
	}

	bearer := "Bearer " + jwt
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (a *APIClient) GetCandlesticks(product_id string, start int64, end int64, granularity string, limit int64) ([]Candlestick, error) {
	candlestick_url := util.GetProductCandlestickUrl(product_id, start, end, granularity, limit)
	fmt.Printf("Encoded URL is %q\n", candlestick_url.String())
	req, err := a.NewRequest("GET", candlestick_url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Request succeeded!")
	} else {
		fmt.Println("Request failed with status:", resp.StatusCode)
		panic(body)
	}

	var candlesticks map[string][]Candlestick
	err = json.Unmarshal([]byte(body), &candlesticks)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	return candlesticks["candles"], nil
}
