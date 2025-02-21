package managers

import "github.com/CalebRose/SimNBA/structs"

type BootstrapData struct {
	AllCollegeTeams      []structs.Team
	CollegeStandings     []structs.CollegeStandings
	CollegeRosterMap     map[uint][]structs.CollegePlayer
	Recruits             []structs.Croot
	TeamProfileMap       map[uint]structs.TeamRecruitingProfile
	PortalPlayers        []structs.CollegePlayer
	CollegeInjuryReport  []structs.CollegePlayer
	CollegeNews          []structs.NewsLog
	CollegeNotifications []structs.Notification
	AllCollegeGames      []structs.Match
	CollegeGameplan      structs.Gameplan

	// Player Profiles by Team?
	// Portal profiles?
	AllProTeams      []structs.NBATeam
	ProStandings     []structs.NBAStandings
	ProRosterMap     map[uint][]structs.NBAPlayer
	CapsheetMap      map[uint]structs.NBACapsheet
	FreeAgency       structs.FreeAgencyResponse
	ProInjuryReport  []structs.NBAPlayer
	ProNews          []structs.NewsLog
	ProNotifications []structs.Notification
	AllProGames      []structs.NBAMatch
	NBAGameplan      structs.NBAGameplan
}

func GetBootstrapData(collegeID, proID string) BootstrapData {

	return BootstrapData{}
}
