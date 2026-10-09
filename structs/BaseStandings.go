package structs

type BaseStandings struct {
	TotalWinPercentage float32
	ConfWinPercentage  float32
	TotalWins          int
	TotalLosses        int
	ConferenceWins     int
	ConferenceLosses   int
	RankedWins         int
	RankedLosses       int
	PointsFor          int
	PointsAgainst      int
	PointsDifferential int
	Streak             int
	HomeWins           int
	AwayWins           int
	Coach              string
}

func (b *BaseStandings) ResetStandings() {
	b.TotalWins = 0
	b.TotalLosses = 0
	b.ConferenceWins = 0
	b.ConferenceLosses = 0
	b.RankedWins = 0
	b.RankedLosses = 0
	b.PointsFor = 0
	b.PointsAgainst = 0
	b.PointsDifferential = 0
	b.Streak = 0
	b.HomeWins = 0
	b.AwayWins = 0
}

func (ns *BaseStandings) CalculatePercentages() {
	totalGames := ns.TotalWins + ns.TotalLosses
	totalConfGames := ns.ConferenceWins + ns.ConferenceLosses
	if totalGames > 0 {
		ns.TotalWinPercentage = float32(ns.TotalWins) / float32(totalGames)
	} else {
		ns.TotalWinPercentage = 0
	}
	if totalConfGames > 0 {
		ns.ConfWinPercentage = float32(ns.ConferenceWins) / float32(totalConfGames)
	} else {
		ns.ConfWinPercentage = 0
	}
}
