// Doost!

package segque

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	//	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	//	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

var (
	TitleFontVariant = font.Variant("Sans")
	TitleFontSize    = font.Length(10)
	TitlePadding     = vg.Inch * font.Length(0.25)
	//	TitleXAlign      = text.XAlignment(-0.5)
	TitleColor = color.RGBA{R: 222, G: 010, B: 111, A: 255}

	LegendFontVariant = font.Variant("Sans")
	LegendFontSize    = font.Length(9)
	LegendYOffset     = (vg.Inch * font.Length(0.5))
	LegendXOffset     = (vg.Inch * font.Length(0.25))

	TickFontVariant = font.Variant("Sans")
)

func newPlot(xmin, xmax, ymin, ymax float64) *plot.Plot {
	p := plot.New()
	p.X.Min = xmin
	p.X.Max = xmax
	p.Y.Min = ymin
	p.Y.Max = ymax

	return p
}

func initPlot(p *plot.Plot) *plot.Plot {

	p.Title.TextStyle.Font.Variant = TitleFontVariant
	p.Title.TextStyle.Font.Size = TitleFontSize
	p.Title.TextStyle.Color = TitleColor
	p.Title.Padding = TitlePadding
	p.Title.TextStyle.XAlign = -0.5

	p.Legend.TextStyle.Font.Variant = LegendFontVariant
	p.Legend.TextStyle.Font.Size = LegendFontSize
	p.Legend.YOffs = vg.Inch * font.Length(0.5)
	p.Legend.XOffs -= (vg.Inch * font.Length(0.25))

	p.X.Tick.Label.Font.Variant = TickFontVariant
	p.Y.Tick.Label.Font.Variant = TickFontVariant
	p.X.Tick.Label.Font.Size = p.X.Tick.Label.Font.Size - 1
	p.Y.Tick.Label.Font.Size = p.Y.Tick.Label.Font.Size - 1

	return p
}

func SavePlot(p *plot.Plot, pfname string, width, height int) {
	w := vg.Inch * font.Length(width)
	h := vg.Inch * font.Length(height)
	fname := fmt.Sprintf("%s.png", pfname)
	if err := p.Save(w, h, fname); err != nil {
		ExitOnError(err, "SavePlot")
	}
}

func PlotHistogram(params *Params, stats *Stats, hbuckets int) *plot.Plot {
	var xmin, xmax, ymin, ymax float64
	p := newPlot(xmin, xmax, ymin, ymax)

	pname := params.CanonicalName()
	p.Title.Text = fmt.Sprintf("%s\nμ:%f\nσ:%f", pname, stats.mean, stats.stddev)
	p.Title.Padding += vg.Inch * font.Length(0.125)

	v := plotter.Values(stats.rnorms)
	h, err := plotter.NewHist(v, hbuckets)
	if err != nil {
		panic(err)
	}
	h.Color = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	h.FillColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Add(h)

	return p
}
