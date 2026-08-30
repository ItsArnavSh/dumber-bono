package utility

import (
	"testing"
	"time"
)

type testMortal struct {
	expiry time.Time
}

func (t testMortal) GetExpiry() time.Time { return t.expiry }

func TestExpiryQueuePopEmpty(t *testing.T) {
	q := NewQueue[testMortal]()
	if _, ok := q.Pop(); ok {
		t.Fatal("expected Pop on empty queue to return false")
	}
}

func TestExpiryQueuePopsUnExpiredItem(t *testing.T) {
	q := NewQueue[testMortal]()
	q.Push(testMortal{expiry: time.Now().Add(time.Hour)})

	item, ok := q.Pop()
	if !ok {
		t.Fatal("expected Pop to return an item")
	}
	if !item.GetExpiry().After(time.Now()) {
		t.Fatal("expected the popped item to be un-expired")
	}
}

func TestExpiryQueueSkipsExpiredItems(t *testing.T) {
	q := NewQueue[testMortal]()
	q.Push(testMortal{expiry: time.Now().Add(-time.Hour)})
	q.Push(testMortal{expiry: time.Now().Add(time.Hour)})

	item, ok := q.Pop()
	if !ok {
		t.Fatal("expected Pop to return the un-expired item")
	}
	if !item.GetExpiry().After(time.Now()) {
		t.Fatal("expected the popped item to be the un-expired one")
	}
}

func TestExpiryQueueAllExpired(t *testing.T) {
	q := NewQueue[testMortal]()
	q.Push(testMortal{expiry: time.Now().Add(-time.Hour)})
	q.Push(testMortal{expiry: time.Now().Add(-time.Minute)})

	if _, ok := q.Pop(); ok {
		t.Fatal("expected Pop to return false when all items are expired")
	}
}

func TestExpiryQueueFIFOOrder(t *testing.T) {
	q := NewQueue[testMortal]()
	now := time.Now()
	first := testMortal{expiry: now.Add(time.Hour)}
	second := testMortal{expiry: now.Add(2 * time.Hour)}
	q.Push(first)
	q.Push(second)

	item, ok := q.Pop()
	if !ok || item.expiry != first.expiry {
		t.Fatal("expected first pushed item to be popped first")
	}
	item, ok = q.Pop()
	if !ok || item.expiry != second.expiry {
		t.Fatal("expected second pushed item to be popped second")
	}
}
