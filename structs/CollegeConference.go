package structs

import "github.com/jinzhu/gorm"

type CollegeConference struct {
	gorm.Model
	ConferenceName              string
	ConferenceAbbr              string
	LatestChampion              string
	LatestRegularSeasonChampion string
}
