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

type Project struct {
	Name                  string
	Instrucment           []string
	basicStrategies       []interface{}
	mapInstrumentStrategy map[string]strategy.IStrategy
}

func NewProject() *Project {
	project := Project{}
	project.mapInstrumentStrategy = map[string]strategy.IStrategy{}

	registerProject(&project)
	return &project
}

func (this *Project) Strategy(strategy interface{}) {
	if debug {
		fmt.Println("Project Strategy", strategy)
	}
	this.basicStrategies = append(this.basicStrategies, strategy)
}

func (this *Project) AddInstrument(symbol string) {
	if debug {
		fmt.Println("Project.AddInstrument", symbol)
	}
	this.Instrucment = append(this.Instrucment, symbol)
}