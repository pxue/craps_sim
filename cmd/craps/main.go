package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/pxue/craps/simulate"
)

type Out struct {
	Rolls  int
	Hits   int
	Amount int
}

func main() {
	s := simulate.SixEightCome{Debug: false}
	// start with 300, play 10 rounds

	perRound := 100
	startAmount := 400
	maxRolls := 3 * 60 // 1 roll per minute, 3h max.
	maxProfit := startAmount * 2

	f, _ := os.Create("out.csv")
	writer := csv.NewWriter(f)
	writer.Write([]string{
		"# of Rounds",
		"# of Hits",
		"# of Rolls",
		"% Hits",
		"$ Bank",
		"$ Profit",
	})

	for i := 0; i < 1000; i++ {
		amount := startAmount
		aRolls := 0
		aRound := 0
		aHits := 0
		for (amount > 100) && aRolls < maxRolls && amount < maxProfit {
			aRound++
			s.Debugf("new shooter coming out! round: %d\n", aRound)

			amount -= perRound
			r := s.Simulate(perRound)

			aRolls += r.Rolls
			amount += r.Amount
			aHits += r.Hits

			s.Debugf("round %d finished with: %d.\n\n", aRound, amount)
			pctHits := float64(aHits*100) / float64(aRolls)
			s.Debugf("\nafter %d rounds, bank $%d/%d, %d/%d (%d%%) rolls was hits.\n", aRound, amount, startAmount, aHits, aRolls, int(pctHits))
		}

		pctHits := float64(aHits) / float64(aRolls)
		writer.Write([]string{
			fmt.Sprintf("%d", aRound),
			fmt.Sprintf("%d", aHits),
			fmt.Sprintf("%d", aRolls),
			fmt.Sprintf("%.2f", pctHits),
			fmt.Sprintf("%d", amount),
			fmt.Sprintf("%d", amount-startAmount),
		})
		writer.Flush()
	}

	// Write any buffered data to the underlying writer (standard output).
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
	f.Close()

}
