package bar

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	debug = false
)

var (
	ErrFieldNoReady = errors.New("bar field not ready")
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

func (this *Bar) Get(which BarField) float64 {
	switch which {
	case Open:
		return this.Open
	case Close:
		return this.Close
	case High:
		return this.High
	case Low:
		return this.Low
	default:
		panic(ErrFieldNoReady)
	}
}

const yearForm = "yyyy-mm-dd"

func NewBar(datetime, open, high, low, close_, volumn, size string) *Bar {

	if datetime == "" || open == "" || high == "" || low == "" || close_ == "" {
		return nil
	}

	if len(datetime) != len(yearForm) {
		panic(datetime)
	}

	date_, _ := time.Parse(time.RFC3339, datetime+"T09:00:00+00:00")

	// date_, _ := time.Time
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

func (this *Bar) ToString() string {
	return fmt.Sprintf("%d-%02d-%02d %f %f %f %f %f %f", this.DateTime.Year(), this.DateTime.Month(), this.DateTime.Day(),
		this.Open, this.High, this.Low, this.Low, this.Volumn, this.Size)
}
