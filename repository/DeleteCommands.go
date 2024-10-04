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
func DeleteProfessionalPlayerRecord(player structs.NBAPlayer, db *gorm.DB) {
	player.Offers = nil
	player.WaiverOffers = nil
	player.Extensions = nil
	player.Contract = structs.NBAContract{}
	player.Stats = nil
	player.SeasonStats = structs.NBAPlayerSeasonStats{}
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

func DeleteContract(contract structs.NBAContract, db *gorm.DB) {
	err := db.Delete(&contract).Error
	if err != nil {
		log.Panicln("Could not delete old college player record.")
	}
}

func DeleteExtension(contract structs.NBAExtensionOffer, db *gorm.DB) {
	err := db.Delete(&contract).Error
	if err != nil {
		log.Panicln("Could not delete old college player record.")
	}
}

func DeleteNotificationRecord(noti structs.Notification, db *gorm.DB) {
	err := db.Delete(&noti).Error
	if err != nil {
		log.Panicln("Could not delete old notification record.")
	}
}
