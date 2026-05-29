package structs

// CreateRecruitPointsDto - Data Transfer Object from UI to API
type UpdateGameplanDto struct {
	CollegeLineups []CollegeLineup
	NBALineups     []NBALineup
	CollegePlayers []CollegePlayerResponse
	NBAPlayers     []NBAPlayer
	TeamID         int
}
