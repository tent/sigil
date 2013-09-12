package gen

import (
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/ajstarks/svgo"
)

type Sigil struct {
	Rows       int
	Foreground []color.NRGBA
	Background color.NRGBA
}

func (s *Sigil) Make(width int, inverted bool, data []byte) image.Image {
	fg, bg := s.colors(data[0], inverted)
	palette := color.Palette{bg, fg}
	img := image.NewPaletted(image.Rect(0, 0, width, width), palette)
	for _, rect := range s.cells(width, data[1:]) {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			for y := rect.Min.Y; y < rect.Max.Y; y++ {
				img.Pix[y*img.Stride+x] = 1
			}
		}
	}
	return img
}

func (s *Sigil) MakeSVG(w io.Writer, width int, inverted bool, data []byte) {
	canvas := svg.New(w)
	fg, bg := s.colors(data[0], inverted)
	fgFill, bgFill := svgFill(fg), svgFill(bg)

	canvas.Start(width, width)
	canvas.Rect(0, 0, width, width, bgFill)
	for _, rect := range s.cells(width, data[1:]) {
		canvas.Rect(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(), fgFill)
	}
	canvas.End()
}

func svgFill(c color.NRGBA) string {
	return fmt.Sprintf("fill:rgba(%d,%d,%d,%.2g);", c.R, c.G, c.B, float64(c.A)*1/255)
}

func (s *Sigil) fill(cell int, data []byte) bool {
	if data[cell/8]>>uint(8-((cell%8)+1))&1 == 0 {
		return false
	}
	return true
}

func (s *Sigil) cells(width int, data []byte) []image.Rectangle {
	width = width / (s.Rows + 1)
	cols := s.Rows/2 + s.Rows%2
	cells := cols * s.Rows
	res := make([]image.Rectangle, 0, s.Rows*s.Rows)
	padding := width / 2
	for i := 0; i < cells; i++ {
		if !s.fill(i, data) {
			continue
		}

		column := i / s.Rows
		row := i % s.Rows

		pt := image.Pt(padding+(column*width), padding+(row*width))
		rect := image.Rectangle{pt, image.Pt(pt.X+width, pt.Y+width)}
		if s.Rows%2 == 0 && column == cols-1 {
			// last/middle column, double width
			rect.Max.X += width
		}
		res = append(res, rect)

		if column < cols-1 {
			// add mirrored column
			rect.Min.X = padding + ((s.Rows - column - 1) * width)
			rect.Max.X = rect.Min.X + width
			res = append(res, rect)
		}
	}
	return res
}

func (s *Sigil) colors(b byte, inverted bool) (color.NRGBA, color.NRGBA) {
	fg := s.Foreground[int(b)%len(s.Foreground)]
	if inverted {
		return s.Background, fg
	}
	return fg, s.Background
}
