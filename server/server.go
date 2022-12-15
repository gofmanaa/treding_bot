package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofmanaa/trading_bot/types"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// Create a new HTTP router
	router := http.NewServeMux()

	// Handle requests to the /ticker endpoint
	router.HandleFunc("/ticker", func(w http.ResponseWriter, r *http.Request) {
		// Get the trading pair from the request query parameters
		pair := r.URL.Query().Get("pair")
		if pair == "" {
			log.Println("error pair is empty")
			// Return an error if the trading pair is not specified
			http.Error(w, "Trading pair must be specified", http.StatusBadRequest)
			return
		}

		// Get the current price of the trading pair
		price, err := getPrice(pair)
		if err != nil {
			// Return an error if the price could not be retrieved
			log.Printf("error getPrice: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the current price as JSON
		w.Header().Set("Content-Type", "application/json")
		var tiker types.Ticker
		tiker.Price = price
		tiker.Symbol = pair
		json.NewEncoder(w).Encode(&tiker)
	})

	// Handle requests to the /buy and /sell endpoints
	router.HandleFunc("/buy", buyHandler)
	router.HandleFunc("/sell", sellHandler)

	// Start the server on port 8080
	http.ListenAndServe(":8080", router)
}

// buyHandler is the handler for requests to the /buy endpoint
func buyHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request
	err := r.ParseForm()
	if err != nil {
		log.Panicf("error parceForm: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the trading pair, amount, and price from the request form data
	pair := r.PostForm.Get("pair")
	if pair == "" {
		log.Println("error pair is empty")
		// Return an error if the trading pair is not specified
		http.Error(w, "Trading pair must be specified", http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(r.Form.Get("amount"), 64)
	if err != nil {
		// Return an error if the amount is not a valid number
		http.Error(w, "Invalid amount: "+r.Form.Get("amount"), http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseFloat(r.Form.Get("price"), 64)
	if err != nil {
		// Return an error if the price is not a valid number
		http.Error(w, "Invalid price: "+r.Form.Get("price"), http.StatusBadRequest)
		return
	}

	// Check if the request is authenticated
	if !isAuthenticated(r) {
		// Return an error if the request is not authenticated
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Execute the buy on the exchange
	err = executeBuy(pair, amount, price)
	if err != nil {
		// Return an error if the buy could not be executed
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response if the buy was executed successfully
	w.WriteHeader(http.StatusOK)
}

// sellHandler is the handler for requests to the /sell endpoint
func sellHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the trading pair, amount, and price from the request form data
	pair := r.Form.Get("pair")
	if pair == "" {
		log.Println("error pair is empty")
		// Return an error if the trading pair is not specified
		http.Error(w, "Trading pair must be specified", http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(r.Form.Get("amount"), 64)
	if err != nil {
		// Return an error if the amount is not a valid number
		http.Error(w, "Invalid amount: "+r.Form.Get("amount"), http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseFloat(r.Form.Get("price"), 64)
	if err != nil {
		// Return an error if the price is not a valid number
		http.Error(w, "Invalid price: "+r.Form.Get("price"), http.StatusBadRequest)
		return
	}

	// Check if the request is authenticated
	if !isAuthenticated(r) {
		// Return an error if the request is not authenticated
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Execute the sell on the exchange
	err = executeSell(pair, amount, price)
	if err != nil {
		// Return an error if the sell could not be executed
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response if the sell was executed successfully
	w.WriteHeader(http.StatusOK)
}

// getPrice retrieves the current price of a trading pair
func getPrice(pair string) (float64, error) {
	// Set up the HTTP request to the exchange API {"symbol":"BTCUSDT","price":"17689.80000000"}
	req, err := http.NewRequest("GET", "https://api.binance.com/api/v3/ticker/price", nil)
	if err != nil {
		return 0, err
	}

	// Add the trading pair to the request query parameters
	q := req.URL.Query()
	q.Add("symbol", pair)
	req.URL.RawQuery = q.Encode()

	// Send the request and get the response from the exchange API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Parse the response from the exchange API
	var tiket struct {
		Price  string `json:"price"`
		Symbol string `json:"symbol"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tiket)
	if err != nil {
		log.Printf("error, can't decode json:%s", err)
		return 0, err
	}

	return strconv.ParseFloat(tiket.Price, 64)
}

// isAuthenticated checks if a request is authenticated
func isAuthenticated(r *http.Request) bool {
	// Get the API key and secret from the request headers
	apiKey := r.Header.Get("X-API-Key")
	apiSecret := r.Header.Get("X-API-Secret")

	// Check if the API key and secret are set
	if apiKey == "" || apiSecret == "" {
		return false
	}

	// Check if the API key and secret are valid
	if !isValidAPIKeyAndSecret(apiKey, apiSecret) {
		return false
	}

	return true
}

// isValidAPIKeyAndSecret checks if an API key and secret are valid
func isValidAPIKeyAndSecret(apiKey, apiSecret string) bool {
	// Check if the API key and secret match the expected values
	if apiKey == "YOUR_API_KEY" && apiSecret == "YOUR_API_SECRET" {
		return true
	}

	return false
}

// executeSell executes a sell on the exchange
func executeSell(pair string, amount, price float64) error {
	// TODO implement Sell
	return nil
	// Set up the HTTP request to the exchange API
	req, err := http.NewRequest("POST", "https://api.exchange.com/sell", nil)
	if err != nil {
		return err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the request form data
	data := url.Values{}
	data.Set("pair", pair)
	data.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	data.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	req.PostForm = data

	// Send the request to the exchange API
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// executeBuy executes a buy on the exchange
func executeBuy(pair string, amount, price float64) error {
	// TODO implement Buy
	return nil
	// Set up the HTTP request to the exchange API
	req, err := http.NewRequest("POST", "https://api.exchange.com/buy", nil)
	if err != nil {
		return err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the request form data
	data := url.Values{}
	data.Set("pair", pair)
	data.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	data.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	req.PostForm = data

	// Send the request to the exchange API
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
