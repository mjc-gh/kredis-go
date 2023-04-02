package kredis

import "time"

type KredisTyped interface {
	~int | ~string | time.Time
}
