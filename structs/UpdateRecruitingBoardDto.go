package structs

// UpdateRecruitingBoardDto - Data Transfer Object from UI to API
type UpdateRecruitingBoardDto struct {
	Profile  RecruitingProfile
	Recruits []RecruitingPoints
	TeamID   int
}
