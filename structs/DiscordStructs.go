package structs

type CollegeTeamResponseData struct {
	TeamData        Team
	TeamStandings   CollegeStandings
	UpcomingMatches []Match
}

type NBATeamResponseData struct {
	TeamData        NBATeam
	TeamStandings   NBAStandings
	UpcomingMatches []NBAMatch
}

type FlexComparisonModel struct {
	TeamOneID      uint
	TeamOne        string
	TeamOneWins    uint
	TeamOneLosses  uint
	TeamOneStreak  uint
	TeamOneMSeason int
	TeamOneMScore  string
	TeamTwoID      uint
	TeamTwo        string
	TeamTwoWins    uint
	TeamTwoLosses  uint
	TeamTwoStreak  uint
	TeamTwoMSeason int
	TeamTwoMScore  string
	CurrentStreak  uint
	LatestWin      string
}
