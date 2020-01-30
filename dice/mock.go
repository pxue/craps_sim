package dice

type Mock struct {
	sequence chan *Pair
}

func NewMock(values []*Pair) *Mock {
	seq := make(chan *Pair, len(values))
	for _, p := range values {
		seq <- p
	}
	return &Mock{seq}
}

func (m *Mock) Roll() *Pair {
	select {
	case v, ok := <-m.sequence:
		if ok {
			return v
		} else {
			return nil
		}
	default:
		return nil
	}
	return nil
}
