package simulate

import (
	"testing"

	"github.com/pxue/craps/dice"
)

type sixeightTest struct {
	name         string
	start        int // starting amount
	initialRound *Round
	rolls        []*dice.Pair
	expected     int // expected profit
}

func TestSixEightSequence(t *testing.T) {
	t.Parallel()

	tests := []sixeightTest{
		{
			name:     "sequence 1",
			start:    100,
			rolls:    []*dice.Pair{{2, 4}, {1, 3}, {1, 2}, {3, 3}, {1, 1}, {4, 2}, {3, 5}, {5, 6}, {3, 4}},
			expected: 109,
		},
		{
			name:     "sequence 2",
			start:    100,
			rolls:    []*dice.Pair{{1, 2}, {2, 4}, {3, 3}, {5, 6}, {2, 2}, {2, 3}, {4, 6}, {4, 4}, {3, 3}, {2, 6}, {3, 6}, {4, 5}, {1, 3}, {3, 6}, {5, 4}, {3, 4}},
			expected: 172,
		},
		{
			name:  "sequence 3",
			start: 172,
			initialRound: &Round{
				placeBets: map[int]int{
					6: 18,
					8: 18,
				},
				passBet:    15,
				minBet:     15,
				initAmount: 172,
				Amount:     172,
				comeBets:   map[int]int{},
				Occurance:  map[int]int{},
				Hits:       map[int]int{},
			},
			rolls:    []*dice.Pair{{3, 4}, {2, 2}, {5, 6}, {2, 4}, {6, 5}, {5, 5}, {4, 4}, {3, 3}, {4, 5}, {2, 5}},
			expected: 259,
		},
		{
			name:     "sequence 4",
			start:    100,
			rolls:    []*dice.Pair{{5, 6}, {2, 2}, {5, 5}, {3, 3}, {2, 2}, {5, 5}, {4, 5}, {5, 4}, {4, 4}, {4, 4}, {1, 2}, {2, 5}},
			expected: 142,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := dice.NewMock(tt.rolls)
			s := NewSixEight(true, src)
			out := s.Simulate(tt.start, tt.initialRound)
			if a := out.Amount; a != tt.expected {
				t.Errorf("test '%s' amount not equal. e: %v got %v", tt.name, tt.expected, a)
			}
		})

	}
}
