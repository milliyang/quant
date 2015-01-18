package bar

type BarData int

const (
	Close BarData = 1 + iota
	Open
	High
	Low
	Median
	Typical
	Weighted
	Average
	Volume
	OpenInt
)

var BarDataName = [...]string{
	"Close",
	"Open",
	"High",
	"Low",
	"Median",
	"Typical",
	"Weighted",
	"Average",
	"Volume",
	"OpenInt",
}
