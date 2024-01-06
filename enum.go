package kredis

import (
	"errors"
	"time"
)

type Enum struct {
	Proxy
	defaultValue string
	values       map[string]bool
}

var EnumEmptyValues = errors.New("values cannot be empty")
var EnumInvalidValue = errors.New("invalid enum value")

func NewEnum(key string, defaultValue string, values []string, opts ...ProxyOption) (*Enum, error) {
	if len(values) == 0 {
		return nil, EnumEmptyValues
	}

	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	// TODO return runtime error if expiresIn option is used -- this option just
	// doesn't fit well into this Kredis data structure

	enum := &Enum{Proxy: *proxy, defaultValue: defaultValue, values: map[string]bool{}}
	for _, value := range values {
		enum.values[value] = true
	}

	err = enum.SetValue(defaultValue)
	if err != nil {
		return nil, err
	}

	return enum, nil
}

func (e *Enum) Is(value string) bool {
	return e.Value() == value
}

func (e *Enum) Value() string {
	value, _ := e.client.Get(e.ctx, e.key).Result()
	return value
}

func (e *Enum) SetValue(value string) error {
	if _, ok := e.values[value]; !ok {
		return EnumInvalidValue
	}

	_, err := e.client.Set(e.ctx, e.key, value, time.Duration(0)).Result()
	return err
}

func (e *Enum) Reset() (err error) {
	pipe := e.client.TxPipeline()
	pipe.Del(e.ctx, e.key)
	pipe.Set(e.ctx, e.key, e.defaultValue, time.Duration(0))

	_, err = pipe.Exec(e.ctx)
	return
}
