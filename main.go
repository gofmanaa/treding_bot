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
	tradingPair = "BTCUSDT"
	// Amount of BTC to buy or sell in each trade
	tradeAmount = 0.01
	// Minimum price difference to consider making a trade
	minPriceDifference = 1.5
)

var totalCryptoAmount float64 = 1.0

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// Set up a connection to the exchange API
	client, err := rpc.NewClient(apiKey, apiSecret)
	if err != nil {
		log.Fatal(err)
	}
	fiatAmount := 0.0
	// Initialize the lastTradePrice variable
	lastTradePrice := 0.0
	// Continuously monitor the market for opportunities to buy and sell
	for {
		// Get the current price of BTC
		ticker, err := client.Ticker(tradingPair)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("Current price of BTC: %v\n", ticker.Price)
		if lastTradePrice == 0 {
			lastTradePrice = ticker.Price
		}
		//fmt.Printf("Abs(%f-%f) >= %f\n", ticker.Price, lastTradePrice, minPriceDifference)
		//fmt.Printf("%f=>%f = %t\n", math.Abs(ticker.Price-lastTradePrice), minPriceDifference, math.Abs(ticker.Price-lastTradePrice) >= minPriceDifference)
		// Check if the price has changed significantly since the last trade
		if math.Abs(ticker.Price-lastTradePrice) >= minPriceDifference {
			// Check if the price has increased or decreased
			if ticker.Price > lastTradePrice {
				// Buy BTC if the price has increased
				err := client.Buy(tradingPair, tradeAmount, ticker.Price)
				if err != nil {
					log.Println(err)
					continue
				}
				fiatAmount -= math.Abs(ticker.Price - lastTradePrice)
				totalCryptoAmount += tradeAmount
				fmt.Println("Bought BTC")
			} else {
				// Sell BTC if the price has decreased
				err := client.Sell(tradingPair, tradeAmount, ticker.Price)
				if err != nil {
					log.Println(err)
					continue
				}
				fiatAmount += math.Abs(ticker.Price - lastTradePrice)
				totalCryptoAmount -= tradeAmount
				fmt.Println("Sold BTC")
			}

			fmt.Printf("totalCryptoAmount: %f\n", totalCryptoAmount)
			fmt.Printf("fiatAmount: %f\n", fiatAmount)
			// Save the last trade price for comparison on the next iteration
			lastTradePrice = ticker.Price
		}

		// Sleep for a short time before checking the price again
		time.Sleep(time.Second)
	}
}
