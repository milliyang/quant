package indicator

import (
	// "container/list"
	"fmt"
	"quant/base/series"
	"quant/base/xbase"

	"time"
)

/*

    moving average convergence divergence (MACD).

    1. The MACD was designed to profit from this divergence by analyzing
        the difference between the two exponential moving averages
	2. Traders will commonly rely on the default settings of 12- and
        26-day periods.

       DIFF = EMA(close,16) - EMA(close,26)

    3. A signal line, also known as the trigger line, is created by taking
        a nine-period exponential moving average(EMA) of the MACD.

        DEM = EMA(DIFF,9)

    4. Zero Center Line

        OSC = DIFF - DEM

    Advantages:

    Drawbacks:
      1.
      2. dollar value, not percentage price oscillator

    refer:
    1. http://www.investopedia.com/articles/technical/082701.asp
    2. http://zh.wikipedia.org/wiki/MACD
	3. http://zh.wikipedia.org/wiki/%E7%A7%BB%E5%8B%95%E5%B9%B3%E5%9D%87

*/
type MACD struct {
	series.FloatSeries // as OSC
	ShortEMA           *EMA
	LongEMA            *EMA
	DIFF               *series.FloatSeries
	DEM                *EMA

	ShortPeriod int
	LongPeriod  int
	DEMPeriod   int
	Length      int
}

func NewMACD(parent xbase.ISeries, short, long, dem int) *MACD {
	s := &MACD{}

	s.Init(parent)
	s.ShortPeriod = short
	s.LongPeriod = long
	s.DEMPeriod = dem
	s.Name = fmt.Sprintf("MACD (%d,%d,%d)", short, long, dem)

	s.Length = s.LongPeriod

	// Manully update short/long term EMA
	// no parent
	s.ShortEMA = NewEMA(nil, short)
	s.LongEMA = NewEMA(nil, long)
	s.DEM = NewEMA(nil, dem)
	s.DIFF = series.NewFloatSeries()

	if parent != nil {
		parent.AddChild(s)
	}
	return s
}

func (this *MACD) IsFake(datetime *time.Time) bool {
	index := this.Index(datetime)
	if index < this.Length-1 || index < 0 {
		return true
	}
	return false
}

// @override ISeries.Append
func (this *MACD) Append(datetime *time.Time, value float64) {
	if debug {
		fmt.Print("MACD.Append:", value, " and ")
	}
	// Manully update short/long term EMA
	this.ShortEMA.Append(datetime, value)
	this.LongEMA.Append(datetime, value)

	// DIFF = EMA(close,16) - EMA(close,26)
	// DEM = EMA(DIFF,9)
	// OSC = DIFF - DEM
	diff := this.ShortEMA.ValueAtTime(datetime) - this.LongEMA.ValueAtTime(datetime)
	this.DIFF.Append(datetime, diff)
	this.DEM.Append(datetime, diff)
	OSC := diff - this.DEM.ValueAtTime(datetime)
	this.FloatSeries.Append(datetime, OSC)
}

func (this *MACD) LongEmaValues() []float64 {
	return this.LongEMA.Values()
}

func (this *MACD) ShortEmaValues() []float64 {
	return this.ShortEMA.Values()
}

func (this *MACD) DemValues() []float64 {
	return this.DEM.Values()
}

func (this *MACD) DiffValues() []float64 {
	return this.DIFF.Values()
}

func (this *MACD) OscValues() []float64 {
	return this.Values()
}

// indicator.IIndicator.OnMeasure

// num of table, X cordirate [datetime count], Y cordirate [min, max, num, percent]
func (this *MACD) OnMeasure() (int, int, float64, float64, int, float64) {

	// put it here now.
	var min, max float64
	var num int

	for idx, value := range this.Data {
		num = idx
		if value < min {
			min = value
		}
		if value > max {
			value = max
		}
	}
	return 1, len(this.DateTime), min, max, num, 100
}

// indicator.IIndicator.OnDraw(ICanvas)
func (this *MACD) OnDraw(canvas xbase.ICanvas) {
	fmt.Println(this.Name, "symbol", this.Symbol, "onDraw")

	canvas.DrawLine(this.DateTime, this.Data, 1)
	// canvas.DrawBar([]time.Time, []bar.Bar, color)

	// canvas.DrawBuy([]time.Time, []float64,color)
	// canvas.DrawSell([]time.Time, []float64,color)
	// canvas.DrawSpark([]time.Time, []float64,color)
	// canvas.DrawShit([]time.Time, []float64,color)
}
