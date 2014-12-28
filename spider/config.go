package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Instructment struct {
	Symbol   string // SZ600001
	Type     string // stock
	YearFrom int    // 2013
	YearTo   int    // 2014
}

type StockConfig struct {
	Instructments []Instructment
}

func (this *Instructment) getSymbolNumber() string {
	if len(this.Symbol) == 8 {
		return this.Symbol[2:]
	} else {
		return this.Symbol
	}
}

func GetInstructmentFromConfigFile(path string) (ins []Instructment) {
	f, ferr := os.Open(path)
	if ferr != nil {
		fmt.Errorf("open failed: %s", ferr.Error())
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Errorf("read failed: %s", err.Error())
	} else {
		stockConfig := StockConfig{}
		err := json.Unmarshal(bytes, &stockConfig)
		if err != nil {
			fmt.Errorf("parse failed: %s", err.Error())
		}
		for _, ins := range stockConfig.Instructments {
			JsonPrint(ins)
		}
		return stockConfig.Instructments
	}
	return nil
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
