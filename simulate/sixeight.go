package simulate

import "fmt"

type SixEightCome struct {
}

/* the strategy:
* assume $15 min bets
*
* always have 6/8 place bets working
* play come to get 4 numbers working at the same time.
* play pass if no come out roll
*
* $30 per roller.
* if profit, play the odds.
 */
func (s *SixEightCome) simulate(r *round) {
	for {
		// stopping conditions.
		// TODO: other?
		roll := NewRoll()
		fmt.Printf("rolled: %s\n", roll)

		switch roll.Value() {
		case 7:
			if r.point == nil {
				r.amount += r.passBet * 2
			}
			r.passBet = 0
			// "lose" clear the board, comeBets and placeBets
			for k, _ := range r.comeBets {
				delete(r.comeBets, k)
			}
			for k, _ := range r.placeBets {
				delete(r.placeBets, k)
			}

			// exit.
			return
		case 11:
			if r.point == nil {
				// clear passBet
				r.amount += r.passBet * 2
				r.passBet = 0
			}
		case 2, 3, 12:
			// clear come and place bets.
			r.comeBet = 0
			if r.point == nil {
				r.passBet = 0
			}

			if r.point != nil && len(r.placeBets)+len(r.comeBets) < 4 {
				// most likely need a come
				r.comeBet = 5
				r.amount -= 5
			}
		default:
			// 4, 5, 6, 8, 9, 10
			// you should never have 6 or 8 vacated.
			if r.point == nil {
				// this was a comeout roll.
				point := roll.Value()
				r.point = &point
				fmt.Printf("\tpoint establisehd! %d\n", point)
			} else {
				// point exists, is it same as rolled value?
				if *r.point == roll.Value() {
					// clear the point
					fmt.Printf("\tpoint won! %d\n", *r.point)
					r.point = nil
				}
			}

			// collect place bet wins
			if v, ok := r.placeBets[roll.Value()]; ok && v > 0 {
				// pays out 7:6
				r.amount += 7
				// keep the bet on.
				fmt.Printf("\twon place(%d): +%d\n", roll.Value(), 7)
			}

			// collect come bet wins
			if v, ok := r.comeBets[roll.Value()]; ok && v > 0 {
				// collect 1:1
				r.amount += v * 2
				delete(r.comeBets, roll.Value())
				fmt.Printf("\twon come(%d): +%d\n", roll.Value(), v*2)
			}

			if r.comeBet == 0 {
				// if the point has not been established, we can't do comeBet
				if r.point != nil && (len(r.placeBets)+len(r.comeBets) < 4) {
					r.comeBet = 5
					r.amount -= 5
				}
				// TODO: use a pass bet
			} else {
				// move the bet up to comeBets
				if r.comeBet > 0 {
					// move it up to point.
					r.comeBets[roll.Value()] = 5 // always 5
					r.comeBet = 0

					// if roll is 6 or 8. pull down the 6 or 8 we don't need it
					// twice

					if roll.Value() == 6 || roll.Value() == 8 {
						r.amount += r.placeBets[roll.Value()]
						fmt.Printf("\t pulled down place bet(%d): +6\n", roll.Value())
						delete(r.placeBets, roll.Value())
					}
				}
				// now check how many number we've got working.
				if r.point != nil && len(r.placeBets)+len(r.comeBets) < 4 {
					// most likely need a come
					r.comeBet = 5
					r.amount -= 5
				}
			}
			fmt.Printf("\tcome bets: %+v\n", r.comeBets)

			// make sure placeBets are on only if there isn't a come bet
			// on 6 or 8
			if v, ok := r.placeBets[6]; !ok || v == 0 {
				if _, ok := r.comeBets[6]; !ok {
					r.placeBets[6] = 6
					r.amount -= 6
				}
			}
			if v, ok := r.placeBets[8]; !ok || v == 0 {
				if _, ok := r.comeBets[8]; !ok {
					r.placeBets[8] = 6
					r.amount -= 6
				}
			}

			// check again, if placeBets + comeBets == 4, make sure
			// we don't have a come bet again.
			if len(r.placeBets)+len(r.comeBets) == 4 {
				r.amount += r.comeBet
				r.comeBet = 0
			}
			fmt.Printf("\tplace bets: %+v\n", r.placeBets)

			fmt.Printf("\tcome: -%d\n", r.comeBet)
			fmt.Printf("\tamount: %d -> profit(%d)\n\n", r.amount, r.amount-30)
		}
	}
}

func (s *SixEightCome) Simulate() int {
	r := &round{
		placeBets: map[int]int{},
		comeBets:  map[int]int{},
		amount:    30, // per shooter
	}
	s.simulate(r)
	fmt.Printf("shooter finished: %d\n", r.amount)

	return r.amount
}
