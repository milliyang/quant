package indicator

import (
	// "container/list"
	"fmt"
	"quant/base/series"
	"time"
)

const (
	debug = false
)

type SMA struct {
	series.FloatSeries
	Length       int
	workingValue []float64
	// workingList  *list.List
}

func NewSMA(parent series.ISeries, length int) *SMA {
	s := &SMA{}

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

func (this *SMA) IsFake(datetime *time.Time) bool {
	index := this.Index(datetime)
	if index < this.Length-1 || index < 0 {
		return true
	}
	return false
}

// @override ISeries.Append
func (this *SMA) Append(datetime *time.Time, value float64) {
	if debug {
		fmt.Print("SMA.Append:", value, " and ")
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
