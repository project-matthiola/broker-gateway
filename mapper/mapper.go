package mapper

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
