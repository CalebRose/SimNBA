package structs

type BaseStandings struct {
	TotalWins        int
	TotalLosses      int
	ConferenceWins   int
	ConferenceLosses int
	RankedWins       int
	RankedLosses     int
	PointsFor        int
	PointsAgainst    int
	Streak           int
	HomeWins         int
	AwayWins         int
	Coach            string
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
	b.Streak = 0
	b.HomeWins = 0
	b.AwayWins = 0
}
