package structs

import "github.com/jinzhu/gorm"

// NBAGameplan - A team's strategy for their weekly gameplan
type NBAGameplan struct {
	gorm.Model
	TeamID               uint
	Game                 string
	Pace                 string
	ThreePointProportion int
	JumperProportion     int
	PaintProportion      int
}
