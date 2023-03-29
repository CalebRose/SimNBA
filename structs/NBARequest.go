package structs

import "github.com/jinzhu/gorm"

type NBARequest struct {
	gorm.Model
	Username            string
	NBATeamID           uint
	NBATeam             string
	NBATeamAbbreviation string
	IsOwner             bool
	IsManager           bool
	IsCoach             bool
	IsAssistant         bool
	IsApproved          bool
}

func (r *NBARequest) ApproveTeamRequest() {
	r.IsApproved = true
}

func (r *NBARequest) RejectTeamRequest() {
	r.IsApproved = false
}
