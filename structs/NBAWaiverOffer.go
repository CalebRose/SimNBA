package structs

import "gorm.io/gorm"

type NBAWaiverOfferDTO struct {
	ID          uint
	PlayerID    uint
	TeamID      uint
	SeasonID    uint
	Team        string
	IsAccepted  bool
	IsRejected  bool
	WaiverOrder uint
}

type NBAWaiverOffer struct {
	gorm.Model
	PlayerID    uint
	TeamID      uint
	SeasonID    uint
	Team        string
	IsAccepted  bool
	IsRejected  bool
	IsActive    bool
	WaiverOrder uint
}

func (n *NBAWaiverOffer) AssignID(id uint) {
	n.ID = id
}

func (wo *NBAWaiverOffer) Map(offer NBAWaiverOfferDTO) {
	wo.TeamID = offer.TeamID
	wo.Team = offer.Team
	wo.PlayerID = offer.PlayerID
	wo.WaiverOrder = offer.WaiverOrder
	wo.IsActive = true
}

func (wo *NBAWaiverOffer) AssignNewWaiverOrder(val uint) {
	wo.WaiverOrder = val
}

func (wo *NBAWaiverOffer) DeactivateWaiverOffer() {
	wo.IsActive = false
}
