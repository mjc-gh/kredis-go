package kredis

import "time"

type Set[T KredisTyped] struct {
	Proxy
	typed *T
}

// Set[bool] type

func NewBoolSet(key string, opts ...ProxyOption) (*Set[bool], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Set[bool]{Proxy: *proxy, typed: new(bool)}, nil
}

func NewBoolSetWithDefault(key string, defaultMembers []bool, opts ...ProxyOption) (s *Set[bool], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &Set[bool]{Proxy: *proxy, typed: new(bool)}
	err = proxy.watch(func() error {
		_, err := s.Add(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// Set[int] type

func NewIntegerSet(key string, opts ...ProxyOption) (*Set[int], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Set[int]{Proxy: *proxy, typed: new(int)}, nil
}

func NewIntegerSetWithDefault(key string, defaultMembers []int, opts ...ProxyOption) (s *Set[int], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &Set[int]{Proxy: *proxy, typed: new(int)}
	err = proxy.watch(func() error {
		_, err := s.Add(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// Set[float64] type

func NewFloatSet(key string, opts ...ProxyOption) (*Set[float64], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Set[float64]{Proxy: *proxy, typed: new(float64)}, nil
}

func NewFloatSetWithDefault(key string, defaultMembers []float64, opts ...ProxyOption) (s *Set[float64], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &Set[float64]{Proxy: *proxy, typed: new(float64)}
	err = proxy.watch(func() error {
		_, err := s.Add(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// Set[string] type

func NewStringSet(key string, opts ...ProxyOption) (*Set[string], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Set[string]{Proxy: *proxy, typed: new(string)}, nil
}

func NewStringSetWithDefault(key string, defaultMembers []string, opts ...ProxyOption) (s *Set[string], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &Set[string]{Proxy: *proxy, typed: new(string)}
	err = proxy.watch(func() error {
		_, err := s.Add(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// Set[time.Time]

func NewTimeSet(key string, opts ...ProxyOption) (*Set[time.Time], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Set[time.Time]{Proxy: *proxy, typed: new(time.Time)}, nil
}

func NewTimeSetWithDefault(key string, defaultMembers []time.Time, opts ...ProxyOption) (s *Set[time.Time], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &Set[time.Time]{Proxy: *proxy, typed: new(time.Time)}
	err = proxy.watch(func() error {
		_, err := s.Add(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// Set[kredisJSON] type

func NewJSONSet(key string, opts ...ProxyOption) (*Set[kredisJSON], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Set[kredisJSON]{Proxy: *proxy, typed: new(kredisJSON)}, nil
}

func NewJSONSetWithDefault(key string, defaultMembers []kredisJSON, opts ...ProxyOption) (s *Set[kredisJSON], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &Set[kredisJSON]{Proxy: *proxy, typed: new(kredisJSON)}
	err = proxy.watch(func() error {
		_, err := s.Add(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// generic Set functions

// TODO this will force an allocation onto the caller -- can this be avoided?
// are there better Go idioms for read all set members?
func (s *Set[T]) Members() ([]T, error) {
	slice, err := s.client.Do(s.ctx, "smembers", s.key).Slice()
	if err != nil {
		return nil, err
	}

	members := make([]T, len(slice))
	copyCmdSliceTo(slice, members)

	return members, nil
}

// TODO return a map? will not work with bool...
//func (s *Set[T]) MembersMap ??

func (s *Set[T]) Add(members ...T) (added int64, err error) {
	if len(members) < 1 {
		return 0, nil
	}

	added, err = s.client.SAdd(s.ctx, s.key, newIter(members).values()...).Result()
	s.RefreshTTL()
	return
}

func (s *Set[T]) Remove(members ...T) (removed int64, err error) {
	if len(members) < 1 {
		return 0, nil
	}

	removed, err = s.client.SRem(s.ctx, s.key, newIter(members).values()...).Result()
	s.RefreshTTL()
	return
}

func (s *Set[T]) Replace(members ...T) (int64, error) {
	pipe := s.client.TxPipeline()
	pipe.Del(s.ctx, s.key)
	add := pipe.SAdd(s.ctx, s.key, newIter(members).values()...)

	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return 0, err
	}

	s.RefreshTTL()
	return add.Val(), nil
}

func (s *Set[T]) Includes(member T) bool {
	return s.client.SIsMember(s.ctx, s.key, typeToInterface(member)).Val()
}

func (s *Set[T]) Size() int64 {
	return s.client.SCard(s.ctx, s.key).Val()
}

func (s *Set[T]) Take() (T, bool) {
	cmd := s.client.SPop(s.ctx, s.key)
	s.RefreshTTL()
	return stringCmdToTyped[T](cmd, s.typed)
}

// TODO func (s *Set[T]) TakeN(memebers []T) (error)

func (s *Set[T]) Clear() error {
	return s.client.Del(s.ctx, s.key).Err()
}

func (s *Set[T]) Sample(members []T) (total int64, err error) {
	if len(members) == 0 {
		return
	}

	slice, err := s.client.Do(s.ctx, "srandmember", s.key, len(members)).Slice()
	if err != nil {
		return
	}

	s.RefreshTTL()
	total = copyCmdSliceTo(slice, members)
	return
}
