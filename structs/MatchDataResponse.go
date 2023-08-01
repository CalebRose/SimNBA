package structs

type MatchStateResponse struct {
	MatchType string
	Week      uint
	Matches   []MatchResponse
}

type MatchResponse struct {
	ID                     uint
	MatchName              string // For Post-Season matchups
	WeekID                 uint
	Week                   uint
	SeasonID               uint
	HomeTeamID             uint
	HomeTeam               string
	AwayTeamID             uint
	AwayTeam               string
	MatchOfWeek            string
	Arena                  string
	City                   string
	State                  string
	IsNeutralSite          bool
	IsNBAMatch             bool
	IsConference           bool
	IsConferenceTournament bool
	IsNITGame              bool
	IsPlayoffGame          bool
	IsNationalChampionship bool
	IsRivalryGame          bool
	IsInvitational         bool
	IsInternational        bool
	Channel                uint
}

type MatchDataResponse struct {
	HomeTeam         MatchTeamResponse
	HomeTeamRoster   []CollegePlayer
	HomeTeamGameplan Gameplan
	AwayTeam         MatchTeamResponse
	AwayTeamRoster   []CollegePlayer
	AwayTeamGameplan Gameplan
	GameID           uint
	Match            string
	WeekID           uint
	SeasonID         uint
	League           string
}

func (mdr *MatchDataResponse) AssignHomeTeam(team MatchTeamResponse, roster []CollegePlayer, gp Gameplan) {
	mdr.HomeTeam = team
	mdr.HomeTeamRoster = roster
	mdr.HomeTeamGameplan = gp
}

func (mdr *MatchDataResponse) AssignAwayTeam(team MatchTeamResponse, roster []CollegePlayer, gp Gameplan) {
	mdr.AwayTeam = team
	mdr.AwayTeamRoster = roster
	mdr.AwayTeamGameplan = gp
}

type MatchTeamResponse struct {
	ID           uint
	TeamName     string
	Mascot       string
	Abbr         string
	Conference   string
	Coach        string
	ConferenceID uint
	LeagueID     uint
}

func (mtr *MatchTeamResponse) Map(team Team) {
	mtr.ID = team.ID
	mtr.TeamName = team.Team
	mtr.Mascot = team.Nickname
	mtr.Abbr = team.Abbr
	mtr.Conference = team.Conference
	mtr.Coach = team.Coach
}
