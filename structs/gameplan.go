package structs

// Gameplan - A team's strategy for their weekly gameplan
type Gameplan struct {
	ID                   int
	TeamID               int
	Game                 string
	Pace                 int
	ThreePointProportion int
	JumperProportion     int
	PaintProportion      int
}
