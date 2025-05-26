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

func SaveContractOfferRecord(contract structs.NBAContractOffer, db *gorm.DB) {
	err := db.Save(&contract).Error
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

func SaveCollegePlayerSeasonStatRecord(stats structs.CollegePlayerSeasonStats, db *gorm.DB) {
	err := db.Save(&stats).Error
	if err != nil {
		log.Panicln("Could not save cbb player season stats record")
	}
}

func SaveNBAPlayerSeasonStatRecord(stats structs.NBAPlayerSeasonStats, db *gorm.DB) {
	err := db.Save(&stats).Error
	if err != nil {
		log.Panicln("Could not save nba player season stats record")
	}
}

func SaveCollegeStandingsRecord(s structs.CollegeStandings, db *gorm.DB) {
	err := db.Save(&s).Error
	if err != nil {
		log.Panicln("Could not save nba player season stats record")
	}
}

func SaveNBAStandingsRecord(s structs.NBAStandings, db *gorm.DB) {
	err := db.Save(&s).Error
	if err != nil {
		log.Panicln("Could not save nba player season stats record")
	}
}

func SaveCollegeTeamSeasonStatRecord(stats structs.TeamSeasonStats, db *gorm.DB) {
	err := db.Save(&stats).Error
	if err != nil {
		log.Panicln("Could not save cbb player season stats record")
	}
}

func SaveNBATeamSeasonStatRecord(stats structs.NBATeamSeasonStats, db *gorm.DB) {
	err := db.Save(&stats).Error
	if err != nil {
		log.Panicln("Could not save nba player season stats record")
	}
}

func SaveNotification(noti structs.Notification, db *gorm.DB) {
	err := db.Save(&noti).Error
	if err != nil {
		log.Panicln("Could not save notification record!")
	}
}

func SaveCBBRecruit(recruit structs.Recruit, db *gorm.DB) {
	recruit.RecruitProfiles = nil
	err := db.Save(&recruit).Error
	if err != nil {
		log.Panicln("Could not save notification record!")
	}
}

func SaveCBBRecruitProfile(profile structs.PlayerRecruitProfile, db *gorm.DB) {
	profile.Recruit = structs.Recruit{}
	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save notification record!")
	}
}

func SaveCBBTeamRecruitingProfile(tp structs.TeamRecruitingProfile, db *gorm.DB) {
	tp.Recruits = nil
	err := db.Save(&tp).Error
	if err != nil {
		log.Panicln("Could not save notification record!")
	}
}

func SaveCollegePromiseRecord(promise structs.CollegePromise, db *gorm.DB) {
	// Save College Promise Record
	err := db.Save(&promise).Error
	if err != nil {
		log.Panicln("Could not save new college recruit record")
	}
}
