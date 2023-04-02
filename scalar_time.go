package kredis

import (
	"time"
)

type ScalarTime struct{ Proxy }

func NewTime(key string, options Options) (*ScalarTime, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarTime{Proxy: *proxy}, nil
}

func (s *ScalarTime) Value() time.Time {
	val, err := s.ValueResult()

	if err != nil || val == nil {
		return s.DefaultValue()
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

func (s *ScalarTime) DefaultValue() time.Time {
	switch s.defaultValue.(type) {
	case time.Time:
		return s.defaultValue.(time.Time)
	default:
		return time.Time{}
	}
}
