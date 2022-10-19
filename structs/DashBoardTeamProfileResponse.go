package structs

type DashboardTeamProfileResponse struct {
	TeamProfile  TeamRecruitingProfile
	TeamNeedsMap map[string]int
}

func (d *DashboardTeamProfileResponse) SetTeamProfile(profile TeamRecruitingProfile) {
	d.TeamProfile = profile
}

func (d *DashboardTeamProfileResponse) SetTeamNeedsMap(obj map[string]int) {
	d.TeamNeedsMap = obj
}

type TeamBoardTeamProfileResponse struct {
	TeamProfile  SimTeamBoardResponse
	TeamNeedsMap map[string]int
}

func (t *TeamBoardTeamProfileResponse) SetTeamProfile(profile SimTeamBoardResponse) {
	t.TeamProfile = profile
}

func (t *TeamBoardTeamProfileResponse) SetTeamNeedsMap(obj map[string]int) {
	t.TeamNeedsMap = obj
}
