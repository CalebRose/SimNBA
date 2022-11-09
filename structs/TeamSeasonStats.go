package structs

import "github.com/jinzhu/gorm"

type TeamSeasonStats struct {
	gorm.Model
	TeamID                    uint
	SeasonID                  uint
	GamesPlayed               uint
	Points                    int
	PointsAgainst             int
	PPG                       float64
	PAPG                      float64
	Possessions               int
	PossessionsPerGame        float64
	FGM                       int
	FGA                       int
	FGPercent                 float64
	FGMPG                     float64
	FGAPG                     float64
	FGMAgainst                int
	FGAAgainst                int
	FGPercentAgainst          float64
	FGMAPG                    float64
	FGAAPG                    float64
	ThreePointsMade           int
	ThreePointAttempts        int
	ThreePointPercent         float64
	ThreePointsMadeAgainst    int
	ThreePointAttemptsAgainst int
	ThreePointPercentAgainst  float64
	TPMPG                     float64
	TPAPG                     float64
	TPMAPG                    float64
	TPAAPG                    float64
	FTM                       int
	FTA                       int
	FTPercent                 float64
	FTMPG                     float64
	FTAPG                     float64
	FTMAgainst                int
	FTAAgainst                int
	FTMAPG                    float64
	FTAAPG                    float64
	FTPercentAgainst          float64
	Rebounds                  int
	OffRebounds               int
	DefRebounds               int
	ReboundsPerGame           float64
	OffReboundsPerGame        float64
	DefReboundsPerGame        float64
	ReboundsAllowed           int
	OffReboundsAllowed        int
	DefReboundsAllowed        int
	ReboundsAllowedPerGame    float64
	OffReboundsAllowedPerGame float64
	DefReboundsAllowedPerGame float64
	Assists                   int
	AssistsAllowed            int
	AssistsPerGame            float64
	AssistsAllowedPerGame     float64
	Steals                    int
	StealsAllowed             int
	StealsPerGame             float64
	StealsAllowedPerGame      float64
	Blocks                    int
	BlocksAllowed             int
	BlocksPerGame             float64
	BlocksAllowedPerGame      float64
	TotalTurnovers            int
	TurnoversAllowed          int
	TurnoversPerGame          float64
	TurnoversAllowedPerGame   float64
	Fouls                     int
	FoulsPerGame              float64
}

func (s *TeamSeasonStats) AddStatsToSeasonRecord(stat TeamStats) {
	s.GamesPlayed++
	s.Possessions += stat.Possessions
	s.FGM += stat.FGM
	s.FGA += stat.FGA
	if s.FGA > 0 {
		s.FGPercent = float64(s.FGM / s.FGA)
	}
	s.ThreePointsMade += stat.ThreePointsMade
	s.ThreePointAttempts += stat.ThreePointAttempts
	if s.ThreePointAttempts > 0 {
		s.ThreePointPercent = float64(s.ThreePointsMade / s.ThreePointAttempts)
	}
	s.FTM += stat.FTM
	s.FTA += stat.FTA
	if s.FTA > 0 {
		s.FTPercent = float64(s.FTM / s.FTA)
	}
	s.Points += stat.Points
	s.Rebounds += stat.Rebounds
	s.OffRebounds += stat.OffRebounds
	s.DefRebounds += stat.DefRebounds
	s.Assists += stat.Assists
	s.Steals += stat.Steals
	s.Blocks += stat.Blocks
	s.TotalTurnovers += stat.TotalTurnovers
	s.Fouls += stat.Fouls

	s.PPG = float64(s.Points) / float64(s.GamesPlayed)
	s.PossessionsPerGame = float64(s.Possessions) / float64(s.GamesPlayed)
	s.FGMPG = float64(s.FGM) / float64(s.GamesPlayed)
	s.FGAPG = float64(s.FGA) / float64(s.GamesPlayed)
	s.TPMPG = float64(s.ThreePointsMade) / float64(s.GamesPlayed)
	s.TPAPG = float64(s.ThreePointAttempts) / float64(s.GamesPlayed)
	s.FTMPG = float64(s.FTM) / float64(s.GamesPlayed)
	s.FTAPG = float64(s.FTA) / float64(s.GamesPlayed)
	s.ReboundsPerGame = float64(s.Rebounds) / float64(s.GamesPlayed)
	s.OffReboundsPerGame = float64(s.OffRebounds) / float64(s.GamesPlayed)
	s.DefReboundsPerGame = float64(s.DefRebounds) / float64(s.GamesPlayed)
	s.AssistsPerGame = float64(s.Assists) / float64(s.GamesPlayed)
	s.StealsPerGame = float64(s.Steals) / float64(s.GamesPlayed)
	s.BlocksPerGame = float64(s.Blocks) / float64(s.GamesPlayed)
	s.TurnoversPerGame = float64(s.TotalTurnovers) / float64(s.GamesPlayed)
	s.FoulsPerGame = float64(s.Fouls) / float64(s.GamesPlayed)
	s.PAPG = float64(s.Points) / float64(s.GamesPlayed)
	s.FGMAPG = float64(s.FGM) / float64(s.GamesPlayed)
	s.FGAAPG = float64(s.FGA) / float64(s.GamesPlayed)
	s.FGPercentAgainst = float64(s.FGMAgainst) / float64(s.FGAAgainst)
	s.TPMAPG = float64(s.ThreePointsMade) / float64(s.GamesPlayed)
	s.TPAAPG = float64(s.ThreePointAttempts) / float64(s.GamesPlayed)
	s.ThreePointPercentAgainst = float64(s.ThreePointsMadeAgainst) / float64(s.ThreePointAttemptsAgainst)
	s.FTMAPG = float64(s.FTM) / float64(s.GamesPlayed)
	s.FTAAPG = float64(s.FTA) / float64(s.GamesPlayed)
	s.FTPercentAgainst = float64(s.FTMAgainst) / float64(s.FTAAgainst)
	s.ReboundsAllowedPerGame = float64(s.Rebounds) / float64(s.GamesPlayed)
	s.OffReboundsAllowedPerGame = float64(s.OffRebounds) / float64(s.GamesPlayed)
	s.DefReboundsAllowedPerGame = float64(s.DefRebounds) / float64(s.GamesPlayed)
	s.AssistsAllowedPerGame = float64(s.Assists) / float64(s.GamesPlayed)
	s.StealsAllowedPerGame = float64(s.Steals) / float64(s.GamesPlayed)
	s.BlocksAllowedPerGame = float64(s.Blocks) / float64(s.GamesPlayed)
	s.TurnoversAllowedPerGame = float64(s.TotalTurnovers) / float64(s.GamesPlayed)
	s.FoulsPerGame = float64(s.Fouls) / float64(s.GamesPlayed)
}
