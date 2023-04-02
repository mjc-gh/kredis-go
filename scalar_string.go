package kredis

type ScalarString struct{ Proxy }

func NewString(key string, options Options) (*ScalarString, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarString{Proxy: *proxy}, err
}

func (s *ScalarString) Value() string {
	val, err := s.ValueResult()

	if err != nil || val == nil {
		return s.DefaultValue()
	}

	return *val
}

func (s *ScalarString) ValueResult() (*string, error) {
	val, err := s.client.Get(s.ctx, s.key).Result()

	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (s *ScalarString) SetValue(v string) error {
	return s.client.Set(s.ctx, s.key, v, s.expiresIn).Err()
}

func (s *ScalarString) DefaultValue() string {
	switch s.defaultValue.(type) {
	case string:
		return s.defaultValue.(string)
	default:
		return ""
	}
}
