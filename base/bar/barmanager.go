package bar

import (
	"fmt"
)

func init() {
	if debug {
		fmt.Println("quant/base/bar/barmanager")
	}
}

var (
	BarManager = map[string]*BarSeries{}
)

func StoreBarToManager(symbol string, bar_ Bar) {
	barSeries, ok := BarManager[symbol]
	if ok {
		barSeries.AppendBar(bar_)
		return
	} else {
		panic("symbol not found:" + symbol)
	}
}

func InitBarManagerWithSymbol(symbol string) *BarSeries {
	oldBarSeries, ok := BarManager[symbol]
	if ok {
		return oldBarSeries
	}

	// new series
	barSeries := NewBarSeries()
	barSeries.Symbol = symbol
	BarManager[symbol] = barSeries
	return barSeries
}
