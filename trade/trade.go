package trade

import (
	"fmt"
	"log"
	"math"

	"github.com/gofmanaa/trading_bot/rpc"
	"github.com/gofmanaa/trading_bot/types"
)

const pricesHistoryLen = 20

type Trade struct {
	initialCryptoAmount float64 //initial crypto amount
	currentCryptoAmount float64 //how many crypto do we have at the moment
	tradingPair         string  //treding pair
	minPriceDifference  float64 //minimum profit

	tradeAmount float64 // hpw many crypto we wont buy/sell

	currentFiatAmount float64 // current fiat profit
	lastTradePrice    float64 // last crypto price

	pricesHistory []float64 // prices history for prediction

	generations int64 // trade tick
}

func NewTade(initialCryptoAmount float64, tradeAmount float64, minPriceDifference float64, tradingPair string) *Trade {
	return &Trade{
		initialCryptoAmount: initialCryptoAmount,
		minPriceDifference:  minPriceDifference,
		tradingPair:         tradingPair,
		tradeAmount:         tradeAmount,
		currentCryptoAmount: initialCryptoAmount,
		pricesHistory:       make([]float64, 0, pricesHistoryLen),
	}
}

func (t *Trade) Trade(tradeClient rpc.Client) error {

	if stopTradeRule(t.initialCryptoAmount, t.currentCryptoAmount) {
		return fmt.Errorf("stop trading curent crypto amount %f", t.currentCryptoAmount)
	}
	var ticker *types.Ticker
	var err error

	// Get the current price of BTC
	ticker, err = tradeClient.Ticker(t.tradingPair)
	if err != nil {
		log.Println(err)
		return err
	}

	if t.lastTradePrice == 0 {
		t.lastTradePrice = ticker.Price
	}

	// predict := t.predictTrend(ticker.Price)

	// if len(t.pricesHistory) >= pricesHistoryLen && predict > 0 {
	// 	fmt.Printf("Price Prediction: %f\n", predict)
	// 	t.lastTradePrice = predict
	// }

	deltaPrice := math.Abs(ticker.Price - t.lastTradePrice)
	// Check if the price has changed significantly since the last trade
	if deltaPrice >= t.minPriceDifference {
		fmt.Printf("Current price of BTC: %v\n", ticker.Price)
		// Check if the price has increased or decreased
		if ticker.Price > t.lastTradePrice {
			// Buy BTC if the price has increased
			err := tradeClient.Buy(t.tradingPair, t.tradeAmount, ticker.Price)
			if err != nil {
				log.Println(err)
				return err
			}
			t.currentFiatAmount -= deltaPrice
			t.currentCryptoAmount += t.tradeAmount
			fmt.Println("Bought BTC")
		} else {
			// Sell BTC if the price has decreased
			err := tradeClient.Sell(t.tradingPair, t.tradeAmount, ticker.Price)
			if err != nil {
				log.Println(err)
				return err
			}
			t.currentFiatAmount += deltaPrice
			t.currentCryptoAmount -= t.tradeAmount
			fmt.Println("Sold BTC")
		}

		fmt.Printf("cryptoAmount: %f\n", t.currentCryptoAmount)
		fmt.Printf("fiatAmount: %f\n", t.currentFiatAmount)
		// Save the last trade price for comparison on the next iteration
		t.lastTradePrice = ticker.Price
		t.generations++

	}
	return nil
}

func stopTradeRule(totalCryptoAmount, cryptoAmount float64) bool {
	return totalCryptoAmount/2 >= cryptoAmount
}

func (t *Trade) predictTrend(currentPrice float64) float64 {
	t.pricesHistory = append(t.pricesHistory, currentPrice)
	if len(t.pricesHistory) < pricesHistoryLen {
		fmt.Println("Not enough data to predict trend")
		return -1
	}
	if len(t.pricesHistory) > pricesHistoryLen {
		t.pricesHistory = t.pricesHistory[1:len(t.pricesHistory)]
	}

	// Calculate the moving average of the past prices
	var sum float64
	for i := 0; i < pricesHistoryLen; i++ {
		sum += t.pricesHistory[len(t.pricesHistory)-i-1]
	}
	movingAverage := sum / pricesHistoryLen
	if currentPrice > movingAverage {
		fmt.Println("Upward trend")
	} else if currentPrice < movingAverage {
		fmt.Println("Downward trend")
	} else {
		fmt.Println("No trend")
	}
	return movingAverage
}
