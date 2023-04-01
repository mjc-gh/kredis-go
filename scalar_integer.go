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

func (s *ScalarInteger) Value() int {
	val, err := s.ValueWithErr()

	if err != nil || val == nil {
		return s.DefaultValue()
	}

	return *val
}

func (s *ScalarInteger) ValueWithErr() (*int, error) {
	val, err := s.client.Get(s.ctx, s.key).Int()

	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (s *ScalarInteger) SetValue(v int) error {
	return s.client.Set(s.ctx, s.key, fmt.Sprintf("%d", v), s.expiresIn).Err()
}

func (s *ScalarInteger) DefaultValue() int {
	switch s.defaultValue.(type) {
	case int:
		return s.defaultValue.(int)
	default:
		return 0
	}
}
