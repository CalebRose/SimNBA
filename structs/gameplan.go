package structs

import "github.com/jinzhu/gorm"

// Gameplan - A team's strategy for their weekly gameplan
type Gameplan struct {
	gorm.Model
	TeamID               int
	Game                 string
	Pace                 int
	ThreePointProportion int
	JumperProportion     int
	PaintProportion      int
}

// UpdatePace - Update the Pace of the Gameplan
func (g *Gameplan) UpdatePace(pace int) {
	g.Pace = pace
}

// Update3PtProportion
func (g *Gameplan) Update3PtProportion(ratio int) {
	g.ThreePointProportion = ratio
}

func (g *Gameplan) UpdateJumperProportion(ratio int) {
	g.JumperProportion = ratio
}

func (g *Gameplan) UpdatePaintProportion(ratio int) {
	g.PaintProportion = ratio
}
