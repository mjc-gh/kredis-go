package kredis

import (
	"strconv"
	"time"
)

type List[T KredisTyped] struct {
	Proxy
}

// TODO add support for default values
// integer_list = Kredis.list "myintegerlist",
//   typed: :integer,
//   default: [ 1, 2, 3 ] # => EXISTS? myintegerlist, RPUSH myintegerlist "1" "2" "3"

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

func NewIntListWithDefault(key string, defaultElements []int, opts ...ProxyOption) (l *List[int], err error) {
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

func NewIntegerList(key string, opts ...ProxyOption) (*List[int], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[int]{Proxy: *proxy}, nil
}

func NewStringList(key string, opts ...ProxyOption) (*List[string], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[string]{Proxy: *proxy}, nil
}

func NewTimeList(key string, opts ...ProxyOption) (*List[time.Time], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[time.Time]{Proxy: *proxy}, nil
}

func NewJSONList(key string, opts ...ProxyOption) (*List[kredisJSON], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &List[kredisJSON]{Proxy: *proxy}, nil
}

func (l *List[T]) Elements(elements []T, opts ...RangeOption) (total int, err error) {
	rangeOptions := RangeOptions{0}
	for _, opt := range opts {
		opt(&rangeOptions)
	}

	// TODO should we read all elements?
	stop := rangeOptions.start + int64(len(elements))
	lrange, err := l.client.LRange(l.ctx, l.key, rangeOptions.start, stop).Result()
	if err != nil {
		return
	}

	for i, e := range lrange {
		if i == len(elements) {
			break
		}

		switch any(elements[i]).(type) {
		case bool:
			b, _ := strconv.ParseBool(e)

			elements[i] = any(b).(T)
		case int:
			n, _ := strconv.Atoi(e)

			elements[i] = any(n).(T)
		case kredisJSON:
			j := kredisJSON(e)

			elements[i] = any(j).(T)
		case time.Time:
			t, _ := time.Parse(time.RFC3339Nano, e)

			elements[i] = any(t).(T)
		default:
			elements[i] = any(e).(T)
		}

		total = total + 1
	}

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

// TODO add function last(n = 1) ??
// https://github.com/rails/kredis/blob/2ccc5c6bf59e5d38870de45a03e9491a3dc8c397/lib/kredis/types/list.rb#L32-L34
