package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

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
