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
	ext := path.Ext(r.URL.Path)
	if ext != "" && ext != ".png" && ext != ".svg" {
		http.Error(w, "Unknown file extension", http.StatusNotFound)
		return
	}

	width := 240
	if ws := r.URL.Query().Get("w"); ws != "" {
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
		config.MakeSVG(w, data)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, config.Make(width, data))
}

func md5hash(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func main() {
	http.HandleFunc("/", imageHandler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
