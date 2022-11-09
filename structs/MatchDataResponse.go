package structs

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
	ID         uint
	TeamName   string
	Mascot     string
	Abbr       string
	Conference string
	Coach      string
}

func (mtr *MatchTeamResponse) Map(team Team) {
	mtr.ID = team.ID
	mtr.TeamName = team.Team
	mtr.Mascot = team.Nickname
	mtr.Abbr = team.Abbr
	mtr.Conference = team.Conference
	mtr.Coach = team.Coach
}
