package simulate

// a round is bets before a 7 rolls again
type round struct {
	point *int // the current point
	n     int  // number of rolls we're on

	passBet   int         // the current pass bet
	oddBet    int         // the current odd bets
	comeBet   int         // the current come bet
	placeBets map[int]int // place bets
	comeBets  map[int]int // active come bets

	// helpers
	takeOdds bool
	won      bool

	amount int
}
