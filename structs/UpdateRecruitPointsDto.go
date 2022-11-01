package structs

// CreateRecruitPointsDto - Data Transfer Object from UI to API
type UpdateRecruitPointsDto struct {
	RecruitPointsId   int
	ProfileId         int
	PlayerId          int
	SpentPoints       int
	Team              string
	RewardScholarship bool
	RevokeScholarship bool
}
