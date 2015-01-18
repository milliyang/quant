package color

// A Color specifies a Color of html/css
// only cover a few at this time [ToDoList]
type Color int

const (
	Red Color = 1 + iota
	Deeppink
	Red
	Hotpink
	Pink
	Lightpink
	Indigo
	Purple
	Violet
	Orangered
	Tomato
	Gold
	Yellow
	Olive
	Green
	Yellowgreen
	Navy
	Black
	Gray
)

var colors = [...]string{
	"red",
	"deeppink",
	"hotpink",
	"pink",
	"lightpink",
	"indigo",
	"purple",
	"violet",
	"orangered",
	"tomato",
	"gold",
	"yellow",
	"olive",
	"green",
	"yellowgreen",
	"navy",
	"black",
	"gray",
}
