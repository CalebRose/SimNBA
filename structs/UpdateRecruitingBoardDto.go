package structs

// UpdateRecruitingBoardDto - Data Transfer Object from UI to API
type UpdateRecruitingBoardDto struct {
	Profile  TeamRecruitingProfile
	Recruits []RecruitPlayerProfile
	TeamID   int
}
