package structs

import "github.com/jinzhu/gorm"

type NBATradeOffer struct {
	gorm.Model
	TradeProposalID uint
	PlayerID        uint
	DraftPickID     uint
	TradeValue      uint
	TeamID          uint
	Team            string
	Player          NBAPlayer
	DraftPick       DraftPick
}
