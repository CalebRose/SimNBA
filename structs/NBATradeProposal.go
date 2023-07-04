package structs

import "github.com/jinzhu/gorm"

type NBATradeProposal struct {
	gorm.Model
	NBATeamID                 uint
	NBATeam                   string
	RecepientTeamID           uint
	RecepientTeam             string
	IsTradeAccepted           bool
	IsTradeRejected           bool
	IsSynced                  bool
	NBATeamTradeOptions       []NBATradeOption `gorm:"foreignKey:TradeProposalID"`
	RecepientTeamTradeOptions []NBATradeOption `gorm:"foreignKey:TradeProposalID"`
}

func (p *NBATradeProposal) ToggleSyncStatus() {
	p.IsSynced = true
}

func (p *NBATradeProposal) AssignID(id uint) {
	p.ID = id
}

func (p *NBATradeProposal) AcceptTrade() {
	p.IsTradeAccepted = true
}

func (p *NBATradeProposal) RejectTrade() {
	p.IsTradeRejected = true
}

type NBATradeOption struct {
	gorm.Model
	TradeProposalID  uint
	NBATeamID        uint
	NBAPlayerID      uint
	NBADraftPickID   uint
	OptionType       string
	CashTransfer     float64
	SalaryPercentage float64 // Will be a percentage that the recepient team (TEAM B) will pay for Y1. Will be between 0 and 100.
}

type NBATradeOptionObj struct {
	ID               uint
	TradeProposalID  uint
	NBATeamID        uint
	NBAPlayerID      uint
	NBADraftPickID   uint
	OptionType       string
	CashTransfer     float64
	SalaryPercentage float64   // Will be a percentage that the recepient team (TEAM B) will pay. Will be between 0 and 100.
	Player           NBAPlayer // If the NBAPlayerID is greater than 0, it will return a player.
	Draftpick        DraftPick // If the NBADraftPickID is greater than 0, it will return a draft pick.
}

func (to *NBATradeOptionObj) AssignPlayer(player NBAPlayer) {
	to.Player = player
	to.NBAPlayerID = player.ID
}

func (to *NBATradeOptionObj) AssignPick(pick DraftPick) {
	to.Draftpick = pick
	to.NBADraftPickID = pick.ID
}

type NBATradeProposalDTO struct {
	ID                        uint
	NBATeamID                 uint
	NBATeam                   string
	RecepientTeamID           uint
	RecepientTeam             string
	IsTradeAccepted           bool
	IsTradeRejected           bool
	NBATeamTradeOptions       []NBATradeOptionObj
	RecepientTeamTradeOptions []NBATradeOptionObj
}

type NBATeamProposals struct {
	SentTradeProposals     []NBATradeProposalDTO
	ReceivedTradeProposals []NBATradeProposalDTO
}
