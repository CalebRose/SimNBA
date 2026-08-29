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

type SearchStatsResponse struct {
	CBBPlayerGameStats   []CollegePlayerStats
	CBBPlayerSeasonStats []CollegePlayerSeasonStats
	CBBTeamGameStats     []TeamStats
	CBBTeamSeasonStats   []TeamSeasonStats
	NBAPlayerGameStats   []NBAPlayerStats
	NBAPlayerSeasonStats []NBAPlayerSeasonStats
	NBATeamGameStats     []NBATeamStats
	NBATeamSeasonStats   []NBATeamSeasonStats
}
