package main

import (
	"crypto/md5"
	"image/color"
	"image/png"
	"log"
	"net/http"

	"github.com/cupcake/sigil/gen"
)

var config = gen.Sigil{
	Width:   840,
	Columns: 3,
	Rows:    5,
	Foreground: []color.Color{
		rgb(45, 79, 255),
		rgb(44, 172, 0),
		rgb(254, 180, 44),
		rgb(226, 121, 234),
		rgb(30, 179, 253),
		rgb(232, 77, 65),
		rgb(49, 203, 115),
		rgb(141, 69, 170),
		rgb(252, 125, 31),
	},
	Background: rgb(224, 224, 224),
}

func rgb(r, g, b uint8) color.NRGBA { return color.NRGBA{r, g, b, 255} }

func imageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, config.Make(md5hash(r.URL.Path[1:])))
}

func md5hash(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func main() {
	http.HandleFunc("/", imageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
