package robo

type OType int8

const (
	OGround OType = 0
	OTrail        = 0b00000001 << iota
	ORobot
)

type Snapshot [][]OType

func (s Snapshot) Copy() Snapshot {
	cpy := make(Snapshot, len(s))
	for i := range cpy {
		inner := make([]OType, len(s[i]))
		for j := range inner {
			inner[j] = s[i][j]
			if inner[j] == ORobot {
				inner[j] = OTrail
			}
		}
		cpy[i] = inner
	}
	return cpy
}

type World struct {
	grid      *Grid2D
	robot     *Robot
	snapshots []Snapshot
}

func NewWorld(size int) *World {
	return &World{
		grid:  NewGrid2D(size, size),
		robot: NewRobot(size/2, size/2, FaceNorth),
	}
}

func (w *World) moveRobot(cmd MoveCommand) {
	for i := 0; i < cmd.Displacement; i++ {
		w.robot.Face = cmd.Direction
		w.robot.Move()
		newX, newY := w.robot.Pos.Get()
		newX, newY = newX%w.grid.W, newY%w.grid.H
		if newX < 0 {
			newX = w.grid.W + newX
		}
		if newY < 0 {
			newY = w.grid.H + newY
		}
		w.robot.Pos.Set(newX, newY)
		w.save()
	}
}

func (w *World) save() {
	var snapshot Snapshot
	if len(w.snapshots) == 0 {
		snapshot = make(Snapshot, w.grid.W)
		for i := range snapshot {
			snapshot[i] = make([]OType, w.grid.H)
		}
	} else {
		snapshot = w.snapshots[len(w.snapshots)-1].Copy()
	}
	x, y := w.robot.Pos.Get()
	snapshot[x][y] = ORobot
	w.snapshots = append(w.snapshots, snapshot)
}

func (w *World) Run(program string) error {
	var p Program = Program(program)
	cmds, err := p.Compile()
	if err != nil {
		return err
	}

	for _, cmd := range cmds {
		w.moveRobot(cmd)
	}

	return nil
}

func (w *World) Snapshots() []Snapshot {
	return w.snapshots
}
