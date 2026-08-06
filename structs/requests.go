package structs

import "gorm.io/gorm"

// Request - A player request to sign for a team
type Request struct {
	gorm.Model
	TeamID                  uint
	Username                string
	IsApproved              bool
	DiscordUsername         string
	HowMuchTimeAnswer       string
	HowDidYouHearAboutSimSN string
	CommunityReference      string
	AboutYourself           string
}

func (r *Request) ApproveTeamRequest() {
	r.IsApproved = true
}

func (r *Request) RejectTeamRequest() {
	r.IsApproved = false
}
