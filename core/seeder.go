package core

type Seeder interface {
	Seed(...interface{}) error
	Seeds() ([]interface{}, error)
}
