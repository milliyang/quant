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
	DefaultQuant = Quant{}

	ErrNoProject = errors.New("no project")
)

type Quant struct {
	Name     string
	Projects []*BaseProject
}

func init() {
	if debug {
		fmt.Println("quant.project.init")
	}
}

type BaseProject struct {
	Name                  string
	Instrucment           []string
	BasicStrategies       []interface{}
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
		fmt.Println("quant.project", strategy)
	}
	this.BasicStrategies = append(this.BasicStrategies, strategy)
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

		oneProject.Instrucment = append(oneProject.Instrucment, "stock0000001")
		oneProject.Instrucment = append(oneProject.Instrucment, "stock0000002")
		oneProject.Instrucment = append(oneProject.Instrucment, "stock0000003")

		for _, ins := range oneProject.Instrucment {
			for _, bs := range oneProject.BasicStrategies {

				// newObjPtr := reflect.New(reflect.TypeOf(bs))
				// c := newObjPtr.Interface().(strategy.IStrategy)

				t := reflect.TypeOf(bs)
				v := reflect.New(t)
				c := v.Interface().(strategy.IStrategy)
				c.Init()

				fmt.Println("final:", v.Interface())

				oneProject.mapInstrumentStrategy[ins] = c
			}
		}

		for _, oneStrategy := range oneProject.mapInstrumentStrategy {
			oneStrategy.OnStrategyStart()
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
