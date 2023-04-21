package structs

import "gorm.io/gorm"

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
	SigningStatus         string
	TopRankModifier       float64
	RivalsRank            float64
	ESPNRank              float64
	Rank247               float64
	RecruitModifier       int
	IsSigned              bool
	IsTransfer            bool
	IsCustomCroot         bool
	CreatedFor            string
	RecruitProfiles       []PlayerRecruitProfile `gorm:"foreignKey:RecruitID"`
	// RecruitPoints         []RecruitPointAllocation `gorm:"foreignKey:RecruitID"`
}

func (r *Recruit) Map(createRecruitDTO CreateRecruitDTO, lastPlayerID uint, expectations int) {
	r.ID = lastPlayerID
	r.FirstName = createRecruitDTO.FirstName
	r.LastName = createRecruitDTO.LastName
	r.Position = createRecruitDTO.Position
	r.Age = 18
	r.Height = createRecruitDTO.Height
	r.State = createRecruitDTO.State
	r.Country = createRecruitDTO.Country
	r.Stars = createRecruitDTO.Stars
	r.Overall = createRecruitDTO.Overall
	r.Stamina = createRecruitDTO.Stamina
	r.Ballwork = createRecruitDTO.Ballwork
	r.Defense = createRecruitDTO.Defense
	r.Finishing = createRecruitDTO.Finishing
	r.Rebounding = createRecruitDTO.Rebounding
	r.Shooting2 = createRecruitDTO.Shooting2
	r.Shooting3 = createRecruitDTO.Shooting3
	r.Potential = createRecruitDTO.Potential
	r.PotentialGrade = createRecruitDTO.PotentialGrade
	r.PlaytimeExpectations = expectations
	r.WorkEthic = createRecruitDTO.WorkEthic
	r.FreeAgency = createRecruitDTO.FreeAgency
	r.Personality = createRecruitDTO.Personality
	r.RecruitingBias = createRecruitDTO.RecruitingBias
	r.AcademicBias = createRecruitDTO.AcademicBias
	r.IsSigned = false
	r.IsCustomCroot = true
	r.CreatedFor = createRecruitDTO.CreatedFor
	r.SigningStatus = "Not Ready"
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
	r.FreeThrow = attr.FreeThrow
	r.Rebounding = attr.Rebounding
	r.Ballwork = attr.Ballwork
	r.InteriorDefense = attr.InteriorDefense
	r.PerimeterDefense = attr.PerimeterDefense
	r.Finishing = attr.Finishing
}

func (r *Recruit) FixRecruit(grade string, pro int, mod int) {
	r.PotentialGrade = grade
	r.ProPotentialGrade = pro
	r.RecruitModifier = mod
}

func (r *Recruit) SetNewAttributes(ft int, id int, pd int) {
	r.FreeThrow = ft
	r.InteriorDefense = id
	r.PerimeterDefense = pd
}

func (r *Recruit) AssignRankValues(rank247 float64, espnRank float64, rivalsRank float64, modifier float64) {
	r.Rank247 = rank247
	r.ESPNRank = espnRank
	r.RivalsRank = rivalsRank
	r.TopRankModifier = modifier
}

func (r *Recruit) ApplySigningStatus(num float64, threshold float64) {
	percentage := num / threshold

	if threshold == 0 || num == 0 || percentage < 0.26 {
		r.SigningStatus = "Not Ready"
	} else if percentage < 0.51 {
		r.SigningStatus = "Hearing Offers"
	} else if percentage < 0.76 {
		r.SigningStatus = "Narrowing Down Offers"
	} else if percentage < 0.96 {
		r.SigningStatus = "Finalizing Decisions"
	} else if percentage < 1 {
		r.SigningStatus = "Ready to Sign"
	} else {
		r.SigningStatus = "Signed"
	}
}

func (b *Recruit) SetNewPosition(pos string) {
	b.Position = pos
}
