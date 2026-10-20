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
	t.UpdateLatestInstance()
}

func (t *Team) AssignCoach(username string) {
	t.Coach = username
}

func (t *Team) RemoveUser() {
	t.Coach = ""
	t.IsUserCoached = false
	t.UpdateLatestInstance()
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

type TeamGrade struct {
	OffenseGradeNumber float64
	DefenseGradeNumber float64
	OverallGradeNumber float64
	OffenseGradeLetter string
	DefenseGradeLetter string
	OverallGradeLetter string
}

func (tg *TeamGrade) SetOffenseGradeNumber(grade float64) {
	tg.OffenseGradeNumber = grade
}

func (tg *TeamGrade) SetDefenseGradeNumber(grade float64) {
	tg.DefenseGradeNumber = grade
}

func (tg *TeamGrade) SetOverallGradeNumber(grade float64) {
	tg.OverallGradeNumber = grade
}

func (tg *TeamGrade) SetOffenseGradeLetter(grade string) {
	tg.OffenseGradeLetter = grade
}

func (tg *TeamGrade) SetDefenseGradeLetter(grade string) {
	tg.DefenseGradeLetter = grade
}

func (tg *TeamGrade) SetOverallGradeLetter(grade string) {
	tg.OverallGradeLetter = grade
}
