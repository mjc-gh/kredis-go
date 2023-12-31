package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type UniqueList[T KredisTyped] struct {
	Proxy
	limit uint64
}

// TODO finish generic Default factories
// TODO use expiresIn

// UniqueList[bool] type

func NewBoolUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[bool], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &UniqueList[bool]{Proxy: *proxy, limit: limit}, nil
}

func NewBoolUniqueListWithDefault(key string, limit uint64, defaultElements []bool, opts ...ProxyOption) (l *UniqueList[bool], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[bool]{Proxy: *proxy, limit: limit}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// UniqueList[int] type

func NewIntegerUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[int], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &UniqueList[int]{Proxy: *proxy, limit: limit}, nil
}

func NewIntegerUniqueListWithDefault(key string, limit uint64, defaultElements []int, opts ...ProxyOption) (l *UniqueList[int], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[int]{Proxy: *proxy, limit: limit}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// UniqueList[string] type

func NewStringUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[string], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &UniqueList[string]{Proxy: *proxy, limit: limit}, nil
}

func NewStringUniqueListWithDefault(key string, limit uint64, defaultElements []string, opts ...ProxyOption) (l *UniqueList[string], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[string]{Proxy: *proxy, limit: limit}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// UniqueList[time] type

func NewTimeUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[time.Time], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &UniqueList[time.Time]{Proxy: *proxy, limit: limit}, nil
}

func NewTimeUniqueListWithDefault(key string, limit uint64, defaultElements []time.Time, opts ...ProxyOption) (l *UniqueList[time.Time], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[time.Time]{Proxy: *proxy, limit: limit}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// UniqueList[kredisJSON] type

func NewJSONUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[kredisJSON], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &UniqueList[kredisJSON]{Proxy: *proxy, limit: limit}, nil
}

func NewJSONUniqueListWithDefault(key string, limit uint64, defaultElements []kredisJSON, opts ...ProxyOption) (l *UniqueList[kredisJSON], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[kredisJSON]{Proxy: *proxy, limit: limit}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// generic UniqueList functions

// NOTE we're using Do() for the "LRANGE" command instead of LRange() as to
// seemingly avoid an []string allocation from StringSliceCmd#Result()
func (l *UniqueList[T]) Elements(elements []T, opts ...RangeOption) (total int64, err error) {
	rangeOptions := RangeOptions{0}
	for _, opt := range opts {
		opt(&rangeOptions)
	}

	stop := rangeOptions.start + int64(len(elements))
	slice, err := l.client.Do(l.ctx, "lrange", l.key, rangeOptions.start, stop).Slice()
	if err != nil {
		return
	}

	total = int64(copyCmdSliceTo(slice, elements))
	return
}

func (l *UniqueList[T]) Remove(elements ...T) error {
	iter := newIter(elements)

	for val, ok := iter.next(); ok; {
		l.client.LRem(l.ctx, l.key, 0, val)

		val, ok = iter.next()
	}

	return nil
}

func (l UniqueList[T]) Prepend(elements ...T) (int64, error) {
	if len(elements) < 1 {
		return 0, nil
	}

	return l.update(elements, func(pipe redis.Pipeliner, uniq []interface{}) *redis.IntCmd {
		return pipe.LPush(l.ctx, l.key, uniq...)
	})
}

func (l *UniqueList[T]) Append(elements ...T) (int64, error) {
	if len(elements) < 1 {
		return 0, nil
	}

	return l.update(elements, func(pipe redis.Pipeliner, uniq []interface{}) *redis.IntCmd {
		return pipe.RPush(l.ctx, l.key, uniq...)
	})
}

func (l *UniqueList[T]) Clear() error {
	_, err := l.client.Del(l.ctx, l.key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (l *UniqueList[T]) Length() (llen int64, err error) {
	llen, err = l.client.LLen(l.ctx, l.key).Result()
	if err == redis.Nil {
		err = nil
	}

	return
}

func (l *UniqueList[T]) update(elements []T, updateFn func(redis.Pipeliner, []interface{}) *redis.IntCmd) (int64, error) {
	uniq := newIter(elements).unique()
	pipe := l.client.TxPipeline()

	for _, u := range uniq {
		pipe.LRem(l.ctx, l.key, 0, u)
	}

	llen := updateFn(pipe, uniq)
	if l.limit > 0 {
		pipe.LTrim(l.ctx, l.key, -int64(l.limit), -1)
	}

	_, err := pipe.Exec(l.ctx)
	if err != nil {
		return 0, err
	}

	return llen.Val(), nil
}

func (l *UniqueList[T]) SetLimit(limit uint64) {
	l.limit = limit
}

// TODO add function last(n = 1) ??
// https://github.com/rails/kredis/blob/2ccc5c6bf59e5d38870de45a03e9491a3dc8c397/lib/kredis/types/list.rb#L32-L34
