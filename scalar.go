package kredis

import (
	"fmt"
	"time"
)

/* Bool Scalar type */

type ScalarBool struct{ Proxy }

func NewBool(key string, opts ...ProxyOption) (*ScalarBool, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &ScalarBool{Proxy: *proxy}, err
}

func NewBoolWithDefault(key string, defaultValue bool, opts ...ProxyOption) (s *ScalarBool, err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	s = &ScalarBool{Proxy: *proxy}
	err = proxy.watch(func() error {
		return s.SetValue(defaultValue)
	})
	if err != nil {
		return nil, err
	}

	return
}

func (s *ScalarBool) Value() bool {
	val, err := s.ValueResult()
	if err != nil {
		return false
	}

	return *val
}

func (s *ScalarBool) ValueResult() (*bool, error) {
	val, err := s.client.Get(s.ctx, s.key).Bool()
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (s *ScalarBool) SetValue(v bool) error {
	return s.client.Set(s.ctx, s.key, fmt.Sprintf("%v", v), s.expiresIn).Err()
}

/* Integer Scalar type */

type ScalarInteger struct{ Proxy }

func NewInteger(key string, opts ...ProxyOption) (*ScalarInteger, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &ScalarInteger{Proxy: *proxy}, err
}

func NewIntegerWithDefault(key string, defaultValue int, opts ...ProxyOption) (s *ScalarInteger, err error) {
	proxy, err := NewProxy(key, opts...)
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

/* Float Scalar type */

type ScalarFloat struct{ Proxy }

func NewFloat(key string, opts ...ProxyOption) (*ScalarFloat, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &ScalarFloat{Proxy: *proxy}, err
}

func NewFloatWithDefault(key string, defaultValue float64, opts ...ProxyOption) (s *ScalarFloat, err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	s = &ScalarFloat{Proxy: *proxy}
	err = proxy.watch(func() error {
		return s.SetValue(defaultValue)
	})
	if err != nil {
		return nil, err
	}

	return
}

func (s *ScalarFloat) Value() float64 {
	val, err := s.ValueResult()

	if err != nil || val == nil {
		return 0
	}

	return *val
}

func (s *ScalarFloat) ValueResult() (*float64, error) {
	val, err := s.client.Get(s.ctx, s.key).Float64()

	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (s *ScalarFloat) SetValue(v float64) error {
	return s.client.Set(s.ctx, s.key, fmt.Sprintf("%f", v), s.expiresIn).Err()
}

/* String Scalar type */

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

/* Time Scalar type */

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

/* JSON Scalar type */

type ScalarJSON struct{ Proxy }

func NewJSON(key string, opts ...ProxyOption) (*ScalarJSON, error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &ScalarJSON{Proxy: *proxy}, nil
}

func NewJSONWithDefault(key string, defaultValue *KredisJSON, opts ...ProxyOption) (s *ScalarJSON, err error) {
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
func (s *ScalarJSON) Value() KredisJSON {
	val, err := s.ValueResult()
	if err != nil || val == nil {
		return KredisJSON{}
	}

	return *val
}

func (s *ScalarJSON) ValueResult() (*KredisJSON, error) {
	val, err := s.client.Get(s.ctx, s.key).Result()
	if err != nil {
		return nil, err
	}

	return NewKredisJSON(val), nil
}

func (s *ScalarJSON) SetValue(v *KredisJSON) error {
	return s.client.Set(s.ctx, s.key, string(v.s), s.expiresIn).Err()
}
