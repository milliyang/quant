package strategy

import (
	"fmt"
	"quant/base"
	"quant/base/bar"
)

const (
	debug = false
)

func init() {
	if debug {
		fmt.Println("quant/base/stratety init")
	}
}

type Strategy struct {
	Name       string
	Instrument map[string]*base.Instrument
}

func (this *Strategy) Init() {
	this.Name = "Strategy"
	this.Instrument = map[string]*base.Instrument{}
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

type IStrategy interface {
	Init()

	OnStrategyStart()
	OnStrategyStop()

	OnBarOpen(bar.Bar)
	OnBar(bar.Bar)
	OnBarSlice(int)

	// public virtual void OnTrade(Trade trade)
	// public virtual void OnQuote(Quote quote)
}
