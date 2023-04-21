package structs

import "github.com/jinzhu/gorm"

type CollegePlayerSeasonStats struct {
	gorm.Model
	GamesPlayed               uint
	CollegePlayerID           uint
	SeasonID                  uint
	Year                      uint
	Minutes                   int
	MinutesPerGame            float64
	PossessionsPerGame        float64
	Possessions               int
	FGM                       int
	FGA                       int
	FGPercent                 float64
	FGMPG                     float64
	FGAPG                     float64
	ThreePointsMade           int
	ThreePointAttempts        int
	ThreePointPercent         float64
	ThreePointsMadePerGame    float64
	ThreePointAttemptsPerGame float64
	FTM                       int
	FTA                       int
	FTPercent                 float64
	FTMPG                     float64
	FTAPG                     float64
	Points                    int
	PPG                       float64
	TotalRebounds             int
	OffRebounds               int
	DefRebounds               int
	ReboundsPerGame           float64
	OffReboundsPerGame        float64
	DefReboundsPerGame        float64
	Assists                   int
	AssistsPerGame            float64
	Steals                    int
	StealsPerGame             float64
	Blocks                    int
	BlocksPerGame             float64
	Turnovers                 int
	TurnoversPerGame          float64
	Fouls                     int
	FoulsPerGame              float64
}

func (s *CollegePlayerSeasonStats) AddStatsToSeasonRecord(stat CollegePlayerStats) {
	if stat.Minutes > 0 {
		s.GamesPlayed++
	}
	s.CollegePlayerID = stat.CollegePlayerID
	s.SeasonID = stat.SeasonID
	s.Minutes += stat.Minutes
	s.Possessions += stat.Possessions
	s.FGM += stat.FGM
	s.FGA += stat.FGA
	if s.FGA > 0 {
		s.FGPercent = float64(s.FGM) / float64(s.FGA)
	}
	s.ThreePointsMade += stat.ThreePointsMade
	s.ThreePointAttempts += stat.ThreePointAttempts
	if s.ThreePointAttempts > 0 {
		s.ThreePointPercent = float64(s.ThreePointsMade) / float64(s.ThreePointAttempts)
	}
	s.FTM += stat.FTM
	s.FTA += stat.FTA
	if s.FTA > 0 {
		s.FTPercent = float64(s.FTM) / float64(s.FTA)
	}
	s.Points += stat.Points
	s.TotalRebounds += stat.TotalRebounds
	s.OffRebounds += stat.OffRebounds
	s.DefRebounds += stat.DefRebounds
	s.Assists += stat.Assists
	s.Steals += stat.Steals
	s.Blocks += stat.Blocks
	s.Turnovers += stat.Turnovers
	s.Fouls += stat.Fouls

	s.PPG = float64(s.Points) / float64(s.GamesPlayed)
	s.PossessionsPerGame = float64(s.Possessions) / float64(s.GamesPlayed)
	s.MinutesPerGame = float64(s.Minutes) / float64(s.GamesPlayed)
	s.FGMPG = float64(s.FGM) / float64(s.GamesPlayed)
	s.FGAPG = float64(s.FGA) / float64(s.GamesPlayed)
	s.ThreePointsMadePerGame = float64(s.ThreePointsMade) / float64(s.GamesPlayed)
	s.ThreePointAttemptsPerGame = float64(s.ThreePointAttempts) / float64(s.GamesPlayed)
	s.FTMPG = float64(s.FTM) / float64(s.GamesPlayed)
	s.FTAPG = float64(s.FTA) / float64(s.GamesPlayed)
	s.ReboundsPerGame = float64(s.TotalRebounds) / float64(s.GamesPlayed)
	s.OffReboundsPerGame = float64(s.OffRebounds) / float64(s.GamesPlayed)
	s.DefReboundsPerGame = float64(s.DefRebounds) / float64(s.GamesPlayed)
	s.AssistsPerGame = float64(s.Assists) / float64(s.GamesPlayed)
	s.StealsPerGame = float64(s.Steals) / float64(s.GamesPlayed)
	s.BlocksPerGame = float64(s.Blocks) / float64(s.GamesPlayed)
	s.TurnoversPerGame = float64(s.Turnovers) / float64(s.GamesPlayed)
	s.FoulsPerGame = float64(s.Fouls) / float64(s.GamesPlayed)
}

func (s *CollegePlayerSeasonStats) RemoveStatsToSeasonRecord(stat CollegePlayerStats) {
	if stat.Minutes > 0 {
		s.GamesPlayed--
	}
	s.CollegePlayerID = stat.CollegePlayerID
	s.SeasonID = stat.SeasonID
	s.Minutes -= stat.Minutes
	s.Possessions -= stat.Possessions
	s.FGM -= stat.FGM
	s.FGA -= stat.FGA
	if s.FGA > 0 {
		s.FGPercent = float64(s.FGM) / float64(s.FGA)
	}
	s.ThreePointsMade -= stat.ThreePointsMade
	s.ThreePointAttempts -= stat.ThreePointAttempts
	if s.ThreePointAttempts > 0 {
		s.ThreePointPercent = float64(s.ThreePointsMade) / float64(s.ThreePointAttempts)
	}
	s.FTM -= stat.FTM
	s.FTA -= stat.FTA
	if s.FTA > 0 {
		s.FTPercent = float64(s.FTM) / float64(s.FTA)
	}
	s.Points -= stat.Points
	s.TotalRebounds -= stat.TotalRebounds
	s.OffRebounds -= stat.OffRebounds
	s.DefRebounds -= stat.DefRebounds
	s.Assists -= stat.Assists
	s.Steals -= stat.Steals
	s.Blocks -= stat.Blocks
	s.Turnovers -= stat.Turnovers
	s.Fouls -= stat.Fouls

	s.PPG = float64(s.Points) / float64(s.GamesPlayed)
	s.PossessionsPerGame = float64(s.Possessions) / float64(s.GamesPlayed)
	s.MinutesPerGame = float64(s.Minutes) / float64(s.GamesPlayed)
	s.FGMPG = float64(s.FGM) / float64(s.GamesPlayed)
	s.FGAPG = float64(s.FGA) / float64(s.GamesPlayed)
	s.ThreePointsMadePerGame = float64(s.ThreePointsMade) / float64(s.GamesPlayed)
	s.ThreePointAttemptsPerGame = float64(s.ThreePointAttempts) / float64(s.GamesPlayed)
	s.FTMPG = float64(s.FTM) / float64(s.GamesPlayed)
	s.FTAPG = float64(s.FTA) / float64(s.GamesPlayed)
	s.ReboundsPerGame = float64(s.TotalRebounds) / float64(s.GamesPlayed)
	s.OffReboundsPerGame = float64(s.OffRebounds) / float64(s.GamesPlayed)
	s.DefReboundsPerGame = float64(s.DefRebounds) / float64(s.GamesPlayed)
	s.AssistsPerGame = float64(s.Assists) / float64(s.GamesPlayed)
	s.StealsPerGame = float64(s.Steals) / float64(s.GamesPlayed)
	s.BlocksPerGame = float64(s.Blocks) / float64(s.GamesPlayed)
	s.TurnoversPerGame = float64(s.Turnovers) / float64(s.GamesPlayed)
	s.FoulsPerGame = float64(s.Fouls) / float64(s.GamesPlayed)
}
