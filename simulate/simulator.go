package simulate

type Simulator interface {
	Simulate(r *round)
}

func Simulate(t Simulator) {
}
