package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofmanaa/trading_bot/types"
)

func NewClient(apiKey string, apiSecret string) (Client, error) {
	return &MyRPCClient{apiKey: apiKey, apiSecret: apiSecret}, nil
}

// .rpc
type Client interface {
	// Methods for communicating with the exchange API
	Ticker(pair string) (*types.Ticker, error)
	Buy(pair string, amount, price float64) error
	Sell(pair string, amount, price float64) error
	// Additional methods as needed...
}

type MyRPCClient struct {
	// Fields for storing API keys and other necessary information
	apiKey    string
	apiSecret string
	// Additional fields as needed...
}

// Implement the Ticker() method of the rpc.Client interface
func (c *MyRPCClient) Ticker(pair string) (*types.Ticker, error) {
	// Set up the HTTP request to the exchange API
	req, err := http.NewRequest("GET", "http://localhost:8080/ticker", nil)
	if err != nil {
		return nil, err
	}

	// Add the API key and secret to the request as headers
	req.Header.Add("X-API-Key", c.apiKey)
	req.Header.Add("X-API-Secret", c.apiSecret)

	// Add the trading pair to the request as a query parameter
	q := req.URL.Query()
	q.Add("pair", pair)
	req.URL.RawQuery = q.Encode()

	// Send the request and get the response from the exchange API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body as JSON to get the current price of the trading pair
	var ticker types.Ticker
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		log.Printf("error, can't unmarshal tiker: %s", err)
		return nil, err
	}

	// Return the current price and any error that occurred
	return &ticker, nil
}

// Implement the Buy() and Sell() methods of the rpc.Client interface
func (c *MyRPCClient) Buy(pair string, amount, price float64) error {
	// Set up the HTTP request to the exchange API
	req, err := http.NewRequest("POST", "http://localhost:8080/buy", nil)
	if err != nil {
		return err
	}

	// Add the API key and secret to the request as headers
	req.Header.Add("X-API-Key", c.apiKey)
	req.Header.Add("X-API-Secret", c.apiSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Add the trading pair, amount, and price to the request as form data
	f := url.Values{}
	f.Add("pair", pair)
	f.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	f.Add("price", strconv.FormatFloat(price, 'f', -1, 64))
	req.Body = ioutil.NopCloser(strings.NewReader(f.Encode()))

	// Send the request and get the response from the exchange API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code to see if the buy was successful
	if resp.StatusCode != http.StatusOK {
		// Read the response body and return an error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to buy %v %v at %v: %s", amount, pair, price, body)
	}

	// Return nil if the buy was successful
	return nil
}

func (c *MyRPCClient) Sell(pair string, amount, price float64) error {
	// Set up the HTTP request to the exchange API
	req, err := http.NewRequest("POST", "http://localhost:8080/sell", nil)
	if err != nil {
		return err
	}

	// Add the API key and secret to the request as headers
	req.Header.Add("X-API-Key", c.apiKey)
	req.Header.Add("X-API-Secret", c.apiSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Add the trading pair, amount, and price to the request as form data
	f := url.Values{}
	f.Add("pair", pair)
	f.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	f.Add("price", strconv.FormatFloat(price, 'f', -1, 64))
	req.Body = ioutil.NopCloser(strings.NewReader(f.Encode()))

	// Send the request and get the response from the exchange API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check the response status code to see if the sell was successful
	if resp.StatusCode != http.StatusOK {
		// Read the response body and return an error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to sell %v %v at %v: %s", amount, pair, price, body)
	}

	// Return nil if the sell was successful
	return nil
}
