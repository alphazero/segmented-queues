// Doost!

package segque

import (
	"fmt"
	"github.com/gonum/stat"
	"io"
)

type Stats struct {
	p                  *Params
	idx                int
	residencies        []int
	rnorms             []float64
	exact, early, late int
	mean               float64
	geomean            float64
	variance           float64
	stddev             float64
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
	} else {
		s.exact++
	}
	s.idx++
}

func (s *Stats) Compute() {

	s.mean, s.stddev = stat.MeanStdDev(s.rnorms, nil)
	s.geomean = stat.GeometricMean(s.rnorms, nil)
	s.variance = stat.Variance(s.rnorms, nil)

}

func (s *Stats) Print(w io.Writer) {
	var cnt = len(s.residencies)
	fmt.Fprintf(w, "--- stats ------------------------------------------\n")
	fmt.Fprintf(w, "mean:           %f\n", s.mean)
	fmt.Fprintf(w, "geometric-mean: %f\n", s.geomean)
	fmt.Fprintf(w, "variance:       %f\n", s.variance)
	fmt.Fprintf(w, "stddev:         %f\n", s.stddev)
	fmt.Fprintf(w, "exact-evicts:   %f\n", float64(s.exact)/float64(cnt))
	fmt.Fprintf(w, "early-evicts:   %f\n", float64(s.early)/float64(cnt))
	fmt.Fprintf(w, "late-evicts:    %f\n", float64(s.late)/float64(cnt))

}
