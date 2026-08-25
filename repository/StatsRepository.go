package repository

import (
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func CreateCollegeTeamStatsBatch(db *gorm.DB, fds []structs.TeamStats, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBATeamStatsBatch(db *gorm.DB, fds []structs.NBATeamStats, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegePlayerStatsBatch(db *gorm.DB, fds []structs.CollegePlayerStats, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBAPlayerStatsBatch(db *gorm.DB, fds []structs.NBAPlayerStats, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegePlayByPlayBatch(db *gorm.DB, fds []structs.CollegePlayByPlay, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBAPlayByPlayBatch(db *gorm.DB, fds []structs.NBAPlayByPlay, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
