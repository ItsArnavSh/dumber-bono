package utility

import (
	"time"
)

type Mortal interface {
	GetExpiry() time.Time
}

type ExpiryQueue[M Mortal] struct {
	items []M
}

func NewQueue[M Mortal]() ExpiryQueue[M] {
	return ExpiryQueue[M]{items: []M{}}
}

func (e *ExpiryQueue[M]) Push(item M) {
	e.items = append(e.items, item)
}

func (e *ExpiryQueue[M]) Pop() (M, bool) {
	now := time.Now()
	for len(e.items) > 0 {
		item := e.items[0]
		e.items = e.items[1:]
		if now.Before(item.GetExpiry()) {
			return item, true
		}
	}
	var zero M
	return zero, false
}
