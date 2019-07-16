package fib

type Sequence struct {
	series []int64
}

func New() *Sequence {
	return &Sequence{[]int64{1, 1}}
}

func (s Sequence) Till(n uint) int64 {
	dim := len(s.series)
	if uint(dim) > n {
		return s.series[n]
	}

	s.series = append(s.series, s.series[dim-2]+s.series[dim-1])
	return s.Till(n)
}
