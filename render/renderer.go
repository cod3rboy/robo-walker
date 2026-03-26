package render

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"

	"github.com/cod3rboy/robo-walker/robo"
)

type Renderer interface {
	Render(opts DrawOpts) ([]byte, error)
}

type gifRenderer struct {
	ss    []robo.Snapshot
	delay int
}

func (r *gifRenderer) Render(opts DrawOpts) ([]byte, error) {
	frames := make([]*image.Paletted, len(r.ss))
	delays := make([]int, len(frames))
	for i := range frames {
		rgbaFrame := RenderImageFromSnapshot(r.ss[i], opts)
		paletted := image.NewPaletted(rgbaFrame.Bounds(), palette.Plan9)
		draw.FloydSteinberg.Draw(paletted, rgbaFrame.Bounds(), rgbaFrame, image.Point{})
		frames[i] = paletted
		delays[i] = r.delay
	}

	b := bytes.Buffer{}
	err := gif.EncodeAll(&b, &gif.GIF{
		Image:     frames,
		Delay:     delays,
		LoopCount: -1,
	})

	if err != nil {
		return b.Bytes(), fmt.Errorf("error render snapshot: %w", err)
	}

	return b.Bytes(), nil
}

func NewGIFRendererForSnapshots(snapshots []robo.Snapshot, delay int) Renderer {
	return &gifRenderer{
		ss:    snapshots,
		delay: delay,
	}
}
