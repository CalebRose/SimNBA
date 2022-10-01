package structs

import "github.com/jinzhu/gorm"

type Recruit struct {
	gorm.Model
	PlayerID int
	TeamID   int
	TeamAbbr string
	BasePlayer
	UninterestedThreshold int
	LowInterestThreshold  int
	MedInterestThreshold  int
	HighInterestThreshold int
	ReadyToSignThreshold  int
	IsSigned              bool
}
