package simulate

import (
	"fmt"
	"math"
)

type SixEightCome struct {
	Debug bool
}

func (s *SixEightCome) Debugf(format string, args ...interface{}) {
	if s.Debug {
		fmt.Printf(format, args...)
	}
}

/* the strategy:
* assume $15 min bets
*
* always have 6/8 place bets working
* play come to get 4 numbers working at the same time.
* play pass if no come out roll
*
* $60 per roller.
* if profit, play the odds. TODO
* maybe combine with field bets
 */
func (s *SixEightCome) simulate(r *Round) {
	for {
		// stopping conditions.
		roll := NewRoll()
		s.Debugf("rolled: %s\n", roll)

		r.Rolls++

		// point established?
		switch roll.Value() {
		case 7:
			// clear comeBets, comeBets are always active
			for k, _ := range r.comeBets {
				delete(r.comeBets, k)
			}

			if r.point == nil {
				// no point established yet. pay 1x passline bet.
				// and keep the passbet on.
				r.Amount += r.passBet
			}

			if r.point != nil {
				// collect comebet win.
				r.Amount += (r.comeBet * 2)

				// don't really need to clean up because we're exiting
				s.Debugf("\t7 rolled before %d\n\n", *r.point)
				return
			}
		case 11:
			if r.point == nil && r.passBet > 0 {
				// passline winner, pocket and keep passBet on
				r.Amount += r.passBet
			}
			if r.point != nil && r.comeBet > 0 {
				// pay come bet, pocket and keep comeBet on
				r.Amount += r.comeBet
			}
		case 2, 3, 12:
			// clear come and place bets.
			if r.point == nil {
				r.passBet = 0
			}
			if r.point != nil {
				r.comeBet = 0
			}
			if r.point != nil && r.active() < 4 {
				// if no point has established.
				// and we don't have 4 numbers established yet.
				// bet come again.
				r.bet(&r.comeBet, r.minBet)
			}
		default:
			// check place bet wins
			// TODO: press place bets
			if v, ok := r.placeBets[roll.Value()]; ok && v > 0 && r.point != nil {
				pay := 0
				switch roll.Value() {
				case 4, 10:
					// pays out 9 to 5
					pay = (9 * v / 5)
				case 5, 9:
					// pays out 7 to 5
					pay = (7 * v / 5)
				case 6, 8:
					// pays out 7 to 6
					pay = (7 * v / 6)
				}
				r.Amount += pay
				r.Hits++
				s.Debugf("\twon place(%d): +%d\n", roll.Value(), pay)
			}

			// 4, 5, 6, 8, 9, 10
			// you should never have 6 or 8 vacated.
			if r.point == nil {
				// this was a comeout roll.
				r.setPoint(roll.Value())
				s.Debugf("\tpoint establisehd! %d\n", roll.Value())
			} else {
				// point exists, is it same as rolled value?
				if *r.point == roll.Value() {
					// clear the point
					s.Debugf("\tpoint won! %d\n", *r.point)
					r.point = nil
				}
			}

			// come bets pays out 1:1
			// TODO: come odds
			if v, ok := r.comeBets[roll.Value()]; ok && v > 0 {
				r.Amount += (v + v) // pays 1:1, collect original
				r.Hits++

				// remove the come bet.
				// >> let later code handle moving up come bet
				delete(r.comeBets, roll.Value())

				s.Debugf("\twon come(%d): +%d\n", roll.Value(), v+v)
			}

			// move up comeBet
			if r.comeBet != 0 {
				// new come bet point is established.
				// move the bet up to comeBets
				// >> take back the already placed come bet.
				r.comeBets[roll.Value()] = r.comeBet

				// >> let later code handle re-betting
				r.comeBet = 0
			}

			// check there's no duplicate 6/8 come bets
			for _, v := range []int{6, 8} {
				if _, ok := r.comeBets[v]; ok {
					if v, okk := r.placeBets[v]; okk {
						// remove the placebet
						r.Amount += v
						delete(r.placeBets, v)
					}
				} else {
					if _, ok := r.placeBets[v]; !ok {
						// add place bet, nearest multiple of 6 from minBet for 6/8
						bet := int(math.Ceil(float64(r.minBet)/float64(6))) * 6
						r.Amount -= bet
						r.placeBets[v] = bet
					}
				}
			}

			// Let's try something new to increase our hit%
			// if we have 4 numbers already, let's place on 4 and 10 if we are up.
			if r.active() == 4 && r.Amount >= r.initAmount {
				// only do this if we've loaded up.
				// TODO: check pctHits threshold, ie. only do this if our hit
				// % is less than a certain amount
				if _, ok := r.comeBets[4]; !ok {
					r.Amount -= r.minBet
					r.placeBets[4] = r.minBet
				} else {
					if _, ok := r.comeBets[10]; !ok {
						r.Amount -= r.minBet
						r.placeBets[10] = r.minBet
					}
				}
			}

			if r.comeBet == 0 {
				// if the point has not been established, we can't do comeBet
				if r.point != nil && r.active() < 4 {
					r.bet(&r.comeBet, r.minBet)
				}
				// TODO: pass bet?
			}

		}
		s.Debugf("\tplace bets: %+v\n", r.placeBets)
		s.Debugf("\tcome bets: %+v\n", r.comeBets)
		s.Debugf("\tcome: %d\n", r.comeBet)
		s.Debugf("\tamount: %d\n\n", r.Amount)
	}
}

func (s *SixEightCome) Simulate(amount int) *Round {
	r := &Round{
		placeBets:  map[int]int{},
		comeBets:   map[int]int{},
		minBet:     15,
		initAmount: amount,
		Amount:     amount, // amount per shooter
	}
	s.simulate(r)

	s.Debugf("shooter finished with $%d.\nwe hit %d times out of %d rolls.\ncost per roll(%.2f)\n",
		r.Amount, r.Hits, r.Rolls,
		(float64(amount)-float64(r.Amount))/float64(r.Rolls),
	)

	return r
}
