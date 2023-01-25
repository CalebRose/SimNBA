package structs

import "github.com/jinzhu/gorm"

type NBATradeBlock struct {
	gorm.Model
	TeamID             uint
	Team               string
	IsLookingForOffers bool
	TradeProposals     []NBATradeProposal
	NBAPlayers         []NBAPlayer
}
