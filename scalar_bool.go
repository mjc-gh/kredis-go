package kredis

import (
	"fmt"
)

type ScalarBool struct{ Proxy }

func NewBool(key string, opts ...ProxyOption) (*ScalarBool, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &ScalarBool{Proxy: *proxy}, err
}

func NewBoolWithDefault(key string, defaultValue bool, opts ...ProxyOption) (s *ScalarBool, err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	s = &ScalarBool{Proxy: *proxy}
	err = proxy.watch(func() error {
		return s.SetValue(defaultValue)
	})
	if err != nil {
		return nil, err
	}

	return
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
