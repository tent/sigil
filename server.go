package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/cupcake/sigil/gen"
)

var config = gen.Sigil{
	Rows: 5,
	Foreground: []color.NRGBA{
		rgb(45, 79, 255),
		rgb(254, 180, 44),
		rgb(226, 121, 234),
		rgb(30, 179, 253),
		rgb(232, 77, 65),
		rgb(49, 203, 115),
		rgb(141, 69, 170),
	},
	Background: rgb(224, 224, 224),
}

func rgb(r, g, b uint8) color.NRGBA { return color.NRGBA{r, g, b, 255} }

func imageHandler(w http.ResponseWriter, r *http.Request) {
	ext := path.Ext(r.URL.Path)
	if ext != "" && ext != ".png" && ext != ".svg" {
		ext = ""
	}

	width := 240
	q := r.URL.Query()
	if ws := q.Get("w"); ws != "" {
		var err error
		width, err = strconv.Atoi(ws)
		if err != nil {
			http.Error(w, "Invalid w parameter, must be an integer", http.StatusBadRequest)
			return
		}
		if width > 600 {
			http.Error(w, "Invalid w parameter, must be less than 600", http.StatusBadRequest)
			return
		}
		div := (config.Rows + 1) * 2
		if width%div != 0 {
			http.Error(w, "Invalid w parameter, must be evenly divisible by "+strconv.Itoa(div), http.StatusBadRequest)
			return
		}
	}
	inverted := false
	if inv := q.Get("inverted"); inv != "" && inv != "false" && inv != "0" {
		inverted = true
	}

	base := path.Base(r.URL.Path)
	base = base[:len(base)-len(ext)]
	var data []byte
	if len(base) == 32 {
		// try to decode hex MD5
		data, _ = hex.DecodeString(base)
	}
	if data == nil {
		data = md5hash(base)
	}

	etag := `"` + base64.StdEncoding.EncodeToString(data) + `"`
	w.Header().Set("Etag", etag)
	if cond := r.Header.Get("If-None-Match"); cond != "" {
		if strings.Contains(cond, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Cache-Control", "max-age=315360000")
	if ext == ".svg" {
		w.Header().Set("Content-Type", "image/svg+xml")
		config.MakeSVG(w, inverted, data)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, config.Make(width, inverted, data))
}

func md5hash(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func main() {
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/", imageHandler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
