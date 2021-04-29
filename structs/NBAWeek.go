package structs

import "github.com/jinzhu/gorm"

// CollegeWeek - the Week of College Basketball in a Season
type NBAWeek struct {
	gorm.Model
	SeasonID    int
	Games       []Match `gorm:"foreignKey:WeekID"`
	IsOffseason bool
}
