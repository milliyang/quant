package quant

import (
	"errors"
	"fmt"
	"quant/base/bar"
	"quant/base/strategy"
	"quant/provider"
	"reflect"
	"time"
)

const (
	debug = true
)

var (
	DefaultQuant = Quant{instruments: map[string]int{}, AccountDate: time.Time{}}

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
	AccountDate time.Time
	lastBar     *bar.Bar
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
				c.Init(ins)
				oneProject.allStrategy = append(oneProject.allStrategy, c)
			}
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
		for _, ss := range oneProject.allStrategy {
			ss.OnStrategyStart()
		}
	}

	// receive datagram from other modules
	go DefaultQuant.Provider.Receive(datagramChan)

	for datagramPtr := range datagramChan {
		DefaultQuant.handleOneBar(datagramPtr)
	}

	// start all strategy
	for _, oneProject := range DefaultQuant.Projects {
		for _, ss := range oneProject.allStrategy {
			ss.OnStrategyStop()
		}
	}
}

func (this *Quant) handleOneBar(dgram *provider.Datagram) {
	if dgram.Event == provider.EventBarOpen {
		// ToDo:
		// not ready yet, haven't figure out how to implement it
		for _, oneProject := range DefaultQuant.Projects {
			for _, oneStrategy := range oneProject.allStrategy {
				return

				// ToDo:
				oneStrategy.OnBarOpen(bar.Bar{})
			}
		}
	} else if dgram.Event == provider.EventBarSlice {
		for _, oneProject := range DefaultQuant.Projects {
			for _, oneStrategy := range oneProject.allStrategy {
				oneStrategy.OnBarSlice(1)
			}
		}
	} else {
		newBar := bar.NewBar(dgram.Time, dgram.Open, dgram.High, dgram.Low, dgram.Close, dgram.Volumn, dgram.Amount)
		// for each instrument , for each bar
		for _, oneProject := range DefaultQuant.Projects {
			for _, oneStrategy := range oneProject.allStrategy {
				if oneStrategy.Match(dgram.Symbol) {
					oneStrategy.OnBar(*newBar)
				}
			}
		}
	}
}
