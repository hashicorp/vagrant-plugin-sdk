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

func (s *Seeds) AddTyped(v ...interface{}) {
	for _, t := range v {
		s.Typed = append(s.Typed, t)
	}
}

func (s *Seeds) AddNamed(n string, v interface{}) {
	s.Named[n] = v
}

type Seeder interface {
	Seed(*Seeds) error
	Seeds() (*Seeds, error)
}
