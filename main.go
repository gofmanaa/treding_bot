package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofmanaa/trading_bot/rpc"
	"github.com/gofmanaa/trading_bot/trade"
)

const (
	// Replace these values with your own API keys and account information
	apiKey    = "YOUR_API_KEY"
	apiSecret = "YOUR_API_SECRET"

	tradingPair = "BTCUSDT"
	// Amount of BTC to buy or sell in each trade
	tradeAmount = 0.01
	// Minimum price difference to consider making a trade
	minPriceDifference = 0.5

	totalCryptoAmount float64 = 1.0
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	startTime := time.Now()
	// Set up a connection to the exchange API
	tradeClient, err := rpc.NewClient(apiKey, apiSecret)
	if err != nil {
		log.Fatal(err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	trading := trade.NewTade(totalCryptoAmount, tradeAmount, minPriceDifference, tradingPair)

	// Continuously monitor the market for opportunities to buy and sell
	for {
		select {
		case <-sigint:
			fmt.Println("\nWork time:", time.Since(startTime))
			return
		default:
			err := trading.Trade(tradeClient)
			if err != nil {
				log.Println(err)
				fmt.Println("\nWork time:", time.Since(startTime))
				return
			}
		}
		// Sleep for a short time before checking the price again
		time.Sleep(time.Second)
	}

}
