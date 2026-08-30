package monitor

import "testing"

func TestShufflerRegisterAndRunSubset(t *testing.T) {
	s := NewShuffler()

	var calls []string
	s.Register("a", func() { calls = append(calls, "a") })
	s.Register("b", func() { calls = append(calls, "b") })
	s.Register("c", func() { calls = append(calls, "c") })

	s.RunSubset(2)
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
}

func TestShufflerRunSubsetMoreThanRegistered(t *testing.T) {
	s := NewShuffler()
	s.Register("a", func() {})
	s.Register("b", func() {})

	// requesting more than registered should not panic and run all
	s.RunSubset(10)

	if len(s.funcs) != 2 {
		t.Fatalf("expected 2 registered funcs, got %d", len(s.funcs))
	}
}

func TestShufflerRunSubsetEmpty(t *testing.T) {
	s := NewShuffler()
	// running on empty shuffler must not panic
	s.RunSubset(3)
	if s.cycle != 1 {
		t.Fatalf("cycle = %d, want 1", s.cycle)
	}
}

func TestShufflerTracksLastRun(t *testing.T) {
	s := NewShuffler()
	s.Register("a", func() {})
	s.Register("b", func() {})

	s.RunSubset(2)
	if len(s.lastRunIdx) != 2 {
		t.Fatalf("expected both funcs to be marked as run, got %d", len(s.lastRunIdx))
	}
	if s.cycle != 1 {
		t.Fatalf("cycle = %d, want 1", s.cycle)
	}
}
