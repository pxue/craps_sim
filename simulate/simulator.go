package simulate

type Simulator interface {
	Simulate(r *Round) *Round
	Debug(format string, args ...interface{})
}

func Simulate(t Simulator) {
}
