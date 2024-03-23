package repository

import (
	"log"

	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func SaveTimeStamp(ts structs.Timestamp, db *gorm.DB) {
	err := db.Save(&ts).Error
	if err != nil {
		log.Panicln("Could not save timestamp")
	}
}

func SaveCollegePlayerRecord(player structs.CollegePlayer, db *gorm.DB) {
	player.Stats = nil
	player.SeasonStats = structs.CollegePlayerSeasonStats{}
	player.Profiles = nil
	err := db.Save(&player).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveProfessionalPlayerRecord(player structs.NBAPlayer, db *gorm.DB) {
	player.Stats = nil
	player.SeasonStats = structs.NBAPlayerSeasonStats{}
	player.Contract = structs.NBAContract{}
	player.Offers = nil
	player.WaiverOffers = nil
	player.Extensions = nil

	err := db.Save(&player).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveTransferPortalProfile(profile structs.TransferPortalProfile, db *gorm.DB) {
	profile.CollegePlayer = structs.CollegePlayer{}
	profile.Promise = structs.CollegePromise{}

	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}

func SaveProfessionalMatchRecord(match structs.NBAMatch, db *gorm.DB) {
	err := db.Save(&match).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}
