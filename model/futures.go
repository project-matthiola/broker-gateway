package model

// Futures represent financial contracts obligating the buyer to purchase an asset or the seller to sell an asset.
type Futures struct {
	FuturesID   string
	FuturesName string
	Symbol      string
	Period      string
}
