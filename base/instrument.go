package base

import (
	"fmt"
	_ "quant/base/bar"
)

const (
	debug = false
)

func init() {
	if debug {
		fmt.Println("quant/base/instument init")
	}
}

type Instrument struct {
	Symbol string
}
