package strategy

import (
	"fmt"
	"quant/base/bar"
	"quant/base/series"
	"quant/base/xbase"
	_ "quant/svgo"
)

const (
	debug = true
)

func init() {
	if debug {
		fmt.Println("quant/base/stratety init")
	}
}

type Strategy struct {
	Name      string
	Symbol    string
	BarSeries *series.BarSeries

	mapIndicator map[int][]xbase.IIndecator
}

func (this *Strategy) Init(symbol string, barSeries *series.BarSeries) {
	if debug {
		fmt.Println("Strategy.Init()")
	}
	this.Name = "Strategy"
	this.Symbol = symbol
	this.BarSeries = barSeries

	this.mapIndicator = map[int][]xbase.IIndecator{}
}

func (this *Strategy) Match(symbol string) bool {
	if this.Symbol == symbol {
		return true
	} else {
		return false
	}
}

func (this *Strategy) OnStrategyStart() {

}

func (this *Strategy) OnStrategyStop() {

}

func (this *Strategy) OnBarOpen(bar bar.Bar) {

}

func (this *Strategy) OnBar(bar bar.Bar) {

}

func (this *Strategy) OnBarSlice(size int) {

}

func (this *Strategy) Draw(table int, indicator_ xbase.IIndecator) {

	indicatorSlice, ok := this.mapIndicator[table]
	if ok {
		indicatorSlice = append(indicatorSlice, indicator_)
	} else {
		this.mapIndicator[table] = []xbase.IIndecator{indicator_}
	}
}

func (this *Strategy) OnDraw(table int, indicator_ xbase.IIndecator) {

	indicatorSlice, ok := this.mapIndicator[table]
	if ok {
		indicatorSlice = append(indicatorSlice, indicator_)
	} else {
		this.mapIndicator[table] = []xbase.IIndecator{indicator_}
	}

	// svgo.TestDraw(indicator_)
	// this.indicators = append(this.indicators, indicator_)
}

type IStrategy interface {
	Init(string, *series.BarSeries)
	Match(string) bool

	OnStrategyStart()
	OnStrategyStop()

	OnBarOpen(bar.Bar)
	OnBar(bar.Bar)
	OnBarSlice(int)

	// public virtual void OnTrade(Trade trade)
	// public virtual void OnQuote(Quote quote)
}
