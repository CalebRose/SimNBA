package structs

// CreateRecruitProfileDto - Data Transfer Object from UI to API
type CreateRecruitProfileDto struct {
	ProfileId int
	PlayerId  int
	SeasonId  int
	Team      string
}
