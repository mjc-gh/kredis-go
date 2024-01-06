package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type List[T KredisTyped] struct {
	Proxy
	typed *T
}

// List[bool] type

func NewBoolList(key string, opts ...ProxyOption) (*List[bool], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[bool]{Proxy: *proxy, typed: new(bool)}, nil
}

func NewBoolListWithDefault(key string, defaultElements []bool, opts ...ProxyOption) (l *List[bool], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[bool]{Proxy: *proxy, typed: new(bool)}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// List[int] type

func NewIntegerList(key string, opts ...ProxyOption) (*List[int], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[int]{Proxy: *proxy, typed: new(int)}, nil
}

func NewIntegerListWithDefault(key string, defaultElements []int, opts ...ProxyOption) (l *List[int], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[int]{Proxy: *proxy, typed: new(int)}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// List[float64] type

func NewFloatList(key string, opts ...ProxyOption) (*List[float64], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &List[float64]{Proxy: *proxy, typed: new(float64)}, nil
}

func NewFloatListWithDefault(key string, defaultElements []float64, opts ...ProxyOption) (l *List[float64], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[float64]{Proxy: *proxy, typed: new(float64)}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// List[string] type

func NewStringList(key string, opts ...ProxyOption) (*List[string], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &List[string]{Proxy: *proxy, typed: new(string)}, nil
}

func NewStringListWithDefault(key string, defaultElements []string, opts ...ProxyOption) (l *List[string], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[string]{Proxy: *proxy, typed: new(string)}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// List[time.Time] type

func NewTimeList(key string, opts ...ProxyOption) (*List[time.Time], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &List[time.Time]{Proxy: *proxy, typed: new(time.Time)}, nil
}

func NewTimeListWithDefault(key string, defaultElements []time.Time, opts ...ProxyOption) (l *List[time.Time], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[time.Time]{Proxy: *proxy, typed: new(time.Time)}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// List[KredisJSON] type

func NewJSONList(key string, opts ...ProxyOption) (*List[KredisJSON], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &List[KredisJSON]{Proxy: *proxy, typed: new(KredisJSON)}, nil
}

func NewJSONListWithDefault(key string, defaultElements []KredisJSON, opts ...ProxyOption) (l *List[KredisJSON], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[KredisJSON]{Proxy: *proxy, typed: new(KredisJSON)}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// generic List functions

// NOTE we're using Do() for the "LRANGE" command instead of LRange() as to
// seemingly avoid an []string allocation from StringSliceCmd#Result()
func (l *List[T]) Elements(elements []T, opts ...RangeOption) (total int64, err error) {
	rangeOptions := RangeOptions{0}
	for _, opt := range opts {
		opt(&rangeOptions)
	}

	stop := rangeOptions.start + int64(len(elements))
	slice, err := l.client.Do(l.ctx, "lrange", l.key, rangeOptions.start, stop).Slice()
	if err != nil {
		return
	}

	total = copyCmdSliceTo(slice, elements)
	return
}

func (l *List[T]) Remove(elements ...T) (err error) {
	for _, element := range elements {
		value := typeToInterface(element)
		if err = l.client.LRem(l.ctx, l.key, 0, value).Err(); err != nil {
			return
		}
	}

	return
}

func (l List[T]) Prepend(elements ...T) (int64, error) {
	if len(elements) < 1 {
		return 0, nil
	}

	llen, err := l.client.LPush(l.ctx, l.key, newIter(elements).values()...).Result()
	if err != nil {
		return 0, err
	}

	l.RefreshTTL()
	return llen, nil
}

func (l *List[T]) Append(elements ...T) (int64, error) {
	if len(elements) < 1 {
		return 0, nil
	}

	llen, err := l.client.RPush(l.ctx, l.key, newIter(elements).values()...).Result()
	if err != nil {
		return 0, err
	}

	l.RefreshTTL()
	return llen, nil
}

func (l *List[T]) Clear() error {
	return l.client.Del(l.ctx, l.key).Err()
}

func (l *List[T]) Length() (llen int64, err error) {
	llen, err = l.client.LLen(l.ctx, l.key).Result()
	if err == redis.Nil {
		err = nil
	}

	return
}

func (l List[T]) Last() (T, bool) {
	slice, err := l.client.Do(l.ctx, "lrange", l.key, -1, -1).Slice()
	if err != nil || len(slice) < 1 {
		return any(*l.typed).(T), false
	}

	elements := make([]T, 1)
	copyCmdSliceTo(slice, elements)
	return elements[0], true
}

func (l List[T]) LastN(elements []T) (total int64, err error) {
	slice, err := l.client.Do(l.ctx, "lrange", l.key, -len(elements), -1).Slice()
	if err != nil {
		return
	}

	total = copyCmdSliceTo(slice, elements)
	return
}
