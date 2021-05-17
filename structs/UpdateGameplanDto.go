package structs

// CreateRecruitPointsDto - Data Transfer Object from UI to API
type UpdateGameplanDto struct {
	Gameplans []Gameplan
	Players   []Player
	TeamID    int
}
