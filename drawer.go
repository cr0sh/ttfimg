package ttfimg

import (
	"image"
	"image/draw"
	"io/ioutil"
	"math"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type drawer struct {
	Width, Height int
	font          font.Face
	size, dpi     float64
}

// NewDrawer returns a new drawer struct.
func NewDrawer(fname string, w, h int, fsize, dpi float64) (*drawer, error) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	ft, err := freetype.ParseFont(b)
	if err != nil {
		return nil, err
	}
	return &drawer{
		Width:  w,
		Height: h,
		font: truetype.NewFace(ft, &truetype.Options{
			Size: fsize,
			DPI:  dpi,
		}),
		size: fsize,
		dpi:  dpi,
	}, nil
}

// Draw returns RGBA image with given string written.
func (dr *drawer) Draw(str string) *image.RGBA {
	str = strings.Trim(str, "\r\n")
	rgba := image.NewRGBA(image.Rect(0, 0, dr.Width, dr.Height))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.ZP, draw.Src)
	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.Black,
		Face: dr.font,
	}
	init_x := (fixed.I(dr.Width) - d.MeasureString(strings.Split(str, "\n")[0])) / 2
	dy := math.Ceil(dr.size*dr.dpi/72) + 5
	lcnt := strings.Count(str, "\n")
	y := int(dr.Height)/2 + int(dy-5)/2 - lcnt*int(dy)/2
	for _, line := range strings.Split(str, "\n") {
		d.Dot = fixed.Point26_6{
			X: init_x,
			Y: fixed.I(y),
		}
		d.DrawString(line)
		y += int(dy)
	}
	return rgba
}
