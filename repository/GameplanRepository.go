package repository

import (
	"log"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

type GameplanQuery struct {
	TeamID string
}

func FindCollegeLineupRecords(clauses GameplanQuery) []structs.CollegeLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups []structs.CollegeLineup

	query := db.Model(&structs.CollegeLineup{})

	if clauses.TeamID != "" {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	query.Find(&lineups)

	return lineups
}

func FindNBALineupRecords(clauses GameplanQuery) []structs.NBALineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups []structs.NBALineup

	query := db.Model(&structs.NBALineup{})

	if clauses.TeamID != "" {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	query.Find(&lineups)

	return lineups
}

func CreateCollegeLineupsRecordsBatch(db *gorm.DB, fds []structs.CollegeLineup, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBALineupsRecordsBatch(db *gorm.DB, fds []structs.NBALineup, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

// Saves

// Saves CollegeLineup
func SaveCollegeLineupRecord(lineup structs.CollegeLineup, db *gorm.DB) {
	err := db.Save(&lineup).Error
	if err != nil {
		log.Panicln("Could not save gameplan record!")
	}
}

// Saves NBALineup
func SaveNBALineupRecord(lineup structs.NBALineup, db *gorm.DB) {
	err := db.Save(&lineup).Error
	if err != nil {
		log.Panicln("Could not save gameplan record!")
	}
}

func SaveNBAGameplanRecord(team structs.NBAGameplan, db *gorm.DB) {
	err := db.Save(&team).Error
	if err != nil {
		log.Panicln("Could not save gameplan record!")
	}
}
