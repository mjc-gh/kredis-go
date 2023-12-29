package kredis

// Loosely based on https://github.com/polyfloyd/go-iterator
type iterator[T KredisTyped] interface {
	next() (T, bool)
	values() []interface{}
	unique() []interface{}
	uniqueMap() map[T]struct{}
}

type iter[T KredisTyped] struct {
	elements []T // rename "elements" to "items" (to be less list-y)
}

func newIter[T KredisTyped](elements []T) iterator[T] {
	return &iter[T]{elements}
}

func (i *iter[T]) next() (elem T, ok bool) {
	if len(i.elements) == 0 {
		return // no elements - return empty values
	}

	elem = i.elements[0]
	i.elements = i.elements[1:]

	return elem, true
}

func (i *iter[T]) values() []interface{} {
	values := make([]interface{}, len(i.elements))

	for idx, e := range i.elements {
		values[idx] = typeToInterface(e)
	}

	return values
}

func (i *iter[T]) unique() []interface{} {
	m := make(map[T]struct{})
	values := make([]interface{}, 0)

	for _, e := range i.elements {
		if _, exists := m[e]; exists {
			continue
		}

		m[e] = struct{}{}
		values = append(values, typeToInterface(e))
	}

	return values
}

func (i *iter[T]) uniqueMap() (m map[T]struct{}) {
	if len(i.elements) == 0 {
		return
	}

	m = make(map[T]struct{})
	for _, e := range i.elements {
		m[e] = struct{}{}
	}

	return
}
