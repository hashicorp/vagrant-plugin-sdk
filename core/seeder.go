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
	// Generate map of existing values for quick filtering
	exist := make(map[interface{}]struct{})
	for _, t := range s.Typed {
		exist[t] = struct{}{}
	}

	for _, t := range v {
		if _, ok := exist[t]; !ok {
			s.Typed = append(s.Typed, t)
			exist[t] = struct{}{}
		}
	}
}

func (s *Seeds) AddNamed(n string, v interface{}) {
	s.Named[n] = v
}

type Seeder interface {
	Seed(*Seeds) error
	Seeds() (*Seeds, error)
}
