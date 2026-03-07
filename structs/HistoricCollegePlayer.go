package structs

import "github.com/jinzhu/gorm"

type HistoricCollegePlayer struct {
	gorm.Model
	BasePlayer
	IsRedshirt         bool
	IsRedshirting      bool
	HasGraduated       bool
	HasProgressed      bool
	WillDeclare        bool
	TransferStatus     uint8  // 1 == Intends, 2 == Is Transferring
	TransferLikeliness string // Low, Medium, High
	LegacyID           uint
	Stats              []CollegePlayerStats     `gorm:"foreignKey:CollegePlayerID"`
	SeasonStats        CollegePlayerSeasonStats `gorm:"foreignKey:CollegePlayerID"`
}

func (h *HistoricCollegePlayer) Map(cp CollegePlayer) {
	h.ID = cp.ID
	h.BasePlayer = cp.BasePlayer
	h.PlayerID = cp.PlayerID
	h.TeamID = cp.TeamID
	h.Team = cp.Team
	h.State = cp.State
	h.Country = cp.Country
}

type UnsignedPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID           uint
	TeamID             uint
	TeamAbbr           string
	IsRedshirt         bool
	IsRedshirting      bool
	HasGraduated       bool
	HasProgressed      bool
	WillDeclare        bool
	TransferStatus     int    // 1 == Intends, 2 == Is Transferring
	TransferLikeliness string // Low, Medium, High
	LegacyID           uint
	Stats              []CollegePlayerStats     `gorm:"foreignKey:CollegePlayerID"`
	SeasonStats        CollegePlayerSeasonStats `gorm:"foreignKey:CollegePlayerID"`
	Profiles           []TransferPortalProfile  `gorm:"foreignKey:CollegePlayerID"`
}

func (up *UnsignedPlayer) MapFromRecruit(r Recruit) {
	up.ID = r.ID
	up.TeamID = 0
	up.TeamAbbr = ""
	up.PlayerID = r.PlayerID
	up.State = r.State
	up.Year = r.Age - 17
	up.IsRedshirt = false
	up.IsRedshirting = false
	up.HasGraduated = false
	up.Age = r.Age + 1
	up.FirstName = r.FirstName
	up.LastName = r.LastName
	up.Position = r.Position
	up.Archetype = r.Archetype
	up.Height = r.Height
	up.Stars = r.Stars
	up.Country = r.Country
	up.Overall = r.Overall
	up.InsideShooting = r.InsideShooting
	up.MidRangeShooting = r.MidRangeShooting
	up.ThreePointShooting = r.ThreePointShooting
	up.FreeThrow = r.FreeThrow
	up.Ballwork = r.Ballwork
	up.Rebounding = r.Rebounding
	up.InteriorDefense = r.InteriorDefense
	up.PerimeterDefense = r.PerimeterDefense
	up.SpecCount = r.SpecCount
	up.SpecInsideShooting = r.SpecInsideShooting
	up.SpecMidRangeShooting = r.SpecMidRangeShooting
	up.SpecThreePointShooting = r.SpecThreePointShooting
	up.SpecFreeThrow = r.SpecFreeThrow
	up.SpecBallwork = r.SpecBallwork
	up.SpecRebounding = r.SpecRebounding
	up.SpecInteriorDefense = r.SpecInteriorDefense
	up.SpecPerimeterDefense = r.SpecPerimeterDefense
	up.Stamina = r.Stamina
	up.PotentialGrade = r.PotentialGrade
	up.Potential = r.Potential
	up.FreeAgency = r.FreeAgency
	up.Personality = r.Personality
	up.RecruitingBias = r.RecruitingBias
	up.RecruitingBiasValue = r.RecruitingBiasValue
	up.WorkEthic = r.WorkEthic
	up.AcademicBias = r.AcademicBias
}
