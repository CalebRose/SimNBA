package repository

import (
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func CreateNBACombineRecordsBatch(db *gorm.DB, fds []structs.NBACombineResults, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}
