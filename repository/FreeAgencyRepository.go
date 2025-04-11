package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

type FreeAgencyQuery struct {
	PlayerID string
	IsActive bool
	TeamID   string
	OfferID  string
}

func FindAllProContracts(onlyActive bool) []structs.NBAContract {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAContract{}

	query := db.Model(&offers)

	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.NBAContract{}
	}

	return offers
}

func FindAllProExtensions(onlyActive bool) []structs.NBAExtensionOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAExtensionOffer{}

	query := db.Model(&offers)

	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.NBAExtensionOffer{}
	}

	return offers
}

func FindAllFreeAgentOffers(clauses FreeAgencyQuery) []structs.NBAContractOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAContractOffer{}
	query := db.Model(&offers)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}

	if len(clauses.OfferID) > 0 {
		query = query.Where("id = ?", clauses.OfferID)
	}

	if clauses.IsActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.NBAContractOffer{}
	}

	return offers
}

func FindAllWaiverOffers(clauses FreeAgencyQuery) []structs.NBAWaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAWaiverOffer{}
	query := db.Model(&offers)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}

	if len(clauses.OfferID) > 0 {
		query = query.Where("id = ?", clauses.OfferID)
	}

	if clauses.IsActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&offers).Error; err != nil {
		return []structs.NBAWaiverOffer{}
	}

	return offers
}
