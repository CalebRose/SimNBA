package managers

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetNBAUserByUsername(username string) structs.NBAUser {
	db := dbprovider.GetInstance().GetDB()

	var user structs.NBAUser

	err := db.Where("username = ?", username).Find(&user).Error
	if err != nil || user.ID == 0 {
		user = structs.NBAUser{
			Username:    username,
			TeamID:      0,
			TotalWins:   0,
			TotalLosses: 0,
			IsActive:    true,
		}
	}

	return user
}
