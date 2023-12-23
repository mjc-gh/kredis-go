package kredis

import "fmt"

type ScalarInteger struct{ Proxy }

func NewInteger(key string, options Options) (*ScalarInteger, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarInteger{Proxy: *proxy}, err
}

func NewIntegerWithDefault(key string, options Options, defaultValue int) (s *ScalarInteger, err error) {
	proxy, err := NewProxy(key, options)
	if err != nil {
		return
	}

	s = &ScalarInteger{Proxy: *proxy}
	err = proxy.watch(func() error {
		return s.SetValue(defaultValue)
	})
	if err != nil {
		return nil, err
	}

	return
}

func (s *ScalarInteger) Value() int {
	val, err := s.ValueResult()

	if err != nil || val == nil {
		return 0
	}

	return *val
}

func (s *ScalarInteger) ValueResult() (*int, error) {
	val, err := s.client.Get(s.ctx, s.key).Int()

	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (s *ScalarInteger) SetValue(v int) error {
	return s.client.Set(s.ctx, s.key, fmt.Sprintf("%d", v), s.expiresIn).Err()
}
