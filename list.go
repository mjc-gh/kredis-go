package kredis

import (
	"errors"
	"fmt"
	"time"
)

type List[T KredisTyped] struct {
	Proxy
}

func NewStringList(key string, options Options) (*List[string], error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &List[string]{Proxy: *proxy}, nil
}

func NewTimeList(key string, options Options) (*List[time.Time], error) {
	proxy, err := NewProxy(key, options)

	if err != nil {
		return nil, err
	}

	return &List[time.Time]{Proxy: *proxy}, nil
}

func (l *List[T]) Append(elements ...T) (int64, error) {
	values := l.elementsToValues(elements)

	fmt.Println(len(values))
	if len(values) < 1 {
		return 0, errors.New("elements is empty")
	}

	llen, err := l.client.RPush(l.ctx, l.key, values...).Result()

	fmt.Println(values)
	fmt.Println(llen)
	fmt.Println(err)

	if err != nil {
		return 0, err
	}

	return llen, nil
}

func (l List[T]) Prepend(elements ...T) (int64, error) {
	values := l.elementsToValues(elements)
	llen, err := l.client.LPush(l.ctx, l.key, values...).Result()

	if err != nil {
		return 0, err
	}

	return llen, nil
}

//func (l *List[T]) Clear() error {
//}

func (l *List[T]) elementsToValues(elements []T) []interface{} {
	values := make([]interface{}, len(elements))

	for i, e := range elements {
		switch any(e).(type) {
		case time.Time:
			values[i] = any(e).(time.Time).Format(time.RFC3339Nano)
		default:
			values[i] = e
		}
	}

	return values
}

func (l *List[T]) Elements(elements []T) (int, error) {
	var total int

	lrange, err := l.client.LRange(l.ctx, l.key, 0, -1).Result()

	if err != nil {
		return total, err
	}

	for i, e := range lrange {
		if i == len(elements) {
			break
		}

		switch any(elements[i]).(type) {
		case time.Time:
			t, _ := time.Parse(time.RFC3339Nano, e)

			elements[i] = any(t).(T)
		default:
			elements[i] = any(e).(T)
		}

		total = total + 1
	}

	return total, nil
}
