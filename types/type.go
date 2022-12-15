package types

// types.
type Ticker struct {
	Symbol string  `json:"-"`     // The fields include the name of the trading pair
	Price  float64 `json:"price"` //  the current price
}
