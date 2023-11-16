package structs

import "gorm.io/gorm"

type NBAUser struct {
	gorm.Model
	Username         string
	TeamID           uint
	TeamAbbreviation string
	IsOwner          bool
	IsManager        bool
	IsHeadCoach      bool
	IsAssistant      bool
	TotalWins        uint
	TotalLosses      uint
	TotalTies        uint
	IsActive         bool
}

func (u *NBAUser) SetTeam(r NBARequest) {
	u.TeamID = r.NBATeamID
	u.TeamAbbreviation = r.NBATeamAbbreviation
	if r.IsOwner {
		u.IsOwner = true
	}
	if r.IsManager {
		u.IsManager = true
	}
	if r.IsCoach {
		u.IsHeadCoach = true
	}
	if r.IsAssistant {
		u.IsAssistant = true
	}
}

func (u *NBAUser) RemoveOwnership() {
	u.IsOwner = false

	if !u.IsHeadCoach && !u.IsManager && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NBAUser) RemoveManagerPosition() {
	u.IsManager = false

	if !u.IsHeadCoach && !u.IsOwner && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NBAUser) RemoveCoachPosition() {
	u.IsHeadCoach = false

	if !u.IsManager && !u.IsOwner && !u.IsAssistant {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}

func (u *NBAUser) RemoveAssistantPosition() {
	u.IsHeadCoach = false

	if !u.IsManager && !u.IsOwner && !u.IsHeadCoach {
		u.TeamID = 0
		u.TeamAbbreviation = ""
	}
}
