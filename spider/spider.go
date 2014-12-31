package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	file  *string = flag.String("file", "", "download config file")
	stock *string = flag.String("stock", "", "stock number. eg. 600031")

	SPIDER_DEBUG      = false
	SPIDER_MULTI_TASK = 10
)

func main() {
	flag.Parse()

	finishChan := make(chan bool)

	if *file == "" && *stock == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	if *file != "" {
		fmt.Println("file:", *file)
		err := parseConfigFile(*file)
		if err != nil {
			os.Exit(2)
		}
	} else if *stock != "" {
		AddStock(*stock)
	}

	_multiTaskC := make(chan int, SPIDER_MULTI_TASK)
	for _, ins := range Config.Instructments {
		_multiTaskC <- 1
		go downloadInstrucment(ins, finishChan)
		<-_multiTaskC
	}

	finishCounter := len(Config.Instructments)
	for {
		select {
		case _ = <-finishChan:
			finishCounter--
		}
		if finishCounter == 0 {
			break
		}
	}
	fmt.Println("all done")
}

func downloadInstrucment(ins Instructment, outC chan bool) {

	// only the newest seaon ( Now )
	if Config.DownloadFlag.Type == "recent" {
		season := getRecentSeason()
		bars, err := HttpGet(ins, season.Year, season.Quarter)
		if err != nil {
			fmt.Println("err:", err.Error())
		}
		if SPIDER_DEBUG {
			for _, oneBar := range bars {
				fmt.Printf(oneBar.toString())
			}
		}
		SaveInstrumentBars(ins, bars)
	} else {
		// update all data from
		seasons := getAllSeason()
		JsonPrint(seasons)
		bars := []Bar{}
		for _, season := range seasons {
			_bars, err := HttpGet(ins, season.Year, season.Quarter)
			if err != nil {
				fmt.Println("err:", err.Error())
			}
			for _, oneBar := range _bars {
				bars = append(bars, oneBar)
				if SPIDER_DEBUG {
					fmt.Printf(oneBar.toString())
				}
			}
		}
		SaveInstrumentBars(ins, bars)
	}
	outC <- true
}
