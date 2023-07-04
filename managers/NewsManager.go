package managers

import (
	"fmt"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetAllCBBNewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ?", "CBB").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func GetAllNBANewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ?", "NBA").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func CreateNewsLog(league, message, messageType string, teamID int, ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	seasonID := 0
	weekID := 0
	week := 0
	if league == "CFB" {
		seasonID = int(ts.SeasonID)
		weekID = int(ts.CollegeWeekID)
		week = ts.CollegeWeek
	} else {
		seasonID = int(ts.SeasonID)
		weekID = int(ts.NBAWeekID)
		week = ts.NBAWeek
	}

	news := structs.NewsLog{
		League:      league,
		Message:     message,
		MessageType: messageType,
		SeasonID:    uint(seasonID),
		WeekID:      uint(weekID),
		Week:        uint(week),
		TeamID:      uint(teamID),
	}

	db.Create(&news)
}
