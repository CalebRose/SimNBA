package structs

import "github.com/jinzhu/gorm"

// Season - Data Structure for a Season
type Season struct {
	gorm.Model
	Year         int
	CollegeWeeks []CollegeWeek `gorm:"foreignKey:SeasonID"`
	NBAWeeks     []NBAWeek     `gorm:"foreignKey:SeasonID"`
}
