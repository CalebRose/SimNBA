package repository

import (
	"log"

	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func CreateRecruitRecord(croot structs.Recruit, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&croot).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateRecruitRecordsBatch(db *gorm.DB, fds []structs.Recruit, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateGlobalRecordsBatch(db *gorm.DB, fds []structs.GlobalPlayer, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateProfessionalPlayerRecord(player structs.NBAPlayer, db *gorm.DB) {
	// Save NBA Player Record
	err := db.Create(&player).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateGlobalPlayerRecord(player structs.GlobalPlayer, db *gorm.DB) {
	// Save College Player Record
	err := db.Create(&player).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}

func CreateCollegePlayerRecordsBatch(db *gorm.DB, fds []structs.CollegePlayer, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateHistoricCollegePlayerRecordsBatch(db *gorm.DB, fds []structs.HistoricCollegePlayer, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBADrafteesRecordsBatch(db *gorm.DB, fds []structs.NBADraftee, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBAPlayerRecordsBatch(db *gorm.DB, fds []structs.NBAPlayer, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNBARetiredPlayerRecordsBatch(db *gorm.DB, fds []structs.RetiredPlayer, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
