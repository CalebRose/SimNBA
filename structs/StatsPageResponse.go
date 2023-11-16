package structs

type StatsPageResponse struct {
	CollegeConferences []CollegeConference
	CollegePlayers     []CollegePlayerResponse
	CollegeTeams       []CollegeTeamResponse
}

type NBAStatsPageResponse struct {
	NBAConferences []NBAConference
	NBAPlayers     []NBAPlayerResponse
	NBATeams       []NBATeamResponse
}
