// Doost!

package main

import (
	"fmt"
	"github.com/alphazero/segque"
)

func main() {
	fmt.Printf("Salaam Sultan of Love!\n")
	p := segque.ParseParams()
	run(p)
}

func run(p *segque.Params) {
	file, r, size := segque.OpenDataFile(p)
	defer file.Close()

	var n int = int(size) / 8
	var stats = segque.NewStats(p, n)
	for i := 0; i < n; i++ {
		residency, eof := segque.ReadInt(r)
		if eof != nil {
			panic("bug - got EOF!")
		}
		stats.Update(residency)
	}
	stats.Compute()
	//	stats.Print(os.Stdout)

	hbuckets := 100
	plot := segque.PlotHistogram(p, stats, hbuckets)

	plotDir := "plots"
	width := 5
	height := 5
	fname := fmt.Sprintf("%s/%s", plotDir, p.CanonicalName())
	segque.SavePlot(plot, fname, width, height)

}
