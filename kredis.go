package kredis

import (
	"encoding/json"
	"time"
)

// TODO does this need to be exported??
type kredisJSON []byte

type KredisTyped interface {
	~bool | ~int | ~string | kredisJSON | time.Time
}

func NewKredisJSON(jsonStr string) *kredisJSON {
	var kj kredisJSON = kredisJSON(jsonStr)

	return &kj
}

func (kj kredisJSON) String() string {
	return string(kj)
}

func (kj *kredisJSON) Unmarshal(data *interface{}) error {
	err := json.Unmarshal(*kj, data)

	if err != nil {
		return err
	}

	return nil
}
