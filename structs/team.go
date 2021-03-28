package structs

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
)

// Team - The CBB / NBA team
type Team struct {
	gorm.Model
	Team        string
	Nickname    string
	Abbr        string
	City        string
	State       string
	Country     string
	Conference  string
	Division    string
	FirstSeason string
	Coach       string
	IsNBA       bool
}

// GetTeam - retrieve team
func (t *Team) GetTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Println(t)
}
