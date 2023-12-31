package kredis

import (
	"github.com/redis/go-redis/v9"
)

type Hash[T KredisTyped] struct {
	Proxy
	typed *T
}

// Hash[bool] type

func NewBoolHash(key string, opts ...ProxyOption) (*Hash[bool], error) {
	proxy, err := NewProxy(key, opts...)

	if err != nil {
		return nil, err
	}

	return &Hash[bool]{Proxy: *proxy, typed: new(bool)}, nil
}

func NewBoolHashWithDefault(key string, defaultElements map[string]bool, opts ...ProxyOption) (h *Hash[bool], err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	h = &Hash[bool]{Proxy: *proxy, typed: new(bool)}
	err = proxy.watch(func() error {
		_, err := h.Update(defaultElements)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

// generic Hash functions

func (h *Hash[T]) Get(field string) (T, bool) {
	return stringCmdToTyped[T](h.client.HGet(h.ctx, h.key, field), h.typed)
}

func (h *Hash[T]) Set(field string, entry T) (err error) {
	_, err = h.Update(map[string]T{field: entry})
	return
}

func (h *Hash[T]) Clear() error {
	return h.client.Del(h.ctx, h.key).Err()
}

func (h *Hash[T]) Update(entries map[string]T) (int64, error) {
	imap := make(map[string]interface{}, len(entries))

	for key, entry := range entries {
		imap[key] = typeToInterface(entry)
	}

	return h.client.HSet(h.ctx, h.key, imap).Result()
}

func (h Hash[T]) Delete(fields ...string) (int64, error) {
	return h.client.HDel(h.ctx, h.key, fields...).Result()
}

func (h *Hash[T]) ValuesAt(fields ...string) (values []T, err error) {
	slice, err := h.client.HMGet(h.ctx, h.key, fields...).Result()
	if err != nil {
		return
	}

	values = make([]T, len(slice))
	copyCmdSliceTo(slice, values)

	return
}

func (h *Hash[T]) Entries() (entries map[string]T, err error) {
	m, err := h.client.HGetAll(h.ctx, h.key).Result()
	if err != nil && err != redis.Nil {
		return
	}

	entries = make(map[string]T, len(m))
	for field, value := range m {
		if t, ok := stringToTyped(value, h.typed); ok {
			entries[field] = t
		}
	}

	return
}

func (h *Hash[T]) Keys() (keys []string, err error) {
	return h.client.HKeys(h.ctx, h.key).Result()
}

func (h *Hash[T]) Values() (values []T, err error) {
	slice, err := h.client.HVals(h.ctx, h.key).Result()
	if err != nil && err != redis.Nil {
		return
	}

	values = make([]T, len(slice))
	for idx, s := range slice {
		if t, ok := stringToTyped(s, h.typed); ok {
			values[idx] = t
		}
	}

	return
}
