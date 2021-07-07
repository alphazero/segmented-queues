// Doost!

package segque

type Stats struct {
	p           *Params
	idx         int
	residencies []int
	rnorms      []float64
	early, late int
}

func NewStats(p *Params, size int) *Stats {
	return &Stats{
		p:           p,
		idx:         0,
		residencies: make([]int, size),
		rnorms:      make([]float64, size),
	}
}

func (s *Stats) Update(r int) {
	var idx = s.idx
	var capacity = s.p.Capacity
	s.residencies[idx] = r
	// normalize by shifting to 0 for perfect residency and dividing by capacity for comparative analysis
	s.rnorms[idx] = float64(r-capacity) / float64(capacity)
	if r < capacity {
		s.early++
	} else if r > capacity {
		s.late++
	}
	s.idx++
}
