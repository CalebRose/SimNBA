package structs

import "github.com/jinzhu/gorm"

type RetiredPlayer struct {
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
}
