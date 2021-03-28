package structs

// Request - A player request to sign for a team
type Request struct {
	ID         int
	TeamID     int
	Username   string
	IsApproved bool
}
