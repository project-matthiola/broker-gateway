package model

// Firm represents a business organization, such as a corporation, limited liability company or partnership, that sells
// goods or services to make a profit.
type Firm struct {
	FirmID   int    `gorm:"primary_key;AUTO_INCREMENT" json:"firm_id"`
	FirmName string `json:"firm_name"`
	Credit   int    `gorm:"default:100" json:"credit"`
}

func (Firm) TableName() string {
	return "firm"
}
