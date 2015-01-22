package base

import (
	"fmt"
	"quant/base/bar"
	"quant/base/series"
)

func init() {
	if debug {
		fmt.Println("quant/base/bar/barmanager")
	}
}

var (
	BarManager = map[string]*series.BarSeries{}
)

func StoreBarToManager(symbol string, bar_ bar.Bar) {
	barSeries, ok := BarManager[symbol]
	if ok {
		barSeries.AppendBar(bar_)
		return
	} else {
		panic("symbol not found:" + symbol)
	}
}

func InitBarManagerWithSymbol(symbol string) *series.BarSeries {
	oldBarSeries, ok := BarManager[symbol]
	if ok {
		return oldBarSeries
	}

	// new series
	barSeries := series.NewBarSeries()
	barSeries.Symbol = symbol
	BarManager[symbol] = barSeries
	return barSeries
}
