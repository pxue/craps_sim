package simulate

import (
	"fmt"
	"math"

	"github.com/pxue/craps/dice"
)

type SixEightCome struct {
	Debug bool
	gen   dice.Generator
}

func NewSixEight(debug bool, gen ...dice.Generator) *SixEightCome {
	s := &SixEightCome{
		Debug: debug,
		gen:   &dice.Simple{},
	}
	if len(gen) > 0 {
		s.gen = gen[0]
	}
	return s
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
		roll := s.gen.Roll()
		if roll == nil {
			// mock roll. no more rolls.
			return
		}
		s.Debugf("rolled: %s\n", roll)

		r.Occurance[roll.Value()]++
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
				// collect and clean up comebet win.
				s.Debugf("\twon come: +%d\n", r.comeBet*2)
				r.Amount += (r.comeBet * 2)
				r.comeBet = 0

				// clean up passline Bet
				r.passBet = 0

				for k, _ := range r.comeBets {
					// comebets lose
					delete(r.comeBets, k)
				}

				// cleanup place bets
				for k, _ := range r.placeBets {
					delete(r.placeBets, k)
				}

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
				r.Hits[roll.Value()]++
				s.Debugf("\twon place(%d): +%d\n", roll.Value(), pay)

				// NOTE: let's test pressing the bets if we don't have a comeBet
				// active.
				// the check is that, every 2nd hit, we press.
				//if r.Occurance[roll.Value()]%2 == 0 && r.comeBet == 0 {
				//s.Debugf("\tpressing place(%d) by 1 unit.\n", roll.Value())
				//switch roll.Value() {
				//case 6, 8:
				//r.Amount -= 6
				//r.placeBets[roll.Value()] += 6
				//default:
				//r.Amount -= 5
				//r.placeBets[roll.Value()] += 5
				//}
				//}
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
					// pay pass line, reset to zero
					r.Amount += 2 * r.passBet
					r.passBet = 0

					// clear the point
					s.Debugf("\tpoint won! %d\n", *r.point)
					r.point = nil

					// make a pass line bet if applicable.
					if r.active() < 4 {
						r.bet(&r.passBet, r.minBet)
					}
				}
			}

			// come bets pays out 1:1
			// TODO: come odds
			if v, ok := r.comeBets[roll.Value()]; ok && v > 0 {
				r.Amount += (v + v) // pays 1:1, collect original
				r.Hits[roll.Value()]++

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
			for _, roll := range []int{6, 8} {
				if _, ok := r.comeBets[roll]; ok {
					if bet, okk := r.placeBets[roll]; okk {
						// remove the placebet
						r.Amount += bet
						delete(r.placeBets, roll)
					}
				} else {
					if _, ok := r.placeBets[roll]; !ok {
						// add place bet, nearest multiple of 6 from minBet for 6/8
						bet := int(math.Ceil(float64(r.minBet)/float64(6))) * 6
						r.Amount -= bet
						r.placeBets[roll] = bet
					}
				}
			}

			// Let's try something new to increase our hit%
			//if we have 4 numbers already, let's place on 4 and 10 if we are up.
			//if r.active() == 4 && r.Amount >= r.initAmount {
			//// only do this if we've loaded up.
			//// TODO: check pctHits threshold, ie. only do this if our hit
			//// % is less than a certain amount
			//if _, ok := r.comeBets[4]; !ok {
			//r.Amount -= r.minBet
			//r.placeBets[4] = r.minBet
			//} else {
			//if _, ok := r.comeBets[10]; !ok {
			//r.Amount -= r.minBet
			//r.placeBets[10] = r.minBet
			//}
			//}
			//}

			if r.comeBet == 0 {
				// if the point has not been established, we can't do comeBet
				if r.point != nil && r.active() < 4 {
					r.bet(&r.comeBet, r.minBet)
				}
				// TODO: pass bet?
			}

			if r.passBet > 0 && r.point == nil {
				// double check here. we've consolidated all the bets.
				// if we've got 4 active numbers. let's not make a passBet
				if r.active() >= 4 {
					r.Amount += r.passBet
					r.passBet = 0
				}
			}

		}
		s.Debugf("\tplace bets: %+v\n", r.placeBets)
		s.Debugf("\tcome bets: %+v\n", r.comeBets)
		s.Debugf("\tcome: %d\n", r.comeBet)
		s.Debugf("\tpass: %d\n", r.passBet)
		if r.point != nil {
			s.Debugf("\tpoint: %d\n", *r.point)
		}
		s.Debugf("\tamount: %d\n\n", r.Amount)
	}
}

func (s *SixEightCome) Simulate(amount int, init ...*Round) *Round {
	var r *Round
	if len(init) > 0 && init[0] != nil {
		r = init[0]
	} else {
		r = &Round{
			placeBets:  map[int]int{},
			comeBets:   map[int]int{},
			minBet:     15,
			initAmount: amount,
			Amount:     amount, // amount per shooter
			Occurance:  map[int]int{},
			Hits:       map[int]int{},
		}
	}
	s.simulate(r)

	s.Debugf("shooter finished with $%d.\nwe hit %d times out of %d rolls.\ncost per roll(%.2f)\n",
		r.Amount, len(r.Hits), r.Rolls,
		(float64(amount)-float64(r.Amount))/float64(r.Rolls),
	)

	return r
}
