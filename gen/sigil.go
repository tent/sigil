package gen

import (
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/ajstarks/svgo"
)

type Sigil struct {
	Width      int
	Columns    int
	Rows       int
	Foreground []color.NRGBA
	Background color.NRGBA
	Inverted   bool
}

func (s *Sigil) Make(data []byte) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, s.Width, s.Width))
	fg, bg := s.colors(data[0])
	for _, c := range s.cells(data[1:]) {
		for x := c.rect.Min.X; x < c.rect.Max.X; x++ {
			for y := c.rect.Min.Y; y < c.rect.Max.Y; y++ {
				fill := fg
				if !c.fill {
					fill = bg
				}
				img.SetNRGBA(x, y, fill)
			}
		}
	}
	padding := s.colWidth() / 2
	for x := 0; x < s.Width; x++ {
		for y := 0; y < s.Width; y++ {
			if x < padding || y < padding || x > s.Width-padding-1 || y > s.Width-padding-1 {
				img.Set(x, y, bg)
			}
		}
	}
	return img
}

func (s *Sigil) MakeSVG(w io.Writer, data []byte) {
	canvas := svg.New(w)
	fg, bg := s.colors(data[0])
	fgFill, bgFill := svgFill(fg), svgFill(bg)

	canvas.Start(s.Width, s.Width)
	canvas.Rect(0, 0, s.Width, s.Width, bgFill)
	for _, c := range s.cells(data[1:]) {
		if c.fill {
			canvas.Rect(c.rect.Min.X, c.rect.Min.Y, c.rect.Dx(), c.rect.Dy(), fgFill)
		}
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

type cell struct {
	fill bool
	rect image.Rectangle
}

func (s *Sigil) cells(data []byte) []cell {
	res := make([]cell, 0, s.cols()*s.Rows)
	cells := s.Columns * s.Rows
	width := s.colWidth()
	cols := s.cols()
	padding := width / 2
	for i := 0; i < cells; i++ {
		column := i / s.Rows
		row := i % s.Rows
		c := cell{fill: s.fill(i, data)}

		pt := image.Pt(padding+(column*width), padding+(row*width))
		c.rect = image.Rectangle{pt, image.Pt(pt.X+width, pt.Y+width)}
		if s.Rows%2 == 0 && column == s.Columns-1 {
			// last/middle column, double width
			c.rect.Max.X *= 2
		}
		res = append(res, c)

		if column < s.Columns-1 {
			// add mirrored column
			c.rect.Min.X = padding + ((cols - column - 1) * width)
			c.rect.Max.X = c.rect.Min.X + width
			res = append(res, c)
		}
	}
	return res
}

func (s *Sigil) colors(b byte) (color.NRGBA, color.NRGBA) {
	fg := s.Foreground[int(b)%len(s.Foreground)]
	if s.Inverted {
		return s.Background, fg
	}
	return fg, s.Background
}

func (s *Sigil) colWidth() int {
	return s.Width / (s.cols() + 1)
}

func (s *Sigil) cols() int {
	cols := s.Columns * 2
	if s.Rows%2 == 1 {
		cols -= 1
	}
	return cols
}
