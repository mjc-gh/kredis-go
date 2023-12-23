package kredis

import (
	"time"
)

type ScalarTime struct{ Proxy }

func NewTime(key string, opts ...ProxyOption) (*ScalarTime, error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &ScalarTime{Proxy: *proxy}, nil
}

func NewTimeWithDefault(key string, defaultValue time.Time, opts ...ProxyOption) (s *ScalarTime, err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	s = &ScalarTime{Proxy: *proxy}
	err = proxy.watch(func() error {
		return s.SetValue(defaultValue)
	})
	if err != nil {
		return nil, err
	}

	return
}

func (s *ScalarTime) Value() time.Time {
	val, err := s.ValueResult()

	if err != nil || val == nil {
		return time.Time{} // empty value
	}

	return *val
}

func (s *ScalarTime) ValueResult() (*time.Time, error) {
	time, err := s.client.Get(s.ctx, s.key).Time()

	if err != nil {
		return nil, err
	}

	return &time, nil
}

func (s *ScalarTime) SetValue(v time.Time) error {
	s.client.Set(s.ctx, s.key, v.Format(time.RFC3339Nano), s.expiresIn)

	return nil
}
