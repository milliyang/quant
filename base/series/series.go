package series

import (
	"time"
)

const (
	enable_safe_mode = true
	debug            = true
)

type ISeries interface {
	Keys() []time.Time
	Values() []float64
	Count() int
	//
	Index(datetime *time.Time) int
	ValueAtIndex(index int) float64
	ValueAtTime(datetime *time.Time) float64

	// event driven interface
	Match(symbol string) bool

	AddChild(child ISeries)
	Append(datetime *time.Time, value float64)
}
