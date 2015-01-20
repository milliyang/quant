package indicator

import (
	// "container/list"
	"fmt"
	"quant/base/series"
	"time"
)

/*

	exponential moving average (EMA).

	1. EMA = Price(t) * k + EMA(y) * (1 – k)
		t = today, y = yesterday, N = number of days in EMA, k = 2/(N+1)

*/

type EMA struct {
	series.FloatSeries
	Length           int     // N
	Alpha            float64 // k
	workingYesterday float64
	// workingList  *list.List
}

func NewEMA(parent series.ISeries, length int) *EMA {
	s := &EMA{}

	s.Init(parent)
	s.Length = length
	s.workingYesterday = 0
	var k float64 = (float64)(length + 1)
	s.Alpha = 2.0 / k

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
