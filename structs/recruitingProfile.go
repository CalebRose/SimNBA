package structs

import "github.com/jinzhu/gorm"

// RecruitingProfile - The profile for a team for recruiting
type RecruitingProfile struct {
	gorm.Model
	TeamID                int
	State                 string
	Region                string
	ScholarshipsAvailable int
	WeeklyPoints          int
	BonusPoints           int
}
