package kredis

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO does this need to be exported??
// type kredisJSON []byte
type kredisJSON struct {
	s string
}

type KredisTyped interface {
	comparable
	~bool | ~int | ~string | kredisJSON | time.Time
}

// kredisJSON is a small struct wrapper for dealing with JSON strings
func NewKredisJSON(jsonStr string) *kredisJSON {
	var kj kredisJSON = kredisJSON{jsonStr}

	return &kj
}

func (kj kredisJSON) String() string {
	return kj.s
}

func (kj *kredisJSON) Unmarshal(data *interface{}) error {
	err := json.Unmarshal([]byte(kj.s), data)
	if err != nil {
		return err
	}

	return nil
}

// convert an interface{} value to a KredisTyped value
func typeToInterface[T KredisTyped](t T) interface{} {
	switch any(t).(type) {
	case time.Time:
		return any(t).(time.Time).Format(time.RFC3339Nano)
	case kredisJSON:
		return any(t).(kredisJSON).String()
	default:
		return t
	}
}

// convert a string value to a KredisTyped value
func stringToTyped[T KredisTyped](value string, typed *T) (T, bool) {
	switch any(*typed).(type) {
	case bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return any(false).(T), false
		}
		return any(b).(T), true

	case int:
		n, err := strconv.Atoi(value)
		if err != nil {
			return any(0).(T), false
		}
		return any(n).(T), true

	case time.Time:
		t, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return any(time.Time{}).(T), false
		}
		return any(t).(T), true

	case string:
		return any(value).(T), true

	case kredisJSON:
		return any(NewKredisJSON(value)).(T), true
	}

	return any(*typed).(T), false
}

// redis.StringCmd has most of the conversion functions we need for converting
// to a KredisTyped. this is only used with Set or Hash.
func stringCmdToTyped[T KredisTyped](cmd *redis.StringCmd, typed *T) (T, bool) {
	if cmd.Err() == redis.Nil {
		goto Empty
	}

	switch any(*typed).(type) {
	case bool:
		b, err := cmd.Bool()
		if err != nil {
			return any(false).(T), false
		}
		return any(b).(T), true

	case int:
		n, err := cmd.Int()
		if err != nil {
			return any(0).(T), false
		}
		return any(n).(T), true

	case time.Time:
		t, err := cmd.Time()
		if err != nil {
			return any(time.Time{}).(T), false
		}
		return any(t).(T), true

	case string:
		return any(cmd.Val()).(T), true

	case kredisJSON:
		return any(NewKredisJSON(cmd.Val())).(T), true
	}

Empty:
	return any(*typed).(T), false
}

// used in most collection types for copying a slice of interfaces to a slice
// of KredisTyped.
func copyCmdSliceTo[T KredisTyped](slice []interface{}, dst []T) (total int64) {
	for i, e := range slice {
		if i == len(dst) {
			break
		}

		switch any(dst[i]).(type) {
		case bool:
			b, _ := strconv.ParseBool(e.(string))

			dst[i] = any(b).(T)
		case int:
			n, _ := strconv.Atoi(e.(string))

			dst[i] = any(n).(T)
		case kredisJSON:
			j := NewKredisJSON(e.(string))

			dst[i] = any(*j).(T)
		case time.Time:
			t, _ := time.Parse(time.RFC3339Nano, e.(string))

			dst[i] = any(t).(T)
		default:
			dst[i] = any(e).(T)
		}

		total += 1
	}

	return
}
