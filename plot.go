// Doost!

package segque

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

/// general plot //////////////////////////////////////////////////////////////

var colidx = 0
var colors = []color.RGBA{
	color.RGBA{R: 150, G: 150, B: 150, A: 255},
	color.RGBA{R: 255, G: 150, B: 150, A: 255},
	color.RGBA{R: 150, G: 255, B: 150, A: 255},
	color.RGBA{R: 150, G: 150, B: 255, A: 255},
	color.RGBA{R: 000, G: 000, B: 000, A: 255},
}

var (
	// DefaultLineStyle is the default style for drawing
	// lines.
	DefaultLineStyle = draw.LineStyle{
		Color:    color.Black,
		Width:    vg.Points(1),
		Dashes:   []vg.Length{},
		DashOffs: 0,
	}

	// DefaultGlyphStyle is the default style used
	// for gyph marks.
	DefaultGlyphStyle = draw.GlyphStyle{
		Color:  color.Black,
		Radius: vg.Points(2.5),
		Shape:  draw.RingGlyph{},
	}
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

func NewPlot(xmin, xmax, ymin, ymax float64) *plot.Plot {
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

/// distribution //////////////////////////////////////////////////////////////

type Distribution struct {
	cnt        int
	xarr, yarr []float64
	draw.LineStyle
}

// NewDistribution returns a Distribution that plots F using
// the default line style with 50 samples.
func NewDistribution(stats *Stats) *Distribution {
	xarr, yarr := ToSortedArrays(stats.pdist)
	ls := DefaultLineStyle
	ls.Color = colors[colidx]
	colidx++
	return &Distribution{
		xarr:      xarr,
		yarr:      yarr,
		cnt:       len(xarr),
		LineStyle: ls,
		//		LineStyle: DefaultLineStyle,
	}
}

// Distribution.Plot implements the Plotter interface,
// drawing a line that connects each point in the Line.
func (d *Distribution) Plot(canvas draw.Canvas, plot *plot.Plot) {
	trX, trY := plot.Transforms(&canvas)
	line := make([]vg.Point, d.cnt)

	for i := range line {
		x := d.xarr[i]
		y := d.yarr[i]
		line[i].X = trX(x)
		line[i].Y = trY(y)
	}
	canvas.StrokeLines(d.LineStyle, canvas.ClipLinesXY(line)...)
}

// Thumbnail draws a line in the given style down the
// center of a DrawArea as a thumbnail representation
// of the LineStyle of the function.
func (s Distribution) Thumbnail(c *draw.Canvas) {
	y := c.Center().Y
	c.StrokeLine2(s.LineStyle, c.Min.X, y, c.Max.X, y)
}

/// histogram /////////////////////////////////////////////////////////////////

// TODO comparative cache-lines
func PlotHistogramXY(params *Params, stats *Stats, hbuckets int, xmin, xmax, ymin, ymax float64) *plot.Plot {
	//	var xmin, xmax, ymin, ymax float64
	p := NewPlot(xmin, xmax, ymin, ymax)

	pname := params.CanonicalName()
	p.Title.Text = fmt.Sprintf("%s\nμ:%f\nσ:%f\nvar:%f", pname, stats.mean, stats.stddev, stats.variance)
	p.Title.Padding += vg.Inch * font.Length(0.125)

	v := plotter.Values(stats.rnorms)
	h, err := plotter.NewHist(v, hbuckets)
	if err != nil {
		panic(err)
	}
	h.Color = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	h.FillColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	//	hxmin, hxmax, hymin, hymax := h.DataRange()
	p.Add(h)

	stdlines := stdDevLines(stats, xmin, xmax, ymin, ymax)
	plotutil.AddLines(p,
		"μ", stdlines[2],
		"μ-σ", stdlines[0],
		"μ-2σ", stdlines[1],
		//		"μ-3σ", stdlines[2], // DON"T use this
		"C LRU", lruResLine(params, C_LRU_ENTRY_SIZE, ymax),
		"Go LRU", lruResLine(params, GO_LRU_ENTRY_SIZE, ymax),
	)
	return p
}

// Adds lines for std-devi (1 and 2) in the early evict side of histogram
// plus the mean itself.
func stdDevLines(stats *Stats, xmin, xmax, ymin, ymax float64) []plotter.XYs {

	stdlines := make([]plotter.XYs, 3)

	sd0line := make(plotter.XYs, 2)
	sd0line[0].X = stats.mean - stats.stddev
	sd0line[0].Y = 0.0
	sd0line[1].X = sd0line[0].X
	sd0line[1].Y = ymax

	sd1line := make(plotter.XYs, 2)
	sd1line[0].X = stats.mean - (2.0 * stats.stddev)
	sd1line[0].Y = 0.0
	sd1line[1].X = sd1line[0].X
	sd1line[1].Y = ymax

	meanline := make(plotter.XYs, 2)
	meanline[0].X = stats.mean
	meanline[0].Y = 0.0
	meanline[1].X = stats.mean
	meanline[1].Y = ymax

	stdlines[0] = sd0line
	stdlines[1] = sd1line
	stdlines[2] = meanline
	return stdlines
}

// An LRU of equivalent memory footprint will have smaller capacity than a CLC LRU.
// This function computes the normalized residency of a given LRU with a given refsize
// (in bytes) per entry. A CLC is a cache line, so it is 64 Bytes. Depending on the number
// of slots per CLC (typically 7 for 64bit keys, and 13 for 32bit keys), and given the
// reference LRU size, we compute the (constant) normalized residency value for that type.
// For C 64bit LRUs, we use 40 Bytes. For Go, 56 Bytes, per LRU entry.
func lruResLine(params *Params, refsize int, ymax float64) plotter.XYs {
	var clcSize = float64(params.CLSize) / float64(params.Slots) // bytes per entry
	var rfactor = clcSize / float64(refsize)
	var rnormal = rfactor - 1.0

	line := make(plotter.XYs, 2)
	line[0].X = rnormal
	line[0].Y = 0.0
	line[1].X = rnormal
	line[1].Y = ymax

	return line
}
