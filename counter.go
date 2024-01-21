package kredis

import "github.com/redis/go-redis/v9"

type Counter struct {
	Proxy
}

func NewCounter(key string, opts ...ProxyOption) (*Counter, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Counter{Proxy: *proxy}, nil
}

func NewCounterWithDefault(key string, defaultValue int64, opts ...ProxyOption) (c *Counter, err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	c = &Counter{Proxy: *proxy}
	err = proxy.watch(func() error {
		_, err := c.Increment(defaultValue)
		return err
	})

	return
}

func (c *Counter) Increment(by int64) (int64, error) {
	pipe := c.client.TxPipeline()
	if c.expiresIn > 0 {
		pipe.SetNX(c.ctx, c.key, 0, c.expiresIn)
	}

	incr := pipe.IncrBy(c.ctx, c.key, by)
	_, err := pipe.Exec(c.ctx)
	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}

func (c *Counter) Decrement(by int64) (int64, error) {
	pipe := c.client.TxPipeline()
	if c.expiresIn > 0 {
		pipe.SetNX(c.ctx, c.key, 0, c.expiresIn)
	}

	decr := pipe.DecrBy(c.ctx, c.key, by)
	_, err := pipe.Exec(c.ctx)
	if err != nil {
		return 0, err
	}

	return decr.Val(), nil
}

// An empty value returned when there is a Redis error as a failsafe
func (c *Counter) Value() (v int64) {
	v, err := c.client.Get(c.ctx, c.key).Int64()
	if err != nil && err != redis.Nil {
		// TODO debug logging
	}

	return
}

func (c *Counter) Reset() (err error) {
	_, err = c.client.Del(c.ctx, c.key).Result()
	return
}
