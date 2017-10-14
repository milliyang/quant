package color

// A Color specifies a Color of html/css
// only cover a few at this time [ToDoList]

const (
	Red int = 0 + iota
	Deeppink
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

var css_colors = [...]string{
	"fill:none;stroke:red",
	"fill:none;stroke:deeppink",
	"fill:none;stroke:hotpink",
	"fill:none;stroke:pink",
	"fill:none;stroke:lightpink",
	"fill:none;stroke:indigo",
	"fill:none;stroke:purple",
	"fill:none;stroke:violet",
	"fill:none;stroke:orangered",
	"fill:none;stroke:tomato",
	"fill:none;stroke:gold",
	"fill:none;stroke:yellow",
	"fill:none;stroke:olive",
	"fill:none;stroke:green",
	"fill:none;stroke:yellowgreen",
	"fill:none;stroke:navy",
	"fill:none;stroke:black",
	"fill:none;stroke:gray",
}

func GetCssColor(idx int) string {
	str := css_colors[idx]
	return str
}
