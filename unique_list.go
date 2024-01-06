package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type UniqueList[T KredisTyped] struct {
	Proxy
	limit uint64
	typed *T
}

// UniqueList[bool] type

func NewBoolUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[bool], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &UniqueList[bool]{Proxy: *proxy, limit: limit, typed: new(bool)}, nil
}

func NewBoolUniqueListWithDefault(key string, limit uint64, defaultElements []bool, opts ...ProxyOption) (l *UniqueList[bool], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[bool]{Proxy: *proxy, limit: limit, typed: new(bool)}
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

	return &UniqueList[int]{Proxy: *proxy, limit: limit, typed: new(int)}, nil
}

func NewIntegerUniqueListWithDefault(key string, limit uint64, defaultElements []int, opts ...ProxyOption) (l *UniqueList[int], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[int]{Proxy: *proxy, limit: limit, typed: new(int)}
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

func NewFloatUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[float64], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &UniqueList[float64]{Proxy: *proxy, limit: limit, typed: new(float64)}, nil
}

func NewFloatUniqueListWithDefault(key string, limit uint64, defaultElements []float64, opts ...ProxyOption) (l *UniqueList[float64], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[float64]{Proxy: *proxy, limit: limit, typed: new(float64)}
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

	return &UniqueList[string]{Proxy: *proxy, limit: limit, typed: new(string)}, nil
}

func NewStringUniqueListWithDefault(key string, limit uint64, defaultElements []string, opts ...ProxyOption) (l *UniqueList[string], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[string]{Proxy: *proxy, limit: limit, typed: new(string)}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// UniqueList[time.Time] type

func NewTimeUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[time.Time], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &UniqueList[time.Time]{Proxy: *proxy, limit: limit, typed: new(time.Time)}, nil
}

func NewTimeUniqueListWithDefault(key string, limit uint64, defaultElements []time.Time, opts ...ProxyOption) (l *UniqueList[time.Time], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[time.Time]{Proxy: *proxy, limit: limit, typed: new(time.Time)}
	err = proxy.watch(func() error {
		_, err := l.Append(defaultElements...)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// UniqueList[KredisJSON] type

func NewJSONUniqueList(key string, limit uint64, opts ...ProxyOption) (*UniqueList[KredisJSON], error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &UniqueList[KredisJSON]{Proxy: *proxy, limit: limit, typed: new(KredisJSON)}, nil
}

func NewJSONUniqueListWithDefault(key string, limit uint64, defaultElements []KredisJSON, opts ...ProxyOption) (l *UniqueList[KredisJSON], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &UniqueList[KredisJSON]{Proxy: *proxy, limit: limit, typed: new(KredisJSON)}
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

func (l *UniqueList[T]) Remove(elements ...T) (err error) {
	for _, element := range elements {
		value := typeToInterface(element)
		if err = l.client.LRem(l.ctx, l.key, 0, value).Err(); err != nil {
			return
		}
	}

	return
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

func (l UniqueList[T]) Last() (T, bool) {
	slice, err := l.client.Do(l.ctx, "lrange", l.key, -1, -1).Slice()
	if err != nil || len(slice) < 1 {
		return any(*l.typed).(T), false
	}

	elements := make([]T, 1)
	copyCmdSliceTo(slice, elements)
	return elements[0], true
}

func (l UniqueList[T]) LastN(elements []T) (total int64, err error) {
	slice, err := l.client.Do(l.ctx, "lrange", l.key, -len(elements), -1).Slice()
	if err != nil {
		return
	}

	total = copyCmdSliceTo(slice, elements)
	return
}

func (l *UniqueList[T]) SetLimit(limit uint64) {
	l.limit = limit
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

	l.RefreshTTL()
	return llen.Val(), nil
}
