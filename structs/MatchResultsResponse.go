package structs

type MatchResultsResponse struct {
	HomePlayers []CollegePlayer
	AwayPlayers []CollegePlayer
	HomeStats   TeamStats
	AwayStats   TeamStats
}
