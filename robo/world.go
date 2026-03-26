package robo

type Snapshot [][]bool

func (s Snapshot) Copy() Snapshot {
	cpy := make(Snapshot, len(s))
	for i := range cpy {
		inner := make([]bool, len(s[i]))
		copy(inner, s[i])
		cpy = append(cpy, inner)
	}
	return cpy
}

type World struct {
	grid      *Grid2D
	robot     *Robot
	snapshots []Snapshot
}

func NewWorld(width, height int) *World {
	return &World{
		grid:  NewGrid2D(width, height),
		robot: NewRobot(width/2, height/2, FaceNorth),
	}
}

func (w *World) moveRobot(cmd MoveCommand) {
	for i := 0; i < cmd.Displacement; i++ {
		w.robot.Face = cmd.Direction
		w.robot.Move()
		newX, newY := w.robot.Pos.Get()
		w.robot.Pos.Set(newX%w.grid.W, newY%w.grid.H)
		w.save()
	}
}

func (w *World) save() {
	var snapshot Snapshot
	if len(w.snapshots) == 0 {
		snapshot = make(Snapshot, w.grid.W)
		for i := range snapshot {
			snapshot[i] = make([]bool, w.grid.H)
		}
	} else {
		snapshot = w.snapshots[len(w.snapshots)-1].Copy()
	}
	x, y := w.robot.Pos.Get()
	snapshot[x][y] = true
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
