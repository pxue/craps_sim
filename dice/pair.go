package dice

import "fmt"

type Pair struct {
	D1, D2 int
}

func (p *Pair) String() string {
	return fmt.Sprintf("(%d, %d)->%d", p.D1, p.D2, p.D1+p.D2)
}

func (p *Pair) Value() int {
	return p.D1 + p.D2
}
