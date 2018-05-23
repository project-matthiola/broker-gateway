package model

// Firm represents a business organization, such as a corporation, limited liability company or partnership, that sells
// goods or services to make a profit.
type Firm struct {
	FirmID   int `gorm:"primary_key"`
	FirmName string
	Credit   int `gorm:"default:100"`
}

func (Firm) TableName() string {
	return "firm"
}
