package structs

import "github.com/jinzhu/gorm"

type HistoricCollegePlayer struct {
	gorm.Model
	BasePlayer
	PlayerID      uint
	TeamID        uint
	TeamAbbr      string
	IsRedshirt    bool
	IsRedshirting bool
	HasGraduated  bool
	HasProgressed bool
	WillDeclare   bool
	Stats         []CollegePlayerStats     `gorm:"foreignKey:CollegePlayerID"`
	SeasonStats   CollegePlayerSeasonStats `gorm:"foreignKey:CollegePlayerID"`
}

func (h *HistoricCollegePlayer) Map(cp CollegePlayer) {
	h.BasePlayer = cp.BasePlayer
	h.PlayerID = cp.PlayerID
	h.TeamID = cp.TeamID
	h.TeamAbbr = cp.TeamAbbr
	h.State = cp.State
	h.Country = cp.Country
}
