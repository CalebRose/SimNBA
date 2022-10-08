package structs

import "github.com/jinzhu/gorm"

type Recruit struct {
	gorm.Model
	PlayerID uint
	TeamID   uint
	TeamAbbr string
	BasePlayer
	UninterestedThreshold int
	LowInterestThreshold  int
	MedInterestThreshold  int
	HighInterestThreshold int
	ReadyToSignThreshold  int
	IsSigned              bool
	IsTransfer            bool
	PlayerRecruitProfiles []PlayerRecruitProfile `gorm:"foreignKey:RecruitID"`
	// RecruitPoints         []RecruitPointAllocation `gorm:"foreignKey:RecruitID"`
}

func (r *Recruit) SetID(id uint) {
	r.ID = uint(id)
}

func (r *Recruit) UpdateTeamID(id uint) {
	r.TeamID = id
	r.IsSigned = true
}

func (r *Recruit) AssignCollege(abbr string) {
	r.TeamAbbr = abbr
}

func (r *Recruit) UpdateSigningStatus() {
	r.IsSigned = true
}

func (r *Recruit) ProgressUnsignedRecruit(attr CollegePlayerProgressions) {
	r.Age++
	r.Shooting2 = attr.Shooting2
	r.Shooting3 = attr.Shooting3
	r.Rebounding = attr.Rebounding
	r.Ballwork = attr.Ballwork
	r.Defense = attr.Defense
	r.Finishing = attr.Finishing
}
