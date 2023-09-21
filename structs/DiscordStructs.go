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
