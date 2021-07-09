// Doost!

package main

import (
	"flag"
	"fmt"
	"github.com/alphazero/segque"
	//	"math"
)

var dparams = struct {
	width   int
	height  int
	plotDir string
}{
	width:   5,
	height:  5,
	plotDir: "plot",
}

func init() {
	flag.IntVar(&dparams.width, "width", dparams.width, "plot width")
	flag.IntVar(&dparams.height, "height", dparams.height, "plot height")
	flag.StringVar(&dparams.plotDir, "plotdir", dparams.plotDir, "plot file dir")
}

func main() {
	fmt.Printf("Salaam Sultan of Love!\n")
	p := segque.ParseParams()
	run(p)
}

func run(params *segque.Params) {
	file, r, size := segque.OpenDataFile(params)
	defer file.Close()

	var n int = int(size) / 8
	var stats = segque.NewStats(params, n)
	for i := 0; i < n; i++ {
		residency, eof := segque.ReadInt(r)
		if eof != nil {
			panic("bug - got EOF!")
		}
		stats.Update(residency)
	}
	stats.Compute()

	//	var precision = 2 // REVU TODO export this from stats
	//	var scale = 2.5
	//	if !params.Ctype.UseCSeqnum() {
	//		scale = 4.5
	//	}
	//	maxY := scale / math.Pow(10., float64(precision))
	maxY := 600000.0 // scale / math.Pow(10., float64(precision))
	plot := segque.NewPlot(-1.0, 3.0, 0, maxY)
	distribution := segque.NewDistribution(stats)
	plot.Add(distribution)

	fname := fmt.Sprintf("%s/%s-distribution", dparams.plotDir, params.CanonicalName())
	segque.SavePlot(plot, fname, dparams.width, dparams.height)
}
