package core

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
)

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

func (s *Seeds) AddTyped(l hclog.Logger, v ...interface{}) {
	log := l.Named("seed-saver")
	// Generate map of existing values for quick filtering
	exist := make(map[interface{}]struct{})
	for _, t := range s.Typed {
		exist[t] = struct{}{}
	}
	log.Debug("existing seeds", "exist", fmt.Sprintf("%#v", exist))

	for _, t := range v {
		log.Debug("attempting to add typed seed", "value", fmt.Sprintf("%#v", t))
		if _, ok := exist[t]; !ok {
			log.Debug("it does not exist so adding it")
			s.Typed = append(s.Typed, t)
			exist[t] = struct{}{}
		} else {
			log.Debug("it already exists so skipping it")
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
