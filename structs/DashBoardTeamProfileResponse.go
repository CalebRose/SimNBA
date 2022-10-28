package structs

type DashboardTeamProfileResponse struct {
	TeamProfile TeamRecruitingProfile
}

func (d *DashboardTeamProfileResponse) SetTeamProfile(profile TeamRecruitingProfile) {
	d.TeamProfile = profile
}

type TeamBoardTeamProfileResponse struct {
	TeamProfile SimTeamBoardResponse
}

func (t *TeamBoardTeamProfileResponse) SetTeamProfile(profile SimTeamBoardResponse) {
	t.TeamProfile = profile
}
