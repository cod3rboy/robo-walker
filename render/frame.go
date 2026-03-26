package render

import (
	"image"
	"image/color"

	"github.com/cod3rboy/robo-walker/robo"
)

type Color uint32

func (c Color) RGBA() color.RGBA {
	return color.RGBA{
		uint8(0xFF000000 & c),
		uint8(0x00FF0000 & c),
		uint8(0x0000FF00 & c),
		uint8(0x000000FF & c),
	}
}

type DrawOpts struct {
	FgColor  Color
	BgColor  Color
	PosColor Color
}

func RenderImageFromSnapshot(s robo.Snapshot, opts DrawOpts) *image.RGBA {
	w := len(s)
	h := 0
	if len(s) > 0 {
		h = len(s[0])
	}
	upperLeft := image.Point{0, 0}
	lowerRight := image.Point{w, h}

	img := image.NewRGBA(image.Rectangle{upperLeft, lowerRight})

	lastFgPoint := image.Point{0, 0}
	for x := range s {
		for y := range s[x] {
			isBg := !s[x][y]
			drawColor := opts.BgColor.RGBA()
			if !isBg {
				drawColor = opts.FgColor.RGBA()
				lastFgPoint.X = x
				lastFgPoint.Y = y
			}

			img.Set(x, y, drawColor)
		}
	}

	img.Set(
		lastFgPoint.X,
		lastFgPoint.Y,
		opts.PosColor.RGBA(),
	)

	return img
}
