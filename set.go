package kredis

import "time"

type Set[T KredisTyped] struct {
	Proxy
}

func NewIntegerSet(key string, options Options) (*Set[int], error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &Set[int]{Proxy: *proxy}, nil
}

func NewStringSet(key string, options Options) (*Set[string], error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &Set[string]{Proxy: *proxy}, nil
}

func NewTimeSet(key string, options Options) (*Set[time.Time], error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &Set[time.Time]{Proxy: *proxy}, nil
}

func (s *Set[T]) Members(members []T) (int, error) {
	return 0, nil
}
