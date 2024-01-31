package structs

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
)

// Team - The CBB / NBA team
type Team struct {
	gorm.Model
	Team              string
	Nickname          string
	Abbr              string
	City              string
	State             string
	Country           string
	ConferenceID      uint
	Conference        string
	Arena             string
	FirstSeason       string
	Coach             string
	OverallGrade      string
	OffenseGrade      string
	DefenseGrade      string
	IsNBA             bool
	IsActive          bool
	IsUserCoached     bool
	Gameplan          []Gameplan            `gorm:"foreignKey:TeamID"`
	TeamStats         []TeamStats           `gorm:"foreignKey:TeamID"`
	TeamSeasonStats   TeamSeasonStats       `gorm:"foreignKey:TeamID"`
	RecruitingProfile TeamRecruitingProfile `gorm:"foreignKey:TeamID"`
}

// GetTeam - retrieve team
func (t *Team) GetTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Println(t)
}

func (t *Team) AssignUserToTeam(username string) {
	t.AssignCoach(username)
	t.IsUserCoached = true
}

func (t *Team) AssignCoach(username string) {
	t.Coach = username
}

func (t *Team) RemoveUser() {
	t.Coach = ""
	t.IsUserCoached = false
}

func (t *Team) ToggleUserCoach() {
	t.IsUserCoached = !t.IsUserCoached
}

func (t *Team) AssignRatings(off string, def string, ovr string) {
	t.OffenseGrade = off
	t.DefenseGrade = def
	t.OverallGrade = ovr
}
