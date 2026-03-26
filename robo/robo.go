package robo

type FaceDirection int

const (
	FaceNorth FaceDirection = iota
	FaceEast
	FaceSouth
	FaceWest
)

type Robot struct {
	Pos  Point2D
	Face FaceDirection
}

func NewRobot(posX, posY int, face FaceDirection) *Robot {
	return &Robot{
		Point2D{posX, posY},
		face,
	}
}

func (r *Robot) Move() {
	switch r.Face {
	case FaceNorth:
		r.stepUp()
	case FaceEast:
		r.stepRight()
	case FaceSouth:
		r.stepDown()
	case FaceWest:
		r.stepLeft()
	}
}

func (r *Robot) stepLeft() {
	x, y := r.Pos.Get()
	r.Pos.Set(x-1, y)
}

func (r *Robot) stepRight() {
	x, y := r.Pos.Get()
	r.Pos.Set(x+1, y)
}

func (r *Robot) stepUp() {
	x, y := r.Pos.Get()
	r.Pos.Set(x, y+1)
}

func (r *Robot) stepDown() {
	x, y := r.Pos.Get()
	r.Pos.Set(x, y-1)
}
