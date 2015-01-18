package indicator

import (
	"fmt"
)

func init() {
	if debug {
		fmt.Println("quant/base/indicator package init")
	}
}

type IIndecator interface {
	Draw(int)
}
