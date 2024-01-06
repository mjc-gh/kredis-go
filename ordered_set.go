package kredis

import (
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

// Backed by Sorted Sets in redis
type OrderedSet[T KredisTyped] struct {
	Proxy
	limit uint64
	base  time.Time
	typed *T
}

// OrderedSet[bool] type

func NewBoolOrderedSet(key string, limit uint64, opts ...ProxyOption) (*OrderedSet[bool], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &OrderedSet[bool]{Proxy: *proxy, limit: limit, typed: new(bool)}, nil
}

func NewBoolOrderedSetWithDefault(key string, limit uint64, defaultMembers []bool, opts ...ProxyOption) (s *OrderedSet[bool], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &OrderedSet[bool]{Proxy: *proxy, limit: limit, typed: new(bool)}
	err = proxy.watch(func() error {
		_, _, err := s.Append(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// OrderedSet[int] type

func NewIntegerOrderedSet(key string, limit uint64, opts ...ProxyOption) (*OrderedSet[int], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &OrderedSet[int]{Proxy: *proxy, limit: limit, typed: new(int)}, nil
}

func NewIntegerOrderedSetWithDefault(key string, limit uint64, defaultMembers []int, opts ...ProxyOption) (s *OrderedSet[int], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &OrderedSet[int]{Proxy: *proxy, limit: limit, typed: new(int)}
	err = proxy.watch(func() error {
		_, _, err := s.Append(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// OrderedSet[int] type

func NewFloatOrderedSet(key string, limit uint64, opts ...ProxyOption) (*OrderedSet[float64], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &OrderedSet[float64]{Proxy: *proxy, limit: limit, typed: new(float64)}, nil
}

func NewFloatOrderedSetWithDefault(key string, limit uint64, defaultMembers []float64, opts ...ProxyOption) (s *OrderedSet[float64], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &OrderedSet[float64]{Proxy: *proxy, limit: limit, typed: new(float64)}
	err = proxy.watch(func() error {
		_, _, err := s.Append(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// OrderedSet[string] type

func NewStringOrderedSet(key string, limit uint64, opts ...ProxyOption) (*OrderedSet[string], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &OrderedSet[string]{Proxy: *proxy, limit: limit, typed: new(string)}, nil
}

func NewStringOrderedSetWithDefault(key string, limit uint64, defaultMembers []string, opts ...ProxyOption) (s *OrderedSet[string], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &OrderedSet[string]{Proxy: *proxy, limit: limit, typed: new(string)}
	err = proxy.watch(func() error {
		_, _, err := s.Append(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// OrderedSet[time.Time]

func NewTimeOrderedSet(key string, limit uint64, opts ...ProxyOption) (*OrderedSet[time.Time], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &OrderedSet[time.Time]{Proxy: *proxy, limit: limit, typed: new(time.Time)}, nil
}

func NewTimeOrderedSetWithDefault(key string, limit uint64, defaultMembers []time.Time, opts ...ProxyOption) (s *OrderedSet[time.Time], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &OrderedSet[time.Time]{Proxy: *proxy, limit: limit, typed: new(time.Time)}
	err = proxy.watch(func() error {
		_, _, err := s.Append(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// OrderedSet[KredisJSON] type

func NewJSONOrderedSet(key string, limit uint64, opts ...ProxyOption) (*OrderedSet[KredisJSON], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &OrderedSet[KredisJSON]{Proxy: *proxy, limit: limit, typed: new(KredisJSON)}, nil
}

func NewJSONOrderedSetWithDefault(key string, limit uint64, defaultMembers []KredisJSON, opts ...ProxyOption) (s *OrderedSet[KredisJSON], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	s = &OrderedSet[KredisJSON]{Proxy: *proxy, limit: limit, typed: new(KredisJSON)}
	err = proxy.watch(func() error {
		_, _, err := s.Append(defaultMembers...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// generic Set functions

func (s *OrderedSet[T]) Members() ([]T, error) {
	slice, err := s.client.Do(s.ctx, "zrange", s.key, 0, -1).Slice()
	if err != nil {
		return nil, err
	}

	members := make([]T, len(slice))
	copyCmdSliceTo(slice, members)

	return members, nil
}

func (s *OrderedSet[T]) Append(members ...T) (added int64, removed int64, err error) {
	if len(members) < 1 {
		return
	}

	pipe := s.client.TxPipeline()
	add := pipe.ZAdd(s.ctx, s.key, newIter(members).valuesWithScoring(s, false)...)
	if s.limit > 0 {
		rem := pipe.ZRemRangeByRank(s.ctx, s.key, 0, -int64(s.limit+1))
		pipe.Exec(s.ctx)
		removed = rem.Val()
	} else {
		pipe.Exec(s.ctx)
	}

	added = add.Val()
	s.RefreshTTL()
	return
}

func (s *OrderedSet[T]) Prepend(members ...T) (added int64, removed int64, err error) {
	if len(members) < 1 {
		return 0, 0, nil
	}

	pipe := s.client.TxPipeline()
	add := pipe.ZAdd(s.ctx, s.key, newIter(members).valuesWithScoring(s, true)...)
	if s.limit > 0 {
		rem := pipe.ZRemRangeByRank(s.ctx, s.key, int64(s.limit), -1)
		pipe.Exec(s.ctx)
		removed = rem.Val()
	} else {
		pipe.Exec(s.ctx)
	}

	added = add.Val()
	s.RefreshTTL()
	return
}

func (s *OrderedSet[T]) Remove(members ...T) (removed int64, err error) {
	removed, err = s.client.ZRem(s.ctx, s.key, newIter(members).values()...).Result()
	s.RefreshTTL()

	return
}

func (s *OrderedSet[T]) Includes(member T) bool {
	err := s.client.ZScore(s.ctx, s.key, typeToInterface(member).(string)).Err()
	if err != nil {
		return false
	}

	return true
}

func (s *OrderedSet[T]) Clear() error {
	return s.client.Del(s.ctx, s.key).Err()
}

func (s *OrderedSet[T]) Size() int64 {
	return s.client.ZCard(s.ctx, s.key).Val()
}

func (s *OrderedSet[T]) SetLimit(limit uint64) {
	s.limit = limit
}

// TODO
//func (s OrderedSet[T]) Rank(member T) int64 {
//}

func (s *OrderedSet[T]) appendScore(index int) float64 {
	baseScore := s.baseScore()
	incrementalScore := float64(index) * 0.000001

	return baseScore + incrementalScore
}

func (s *OrderedSet[T]) prependScore(index int) float64 {
	baseScore := s.baseScore()
	incrementalScore := float64(index) * 0.000001

	return -baseScore - incrementalScore
}

func (s *OrderedSet[T]) baseScore() float64 {
	if s.base.IsZero() {
		s.base = s.client.Time(s.ctx).Val()
	}

	ts := s.base.Add(time.Since(s.base)).UnixNano()

	return float64(ts) / float64(1e9)
}
