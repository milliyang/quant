package main

import (
	"flag"
	"fmt"
)

var (
	file *string = flag.String("file", "", "input download file")
)

func main() {
	flag.Parse()

	if *file == "" {
		flag.PrintDefaults()
		return
	}

	fmt.Println("file:", *file)
	err := parseConfigFile(*file)
	if err == nil {
		for _, ins := range Config.Instructments {
			downloadInstrucment(ins)
		}
	}
}

func downloadInstrucment(ins Instructment) {
	// only the newest seaon ( Now )
	if Config.DownloadFlag.Type == "recent" {
		season := getRecentSeason()
		bars, err := HttpGet(ins, season.Year, season.Quarter)
		if err != nil {
			fmt.Println("err:", err.Error())
		}
		for _, oneBar := range bars {
			fmt.Println(oneBar.Date, oneBar.Items)
		}
	} else {
		// update all data from
		seasons := getAllSeason()
		JsonPrint(seasons)
		for _, season := range seasons {
			bars, err := HttpGet(ins, season.Year, season.Quarter)
			if err != nil {
				fmt.Println("err:", err.Error())
			}
			for _, oneBar := range bars {
				fmt.Println(oneBar.Date, oneBar.Items)
			}
		}
	}
}
