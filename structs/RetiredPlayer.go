package structs

import "github.com/jinzhu/gorm"

type RetiredPlayer struct {
	gorm.Model
	BasePlayer
	PlayerID            uint
	TeamID              uint
	TeamAbbr            string
	CollegeID           uint
	College             string
	DraftPickID         uint
	DraftedRound        uint
	DraftPick           uint
	DraftedTeamID       uint
	DraftedTeamAbbr     string
	PrimeAge            uint
	IsNBA               bool
	MaxRequested        bool
	IsSuperMaxQualified bool
	IsFreeAgent         bool
	IsGLeague           bool
	IsTwoWay            bool
	IsWaived            bool
	IsOnTradeBlock      bool
	IsFirstTeamANBA     bool
	IsDPOY              bool
	IsMVP               bool
	IsInternational     bool
	IsRetiring          bool
	IsAcceptingOffers   bool
	IsNegotiating       bool
	MinimumValue        float64
	SigningRound        uint
	NegotiationRound    uint
	Offers              []NBAContractOffer   `gorm:"foreignKey:PlayerID"`
	WaiverOffers        []NBAWaiverOffer     `gorm:"foreignKey:PlayerID"`
	Contract            NBAContract          `gorm:"foreignKey:PlayerID"`
	Stats               []NBAPlayerStats     `gorm:"foreignKey:NBAPlayerID"`
	SeasonStats         NBAPlayerSeasonStats `gorm:"foreignKey:NBAPlayerID"`
}
