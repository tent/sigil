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
	img := image.NewNRGBA(image.Rect(0, 0, width, width))
	fg, bg := s.colors(data[0], inverted)
	colWidth := s.colWidth(width)
	for _, c := range s.cells(colWidth, data[1:]) {
		for x := c.rect.Min.X; x < c.rect.Max.X; x++ {
			for y := c.rect.Min.Y; y < c.rect.Max.Y; y++ {
				fill := fg
				if !c.fill {
					fill = bg
				}
				i := y*img.Stride + x*4
				img.Pix[i+0] = fill.R
				img.Pix[i+1] = fill.G
				img.Pix[i+2] = fill.B
				img.Pix[i+3] = fill.A
			}
		}
	}
	padding := colWidth / 2
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			if x < padding || y < padding || x > width-padding-1 || y > width-padding-1 {
				i := y*img.Stride + x*4
				img.Pix[i+0] = bg.R
				img.Pix[i+1] = bg.G
				img.Pix[i+2] = bg.B
				img.Pix[i+3] = bg.A
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
	for _, c := range s.cells(s.colWidth(width), data[1:]) {
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

func (s *Sigil) cells(width int, data []byte) []cell {
	cols := s.Rows/2 + s.Rows%2
	cells := cols * s.Rows
	res := make([]cell, 0, s.Rows*s.Rows)
	padding := width / 2
	for i := 0; i < cells; i++ {
		column := i / s.Rows
		row := i % s.Rows
		c := cell{fill: s.fill(i, data)}

		pt := image.Pt(padding+(column*width), padding+(row*width))
		c.rect = image.Rectangle{pt, image.Pt(pt.X+width, pt.Y+width)}
		if s.Rows%2 == 0 && column == cols-1 {
			// last/middle column, double width
			c.rect.Max.X += width
		}
		res = append(res, c)

		if column < cols-1 {
			// add mirrored column
			c.rect.Min.X = padding + ((s.Rows - column - 1) * width)
			c.rect.Max.X = c.rect.Min.X + width
			res = append(res, c)
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

func (s *Sigil) colWidth(w int) int {
	return w / (s.Rows + 1)
}
