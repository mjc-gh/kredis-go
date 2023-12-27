package kredis

import "time"

// Loosely based on https://github.com/polyfloyd/go-iterator
type iterator[T KredisTyped] interface {
	next() (T, bool)
	values() []interface{}
}

type iter[T KredisTyped] struct {
	elements []T
}

func newIter[T KredisTyped](elements []T) iterator[T] {
	return &iter[T]{elements}
}

func (i *iter[T]) next() (val T, ok bool) {
	if len(i.elements) == 0 {
		return // no elements - return empty values
	}

	elem := i.elements[0]
	i.elements = i.elements[1:]

	return elem, true
}

func (i *iter[T]) values() []interface{} {
	values := make([]interface{}, len(i.elements))

	for i, e := range i.elements {
		switch any(e).(type) {
		case time.Time:
			values[i] = any(e).(time.Time).Format(time.RFC3339Nano)
		case kredisJSON:
			values[i] = any(e).(kredisJSON).String()
		default:
			values[i] = e
		}
	}

	return values
}
