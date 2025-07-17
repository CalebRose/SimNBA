package structs

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// Team - The CBB / NBA team
type Team struct {
	gorm.Model
	Team               string
	Nickname           string
	Abbr               string
	City               string
	State              string
	Country            string
	ConferenceID       uint
	Conference         string
	Arena              string
	FirstSeason        string
	Coach              string
	OverallGrade       string
	OffenseGrade       string
	DefenseGrade       string
	IsNBA              bool
	IsActive           bool
	IsUserCoached      bool
	ColorOne           string
	ColorTwo           string
	ColorThree         string
	DiscordID          string
	PenaltyMark        uint8
	NoLongerActiveUser bool
	LastLogin          time.Time
	Gameplan           []Gameplan            `gorm:"foreignKey:TeamID"`
	TeamStats          []TeamStats           `gorm:"foreignKey:TeamID"`
	TeamSeasonStats    TeamSeasonStats       `gorm:"foreignKey:TeamID"`
	RecruitingProfile  TeamRecruitingProfile `gorm:"foreignKey:TeamID"`
}

func (bt *Team) UpdateLatestInstance() {
	bt.LastLogin = time.Now()
}

func (t *Team) AssignDiscordID(id string) {
	t.DiscordID = id
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

func (t *Team) MarkTeamPenalty() {
	t.PenaltyMark++
}

func (t *Team) ResetPenaltyMark() {
	t.PenaltyMark = 0
}

func (t *Team) CheckUserActivity() {
	if t.PenaltyMark >= 3 {
		t.NoLongerActiveUser = true
	}
	if t.LastLogin.Add(4 * 7 * 24 * time.Hour).Before(time.Now()) {
		t.NoLongerActiveUser = true
	}
}
