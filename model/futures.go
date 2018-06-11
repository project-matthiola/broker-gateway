package model

// Futures represent financial contracts obligating the buyer to purchase an asset or the seller to sell an asset.
type Futures struct {
	FuturesID   string `gorm:"primary_key" json:"futures_id"`
	FuturesName string `json:"futures_name"`
	Symbol      string `json:"symbol"`
	Period      string `json:"period"`
}

func (Futures) TableName() string {
	return "futures"
}
