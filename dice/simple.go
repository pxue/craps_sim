package dice

import (
	"math/rand"
)

// Simple dice uses global rand source to generate dice rolls
type Simple struct {
	rand rand.Rand
}

func (s *Simple) Roll() *Pair {
	d1 := d[rand.Intn(len(d))]
	d2 := d[rand.Intn(len(d))]
	return &Pair{d1, d2}
}
