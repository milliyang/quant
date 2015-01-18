package bar

import (
	"time"
)

type BarSeries struct {
	DateTime time.Time
	Close    float64
	Open     float64
	High     float64
	Low      float64
	Volumn   float64
	Size     float64
	Type     BarType
}
