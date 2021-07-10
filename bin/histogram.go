// Doost!

package main

import (
	"flag"
	"fmt"
	"github.com/alphazero/segque"
)

var hparams = struct {
	width    int
	height   int
	hbuckets int
	plotDir  string
}{
	width:    5,
	height:   5,
	hbuckets: 100,
	plotDir:  "plot",
}

func init() {
	flag.IntVar(&hparams.width, "width", hparams.width, "plot width")
	flag.IntVar(&hparams.height, "height", hparams.height, "plot height")
	flag.IntVar(&hparams.hbuckets, "buckets", hparams.hbuckets, "histogram buckets")
	flag.StringVar(&hparams.plotDir, "plotdir", hparams.plotDir, "plot file dir")
}

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

	plot := segque.PlotHistogramXY(p, stats, hparams.hbuckets, -1.000, 3.0, 0.0, 600000.0)

	fname := fmt.Sprintf("%s/%s", hparams.plotDir, p.CanonicalName())
	segque.SavePlot(plot, fname, hparams.width, hparams.height)
}
