package server

type Mapper interface {
	New(model interface{})
}

type mapper struct {
}

func (m *mapper) New(model interface{}) {

}
