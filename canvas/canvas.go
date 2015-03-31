package canvas

import (
	"bytes"
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/milliyang/dice"
	"github.com/milliyang/dice/utils"
	"quant/base/bar"
	"quant/base/color"
	"quant/base/xbase"
	"strconv"
	"time"
)

var (
	MAGNIFY_FACTOR   int     = 100
	MAGNIFY_FACTOR_F float64 = 100
	TIME_ORIGIN      time.Time

	DAY100                = 100
	WIDTH_FOR_100_DAY_BAR = 1200
	HEIGTH_DEFAULT        = 500
)

func init() {
	TIME_ORIGIN = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local)

	if xbase.CasinoDicingGame {
		WIDTH_FOR_100_DAY_BAR = 2000
		HEIGTH_DEFAULT = 2000

		WIDTH_FOR_100_DAY_BAR = WIDTH_FOR_100_DAY_BAR * 4
	}
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

	self.rectWidth = WIDTH_FOR_100_DAY_BAR
	self.rectHeight = HEIGTH_DEFAULT

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

	// Find which indicator has the most fullset of DateTime (X cordirate)
	for idx, ind := range this.indicators {

		// TODO
		fmt.Println("[Warn]: ind.OnMeasure(TIME_ORIGIN, time.Now())", TIME_ORIGIN, time.Now())
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

	// Calc Rect Width
	this.rectWidth = (len(this.xBaseTime)/DAY100 + 1) * WIDTH_FOR_100_DAY_BAR

	this.drawItemX = len(this.xBaseTime) + 1
	this.drawStepX = this.rectWidth / (this.drawItemX)

	for idx, oneTime := range this.xBaseTime {
		num := oneTime.Year()*1000 + oneTime.YearDay()
		this.mapTimeIdx[num] = idx + 1
	}

	fmt.Println("yMax:", this.yMax, " yMin:", this.yMin)

	// Hacking for Casino Dicing Game
	if xbase.CasinoDicingGame {
		this.yMin = 1
	}

	priceRange := (this.yMax*1.1 - this.yMin*0.9) * 100.0

	// Calc Rect Height
	this.drawStepY = 0
	for {
		if this.drawStepY != 0 {
			break
		} else {
			this.rectHeight += HEIGTH_DEFAULT
		}
		this.drawStepY = int(float64(this.rectHeight) / priceRange)
	}
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

	for i := 0; i < this.drawItemX; i++ {
		this.svgCanvas.Line(i*this.drawStepX, 0, i*this.drawStepX, this.rectHeight, "fill:none;stroke:white")
	}

	// CasinoDicingGame Line for [1,2,3,4,5,6]
	if xbase.CasinoDicingGame {
		for i := 1.0; i < 19; i = i + 1.0 {

			// yOffset := this.rectHeight - int((i-this.yMin*0.9)*100*float64(this.drawStepY))
			yOffset := this.calcYOffsetByPrice(i)
			if i <= 10 {
				this.svgCanvas.Line(0, yOffset, this.rectWidth, yOffset, "fill:none;stroke:green")
			} else {
				this.svgCanvas.Line(0, yOffset, this.rectWidth, yOffset, "fill:none;stroke:red")
			}
		}
	} else {
		for i := 0; i < this.drawItemY; i += 10 {
			this.svgCanvas.Line(0, i*this.drawStepY, this.rectWidth, i*this.drawStepY, "fill:none;stroke:white")
		}

		// Cross Line for Mouse
		this.svgCanvas.Line(110, 0, 110, this.rectHeight, "fill:none;stroke:red")
		this.svgCanvas.Line(0, 110, this.rectWidth, 110, "fill:none;stroke:red")
	}

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

func (this *Canvas) DrawLine(times []time.Time, values []float64, color_ int) {
	if this.finished {
		return
	}

	xOffsets := []int{}

	for _, oneTime := range times {
		xIdx := this.calcXIdxByDatetime(&oneTime)
		xOffsets = append(xOffsets, this.calcXOffsetByIdx(xIdx))
	}

	yOffsets := []int{}
	for _, onePrice := range values {
		yOffset := this.calcYOffsetByPrice(onePrice)
		yOffsets = append(yOffsets, yOffset)
	}

	this.svgCanvas.Polyline(xOffsets, yOffsets, color.GetCssColor(color_))
}

/*
	for xbase.CasinoDicingGame
*/
func (this *Canvas) DrawPoints(times []time.Time, values []float64, color_ int, withContinue bool) {
	if this.finished {
		return
	}

	if len(times) != len(values) {
		panic("invalid (x,y) pairs")
	}

	var xIdx, xOffset, yOffset int

	var lastX, lastY int
	var lastValue float64

	for i := 0; i < len(times); i++ {
		oneTime := times[i]
		oneValue := values[i]

		xIdx = this.calcXIdxByDatetime(&oneTime)
		xOffset = this.calcXOffsetByIdx(xIdx)
		yOffset = this.calcYOffsetByPrice(oneValue)

		this.svgCanvas.Circle(xOffset, yOffset, 5, "fill:pink;stroke:yellow")

		if withContinue && lastX > 0 {
			if isSebaoContinue(lastValue, oneValue) {
				this.svgCanvas.Line(lastX, lastY, xOffset, yOffset, "fill:pink;stroke:red;stroke-width:5")
			}
		}

		lastValue = oneValue
		lastX = xOffset
		lastY = yOffset
	}
}

// continue means:
//  [4 ,10]
//  [11,18]
func isSebaoContinue(value1, value2 float64) bool {
	if value1 <= 10 && value2 <= 10 || value1 >= 11 && value2 >= 11 {
		return true
	} else {
		return false
	}
}

func (this *Canvas) DrawBar(bars []bar.Bar, color int) {
	if this.finished {
		return
	}

	var acc []*dice.DiceRoll

	for _, oneBar := range bars {
		if oneBar.Dice != nil {
			this.drawOneDiceBar(oneBar)
			acc = append(acc, oneBar.Dice)
		} else {
			this.drawOneBar(oneBar)
		}
	}

	if xbase.CasinoDicingGame && acc != nil {
		utils.CheckCasinoPoint(acc)
		utils.CheckRandom(acc)
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

func (this *Canvas) DrawTextAtPrice(times []time.Time, values []string, price float64, color int) {
	if this.finished {
		return
	}
	fmt.Println("DrawTextAtPrice:", "times:", len(times), "values:", values, "color:", color)
}

func (this *Canvas) drawOneBar(bar bar.Bar) {

	// Android Style: (0,0) is at (Left,Top)
	var barHeight int

	xIdx := this.calcXIdxByDatetime(&bar.DateTime)
	xOffset := this.calcXOffsetByIdx(xIdx)

	yOffset := 0
	barColor := ""

	yOpen := this.calcYOffsetByPrice(bar.Open)
	yClose := this.calcYOffsetByPrice(bar.Close)

	if yOpen > yClose {
		barHeight = yOpen - yClose
		yOffset = yClose
		barColor = "fill:red;stroke:red"
	} else {
		barHeight = yClose - yOpen
		yOffset = yOpen
		barColor = "fill:green;stroke:green"
	}
	if barHeight == 0 {
		barHeight = 1
	}

	this.svgCanvas.Rect(xOffset-this.barWidthHalf, yOffset, this.barWidth, barHeight, barColor)

	lowPriceOffset := this.calcYOffsetByPrice(bar.Low)
	highPriceOffset := this.calcYOffsetByPrice(bar.High)
	this.svgCanvas.Line(xOffset, lowPriceOffset, xOffset, highPriceOffset, barColor)

	if false {
		oo := this.calcYOffsetByPrice(bar.Open)
		this.svgCanvas.Circle(xOffset, oo, 5, "fill:pink;stroke:yellow")

		cc := this.calcYOffsetByPrice(bar.Close)
		this.svgCanvas.Circle(xOffset, cc, 5, "fill:yellow;stroke:yellow")
	}
}

func (this *Canvas) drawOneDiceBar(bar bar.Bar) {

	// Android Style: (0,0) is at (Left,Top)

	// A <= B <= C
	var A, B, C float64
	A = float64(bar.Dice.Rolls[0])
	B = float64(bar.Dice.Rolls[1])
	C = float64(bar.Dice.Rolls[2])

	var barHeight int
	var lowerPrice float64

	barColor := "fill:red;stroke:red"

	barHeight = int((C - A) * 100 * float64(this.drawStepY))
	if barHeight == 0 {
		barHeight = 1
	}
	lowerPrice = A
	yOffset := this.calcYOffsetByPrice(lowerPrice)

	xIdx := this.calcXIdxByDatetime(&bar.DateTime)
	xOffset := this.calcXOffsetByIdx(xIdx)
	this.svgCanvas.Rect(xOffset-this.barWidthHalf, yOffset-barHeight, this.barWidth, barHeight, barColor)

	// draw circle
	AA := this.calcYOffsetByPrice(A)
	this.svgCanvas.Circle(xOffset, AA, 5, "fill:red;stroke:yellow")

	BB := this.calcYOffsetByPrice(B)
	this.svgCanvas.Circle(xOffset, BB, 5, "fill:green;stroke:yellow")

	CC := this.calcYOffsetByPrice(C)
	this.svgCanvas.Circle(xOffset, CC, 5, "fill:blue;stroke:yellow")

	// (1,2,3)

	comment := fmt.Sprintf("(%d,%d,%d)", bar.Dice.Rolls[0], bar.Dice.Rolls[1], bar.Dice.Rolls[2])

	this.drawCasinoBigSmallText(xOffset, int64(A+B+C), comment)
}

func (this *Canvas) calcYOffsetByPrice(price float64) int {
	// Hacking
	// To Make more space between 10 - 11
	// if xbase.CasinoDicingGame {
	// 	if price > 10.5 {
	// 		price = price + 1
	// 	}
	// }
	yOffset := this.rectHeight - int((price-this.yMin*0.9)*100*float64(this.drawStepY))
	return yOffset
}

func (this *Canvas) calcXOffsetByIdx(idx int) int {
	return idx * this.drawStepX
}

func (this *Canvas) calcXIdxByDatetime(time_ *time.Time) int {
	num := time_.Year()*1000 + time_.YearDay()
	xIdx, ok := this.mapTimeIdx[num]
	if ok {
		return xIdx
	} else {
		panic(time_)
	}
}

func (this *Canvas) drawCasinoBigSmallText(xOffset int, total int64, comment string) {
	_color := "fill:red;stroke:red"
	yOffset := 0

	if total >= 11 {
		yOffset = this.calcYOffsetByPrice(10.8)
	} else {
		_color = "fill:green;stroke:green"
		yOffset = this.calcYOffsetByPrice(10.1)
	}
	this.svgCanvas.Text(xOffset, yOffset, strconv.FormatInt(total, 10)+comment, _color)
}
