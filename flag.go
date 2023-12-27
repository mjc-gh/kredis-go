package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Flag struct {
	Proxy
}

type FlagMarkOptions struct {
	expiresIn time.Duration
	force     bool
}

type FlagMarkOption func(*FlagMarkOptions)

// TODO add default value factory

func NewFlag(key string, opts ...ProxyOption) (*Flag, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Flag{Proxy: *proxy}, nil
}

func (f *Flag) Mark(opts ...FlagMarkOption) error {
	options := FlagMarkOptions{force: false}
	for _, opt := range opts {
		opt(&options)
	}

	if options.force {
		f.client.Set(f.ctx, f.key, 1, options.expiresIn)
	} else {
		f.client.SetNX(f.ctx, f.key, 1, options.expiresIn)
	}

	return nil
}

func (f *Flag) IsMarked() bool {
	n, err := f.client.Exists(f.ctx, f.key).Result()
	if err != nil && err != redis.Nil {
		// TODO debug logging
	}

	return n > 0
}

func (f *Flag) Remove() (err error) {
	_, err = f.client.Del(f.ctx, f.key).Result()
	return
}

// Mark() function optional configuration functions

func WithFlagMarkExpiry(expires string) FlagMarkOption {
	return func(o *FlagMarkOptions) {
		duration, err := time.ParseDuration(expires)
		if err != nil {
			return
		}

		o.expiresIn = duration
	}
}

func WithFlagMarkForced() FlagMarkOption {
	return func(o *FlagMarkOptions) {
		o.force = true
	}
}
