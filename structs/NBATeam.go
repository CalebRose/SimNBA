package structs

import "github.com/jinzhu/gorm"

type NBATeam struct {
	gorm.Model
	Team              string
	Nickname          string
	Abbr              string
	City              string
	State             string
	Country           string
	LeagueID          uint
	League            string
	ConferenceID      uint
	Conference        string
	DivisionID        uint
	Division          string
	ArenaID           uint
	Arena             string
	NBAOwnerID        uint
	NBAOwnerName      string
	NBACoachID        uint
	NBACoachName      string
	NBAGMID           uint
	NBAGMName         string
	NBAAssistantID    uint
	NBAAssistantName  string
	OverallGrade      string
	OffenseGrade      string
	DefenseGrade      string
	IsActive          bool
	WaiverOrder       uint
	Gameplan          NBAGameplan           `gorm:"foreignKey:TeamID"`
	TeamStats         []TeamStats           `gorm:"foreignKey:TeamID"`
	TeamSeasonStats   TeamSeasonStats       `gorm:"foreignKey:TeamID"`
	RecruitingProfile TeamRecruitingProfile `gorm:"foreignKey:TeamID"`
	Capsheet          NBACapsheet           `gorm:"foreignKey:TeamID"`
	Contracts         []NBAContract         `gorm:"foreignKey:TeamID"`
	DraftPicks        []DraftPick           `gorm:"foreignKey:TeamID"`
}

func (t *NBATeam) AssignID(id uint) {
	t.ID = id
}

func (t *NBATeam) AssignNBAUserToTeam(r NBARequest, u NBAUser) {
	if r.IsOwner {
		t.NBAOwnerID = u.ID
		t.NBAOwnerName = r.Username
	}
	if r.IsManager {
		t.NBAGMID = u.ID
		t.NBAGMName = r.Username
	}
	if r.IsCoach {
		t.NBACoachID = u.ID
		t.NBACoachName = r.Username
	}
	if r.IsAssistant {
		t.NBAAssistantID = u.ID
		t.NBAAssistantName = r.Username
	}
}

func (t *NBATeam) RemoveUser(username string) {
	if t.NBAOwnerName == username {
		t.NBAOwnerName = ""
	} else if t.NBAGMName == username {
		t.NBAGMName = ""
	} else if t.NBACoachName == username {
		t.NBACoachName = ""
	} else {
		t.NBAAssistantName = ""
	}
}

func (t *NBATeam) AssignRatings(og string, dg string, ov string) {
	t.OffenseGrade = og
	t.DefenseGrade = dg
	t.OverallGrade = ov
}

func (t *NBATeam) AssignWaiverOrder(order uint) {
	t.WaiverOrder = order
}
