package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

type SchedulerQuery struct {
	ID       string
	TeamID   string
	SeasonID string
	WeekID   string
}

func FindCBBGameRequestRecord(q SchedulerQuery) structs.CBBGameRequest {
	db := dbprovider.GetInstance().GetDB()
	var request structs.CBBGameRequest
	query := db.Model(&structs.CBBGameRequest{})
	if q.ID != "" {
		query = query.Where("id = ?", q.ID)
	}
	query.First(&request)
	return request
}

func FindCBBGameRequestRecords(q SchedulerQuery) []structs.CBBGameRequest {
	db := dbprovider.GetInstance().GetDB()
	var requests []structs.CBBGameRequest
	query := db.Model(&structs.CBBGameRequest{})
	if q.TeamID != "" {
		query = query.Where("home_team_id = ? OR away_team_id = ?", q.TeamID, q.TeamID)
	}
	if q.SeasonID != "" {
		query = query.Where("season_id = ?", q.SeasonID)
	}
	if q.WeekID != "" {
		query = query.Where("week_id = ?", q.WeekID)
	}
	query.Find(&requests)
	return requests
}

func CreateCBBGameRequest(request structs.CBBGameRequest, db *gorm.DB) {
	db.Create(&request)
}

func SaveCBBGameRequest(request structs.CBBGameRequest, db *gorm.DB) {
	db.Save(&request)
}

func DeleteCBBGameRequest(request structs.CBBGameRequest, db *gorm.DB) {
	db.Delete(&request)
}

func FindArenaByID(id uint) structs.Arena {
	db := dbprovider.GetInstance().GetDB()
	var arena structs.Arena
	db.First(&arena, id)
	return arena
}
