package structs

type StatsPageResponse struct {
	CollegeConferences []CollegeConference
	CollegePlayers     []CollegePlayerResponse
	CollegeTeams       []CollegeTeamResponse
}
