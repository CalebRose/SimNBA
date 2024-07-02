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
	Capacity               int
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

type CBBMatchDataResponse struct {
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

func (mdr *CBBMatchDataResponse) AssignHomeTeam(team MatchTeamResponse, roster []CollegePlayer, gp Gameplan) {
	mdr.HomeTeam = team
	mdr.HomeTeamRoster = roster
	mdr.HomeTeamGameplan = gp
}

func (mdr *CBBMatchDataResponse) AssignAwayTeam(team MatchTeamResponse, roster []CollegePlayer, gp Gameplan) {
	mdr.AwayTeam = team
	mdr.AwayTeamRoster = roster
	mdr.AwayTeamGameplan = gp
}

type NBAMatchDataResponse struct {
	HomeTeam         MatchTeamResponse
	HomeTeamRoster   []NBAPlayer
	HomeTeamGameplan NBAGameplan
	AwayTeam         MatchTeamResponse
	AwayTeamRoster   []NBAPlayer
	AwayTeamGameplan NBAGameplan
	GameID           uint
	Match            string
	WeekID           uint
	SeasonID         uint
	League           string
}

func (mdr *NBAMatchDataResponse) AssignHomeTeam(team MatchTeamResponse, roster []NBAPlayer, gp NBAGameplan) {
	mdr.HomeTeam = team
	mdr.HomeTeamRoster = roster
	mdr.HomeTeamGameplan = gp
}

func (mdr *NBAMatchDataResponse) AssignAwayTeam(team MatchTeamResponse, roster []NBAPlayer, gp NBAGameplan) {
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

func (mtr *MatchTeamResponse) MapNBATeam(team NBATeam) {
	mtr.ID = team.ID
	mtr.TeamName = team.Team
	mtr.Mascot = team.Nickname
	mtr.Abbr = team.Abbr
	mtr.Conference = team.Conference
	if len(team.NBACoachName) > 0 {
		mtr.Coach = team.NBACoachName
	} else {
		mtr.Coach = team.NBAOwnerName
	}

}

type PollDataResponse struct {
	Poll          CollegePollSubmission
	Matches       []Match
	Standings     []CollegeStandings
	OfficialPolls []CollegePollOfficial
}
