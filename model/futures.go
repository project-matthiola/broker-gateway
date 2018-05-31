package model

// Futures represent financial contracts obligating the buyer to purchase an asset or the seller to sell an asset.
type Futures struct {
	FuturesID   string `gorm:"primary_key"`
	FuturesName string
	Symbol      string
	Period      string
}

func (Futures) TableName() string {
	return "futures"
}
