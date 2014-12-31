package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var (
	OUTPUT = "./download/"
	EXPORT = "./export/"

	FILE_DEBUG = false
)

// 1. read old data
// 2. merge all data, overide duplicate
// 3. write to file
func SaveInstrumentBars(ins Instructment, bars []Bar) {
	// read
	oldBars := readBars(ins)

	if FILE_DEBUG {
		fmt.Println("<oldBars>")
		for _, bar_ := range oldBars {
			fmt.Printf(bar_.toString())
		}
		fmt.Println("</oldBars>")

		fmt.Println("<downloadBars>")
		for _, bar_ := range bars {
			fmt.Printf(bar_.toString())
		}
		fmt.Println("</downloadBars>")
	}
	// merge
	barMap := map[string]Bar{}
	for _, oneBar := range oldBars {
		barMap[oneBar.Date] = oneBar
	}
	for _, oneBar := range bars {
		barMap[oneBar.Date] = oneBar
	}

	// sort
	keys := []string{}
	for k, _ := range barMap {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))
	if FILE_DEBUG {
		JsonPrint(keys)
	}

	newBars := []Bar{}
	for _, k := range keys {
		newBars = append(newBars, barMap[k])
	}

	if FILE_DEBUG {
		fmt.Println("<newBars>")
		for _, bar_ := range newBars {
			fmt.Printf(bar_.toString())
		}
		fmt.Println("</newBars>")
	}
	// save
	saveBars(ins, newBars)
	exportBars(ins, newBars)
}

func readBars(ins Instructment) []Bar {

	outfile := ins.getFileName()

	bars := []Bar{}

	f, err := os.Open(outfile)
	if err != nil {
		fmt.Println(err, ins.getFileName())
		return bars
	}
	defer f.Close()

	all, err := ioutil.ReadAll(f)
	if err == nil {
		allLine := string(all)
		lines := strings.Split(allLine, "\n")
		for _, oneLine := range lines {
			if strings.Contains(oneLine, "#") {
				continue
			}
			items := strings.Split(oneLine, ",")
			oneBar := Bar{}

			// 2007-08-08,48.80,58.95,50.80,48.80,35241892.00,1861702656.00,4.00
			if len(items) == 8 {
				for idx, item := range items {
					if idx == 0 {
						oneBar.Date = item
					} else {
						oneBar.Items = append(oneBar.Items, item)
					}
				}
				bars = append(bars, oneBar)
			}
		}
	}
	return bars
}

func saveBars(ins Instructment, bars []Bar) {
	outfile := ins.getFileName()

	os.Remove(outfile)

	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	f.WriteString("#date,open,high,low,close,volumn(share),amount(CNY),factor\n")
	for _, bar := range bars {
		f.WriteString(bar.toString())
	}
}

// no factor
// Open Quant Format
func exportBars(ins Instructment, bars []Bar) {
	outfile := ins.getExportFileName()
	os.Remove(outfile)

	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// newest factor
	var factor float64 = 1.0
	if len(bars) > 0 {
		lastBar := bars[len(bars)-1]
		if lastBar.hasFactor() {
			factor = lastBar.getFactor()
		} else {
			fmt.Println("invalid bar. no factor")
			os.Exit(2)
		}
	}

	f.WriteString("#date,open,high,low,close,volumn(share/10K),amount(CNY/10K)\n")
	for _, bar := range bars {
		numBar := bar.toNumBar(factor)
		f.WriteString(numBar.toString())
	}
}
