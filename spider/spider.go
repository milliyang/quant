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
	instruments := GetInstructmentFromConfigFile(*file)

	for _, ins := range instruments {
		HttpGet(ins)
	}
}
