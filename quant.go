package quant

import (
	"errors"
	"fmt"
	"quant/base/bar"
	"quant/base/strategy"
	"reflect"
)

const (
	debug = true
)

var (
	DefaultQuant = Quant{instruments: map[string]int{}}

	ErrNoProject = errors.New("no project")
)

func init() {
	if debug {
		fmt.Println("quant.quant.init")
	}
}

type Quant struct {
	Name        string
	Projects    []*BaseProject
	instruments map[string]int
}

func (this *Quant) addInstrument(symbol string) {
	count, ok := this.instruments[symbol]
	if ok {
		this.instruments[symbol] = count + 1
	} else {
		this.instruments[symbol] = 1
	}
}

func registerProject(pro *BaseProject) {
	if debug {
		fmt.Println("quant.Project Add:", pro.Name)
	}
	DefaultQuant.Projects = append(DefaultQuant.Projects, pro)
}

func Run() {
	// check
	if len(DefaultQuant.Projects) == 0 {
		panic(ErrNoProject)
	}

	// setup
	for _, oneProject := range DefaultQuant.Projects {
		for _, ins := range oneProject.Instrucment {

			DefaultQuant.addInstrument(ins)
			for _, bs := range oneProject.basicStrategies {
				t := reflect.TypeOf(bs)
				v := reflect.New(t)
				/*
					hack to set Symbol of Strategy, but fail  -_-!!
					use Init(symbol) instead
				*/
				/*
					vv := reflect.Indirect(v)
					typeOfT := vv.Type()
					for i := 0; i < vv.NumField(); i++ {
						if typeOfT.Field(i).Name == "Symbol" {
							vv.SetString(ins)
						}
						fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, vv.Type(), v.Interface())
					}
				*/
				c := v.Interface().(strategy.IStrategy)
				oneProject.mapInstrumentStrategy[ins] = c
			}
		}

		for ins, ss := range oneProject.mapInstrumentStrategy {
			ss.Init(ins)
		}
	}

	// connect to bar provider

	// request our stock bar data
	if debug {
		fmt.Println("request stock:", DefaultQuant.instruments)
	}

	// get bar data from other modules

	// start all strategy
	for _, oneProject := range DefaultQuant.Projects {
		for _, ss := range oneProject.mapInstrumentStrategy {
			ss.OnStrategyStart()
		}
	}

	ins := "stock000001"
	bar := bar.Bar{}

	// for each instrument , for each bar
	for _, oneProject := range DefaultQuant.Projects {
		oneStrategy, ok := oneProject.mapInstrumentStrategy[ins]
		if ok {
			oneStrategy.OnBar(bar)
			fmt.Println(oneStrategy, bar)
		}
	}
}
