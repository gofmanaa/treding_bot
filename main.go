package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/gofmanaa/trading_bot/rpc"
)

const (
	// Replace these values with your own API keys and account information
	apiKey      = "YOUR_API_KEY"
	apiSecret   = "YOUR_API_SECRET"
	tradingPair = "BTC-USD"
	// Amount of BTC to buy or sell in each trade
	tradeAmount = 0.01
	// Minimum price difference to consider making a trade
	minPriceDifference = 0.01
)

func main() {
	// Set up a connection to the exchange API
	client, err := rpc.NewClient(apiKey, apiSecret)
	if err != nil {
		log.Fatal(err)
	}

	// Continuously monitor the market for opportunities to buy and sell
	for {
		// Get the current price of BTC
		ticker, err := client.Ticker(tradingPair)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("Current price of BTC: %v\n", ticker.Last)

		// Initialize the lastTradePrice variable
		lastTradePrice := 0.0

		// Check if the price has changed significantly since the last trade
		if math.Abs(ticker.Last-lastTradePrice) >= minPriceDifference {
			// Check if the price has increased or decreased
			if ticker.Last > lastTradePrice {
				// Buy BTC if the price has increased
				err := client.Buy(tradingPair, tradeAmount, ticker.Last)
				if err != nil {
					log.Println(err)
					continue
				}
				fmt.Println("Bought BTC")
			} else {
				// Sell BTC if the price has decreased
				err := client.Sell(tradingPair, tradeAmount, ticker.Last)
				if err != nil {
					log.Println(err)
					continue
				}
				fmt.Println("Sold BTC")
			}

			// Save the last trade price for comparison on the next iteration
			lastTradePrice = ticker.Last
		}

		// Sleep for a short time before checking the price again
		time.Sleep(time.Second)
	}
}
