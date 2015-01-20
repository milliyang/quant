package indicator

import (
	// "container/list"
	"fmt"
	"quant/base/series"
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
	DEM                *EMA

	ShortPeriod int
	LongPeriod  int
	DEMPeriod   int
	Length      int
}

func NewMACD(parent series.ISeries, short, long, dem int) *MACD {
	s := &MACD{}

	s.Init(parent)
	s.ShortPeriod = short
	s.LongPeriod = long
	s.DEMPeriod = dem

	s.Length = s.LongPeriod

	// Manully update short/long term EMA
	// no parent
	s.ShortEMA = NewEMA(nil, short)
	s.LongEMA = NewEMA(nil, long)
	s.DEM = NewEMA(nil, dem)

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
	DIFF := this.ShortEMA.ValueAtTime(datetime) - this.LongEMA.ValueAtTime(datetime)
	this.DEM.Append(datetime, DIFF)
	OSC := DIFF - this.DEM.ValueAtTime(datetime)
	this.FloatSeries.Append(datetime, OSC)
}

func (this *MACD) LongEmaValues() []float64 {
	return this.LongEMA.Values()
}

func (this *MACD) ShortEmaValues() []float64 {
	return this.ShortEMA.Values()
}

// func (this *MACD) DemValues() []float64 {
// 	return this.DEM.Values()
// }

func (this *MACD) OscValues() []float64 {
	return this.Values()
}
