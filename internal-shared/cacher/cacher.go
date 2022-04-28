package cacher

type Cache interface {
	Register(key string, value interface{})
	Get(key string) interface{}
	Keys() []string
	Values() []interface{}
}

type HasCache interface {
	SetCache(Cache)
}

func New() Cache {
	return &cache{
		registry: map[string]interface{}{},
	}
}

// Used for caching local conversions
type cache struct {
	registry map[string]interface{}
}

func (c *cache) Register(key string, value interface{}) {
	c.registry[key] = value
}

func (c *cache) Get(key string) interface{} {
	return c.registry[key]
}

func (c *cache) Keys() []string {
	keys := make([]string, 0, len(c.registry))
	for k, _ := range c.registry {
		keys = append(keys, k)
	}

	return keys
}

func (c *cache) Values() []interface{} {
	values := make([]interface{}, 0, len(c.registry))
	for _, v := range c.registry {
		values = append(values, v)
	}

	return values
}
