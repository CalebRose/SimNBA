package structs

import "github.com/jinzhu/gorm"

// Timestamp - The Global Timestamp for the Season
type Timestamp struct {
	gorm.Model
	SeasonID          int
	CollegeWeekID     int
	NBAWeekID         int
	CollegeWeek       int
	NBAWeek           int
	GamesARan         bool
	GamesBRan         bool
	RecruitingSynced  bool
	GMActionsComplete bool
	IsOffSeason       bool
}
