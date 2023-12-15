package kredis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ScalarBool struct{ Proxy }

func NewBool(key string, options Options) (*ScalarBool, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarBool{Proxy: *proxy}, err
}

func NewBoolWithDefault(key string, options Options, defaultValue bool) (*ScalarBool, error) {
	s, err := NewBool(key, options)
	if err != nil {
		return nil, err
	}

	// set default value with a transaction
	err = s.client.Watch(s.ctx, func(tx *redis.Tx) error {
		n, err := tx.Exists(s.ctx, s.key).Result()
		if err != nil {
			return err
		} else if n > 0 { // already exists
			return nil
		}

		return s.SetValue(defaultValue)
	}, key)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *ScalarBool) Value() bool {
	val, err := s.ValueResult()

	if err != nil {
		return false
	}

	return *val
}

func (s *ScalarBool) ValueResult() (*bool, error) {
	val, err := s.client.Get(s.ctx, s.key).Bool()

	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (s *ScalarBool) SetValue(v bool) error {
	return s.client.Set(s.ctx, s.key, fmt.Sprintf("%v", v), s.expiresIn).Err()
}
