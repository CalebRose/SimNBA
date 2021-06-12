package managers

import (
	"log"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/jinzhu/gorm"
)

// TEAM Functions
func GetTeamByTeamID(teamId string) structs.Team {
	var team structs.Team
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("RecruitingProfile").Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func RemoveUserFromTeam(team structs.Team) {
	db := dbprovider.GetInstance().GetDB()
	db.Save(&team)
}

// PLAYER Functions
func GetPlayerByPlayerId(playerId string) structs.Player {
	// Test
	var player structs.Player
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("id = ?", playerId).Find(&player).Error
	if err != nil {
		log.Fatal(err)
	}
	return player
}

func UpdatePlayer(p structs.Player) {
	db := dbprovider.GetInstance().GetDB()
	err := db.Save(&p).Error
	if err != nil {
		log.Fatal(err)
	}
}

func GetRecruitingProfileByTeamId(teamId string) structs.RecruitingProfile {
	var profile structs.RecruitingProfile
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("Recruits.Recruit.RecruitingPoints", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points_spent DESC")
	}).Where("id = ?", teamId).Find(&profile).Error
	if err != nil {
		log.Fatal(err)
	}
	return profile
}

func GetRecruitingPointsProfileByPlayerId(db *gorm.DB, playerId string, profileId string) structs.RecruitingPoints {
	var recruitingPoints structs.RecruitingPoints
	err := db.Where("player_id = ? AND profile_id = ?", playerId, profileId).Find(&recruitingPoints).Error
	if err != nil {
		log.Fatal(err)
	}
	return recruitingPoints
}

func GetTimestamp(db *gorm.DB) structs.Timestamp {
	var timeStamp structs.Timestamp
	err := db.Find(&timeStamp).Error
	if err != nil {
		log.Fatal(err)
	}
	return timeStamp
}

func GetPlayersByConference(db *gorm.DB, seasonId string, conference string) []structs.Player {
	var players []structs.Player
	db.Preload("PlayerStats", "season_id = ?", seasonId).Joins("Team").Where("Team.Conference = ?", conference).Find(&players)
	return players
}

func GetTeamsInConference(db *gorm.DB, conference string) []structs.Team {
	var teams []structs.Team
	err := db.Where("conference = ?", conference).Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}

	return teams
}

func GetGameplansByTeam(db *gorm.DB, teamId string) []structs.Gameplan {
	var gameplans []structs.Gameplan
	db.Where("team_id = ?", teamId).Order("game asc").Find(&gameplans)

	return gameplans
}

func GetPlayersByTeamId(db *gorm.DB, teamId string) []structs.Player {
	var players []structs.Player
	db.Where("team_id = ?", teamId).Find(&players)

	return players
}
