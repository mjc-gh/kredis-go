package kredis

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

func (l *List[T]) Append(elements ...T) (int64, error) {
	values := l.elementsToValues(elements)
	llen, err := l.client.RPush(l.ctx, l.key, values...).Result()

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
		// TODO do typing casting here depending on T
		values[i] = e
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

		elements[i] = any(e).(T)
		total = total + 1
	}

	return total, nil
}
