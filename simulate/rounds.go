package simulate

// a round is bets before a 7 rolls again
type Round struct {
	point *int // the current point

	passBet   int         // the current pass bet
	comeBet   int         // the current come bet
	placeBets map[int]int // place bets
	comeBets  map[int]int // active come bets
	comeOdds  map[int]int // the current come bet odds
	passOdds  map[int]int
	minBet    int // minimum bets

	// helpers
	initAmount int  // initial starting amount
	takeOdds   bool // on/off should the algorithm take odds
	won        bool

	// expose
	Occurance map[int]int // count of times a value occured
	Hits      map[int]int // number of times a roll hit us
	Rolls     int         // number of rolls we're on
	Amount    int         // bank
}

func (r *Round) bet(dest *int, a int) {
	r.Amount -= a
	*dest = *dest + a
}

func (r *Round) active() int {
	c := len(r.placeBets) + len(r.comeBets)
	if r.passBet != 0 && r.point != nil {
		c = c + 1
	}
	return c
}

func (r *Round) setPoint(p int) {
	r.point = &p
}
