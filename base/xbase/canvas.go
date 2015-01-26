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

	DrawBuy([]time.Time, []float64, int)
	DrawSell([]time.Time, []float64, int)

	DrawSpark([]time.Time, []float64, int)
	DrawShit([]time.Time, []float64, int)

	DrawBar([]bar.Bar, int)

	// canvas.DrawLine(table, []time.Time, []float64, color)
	// canvas.DrawBar(table,  []time.Time, []bar.Bar, color)
	// canvas.DrawBuy(table,  []time.Time, []float64,color)
	// canvas.DrawSell(table, []time.Time, []float64,color)
	// canvas.DrawSpark(table,[]time.Time, []float64,color)
	// canvas.DrawShit(table, []time.Time, []float64,color)
}
