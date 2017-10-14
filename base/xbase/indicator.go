package xbase

import (
	"fmt"
	"time"
)

const (
	/*
		Note:
			The greater IndicatorType will override the smaller one's,
			when in canvas drawing
	*/
	IndicatorTypeDayBar = 1 + iota
	IndicatorTypePerformance
	IndicatorTypeDicingGameBar

	//
)

func init() {
	if debug {
		fmt.Println("quant/base/indicator package init")
	}
}

type IIndecator interface {
	GetIndicatorType() int

	OnMeasure(time.Time, time.Time) (int, float64, float64, int, float64) // X cordirate [datetime count], Y cordirate [min, max, num, percent]
	OnDraw(ICanvas)
	OnLayout() []time.Time
}
