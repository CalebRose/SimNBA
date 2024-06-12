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

func SaveProfessionalContractRecord(contract structs.NBAContract, db *gorm.DB) {
	err := db.Save(&contract).Error
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

func SaveCBBTeamRecruitingProfile(profile *structs.TeamRecruitingProfile, db *gorm.DB) {
	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save team profile record")
	}
}

func SaveISLScoutingDeptRecord(dept structs.ISLScoutingDept, db *gorm.DB) {
	if dept.ID == 0 {
		log.Panicln("ID is not set for the scouting dept record")
	}
	err := db.Save(&dept).Error
	if err != nil {
		log.Panicln("Could not save scouting dept record")
	}
}

func SaveISLScoutingReportRecord(report structs.ISLScoutingReport, db *gorm.DB) {
	err := db.Save(&report).Error
	if err != nil {
		log.Panicln("Could not save scouting report record")
	}
}
