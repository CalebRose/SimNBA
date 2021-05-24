package structs

// Request - A player request to sign for a team
type RequestDTO struct {
	ID         int
	TeamID     int
	Team       string
	Abbr       string
	Username   string
	Conference string
	IsNBA      bool
	IsApproved bool
}
