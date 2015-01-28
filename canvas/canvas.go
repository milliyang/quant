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
	TIME_ORIGIN      time.Time
)

func init() {
	TIME_ORIGIN = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local)
}

// implement quant/base/indicator.ICanvas
type Canvas struct {
	finished   bool
	indicators []xbase.IIndecator
	svgCanvas  *svg.SVG
	buffer     *bytes.Buffer

	xNum      int
	yMin      float64
	yMax      float64
	yNum      int
	yPercent  float64
	xBaseIdx  int
	xBaseTime []time.Time

	drawStepX    int
	drawStepY    int
	drawItemX    int
	drawItemY    int
	rectWidth    int
	rectHeight   int
	barWidth     int
	barWidthHalf int

	mapTimeIdx map[int]int
}

func NewCanvas() *Canvas {
	self := &Canvas{}
	self.buffer = new(bytes.Buffer)
	self.svgCanvas = svg.New(self.buffer)
	self.finished = false

	self.rectWidth = 1200
	self.rectHeight = 500

	self.xNum = 0
	self.yNum = 0
	self.yMin = 10000
	self.yMax = -10000

	self.mapTimeIdx = make(map[int]int)
	return self
}

// xbase.ICanvas
func (this *Canvas) Draw(indicator_ xbase.IIndecator) {
	this.indicators = append(this.indicators, indicator_)
}

func (this *Canvas) prepareDrawing() {
	//(int, float64, float64, int, float64) // X cordirate [datetime count], Y cordirate [min, max, num, percent]
	for idx, ind := range this.indicators {
		xnum, ymin, ymax, ynum, ypercent := ind.OnMeasure(TIME_ORIGIN, time.Now())

		if xnum > this.xNum {
			this.xNum = xnum
			this.xBaseIdx = idx
		}
		if ynum > this.yNum {
			this.yNum = ynum
		}
		if ymin < this.yMin {
			this.yMin = ymin
		}
		if ymax > this.yMax {
			this.yMax = ymax
		}
		this.yPercent = ypercent // useless currently
	}

	this.xBaseTime = this.indicators[this.xBaseIdx].OnLayout()
	this.drawItemX = len(this.xBaseTime) + 1

	for idx, oneTime := range this.xBaseTime {
		num := oneTime.Year()*1000 + oneTime.YearDay()
		this.mapTimeIdx[num] = idx + 1
	}

	this.drawStepX = this.rectWidth / (this.drawItemX)

	priceRange := (this.yMax*1.1 - this.yMin*0.9) * 100.0
	this.drawStepY = int(float64(this.rectHeight) / priceRange)
	this.drawItemY = this.rectHeight / this.drawStepY

	this.barWidth = int(float64(this.drawStepX) * 0.66)
	this.barWidthHalf = int(float64(this.drawStepX) * 0.33)

	// fmt.Println("yMax1.1:", this.yMax*1.1, " yMin0.9:", this.yMin*0.9, " priceRange", priceRange)
	// fmt.Println("xStep:", this.drawStepX, " yStep:", this.drawStepY)
}

func (this *Canvas) GetResult() string {
	this.finishDrawing()
	return this.buffer.String()
}

func (this *Canvas) startDrawing() {
	// draw basic framework

	// left,bar,right,buy,self.

	this.svgCanvas.Start(this.rectWidth, this.rectHeight)
	this.svgCanvas.Rect(0, 0, this.rectWidth, this.rectHeight, "fill:gray;stroke:black")

	return

	for i := 0; i < this.drawItemX; i++ {
		this.svgCanvas.Line(i*this.drawStepX, 0, i*this.drawStepX, this.rectHeight, "fill:none;stroke:white")
	}

	for i := 0; i < this.drawItemY; i += 5 {
		this.svgCanvas.Line(0, i*this.drawStepY, this.rectWidth, i*this.drawStepY, "fill:none;stroke:white")
	}

	// Cross Line for Mouse
	this.svgCanvas.Line(110, 0, 110, this.rectHeight, "fill:none;stroke:red")
	this.svgCanvas.Line(0, 110, this.rectWidth, 110, "fill:none;stroke:red")

	/* draw point
	x := []int{}
	y := []int{}
	for idx, value := range values {
		x = append(x, idx*MAGNIFY_FACTOR)
		y = append(y, h-int(value*MAGNIFY_FACTOR_F))
	}
	this.svgCanvas.Polygon(x, y, "fill:none;stroke:black")
	*/

	//this.svgCanvas.Circle(250, 250, 125, "fill:none;stroke:black")
}

func (this *Canvas) doDrawing() {
	for _, ind := range this.indicators {
		ind.OnDraw(this)
	}
}

func (this *Canvas) finishDrawing() {
	if !this.finished {
		this.prepareDrawing()
		this.startDrawing()
		this.doDrawing()
		this.svgCanvas.End()
		this.finished = true
	}
}

func (this *Canvas) DrawLine(times []time.Time, values []float64, color int) {
	fmt.Println("DrawLine tables:", "times:", len(times), "values:", values, "color:", color)
}

func (this *Canvas) DrawBar(bars []bar.Bar, color int) {
	if this.finished {
		return
	}

	for _, oneBar := range bars {
		this.drawOneBar(oneBar)
	}
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

func (this *Canvas) drawOneBar(bar bar.Bar) {

	// Android Style: (0,0) is at (Left,Top)
	var barHeight int
	var higher float64

	barColor := "fill:red;stroke:red"

	if bar.Close > bar.Open {
		barHeight = int((bar.Close - bar.Open) * 100 * float64(this.drawStepY))
		higher = bar.Open
	} else {
		barHeight = int((bar.Open - bar.Close) * 100 * float64(this.drawStepY))
		higher = bar.Close
		barColor = "fill:green;stroke:green"
	}

	yOffset := int((higher - this.yMin*0.9) * 100 * float64(this.drawStepY))

	num := bar.DateTime.Year()*1000 + bar.DateTime.YearDay()
	xIdx, ok := this.mapTimeIdx[num]
	if ok {
		if barHeight == 0 {
			barHeight = 1
		}
		this.svgCanvas.Rect(xIdx*this.drawStepX-this.barWidthHalf, this.rectHeight-yOffset-barHeight, this.barWidth, barHeight, barColor)
	} else {
		panic(bar)
	}

	lowPriceOffset := this.rectHeight - int((bar.High-this.yMin*0.9)*100*float64(this.drawStepY))
	highPriceOffset := this.rectHeight - int((bar.Low-this.yMin*0.9)*100*float64(this.drawStepY))

	this.svgCanvas.Line(xIdx*this.drawStepX, lowPriceOffset, xIdx*this.drawStepX, highPriceOffset, barColor)

	if false {
		oo := this.rectHeight - int((bar.Open-this.yMin*0.9)*100*float64(this.drawStepY))
		this.svgCanvas.Circle(xIdx*this.drawStepX, oo, 5, "fill:pink;stroke:yellow")
		cc := this.rectHeight - int((bar.Close-this.yMin*0.9)*100*float64(this.drawStepY))
		this.svgCanvas.Circle(xIdx*this.drawStepX, cc, 5, "fill:yellow;stroke:yellow")
	}

}
