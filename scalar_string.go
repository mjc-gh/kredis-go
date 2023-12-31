package kredis

type ScalarString struct{ Proxy }

func NewString(key string, opts ...ProxyOption) (*ScalarString, error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &ScalarString{Proxy: *proxy}, err
}

func NewStringWithDefault(key string, defaultValue string, opts ...ProxyOption) (s *ScalarString, err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	s = &ScalarString{Proxy: *proxy}
	err = proxy.watch(func() error {
		return s.SetValue(defaultValue)
	})
	if err != nil {
		return nil, err
	}

	return
}

func (s *ScalarString) Value() string {
	val, err := s.ValueResult()

	if err != nil || val == nil {
		return ""
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
