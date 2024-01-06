package kredis

import "github.com/redis/go-redis/v9"

type Cycle struct {
	Proxy
	values []string
}

func NewCycle(key string, values []string, opts ...ProxyOption) (*Cycle, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Cycle{Proxy: *proxy, values: values}, nil
}

func (c *Cycle) Index() int64 {
	return c.value()
}

func (c *Cycle) Value() string {
	return c.values[c.Index()]
}

func (c *Cycle) Next() (err error) {
	value := (c.value() + 1) % int64(len(c.values))
	_, err = c.client.Set(c.ctx, c.key, value, c.expiresIn).Result()

	return
}

func (c *Cycle) value() (v int64) {
	v, err := c.client.Get(c.ctx, c.key).Int64()
	if err != nil && err != redis.Nil {
		// TODO debug logging
	}

	return
}
