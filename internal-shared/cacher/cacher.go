package cacher

type Cache interface {
	Register(key string, value interface{})
	Get(key string) interface{}
}

type HasCache interface {
	SetCache(Cache)
}

func New() *Cacher {
	return &Cacher{
		registry: map[string]interface{}{},
	}
}

// Used for caching local conversions
type Cacher struct {
	registry map[string]interface{}
}

func (c *Cacher) Register(key string, value interface{}) {
	c.registry[key] = value
}

func (c *Cacher) Get(key string) interface{} {
	return c.registry[key]
}
