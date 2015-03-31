package xbase

import (
	"fmt"
	"quant/base/bar"
	"time"
)

func init() {
	if debug {
		fmt.Println("quant/base/indicator package init")
	}
}

type ICanvas interface {
	DrawLine([]time.Time, []float64, int)
	DrawPoints([]time.Time, []float64, int, bool)

	DrawBuy([]time.Time, []float64, int)
	DrawSell([]time.Time, []float64, int)

	DrawSpark([]time.Time, []float64, int)
	DrawShit([]time.Time, []float64, int)

	DrawBar([]bar.Bar, int)

	// CasinoDicingGame
	DrawTextAtPrice([]time.Time, []string, float64, int)
}
