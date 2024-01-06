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
}

type FlagMarkOption func(*FlagMarkOptions)

func NewFlag(key string, opts ...ProxyOption) (*Flag, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Flag{Proxy: *proxy}, nil
}

func NewMarkedFlag(key string, opts ...ProxyOption) (*Flag, error) {
	flag, err := NewFlag(key, opts...)
	if err != nil {
		return nil, err
	}

	err = flag.Mark()
	if err != nil {
		return nil, err
	}

	return flag, nil
}

func (f *Flag) Mark(opts ...FlagMarkOption) error {
	options := FlagMarkOptions{f.expiresIn}
	for _, opt := range opts {
		opt(&options)
	}

	return f.client.Set(f.ctx, f.key, 1, options.expiresIn).Err()
}

func (f *Flag) SoftMark(opts ...FlagMarkOption) error {
	options := FlagMarkOptions{f.expiresIn}
	for _, opt := range opts {
		opt(&options)
	}

	return f.client.SetNX(f.ctx, f.key, 1, options.expiresIn).Err()
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

func WithFlagExpiry(expires string) FlagMarkOption {
	return func(o *FlagMarkOptions) {
		duration, err := time.ParseDuration(expires)
		if err != nil {
			return
		}

		o.expiresIn = duration
	}
}
