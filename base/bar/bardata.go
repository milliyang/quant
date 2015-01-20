package bar

type BarField int

const (
	Close BarField = 1 + iota
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

var BarFieldName = [...]string{
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
