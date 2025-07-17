package structs

import (
	"time"

	"github.com/jinzhu/gorm"
)

type NBATeam struct {
	gorm.Model
	Team               string
	Nickname           string
	Abbr               string
	City               string
	State              string
	Country            string
	LeagueID           uint
	League             string
	ConferenceID       uint
	Conference         string
	DivisionID         uint
	Division           string
	ArenaID            uint
	Arena              string
	NBAOwnerID         uint
	NBAOwnerName       string
	NBACoachID         uint
	NBACoachName       string
	NBAGMID            uint
	NBAGMName          string
	NBAAssistantID     uint
	NBAAssistantName   string
	OverallGrade       string
	OffenseGrade       string
	DefenseGrade       string
	IsActive           bool
	CanTrade           bool
	WaiverOrder        uint
	ColorOne           string
	ColorTwo           string
	ColorThree         string
	OwnerDiscordID     string
	GMDiscordID        string
	CoachDiscordID     string
	AssistantDiscordID string
	LastLogin          time.Time
	PenaltyMark        uint8
	NoLongerActiveUser bool
	Gameplan           NBAGameplan        `gorm:"foreignKey:TeamID"`
	TeamStats          []NBATeamStats     `gorm:"foreignKey:TeamID"`
	TeamSeasonStats    NBATeamSeasonStats `gorm:"foreignKey:TeamID"`
	Capsheet           NBACapsheet        `gorm:"foreignKey:TeamID"`
	Contracts          []NBAContract      `gorm:"foreignKey:TeamID"`
	DraftPicks         []DraftPick        `gorm:"foreignKey:TeamID"`
}

func (bt *NBATeam) UpdateLatestInstance() {
	bt.LastLogin = time.Now()
}

func (t *NBATeam) AssignID(id uint) {
	t.ID = id
}

func (t *NBATeam) AssignDiscordID(id, username string) {
	if t.NBAOwnerName == username {
		t.OwnerDiscordID = id
	} else if t.NBAGMName == username {
		t.GMDiscordID = id
	} else if t.NBACoachName == username {
		t.CoachDiscordID = id
	} else if t.NBAAssistantName == username {
		t.AssistantDiscordID = id
	}
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
		t.NBAOwnerID = 0
	} else if t.NBAGMName == username {
		t.NBAGMName = ""
		t.NBAGMID = 0
	} else if t.NBACoachName == username {
		t.NBACoachName = ""
		t.NBACoachID = 0
	} else {
		t.NBAAssistantName = ""
		t.NBAAssistantID = 0
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

func (t *NBATeam) ActivateTradeAbility() {
	t.CanTrade = true
}

func (t *NBATeam) MarkTeamPenalty() {
	t.PenaltyMark++
}

func (t *NBATeam) ResetPenaltyMark() {
	t.PenaltyMark = 0
}

func (t *NBATeam) CheckUserActivity() {
	if t.PenaltyMark >= 3 {
		t.NoLongerActiveUser = true
	}
	if t.LastLogin.Add(4 * 7 * 24 * time.Hour).Before(time.Now()) {
		t.NoLongerActiveUser = true
	}
}

type ISLTeamNeeds struct {
	TeamNeedsMap  map[string]bool
	PositionCount map[string]int
	TotalCount    int
}

func (itn *ISLTeamNeeds) IncrementTotalCount() {
	itn.TotalCount += 1
}

func (itn *ISLTeamNeeds) IncrementPositionCount(pos string) {
	itn.PositionCount[pos] += 1
}
