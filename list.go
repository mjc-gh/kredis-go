package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type List[T KredisTyped] struct {
	Proxy
}

// TODO finish generic Default factories
// TODO use expiresIn

// List[bool] type

func NewBoolList(key string, opts ...ProxyOption) (*List[bool], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[bool]{Proxy: *proxy}, nil
}

func NewBoolListWithDefault(key string, defaultElements []bool, opts ...ProxyOption) (l *List[bool], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[bool]{Proxy: *proxy}
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

	return &List[int]{Proxy: *proxy}, nil
}

func NewIntegerListWithDefault(key string, defaultElements []int, opts ...ProxyOption) (l *List[int], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[int]{Proxy: *proxy}
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

	return &List[string]{Proxy: *proxy}, nil
}

func NewStringListWithDefault(key string, defaultElements []string, opts ...ProxyOption) (l *List[string], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[string]{Proxy: *proxy}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// List[time] type

func NewTimeList(key string, opts ...ProxyOption) (*List[time.Time], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[time.Time]{Proxy: *proxy}, nil
}

func NewTimeListWithDefault(key string, defaultElements []time.Time, opts ...ProxyOption) (l *List[time.Time], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[time.Time]{Proxy: *proxy}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// List[kredisJSON] type

func NewJSONList(key string, opts ...ProxyOption) (*List[kredisJSON], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[kredisJSON]{Proxy: *proxy}, nil
}

func NewJSONListWithDefault(key string, defaultElements []kredisJSON, opts ...ProxyOption) (l *List[kredisJSON], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &List[kredisJSON]{Proxy: *proxy}
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
func (l *List[T]) Elements(elements []T, opts ...RangeOption) (total int, err error) {
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

func (l *List[T]) Remove(elements ...T) error {
	iter := newIter(elements)

	for val, ok := iter.next(); ok; {
		l.client.LRem(l.ctx, l.key, 0, val)

		val, ok = iter.next()
	}

	return nil
}

// TODO should Prepend and Append return an int not an int64 for greater ease
// of use??

func (l List[T]) Prepend(elements ...T) (int64, error) {
	if len(elements) < 1 {
		return 0, nil
	}

	llen, err := l.client.LPush(l.ctx, l.key, newIter(elements).values()...).Result()
	if err != nil {
		return 0, err
	}

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

	return llen, nil
}

func (l *List[T]) Clear() error {
	_, err := l.client.Del(l.ctx, l.key).Result()

	if err != nil {
		return err
	}

	return nil
}

func (l *List[T]) Length() (llen int64, err error) {
	llen, err = l.client.LLen(l.ctx, l.key).Result()
	if err == redis.Nil {
		err = nil
	}

	return
}

// TODO add function last(n = 1) ??
// https://github.com/rails/kredis/blob/2ccc5c6bf59e5d38870de45a03e9491a3dc8c397/lib/kredis/types/list.rb#L32-L34
