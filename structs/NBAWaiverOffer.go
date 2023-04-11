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
