package kredis

import "encoding/json"

type ScalarJSON struct{ Proxy }

func NewJSON(key string, options Options) (*ScalarJSON, error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &ScalarJSON{Proxy: *proxy}, nil
}

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
	s.client.Set(s.ctx, s.key, string(*v), s.expiresIn)

	return nil
}

func (s *ScalarJSON) DefaultValue() kredisJSON {
	switch s.defaultValue.(type) {
	case kredisJSON:
		return s.defaultValue.(kredisJSON)
	default:
		return kredisJSON{}
	}
}

type kredisJSON []byte

func NewKredisJSON(jsonStr string) *kredisJSON {
	var kj kredisJSON = kredisJSON(jsonStr)

	return &kj
}

func (kj *kredisJSON) Unmarshal(data *interface{}) error {
	err := json.Unmarshal(*kj, data)

	if err != nil {
		return err
	}

	return nil
}
