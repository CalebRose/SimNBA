package structs

import "github.com/jinzhu/gorm"

type CollegePlayerStats struct {
	gorm.Model
	CollegePlayerID    uint
	MatchID            uint
	SeasonID           uint
	WeekID             uint
	Week               uint
	MatchType          string
	Year               uint
	Minutes            int
	Possessions        int
	FGM                int
	FGA                int
	FGPercent          float64
	ThreePointsMade    int
	ThreePointAttempts int
	ThreePointPercent  float64
	FTM                int
	FTA                int
	FTPercent          float64
	Points             int
	TotalRebounds      int
	OffRebounds        int
	DefRebounds        int
	Assists            int
	Steals             int
	Blocks             int
	Turnovers          int
	Fouls              int
	FouledOut          bool
	IsInjured          bool
	InjuryName         string
	InjuryType         string
	WeeksOfRecovery    uint8
	RevealResults      bool
}

func (s *CollegePlayerStats) MapNewProperties(weekID uint, matchType string) {
	s.WeekID = weekID
	s.MatchType = matchType
}
