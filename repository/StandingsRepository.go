package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

type StandingsQuery struct {
	SeasonID string
	TeamID   string
}

func FindAllCollegeStandingsRecords(clauses StandingsQuery) []structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.CollegeStandings

	query := db.Model(&structs.CollegeStandings{})
	if clauses.SeasonID != "" {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if clauses.TeamID != "" {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	query.Find(&standings)

	return standings
}

func FindAllNBAStandingsRecords(clauses StandingsQuery) []structs.NBAStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.NBAStandings

	query := db.Model(&structs.NBAStandings{})
	if clauses.SeasonID != "" {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if clauses.TeamID != "" {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	query.Find(&standings)

	return standings
}

func CreateCollegeStandingsRecordsBatch(records []structs.CollegeStandings, db *gorm.DB, batchSize int) error {
	total := len(records)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(records[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBAStandingsRecordsBatch(records []structs.NBAStandings, db *gorm.DB, batchSize int) error {
	total := len(records)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(records[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
