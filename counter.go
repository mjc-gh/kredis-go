package kredis

import "github.com/redis/go-redis/v9"

type Counter struct {
	Proxy
}

// TODO add default value factory

func NewCounter(key string, opts ...ProxyOption) (*Counter, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Counter{Proxy: *proxy}, nil
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
