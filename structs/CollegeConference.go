package structs

import "github.com/jinzhu/gorm"

type CollegeConference struct {
	gorm.Model
	ConferenceName              string
	ConferenceAbbr              string
	LatestChampion              string
	LatestRegularSeasonChampion string
}

type NBAConference struct {
	gorm.Model
	ConferenceName              string
	ConferenceAbbr              string
	LatestChampion              string
	LatestRegularSeasonChampion string
}

type NBADivision struct {
	gorm.Model
	DivisionName                string
	DivisionAbbr                string
	LatestChampion              string
	LatestRegularSeasonChampion string
}
