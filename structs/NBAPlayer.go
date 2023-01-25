package structs

import "github.com/jinzhu/gorm"

type NBAPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID             uint
	TeamID               uint
	TeamAbbr             string
	CollegeID            uint
	College              string
	DraftPickID          uint
	DraftPick            uint
	DraftedTeamID        uint
	DraftedTeamAbbr      string
	PrimeAge             uint
	IsNBA                bool
	MaxRequested         bool
	IsSuperMaxQualified  bool
	IsFreeAgent          bool
	IsGLeague            bool
	IsTwoWay             bool
	IsWaived             bool
	IsOnTradeBlock       bool
	IsFirstTeamANBA      bool
	IsDPOY               bool
	IsMVP                bool
	IsInternational      bool
	IsRetiring           bool
	PositionOne          string
	PositionTwo          string
	PositionThree        string
	Position1Minutes     uint
	Position2Minutes     uint
	Position3Minutes     uint
	InsidePercentage     uint
	MidPercentage        uint
	ThreePointPercentage uint
	// Contracts           []NBAContract
	// NBAPlayerStats
	// NBASeasonStats
}

func (n *NBAPlayer) SetID(id uint) {
	n.ID = id
}

func (n *NBAPlayer) SetRetiringStatus() {
	n.IsRetiring = true
}

func (n *NBAPlayer) BecomeFreeAgent() {
	n.TeamAbbr = "FA"
	n.TeamID = 0
}

func (n *NBAPlayer) SignWithTeam(teamID uint, team string) {
	n.TeamAbbr = team
	n.TeamID = teamID
	n.IsFreeAgent = false
}

func (n *NBAPlayer) Progress(p NBAPlayerProgressions) {
	n.Shooting2 = p.Shooting2
	n.Shooting3 = p.Shooting3
	n.FreeThrow = p.FreeThrow
	n.Ballwork = p.Ballwork
	n.Finishing = p.Finishing
	n.Rebounding = p.Rebounding
	n.InteriorDefense = p.InteriorDefense
	n.PerimeterDefense = p.PerimeterDefense
	n.Overall = p.Overall
	n.Age = p.Age
	n.Stamina = p.Stamina
	if n.Stamina < 1 {
		n.Stamina = 1
	}
	n.Year++
}
func (n *NBAPlayer) QualifyForSuperMax() {
	n.IsSuperMaxQualified = true
}

func (n *NBAPlayer) QualifiesForMax() {
	n.MaxRequested = true
}
