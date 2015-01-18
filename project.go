package quant

import (
	"fmt"
	"quant/base/strategy"
)

func init() {
	if debug {
		fmt.Println("quant.project.init")
	}
}

type BaseProject struct {
	Name                  string
	Instrucment           []string
	basicStrategies       []interface{}
	mapInstrumentStrategy map[string]strategy.IStrategy
}

func NewProject() *BaseProject {
	project := BaseProject{}
	project.mapInstrumentStrategy = map[string]strategy.IStrategy{}

	registerProject(&project)
	return &project
}

func (this *BaseProject) Strategy(strategy interface{}) {
	if debug {
		fmt.Println("BaseProject Strategy", strategy)
	}
	this.basicStrategies = append(this.basicStrategies, strategy)
}

func (this *BaseProject) AddInstrument(symbol string) {
	if debug {
		fmt.Println("BaseProject.AddInstrument", symbol)
	}
	this.Instrucment = append(this.Instrucment, symbol)
}
