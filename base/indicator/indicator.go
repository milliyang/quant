package indicator

import (
	"fmt"
)

const (
	debug = false
)

func init() {
	if debug {
		fmt.Println("quant/base/indicator package init")
	}
}

type IIndecator interface {
	Draw(int)
}
