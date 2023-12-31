package kredis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterValues(t *testing.T) {
	i := newIter([]kredisJSON{
		*NewKredisJSON(`{"a":"b"}`),
		*NewKredisJSON(`{"c":"d"}`),
		*NewKredisJSON(`{"a":"b"}`),
	})

	assert.Len(t, i.values(), 3)
}

func TestIterUnique(t *testing.T) {
	i := newIter([]kredisJSON{
		*NewKredisJSON(`{"a":"b"}`),
		*NewKredisJSON(`{"c":"d"}`),
		*NewKredisJSON(`{"a":"b"}`),
	})

	assert.Len(t, i.unique(), 2)
}
