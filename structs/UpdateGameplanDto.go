package structs

// CreateRecruitPointsDto - Data Transfer Object from UI to API
type UpdateGameplanDto struct {
	Gameplan       Gameplan
	CollegePlayers []CollegePlayer
	NBAPlayers     []NBAPlayer
	TeamID         int
}
