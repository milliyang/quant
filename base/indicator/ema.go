package indicator

import (
	// "container/list"
	"fmt"
	"quant/base/series"
	"quant/base/xbase"
	"time"
)

/*

	exponential moving average (EMA).

	1. EMA = Price(t) * k + EMA(y) * (1 – k)
		t = today, y = yesterday, N = number of days in EMA, k = 2/(N+1)


	refer:
	1. http://zh.wikipedia.org/wiki/%E7%A7%BB%E5%8B%95%E5%B9%B3%E5%9D%87
    2. http://www.iexplain.org/ema-how-to-calculate/

*/

type EMA struct {
	series.FloatSeries
	Length           int     // N
	Alpha            float64 // k
	workingYesterday float64
	// workingList  *list.List
}

func NewEMA(parent xbase.ISeries, length int) *EMA {
	s := &EMA{}

	s.Init(parent, 0)
	s.Length = length
	s.workingYesterday = 0
	var k float64 = (float64)(length + 1)
	s.Alpha = 2.0 / k

	s.Name = "EMA" + fmt.Sprintf("%2d", length)

	// [GoBug Tag00001]??
	//
	if parent != nil {
		parent.AddChild(s)
	}
	return s
}

func (this *EMA) IsFake(datetime *time.Time) bool {
	index := this.Index(datetime)
	if index < this.Length-1 || index < 0 {
		return true
	}
	return false
}

// @override ISeries.Append
func (this *EMA) Append(datetime *time.Time, value float64) {
	if debug {
		fmt.Print("EMA.Append:", value, " and ")
	}

	// How about EMA of the first day
	// EMA(0) =? VALUE
	if this.workingYesterday == 0 {
		this.workingYesterday = value
	}

	ema := value*this.Alpha + this.workingYesterday*(1-this.Alpha)
	this.workingYesterday = ema
	this.FloatSeries.Append(datetime, ema)
	return
}
