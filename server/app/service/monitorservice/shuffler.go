package monitor

import (
	"math/rand"
)

// InformFunc is any of your no-arg informational callouts.
type InformFunc func()

// Shuffler runs a random subset of registered functions each cycle,
// avoiding immediate repeats and biasing against recently-run functions.
type Shuffler struct {
	funcs      []InformFunc
	names      []string    // for logging/debugging which ran
	lastRunIdx map[int]int // index -> "cycles ago it last ran"
	cycle      int
}

func NewShuffler() *Shuffler {
	return &Shuffler{
		lastRunIdx: make(map[int]int),
	}
}

func (s *Shuffler) Register(name string, fn InformFunc) {
	s.names = append(s.names, name)
	s.funcs = append(s.funcs, fn)
}

// RunSubset picks `count` functions to run this cycle, weighted to prefer
// ones that haven't run recently (a "shuffle-bag" style pick), and runs them.
func (s *Shuffler) RunSubset(count int) {
	s.cycle++
	if count > len(s.funcs) {
		count = len(s.funcs)
	}

	// Build weighted candidate list: functions not run recently get higher weight.
	type candidate struct {
		idx    int
		weight int
	}
	candidates := make([]candidate, len(s.funcs))
	for i := range s.funcs {
		lastRan, ok := s.lastRunIdx[i]
		age := s.cycle // never run = max weight (very "stale")
		if ok {
			age = s.cycle - lastRan
		}
		// weight grows the longer it's been since this function last ran
		weight := age + 1
		candidates[i] = candidate{idx: i, weight: weight}
	}

	chosen := make(map[int]bool)
	for len(chosen) < count {
		total := 0
		for _, c := range candidates {
			if !chosen[c.idx] {
				total += c.weight
			}
		}
		if total == 0 {
			break
		}
		r := rand.Intn(total)
		for _, c := range candidates {
			if chosen[c.idx] {
				continue
			}
			r -= c.weight
			if r < 0 {
				chosen[c.idx] = true
				break
			}
		}
	}

	for idx := range chosen {
		s.funcs[idx]()
		s.lastRunIdx[idx] = s.cycle
	}
}
