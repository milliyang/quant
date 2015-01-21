package xbase

import (
	"fmt"
)

func init() {
	if debug {
		fmt.Println("quant/base/indicator package init")
	}
}

type IIndecator interface {
	OnMeasure() (int, int, float64, float64, int, float64) // X cordirate [datetime count], Y cordirate [min, max, num, percent]
	OnDraw(ICanvas)
}
