package repository

import (
	"github.com/CalebRose/SimNBA/dbprovider"
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

func FindNBACombineRecords() []structs.NBACombineResults {
	db := dbprovider.GetInstance().GetDB()
	var records []structs.NBACombineResults
	db.Find(&records)
	return records
}
