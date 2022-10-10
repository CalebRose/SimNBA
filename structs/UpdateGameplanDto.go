package structs

// CreateRecruitPointsDto - Data Transfer Object from UI to API
type UpdateGameplanDto struct {
	Gameplan Gameplan
	Players  []CollegePlayer
	TeamID   int
}
