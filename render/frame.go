package render

import (
	"image"
	"image/color"

	"github.com/cod3rboy/robo-walker/robo"
)

type Color uint32

func (c Color) RGBA() color.RGBA {
	return color.RGBA{
		uint8((0xFF000000 & c) >> 24),
		uint8((0x00FF0000 & c) >> 16),
		uint8((0x0000FF00 & c) >> 8),
		uint8((0x000000FF & c) >> 0),
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

	for x := range s {
		for y := range s[x] {
			otype := s[x][y]
			drawColor := opts.BgColor.RGBA()
			switch otype {
			case robo.OTrail:
				drawColor = opts.FgColor.RGBA()
			case robo.ORobot:
				drawColor = opts.PosColor.RGBA()
			}

			img.Set(x, y, drawColor)
		}
	}
	return img
}
