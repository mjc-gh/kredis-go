package kredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO Add Pipelining support
// https://redis.uptrace.dev/guide/go-redis-pipelines.html#pipelines

type Proxy struct {
	ctx       context.Context
	client    *redis.Client
	key       string
	expiresIn time.Duration
}

func NewProxy(key string, opts ...ProxyOption) (*Proxy, error) {
	options := ProxyOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	client, namespace, err := getConnection(options.Config())
	if err != nil {
		return nil, err
	}

	if namespace != nil {
		key = fmt.Sprintf("%s:%s", *namespace, key)
	}

	return &Proxy{
		// TODO figure out the best way to handle context
		ctx:       context.TODO(),
		client:    client,
		key:       key,
		expiresIn: options.ExpiresIn(),
	}, nil
}

func (p *Proxy) watch(setter func() error) error {
	err := p.client.Watch(p.ctx, func(tx *redis.Tx) error {
		n, err := tx.Exists(p.ctx, p.key).Result()
		if err != nil {
			return err
		} else if n > 0 { // already exists
			return nil
		}

		return setter()
	}, p.key)
	if err != nil {
		return err
	}

	return nil
}
