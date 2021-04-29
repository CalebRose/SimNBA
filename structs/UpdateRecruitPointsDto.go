package structs

// CreateRecruitPointsDto - Data Transfer Object from UI to API
type UpdateRecruitPointsDto struct {
	RecruitPointsId   int
	ProfileId         string
	PlayerId          string
	SpentPoints       int
	RewardScholarship bool
	RevokeScholarship bool
}
