package core

func NewSeeds() *Seeds {
	return &Seeds{
		Named: map[string]interface{}{},
		Typed: []interface{}{},
	}
}

type Seeds struct {
	Named map[string]interface{}
	Typed []interface{}
}

func (s *Seeds) AddTyped(v interface{}) {
	s.Typed = append(s.Typed, v)
}

func (s *Seeds) AddNamed(n string, v interface{}) {
	s.Named[n] = v
}

type Seeder interface {
	Seed(*Seeds) error
	Seeds() (*Seeds, error)
}
