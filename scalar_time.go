package kredis

import (
	"time"
)

//t := time.Now()

//fmt.Println(t.Format(time.RFC3339Nano))
//fmt.Println(time.Parse(time.RFC3339Nano, t.Format(time.RFC3339Nano)))

type ScalarTime struct{ Proxy }

func NewTime(key string, options Options) (*ScalarTime, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarTime{Proxy: *proxy}, nil
}

func (s *ScalarTime) Value() time.Time {
	val, err := s.ValueWithErr()

	if err != nil || val == nil {
		return s.DefaultValue()
	}

	return *val
}

func (s *ScalarTime) ValueWithErr() (*time.Time, error) {
	val, redisErr := s.client.Get(s.ctx, s.key).Result()

	if redisErr != nil {
		return nil, redisErr
	}

	time, err := time.Parse(time.RFC3339Nano, val)

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
