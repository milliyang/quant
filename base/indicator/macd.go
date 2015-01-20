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

*/
type MACD struct {
	series.FloatSeries
	Length       int
	workingValue []float64
	// workingList  *list.List
}

func NewMACD(parent series.ISeries, length int) *MACD {
	s := &MACD{}

	s.Init(parent)
	s.workingValue = []float64{}
	s.Length = length

	// [GoBug Tag00001]??
	//
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

	if len(this.workingValue) < this.Length {
		this.workingValue = append(this.workingValue, value)
	} else {
		this.workingValue = append(this.workingValue[1:], value)
	}

	// average
	var total float64
	for _, item := range this.workingValue {
		total += item
	}

	var num float64
	num = float64(len(this.workingValue))
	this.FloatSeries.Append(datetime, total/num)
	return
}
