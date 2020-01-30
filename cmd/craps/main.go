package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pxue/craps/simulate"
)

var (
	debug = flag.Bool("debug", false, "debug mode")
)

func main() {
	var fname string
	flag.StringVar(&fname, "file", "results.csv", "name of output file")

	var iter int
	flag.IntVar(&iter, "iter", 10000, "number of iterations to run")

	flag.Parse()

	s := simulate.NewSixEight(*debug)
	// start with 300, play 10 rounds

	perRound := 100
	startAmount := 400
	maxRolls := 3 * 60 // 1 roll per minute, 3h max.
	maxProfit := startAmount * 2

	ioWriter := io.WriteCloser(os.Stdout)
	if !s.Debug {
		var err error
		ioWriter, err = os.Create(fname)
		if err != nil {
			log.Fatalf("file create: %v", err)
		}
	}

	writer := csv.NewWriter(ioWriter)
	writer.Write([]string{
		"# of Rounds",
		"# of Hits",
		"# of Rolls",
		"% Hits",
		"$ Bank",
		"$ Profit",
	})

	for i := 0; i < iter; i++ {
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
			aHits += len(r.Hits)

			s.Debugf("round %d finished with: %d.\n\n", aRound, amount)
			s.Debugf("numbers we hit: %+v\n", r.Hits)
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
	ioWriter.Close()
}
