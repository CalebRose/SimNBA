package structs

import "github.com/jinzhu/gorm"

type NBATradeProposal struct {
	gorm.Model
	TeamAID      uint
	TeamA        string
	TeamBID      uint
	TeamB        string
	TeamAOffers  bool
	TeamBOffers  bool
	TeamAConfirm bool
	TeamBConfirm bool
	TeamAReject  bool
	TeamBReject  bool
	AdminConfirm bool
	AdminReject  bool
}
