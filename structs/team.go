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
	Conference        string
	Division          string
	FirstSeason       string
	Coach             string
	IsNBA             bool
	Players           []Player          `gorm:"foriegnKey:TeamID"`
	Gameplan          []Gameplan        `gorm:"foreignKey:TeamID"`
	TeamStats         TeamStats         `gorm:"foreignKey:TeamID"`
	RecruitingProfile RecruitingProfile `gorm:"foreignKey:TeamID"`
}

// GetTeam - retrieve team
func (t *Team) GetTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Println(t)
}

func (t *Team) RemoveUser() {
	t.Coach = ""
}
