// Doost!

package segque

import (
	"fmt"
	"github.com/gonum/stat"
	"io"
	"math"
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
	hprecision         int // value of 2 means histogram values are bucketed by 2 precision float values +/-0.nn
	histogram          map[float64]int
	pdist              map[float64]float64
	maxRx              float64 // histogram bucket with max value
	maxRy              int     // corresponding number at maxRx
}

func NewStats(p *Params, size int) *Stats {
	return &Stats{
		p:           p,
		idx:         0,
		residencies: make([]int, size),
		rnorms:      make([]float64, size),
		hprecision:  2,
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

	// compute histogram and probability distribution for normalized values per precision
	s.histogram = make(map[float64]int)
	fprecision := math.Pow(10., float64(s.hprecision))
	n := float64(len(s.rnorms)) // total number of data points
	for _, v0 := range s.rnorms {
		// BUG here is doubling the 0.0 values since +/- values below the precision
		//     end up rounded to 0.0. This basically doubles the expected value
		//     for s.histogram and s.pdist at [0.0] and causes a spike in the dist
		//     plots.
		v := float64(int(v0*fprecision)) / fprecision
		s.histogram[v]++
	}
	s.pdist = make(map[float64]float64)
	for r, cnt := range s.histogram {
		if r == 0.0 {
			// BUG HACK is to just ignore the peak at 0.0
			continue
		}
		s.pdist[r] = float64(cnt) / n
	}

	// determine max of histogram
	for x, y := range s.histogram {
		if s.maxRy < y {
			s.maxRx = x
			s.maxRy = y
		}
	}
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

	// emit histogram / pdistribution
	buckets, _ := ToSortedArrays(s.pdist)
	fmt.Fprintf(w, "hist / pdist:\n")
	for i, bucket := range buckets {
		fmt.Fprintf(w, "[%03d] %+.2f %7d  p: %f\n", i, bucket, s.histogram[bucket], s.pdist[bucket])
	}
	fmt.Fprintf(w, "max @ %+.2f : %d\n", s.maxRx, s.maxRy)
}
