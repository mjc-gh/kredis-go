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

func (c *Cycle) Index() (idx int64) {
	idx, err := c.client.Get(c.ctx, c.key).Int64()
	if err != nil && err != redis.Nil {
		if debugLogger != nil {
			debugLogger.Warn("Cycle#Index", err)
		}
	}

	return
}

func (c *Cycle) IndexResult() (int64, error) {
	return c.client.Get(c.ctx, c.key).Int64()
}

func (c *Cycle) Value() string {
	return c.values[c.Index()]
}

func (c *Cycle) Next() (err error) {
	idx, err := c.client.Get(c.ctx, c.key).Int64()
	if err != nil && err != redis.Nil {
		return // GET error
	}

	value := (idx + 1) % int64(len(c.values))
	_, err = c.client.Set(c.ctx, c.key, value, c.expiresIn).Result()

	return
}
