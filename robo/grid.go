package robo

type Point2D [2]int

func (p *Point2D) Set(x, y int) {
	p[0], p[1] = x, y
}

func (p *Point2D) Get() (int, int) {
	return p[0], p[1]
}

type Grid2D struct {
	O Point2D
	W int
	H int
}

func NewGrid2D(w, h int) *Grid2D {
	grid := &Grid2D{
		Point2D{},
		w,
		h,
	}
	grid.O.Set(0, 0)
	return grid
}
