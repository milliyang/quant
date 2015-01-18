package bar

type BarType int

const (
	Time BarType = 1 + iota
	Tick
	// Volume
	Range
)
