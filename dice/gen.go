package dice

import (
	"math/rand"
	"sync"
	"time"
)

var (
	d        = []int{1, 2, 3, 4, 5, 6}
	onlyOnce sync.Once
)

type Generator interface {
	Roll() *Pair
}

func init() {
	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano()) // only run once
	})
}
