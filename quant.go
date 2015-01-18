package quant

import (
	"errors"
	"fmt"
	"quant/base/bar"
	"quant/base/strategy"
	"quant/provider"
	"reflect"
)

const (
	debug = true
)

var (
	DefaultQuant = Quant{instruments: map[string]int{}}

	ErrNoProject  = errors.New("no project")
	ErrNoProvider = errors.New("no provider")
)

func init() {
	if debug {
		fmt.Println("quant.quant.init")
	}
}

type Quant struct {
	Name        string
	Projects    []*Project
	instruments map[string]int
	Provider    provider.IProvider
}

func (this *Quant) addInstrument(symbol string) {
	count, ok := this.instruments[symbol]
	if ok {
		this.instruments[symbol] = count + 1
	} else {
		this.instruments[symbol] = 1
	}
}

func registerProject(pro *Project) {
	if debug {
		fmt.Println("quant.Project Add:", pro.Name)
	}
	DefaultQuant.Projects = append(DefaultQuant.Projects, pro)
}

func RegisterProvider(provider_ provider.IProvider) {
	DefaultQuant.Provider = provider_
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
	if DefaultQuant.Provider == nil {
		panic(ErrNoProvider)
	}

	connectErr := DefaultQuant.Provider.Connect()
	if connectErr != nil {
		panic(connectErr)
	}

	symbols := []string{}
	for key, _ := range DefaultQuant.instruments {
		symbols = append(symbols, key)
	}
	requestErr := DefaultQuant.Provider.RequestInstrument(symbols)
	if requestErr != nil {
		panic(requestErr)
	}

	// request our stock bar data
	if debug {
		fmt.Println("request stock:", DefaultQuant.instruments)
	}

	datagramChan := make(chan *provider.Datagram, 100)

	// start all strategy
	for _, oneProject := range DefaultQuant.Projects {
		for _, ss := range oneProject.mapInstrumentStrategy {
			ss.OnStrategyStart()
		}
	}

	// get bar data from other modules
	go DefaultQuant.Provider.Receive(datagramChan)

	// select {

	// 	case datagram := <- datagramChan

	// }

	datagram := &provider.Datagram{}
	DefaultQuant.HandleOneBar(datagram)

}

func (this *Quant) HandleOneBar(datagram *provider.Datagram) {

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
