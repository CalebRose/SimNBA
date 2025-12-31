package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func FindAllTransferPortalProfiles() []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	db.Where("removed_from_board = ?", false).Find(&profiles)

	return profiles
}
