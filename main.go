package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"mjlee.dev/btc-analysis/api"
)

const dollar_to_cents uint8 = 100

// money represented in cents
type Candlestick struct {
	Start  StringUInt64  `json:"Start"`
	Low    StringUInt32  `json:"Low"`
	High   StringUInt32  `json:"High"`
	Open   StringUInt32  `json:"Open"`
	Close  StringUInt32  `json:"Close"`
	Volume StringFloat32 `json:"Volume"`
}

type StringUInt32 uint32

func (s *StringUInt32) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return err
	}
	*s = StringUInt32(uint32(value * float64(dollar_to_cents)))
	return nil
}

type StringUInt64 uint64

func (s *StringUInt64) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*s = StringUInt64(uint64(value))
	return nil
}

type StringFloat32 float32

func (s *StringFloat32) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*s = StringFloat32(value)
	return nil
}

func getRequest() (*http.Request, error) {
	jwt, err := api.BuildJWT()
	if err != nil {
		log.Fatal("error building jwt: %v", err)
	}
	fmt.Println(jwt)

	bearer := "Bearer " + jwt
	candlestick_url, err := url.Parse("https://api.coinbase.com/api/v3/brokerage/products/BTC-USD/candles")
	if err != nil {
		return nil, fmt.Errorf("jwt: %w", err)
	}
	// // Query params
	query := candlestick_url.Query()
	query.Set("start", "1735711200")
	query.Set("end", "1735711320")
	query.Set("granularity", "ONE_MINUTE")
	query.Set("limit", "2")
	candlestick_url.RawQuery = query.Encode()
	fmt.Printf("Encoded URL is %q\n", candlestick_url.String())
	req, err := http.NewRequest("GET", candlestick_url.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
func getCandlesticks() ([]Candlestick, error) {
	req, err := getRequest()
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
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
	fmt.Printf("%+v\n", body)

	var candlesticks map[string][]Candlestick
	err = json.Unmarshal([]byte(body), &candlesticks)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	return candlesticks["candles"], nil
}

func main() {
	candlesticks, err := getCandlesticks()
	if err != nil {
		log.Fatal(err)
	}
	for _, candlestick := range candlesticks {
		fmt.Printf("%+v\n", candlestick)
	}
}
