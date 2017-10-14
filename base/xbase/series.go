package xbase

import (
	"time"
)

const (
	enable_safe_mode = true
	debug            = false
)

type ISeries interface {
	Keys() []time.Time
	Values() []float64
	Count() int

	Now() time.Time

	//
	Index(datetime *time.Time) int
	ValueAtIndex(index int) float64
	ValueAtTime(datetime *time.Time) float64

	// event driven interface
	Match(symbol string) bool

	AddChild(child ISeries)
	Append(datetime *time.Time, value float64)
}
