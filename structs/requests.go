package structs

import "github.com/jinzhu/gorm"

// Request - A player request to sign for a team
type Request struct {
	gorm.Model
	TeamID     int
	Username   string
	IsApproved bool
}

func (r *Request) ApproveTeamRequest() {
	r.IsApproved = true
}

func (r *Request) RejectTeamRequest() {
	r.IsApproved = false
}
