package simulate

import (
	"fmt"
	"math/rand"
	"time"
)

type Naive struct {
	minBet   int
	maxLoss  int
	maxRolls int
}

func NewNaive() *Naive {
	return &Naive{
		minBet:   15,
		maxLoss:  -500,
		maxRolls: 500,
	}
}

// change up odd bet depending on the point value
func (n *Naive) getOddBet(point int) int {
	switch point {
	case 4, 10:
		return n.minBet * 3
	case 5, 9:
		return n.minBet * 1
	case 6, 8:
		return n.minBet * 1
	}
	return 0
}

func (n *Naive) simulate(r *Round) int {
	roll := NewRoll()
	//fmt.Printf("round %d: %s (point %d)\n", r.n, roll, point)

	switch roll.Value() {
	case 7, 11:
		// 22.22% -> 8/36 -> 2/9
		if r.point != nil {
			// point is established we lose our minBet and oddBet
			if r.takeOdds {
				return -(n.minBet + n.getOddBet(*r.point))
			}
			// lose just the pass bet
			return -1 * n.minBet
		}
		// on comeout roll. we win.
		return n.minBet
	case 2, 3, 12:
		// 11.11% -> 4/36 -> 1/9
		// lose on come out roll
		return -1 * n.minBet
	default:
		// other: 4, 5, 6, 8, 9, 10
		// odds: 66.67 -> 2/3
		if r.point != nil && roll.Value() == *r.point {
			if r.takeOdds {
				mod := 0
				// pass with odds
				switch *r.point {
				case 4, 10:
					mod = 2 // risk 40 for 70
				case 5, 9:
					mod = 3 / 2 // risk 40 for 55
				case 6, 8:
					mod = 6 / 5 // risk 40 for 46
				}
				return n.minBet + n.getOddBet(*r.point)*mod
			}
			// just pass winning
			return n.minBet
		}
		if r.point == nil {
			v := roll.Value()
			r.point = &v
		}
		r.Rolls++
		return n.simulate(r)
	}

	// unreachable
	return 0
}

func (n *Naive) Simulate() {
	g := &boolgen{src: rand.NewSource(time.Now().UnixNano())}
	var rounds []*Round

	rolls := 0
	comeout := 0
	w := 0  // wins
	l := 0  // loses
	vv := 0 // accumulative $$
	to := 0 // # of odds took

	for comeout < n.maxRolls && rolls < n.maxRolls {
		r := &Round{
			takeOdds: g.Bool(),
		}
		if r.takeOdds {
			to++
		}

		v := n.simulate(r)
		if v > 0 {
			r.won = true
			w++
		} else if v < 0 {
			l++
		}
		rolls += r.Rolls
		vv += v

		if vv < n.maxLoss {
			break
		}
		comeout++

		rounds = append(rounds, r)
	}
	wp := float64(w) / float64(rolls)
	vp := float64(vv) / float64(rolls)

	fmt.Printf("won %d(%.2f) vs lose %d\n", w, wp, l)
	fmt.Printf("played %d rolls across %d comeouts, for $%d\n", rolls, comeout, vv)
	fmt.Printf("avg %.2f rolls per comeout\n", float64(rolls)/float64(comeout))
	fmt.Printf("%d rolls took the odds\n", to)
	fmt.Printf("avg loss per roll: %.2f\n", vp)

	//oW := 0
	//oL := 0
	//for _, v := range rounds {
	//s := "L"
	//if v.won {
	//s = "W"
	//}
	//fmt.Printf("%d%s", v.n, s)

	//if v.n == 0 {
	//if v.won {
	//oW++
	//} else {
	//oL++
	//}
	//}
	//}

	//fmt.Printf("\n0W %.2f, 0L %.2f", float64(oW)/float64(comeout), float64(oL)/float64(comeout))
	fmt.Println("\n")
}
