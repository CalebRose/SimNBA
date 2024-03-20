package repository

import (
	"log"

	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func DeleteCollegePlayerRecord(player structs.CollegePlayer, db *gorm.DB) {
	err := db.Delete(&player).Error
	if err != nil {
		log.Panicln("Could not delete old college player record.")
	}
}

func DeleteCollegeRecruitRecord(player structs.Recruit, db *gorm.DB) {
	err := db.Delete(&player).Error
	if err != nil {
		log.Panicln("Could not delete old college recruit record.")
	}
}

func DeleteCollegePromise(promise structs.CollegePromise, db *gorm.DB) {
	err := db.Delete(&promise).Error
	if err != nil {
		log.Panicln("Could not delete old college player record.")
	}
}
