package kredis

type ScalarJSON struct{ Proxy }

func NewJSON(key string, options Options) (*ScalarJSON, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarJSON{Proxy: *proxy}, nil
}

// TODO should this be returning a pointer instead struct value itself??
func (s *ScalarJSON) Value() kredisJSON {
	val, err := s.ValueResult()

	if err != nil || val == nil {
		return s.DefaultValue()
	}

	return *val
}

func (s *ScalarJSON) ValueResult() (*kredisJSON, error) {
	val, err := s.client.Get(s.ctx, s.key).Result()

	if err != nil {
		return nil, err
	}

	kjson := kredisJSON(val)

	return &kjson, nil
}

func (s *ScalarJSON) SetValue(v *kredisJSON) error {
	return s.client.Set(s.ctx, s.key, string(*v), s.expiresIn).Err()
}

func (s *ScalarJSON) DefaultValue() kredisJSON {
	switch s.defaultValue.(type) {
	case kredisJSON:
		return s.defaultValue.(kredisJSON)
	default:
		return kredisJSON{}
	}
}
