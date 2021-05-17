package managers

import (
	"fmt"

	"github.com/CalebRose/SimNBA/config"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var c = config.Config()

// TEAM Functions
func GetTeamByTeamID(teamId string) structs.Team {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var team structs.Team
	db.Preload("RecruitingProfile").Where("id = ?", teamId).Find(&team)

	return team
}

func RemoveUserFromTeam(team structs.Team) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	db.Model(&team).Update("coach", nil)
}

// PLAYER Functions
func GetPlayerByPlayerId(playerId string) structs.Player {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()
	// Test
	var player structs.Player
	db.Where("id = ?", playerId).Find(&player)

	return player
}

func UpdatePlayer(p structs.Player) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	db.Save(&p)
}

func GetRecruitingProfileByTeamId(db *gorm.DB, teamId string) structs.RecruitingProfile {
	var profile structs.RecruitingProfile
	db.Preload("Recruits").Where("id = ?", teamId).Find(&profile)
	return profile
}

func GetRecruitingPointsProfileByPlayerId(db *gorm.DB, playerId string, profileId string) structs.RecruitingPoints {
	var recruitingPoints structs.RecruitingPoints
	db.Where("player_id = ? AND profile_id = ?", playerId, profileId).Find(&recruitingPoints)

	return recruitingPoints
}

func GetTimestamp(db *gorm.DB) structs.Timestamp {
	var timeStamp structs.Timestamp
	db.Find(&timeStamp)
	return timeStamp
}

func GetPlayersByConference(db *gorm.DB, seasonId string, conference string) []structs.Player {
	var players []structs.Player
	db.Preload("PlayerStats", "season_id = ?", seasonId).Joins("Team").Where("Team.Conference = ?", conference).Find(&players)
	return players
}

func GetTeamsInConference(db *gorm.DB, conference string) []structs.Team {
	var teams []structs.Team
	db.Where("conference = ?", conference).Find(&teams)

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
