package kredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Proxy struct {
	ctx       context.Context
	client    *redis.Client
	key       string
	expiresIn time.Duration
}

func NewProxy(key string, opts ...ProxyOption) (*Proxy, error) {
	options := ProxyOptions{context: context.Background()}
	for _, opt := range opts {
		opt(&options)
	}

	client, namespace, err := getConnection(options.Config())
	if err != nil {
		return nil, err
	}

	if namespace != "" {
		key = fmt.Sprintf("%s:%s", namespace, key)
	}

	return &Proxy{
		ctx:       options.context,
		client:    client,
		key:       key,
		expiresIn: options.ExpiresIn(),
	}, nil
}

// Used when setting defaults
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

// Get the key's current TTL. Redis is only called if the type was configured
// WithExpiry(). If no expiry is configured, a zero value Duration is returned
func (p *Proxy) TTL() (time.Duration, error) {
	if p.expiresIn == 0 {
		return time.Duration(0), nil
	}

	return p.client.TTL(p.ctx, p.key).Result()
}

// Set the key's EXPIRE using the configured expiresIn. If there is no
// value configured, nothing happens.
func (p *Proxy) RefreshTTL() (bool, error) {
	if p.expiresIn == 0 {
		return false, nil
	}

	return p.client.ExpireXX(p.ctx, p.key, p.expiresIn).Result()
}
