package managers

import (
	"fmt"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetAllNewsLogs(seasonID string) []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("season_id = ?", seasonID).Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}
