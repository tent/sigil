package gen

import (
	"encoding/hex"
	"image/color"
	"image/png"
	"io/ioutil"
	"testing"
)

var data, _ = hex.DecodeString("0447550ed53b522827e17b7d7976dd3f")
var config = Sigil{
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

func BenchmarkMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config.Make(420, false, data)
	}
}

func BenchmarkMakeEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		png.Encode(ioutil.Discard, config.Make(420, false, data))
	}
}

func BenchmarkMakeSVG(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config.MakeSVG(ioutil.Discard, 420, false, data)
	}
}
