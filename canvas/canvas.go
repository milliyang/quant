package canvas

import (
	"bytes"
	"fmt"
	"github.com/ajstarks/svgo"
	"quant/base/bar"
	"quant/base/xbase"

	"time"
)

var (
	MAGNIFY_FACTOR   int     = 100
	MAGNIFY_FACTOR_F float64 = 100
)

// implement quant/base/indicator.ICanvas
type Canvas struct {
	finished   bool
	indicators []xbase.IIndecator
	svgCanvas  *svg.SVG
	buffer     *bytes.Buffer
}

func NewCanvas() *Canvas {
	self := &Canvas{}
	self.buffer = new(bytes.Buffer)
	self.svgCanvas = svg.New(self.buffer)
	self.finished = false
	return self
}

type Table struct {
	Max     float64
	Min     float64
	NumY    int
	NumX    int
	Percent float64
}

// xbase.ICanvas
func (this *Canvas) Draw(indicator_ xbase.IIndecator) {
	this.indicators = append(this.indicators, indicator_)
}

func (this *Canvas) doRealDrawing() {
	//(int, int, float64, float64, int, float64) // X cordirate [datetime count], Y cordirate [min, max, num, percent]

	for _, ind := range this.indicators {
		// table, timeNum, min, max, num, percent := ind.OnMeasure()
		ind.OnMeasure()
	}

	for _, ind := range this.indicators {
		ind.OnDraw(this)
	}
}

func (this *Canvas) DrawLine(times []time.Time, values []float64, color int) {
	fmt.Println("DrawLine tables:", "times:", len(times), "values:", values, "color:", color)

	w := 1500
	h := 1500

	this.svgCanvas.Start(w, h)
	this.svgCanvas.Rect(0, 0, w, h, "fill:green;stroke:black")

	for i := 0; i < w; i += MAGNIFY_FACTOR {
		this.svgCanvas.Line(i, 0, i, h, "fill:none;stroke:black")
	}

	for i := 0; i < h; i += MAGNIFY_FACTOR {
		this.svgCanvas.Line(0, i, w, i, "fill:none;stroke:black")
	}

	this.svgCanvas.Line(110, 0, 110, h, "fill:none;stroke:red")
	this.svgCanvas.Line(0, 110, w, 110, "fill:none;stroke:red")

	x := []int{}
	y := []int{}

	for idx, value := range values {
		x = append(x, idx*MAGNIFY_FACTOR)
		y = append(y, h-int(value*MAGNIFY_FACTOR_F))
	}
	this.svgCanvas.Polygon(x, y, "fill:none;stroke:black")

	this.svgCanvas.Circle(250, 250, 125, "fill:none;stroke:black")
}

func (this *Canvas) GetResult() string {
	this.finishDrawing()
	return this.buffer.String()
}

func (this *Canvas) finishDrawing() {
	if !this.finished {
		this.doRealDrawing()
		this.svgCanvas.End()
		this.finished = true
	}
}

func (this *Canvas) DrawBar(bars []bar.Bar, color int) {
	if this.finished {
		return
	}
	fmt.Println("DrawBar tables:", "bars:", bars, "color:", color)
}
func (this *Canvas) DrawBuy(times []time.Time, values []float64, color int) {
	if this.finished {
		return
	}
	fmt.Println("DrawBuy tables:", "times:", len(times), "values:", values, "color:", color)
}
func (this *Canvas) DrawSell(times []time.Time, values []float64, color int) {
	if this.finished {
		return
	}
	fmt.Println("DrawSell tables:", "times:", len(times), "values:", values, "color:", color)
}
func (this *Canvas) DrawSpark(times []time.Time, values []float64, color int) {
	if this.finished {
		return
	}
	fmt.Println("DrawSpark tables:", "times:", len(times), "values:", values, "color:", color)
}
func (this *Canvas) DrawShit(times []time.Time, values []float64, color int) {
	if this.finished {
		return
	}
	fmt.Println("DrawShit tables:", "times:", len(times), "values:", values, "color:", color)

}
