package badger

import (
	"path/filepath"
	"testing"
)

func TestBadgerSetGet(t *testing.T) {
	db, err := NewBadger(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open badger: %v", err)
	}
	defer func() { _ = db.Close() }()

	var out int
	if err := db.Set("ns", "key", 42); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if err := db.Get("ns", "key", &out); err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if out != 42 {
		t.Fatalf("got %d, want 42", out)
	}
}

func TestBadgerGetMissing(t *testing.T) {
	db, err := NewBadger(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open badger: %v", err)
	}
	defer func() { _ = db.Close() }()

	var out int
	if err := db.Get("ns", "nope", &out); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestBadgerList(t *testing.T) {
	db, err := NewBadger(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open badger: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.BulkSet("cars", map[string]any{"0": 10, "1": 20}); err != nil {
		t.Fatalf("BulkSet failed: %v", err)
	}

	items, err := db.List("cars")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("List returned %d items, want 2", len(items))
	}
}

func TestBadgerBulkSetAndGet(t *testing.T) {
	db, err := NewBadger(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open badger: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.BulkSet("lap", map[string]any{"0": "first", "1": "second"}); err != nil {
		t.Fatalf("BulkSet failed: %v", err)
	}

	var out string
	if err := db.Get("lap", "1", &out); err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if out != "second" {
		t.Fatalf("got %q, want %q", out, "second")
	}
}
