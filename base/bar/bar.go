package bar

import (
	"fmt"
	"strconv"
	"time"
)

const (
	debug = true
)

type Bar struct {
	DateTime time.Time
	Close    float64
	Open     float64
	High     float64
	Low      float64
	Volumn   float64
	Size     float64
	Type     BarType
}

func (this *Bar) get(which BarField) float64 {
	switch which {
	case Open:
		return this.Open
	case Close:
		return this.Close
	case High:
		return this.High
	case Low:
		return this.Low
	}
	return -1
}

const yearForm = "2006-01-22"

func NewBar(datetime, open, high, low, close_, volumn, size string) *Bar {

	if datetime == "" || open == "" || high == "" || low == "" || close_ == "" {
		return nil
	}

	date_, _ := time.Parse(yearForm, datetime)
	open_, _ := strconv.ParseFloat(open, 32)
	high_, _ := strconv.ParseFloat(high, 32)
	low_, _ := strconv.ParseFloat(low, 32)
	close_x, _ := strconv.ParseFloat(close_, 32)

	newBar := &Bar{}
	newBar.Type = Time

	newBar.DateTime = date_
	newBar.Open = open_
	newBar.High = high_
	newBar.Low = low_
	newBar.Close = close_x

	if volumn != "" {
		v_, _ := strconv.ParseFloat(volumn, 32)
		newBar.Volumn = v_
	}
	if size != "" {
		s_, _ := strconv.ParseFloat(size, 32)
		newBar.Size = s_
	}
	return newBar
}

func (this *Bar) toString() string {
	return fmt.Sprintf("%v %f %f %f %f %f %f", this.DateTime, this.Open, this.High, this.Low, this.Low, this.Volumn, this.Size)

}
