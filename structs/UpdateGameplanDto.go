package structs

// CreateRecruitPointsDto - Data Transfer Object from UI to API
type UpdateGameplanDto struct {
	CollegeLineups          []CollegeLineup
	NBALineups              []NBALineup
	CollegePlayers          []CollegePlayerResponse
	NBAPlayers              []NBAPlayer
	TeamID                  int
	Pace                    string
	OffensiveFormation      string
	DefensiveFormation      string
	FocusPlayer             string
	TimeoutSettingsProvided bool
	PreserveTimeouts        bool
	Trigger1Enabled         bool
	Trigger1Type            uint8
	Trigger1Value           uint
	Trigger2Enabled         bool
	Trigger2Value           uint
	Trigger3Enabled         bool
	Trigger3Value           uint
	Trigger3Exhaustion      uint
	Trigger4Enabled         bool
	Trigger4Value           uint
}
