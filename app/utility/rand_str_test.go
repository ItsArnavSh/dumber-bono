package utility

import (
	"testing"
)

func TestRandomStringSelectsFromItems(t *testing.T) {
	items := []string{"a", "b", "c"}
	got, err := RandomString(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, item := range items {
		if got == item {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("got %q, want one of %v", got, items)
	}
}

func TestRandomStringEmptySlice(t *testing.T) {
	if _, err := RandomString(nil); err == nil {
		t.Fatal("expected error for empty slice")
	}
}

func TestRandomStringSingleItem(t *testing.T) {
	got, err := RandomString([]string{"only"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "only" {
		t.Fatalf("got %q, want %q", got, "only")
	}
}

func TestRandomStringNeverReturnsOutOfSlice(t *testing.T) {
	items := []string{"a", "b", "c", "d"}
	seen := make(map[string]bool)
	for range 1000 {
		got, err := RandomString(items)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		seen[got] = true
	}
	for _, item := range items {
		if !seen[item] {
			t.Fatalf("item %q was never selected over 1000 draws", item)
		}
	}
}
