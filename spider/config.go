package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	CONFIG_YEAR = 2005
)

var (
	Config = StockConfig{}
)

type StockConfig struct {
	DownloadFlag  DownFlag
	Instructments []Instructment
}

type DownFlag struct {
	Type     string // {all, recent}
	YearFrom int    // 2014
}

type Instructment struct {
	Symbol string // SZ600001
	Type   string // stock
}

func (this *Instructment) getSymbolNumber() string {
	if len(this.Symbol) == 8 {
		return this.Symbol[2:]
	} else {
		return this.Symbol
	}
}

func (this *Instructment) getFileName() string {
	var outfile string
	if this.Type == "stock" {
		outfile = OUTPUT + this.Symbol + ".txt"
	} else {
		outfile = OUTPUT + this.Type + this.Symbol + ".txt"
	}
	return outfile
}

func (this *Instructment) getExportFileName() string {
	var outfile string
	if this.Type == "stock" {
		outfile = EXPORT + this.Symbol + ".txt"
	} else {
		outfile = EXPORT + this.Type + this.Symbol + ".txt"
	}
	return outfile
}

func parseConfigFile(path string) error {
	f, ferr := os.Open(path)
	if ferr != nil {
		fmt.Errorf("open failed: %s", ferr.Error())
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Errorf("read failed: %s", err.Error())
	} else {
		err := json.Unmarshal(bytes, &Config)
		if err != nil {
			fmt.Errorf("parse failed: %s", err.Error())
		}
		JsonPrint(Config)
	}
	return nil
}

func AddStock(stock string) {
	Config.DownloadFlag.Type = "all"
	Config.DownloadFlag.YearFrom = CONFIG_YEAR

	ins := Instructment{}
	ins.Symbol = stock
	ins.Type = "stock"
	Config.Instructments = append(Config.Instructments, ins)
}

func JsonPrint(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	os.Stdout.Write(b)
	fmt.Println("")
}
