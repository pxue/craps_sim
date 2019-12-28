package simulate

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var dice = []int{1, 2, 3, 4, 5, 6}
var onlyOnce sync.Once

type Pair struct {
	d1, d2 int
}

func (p *Pair) String() string {
	return fmt.Sprintf("(%d, %d)->%d", p.d1, p.d2, p.d1+p.d2)
}

func (p *Pair) Value() int {
	return p.d1 + p.d2
}

func NewRoll() *Pair {
	d1 := dice[rand.Intn(len(dice))]
	d2 := dice[rand.Intn(len(dice))]
	return &Pair{d1, d2}
}

func init() {
	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano()) // only run once
	})
}
