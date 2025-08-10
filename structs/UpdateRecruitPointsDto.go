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

type UpdateRecruitPointsDtoV2 struct {
	RecruitPointsID   int
	ProfileID         int
	PlayerID          int
	RecruitID         int
	SpentPoints       int
	Team              string
	RewardScholarship bool
	RevokeScholarship bool
}
