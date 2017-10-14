// package main

// import (
// 	"github.com/ajstarks/svgo"
// 	"os"
// )

// func main() {
// 	width := 5000
// 	height := 5000
// 	canvas := svg.New(os.Stdout)
// 	canvas.Start(width, height)
// 	canvas.Circle(width/2, height/2, 1000)
// 	canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
// 	canvas.End()
// }

package main

import (
	"github.com/ajstarks/svgo"
	"log"
	"net/http"
)

func main() {
	DrawData()
	return

	http.Handle("/circle", http.HandlerFunc(circle))
	err := http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func circle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(5000, 5000)
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	s.End()
}
