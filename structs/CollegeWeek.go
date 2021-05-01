package structs

import "github.com/jinzhu/gorm"

// CollegeWeek - the Week of College Basketball in a Season
type CollegeWeek struct {
	gorm.Model
	WeekNumber  int
	SeasonID    int
	IsOffseason bool
	Games       []Match `gorm:"foreignKey:WeekID"`
}
