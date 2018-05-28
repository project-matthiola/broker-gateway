package mapper

import (
	"github.com/jinzhu/gorm"
)

type Mapper struct {
	db *gorm.DB
}

func NewMapper() *Mapper {
	return &Mapper{
		db: NewDB(),
	}
}

func (m *Mapper) Create(model interface{}) error {
	return m.db.Create(model).Error
}

func (m *Mapper) FirstByID(model interface{}, id int) error {
	return m.db.First(model, id).Error
}
