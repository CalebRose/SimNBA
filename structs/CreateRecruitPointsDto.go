package structs

// CreateRecruitProfileDto - Data Transfer Object from UI to API
type CreateRecruitProfileDto struct {
	ProfileId      int
	PlayerId       int
	SeasonId       int
	HasStateBonus  bool
	HasRegionBonus bool
	Team           string
}

type CreateRecruitProfileDtoV2 struct {
	ProfileID      uint
	PlayerID       uint
	RecruitID      uint
	SeasonID       uint
	HasStateBonus  bool
	HasRegionBonus bool
	Team           string
}
