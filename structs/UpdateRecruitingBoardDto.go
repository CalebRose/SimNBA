package structs

// UpdateRecruitingBoardDto - Data Transfer Object from UI to API
type UpdateRecruitingBoardDto struct {
	Profile  SimTeamBoardResponse
	Recruits []CrootProfile
	TeamID   int
}
