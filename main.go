package main

import (
	"fmt"

	"github.com/pxue/craps/simulate"
)

func main() {
	s := simulate.SixEightCome{}
	// start with 300, play 10 rounds
	amount := 0
	for i := 0; i < 10; i++ {
		fmt.Println("new shooter coming out!")
		amount += s.Simulate()
	}

	fmt.Printf("after 10 rounds: net(%d), profit(%d)\n", amount, amount-300)
}
