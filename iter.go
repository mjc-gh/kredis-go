package kredis

import (
	"github.com/redis/go-redis/v9"
)

type iterator[T KredisTyped] interface {
	values() []interface{}
	valuesWithScoring(scoreable, bool) []redis.Z
	unique() []interface{}
	uniqueMap() map[T]struct{}
}

type iter[T KredisTyped] struct {
	elements []T // rename "elements" to "items" (to be less list-y)
}

func newIter[T KredisTyped](elements []T) iterator[T] {
	return &iter[T]{elements}
}

func (i *iter[T]) values() []interface{} {
	values := make([]interface{}, len(i.elements))

	for idx, e := range i.elements {
		values[idx] = typeToInterface(e)
	}

	return values
}

type scoreable interface {
	prependScore(int) float64
	appendScore(int) float64
}

func (i *iter[T]) valuesWithScoring(scorer scoreable, prepended bool) []redis.Z {
	values := make([]redis.Z, len(i.elements))

	for idx, e := range i.elements {
		if prepended {
			values[idx].Score = scorer.prependScore(idx)
		} else {
			values[idx].Score = scorer.appendScore(idx)
		}

		values[idx].Member = typeToInterface(e)
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
