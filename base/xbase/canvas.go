package xbase

import (
	"fmt"
	"time"
)

func init() {
	if debug {
		fmt.Println("quant/base/indicator package init")
	}
}

type ICanvas interface {
	DrawLine(int, []time.Time, []float64, int)
	DrawBar(int, []time.Time, []float64, int)

	DrawBuy(int, []time.Time, []float64, int)
	DrawSell(int, []time.Time, []float64, int)

	DrawSpark(int, []time.Time, []float64, int)
	DrawShit(int, []time.Time, []float64, int)

	// canvas.DrawLine(table, []time.Time, []float64, color)
	// canvas.DrawBar(table,  []time.Time, []bar.Bar, color)
	// canvas.DrawBuy(table,  []time.Time, []float64,color)
	// canvas.DrawSell(table, []time.Time, []float64,color)
	// canvas.DrawSpark(table,[]time.Time, []float64,color)
	// canvas.DrawShit(table, []time.Time, []float64,color)
}
