package kredis

type ScalarJSON struct{ Proxy }

func NewJSON(key string, opts ...ProxyOption) (*ScalarJSON, error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &ScalarJSON{Proxy: *proxy}, nil
}

func NewJSONWithDefault(key string, defaultValue *kredisJSON, opts ...ProxyOption) (s *ScalarJSON, err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	s = &ScalarJSON{Proxy: *proxy}
	err = proxy.watch(func() error {
		return s.SetValue(defaultValue)
	})
	if err != nil {
		return nil, err
	}

	return
}

// TODO should this be returning a pointer instead struct value itself??
func (s *ScalarJSON) Value() kredisJSON {
	val, err := s.ValueResult()

	if err != nil || val == nil {
		return kredisJSON{}
	}

	return *val
}

func (s *ScalarJSON) ValueResult() (*kredisJSON, error) {
	val, err := s.client.Get(s.ctx, s.key).Result()

	if err != nil {
		return nil, err
	}

	return NewKredisJSON(val), nil
}

func (s *ScalarJSON) SetValue(v *kredisJSON) error {
	return s.client.Set(s.ctx, s.key, string(v.s), s.expiresIn).Err()
}
