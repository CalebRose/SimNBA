package managers

import (
	"fmt"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetAllCBBNewsLogs(seasonID string) []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("season_id = ? AND league = ?", seasonID, "CBB").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func GetAllNBANewsLogs(seasonID string) []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("season_id = ? AND league = ?", seasonID, "NBA").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}
