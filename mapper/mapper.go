package mapper

import "strings"

type Mapper struct {
}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Create(model interface{}) error {
	return DB.Create(model).Error
}

func (m *Mapper) FirstByID(model interface{}, id int) error {
	return DB.First(model, id).Error
}

func (m *Mapper) Save(model interface{}) error {
	return DB.Save(model).Error
}

func (m *Mapper) Update(model interface{}, column string, value string) error {
	return DB.Model(model).Update(column, value).Error
}

func (m *Mapper) Delete(model interface{}) error {
	return DB.Delete(model).Error
}

func (m *Mapper) WhereByUUID(model interface{}, column string, uuid string) error {
	return DB.Where(strings.Replace("column = ?", "column", column, -1), uuid).First(model).Error
}

func (m *Mapper) FindWithLimit(models interface{}, limit int) error {
	return DB.Order("created_at desc").Limit(limit).Find(models).Error
}

func (m *Mapper) FutureIDs() ([]string, error) {
	rows, err := DB.Raw("SELECT DISTINCT futures_id FROM trade").Rows()
	defer rows.Close()
	var futuresIDs []string
	for rows.Next() {
		var futuresID string
		rows.Scan(&futuresID)
		futuresIDs = append(futuresIDs, futuresID)
	}
	return futuresIDs, err
}

func (m *Mapper) FindByFuturesID(models interface{}, futuresID string) error {
	return DB.Where("futures_id = ?", futuresID).Limit(50).Find(models).Error
}
