package types

// types.
type Ticker struct {
	Pair          string  // The fields include the name of the trading pair
	Last          float64 //  the current price
	Low           float64 // the lowest price in the past 24 hours
	High          float64 // the highest price in the past 24 hours
	Volume        float64 // the volume of trades in the past 24 hours
	QuoteVolume   float64 // the quote volume in the past 24 hours
	PercentChange float64
	// Additional fields as needed...
}
