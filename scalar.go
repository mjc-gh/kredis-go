package kredis

import "fmt"

type ScalarInteger struct{ Proxy }

func NewInteger(key string, options Options) (*ScalarInteger, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarInteger{Proxy: *proxy}, nil
}

func (s *ScalarInteger) Value() int {
	val, err := s.client.Get(s.ctx, s.key).Int()

	if err != nil {
		return s.DefaultValue()
	}

	return val
}

func (s *ScalarInteger) SetValue(v int) error {
	s.client.Set(s.ctx, s.key, fmt.Sprintf("%d", v), s.expiresIn)

	return nil
}

func (s *ScalarInteger) DefaultValue() int {
	switch s.defaultValue.(type) {
	case int:
		return s.defaultValue.(int)
	default:
		return 0
	}
}

type ScalarString struct{ Proxy }

func NewString(key string, options Options) (*ScalarString, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarString{Proxy: *proxy}, nil
}

func (s *ScalarString) Value() string {
	val, err := s.client.Get(s.ctx, s.key).Result()

	if err != nil {
		// TODO handle defaultValue
		return ""
	}

	return val
}

func (s *ScalarString) SetValue(v string) error {
	s.client.Set(s.ctx, s.key, v, s.expiresIn)

	return nil
}
