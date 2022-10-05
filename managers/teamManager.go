package managers

import (
	"log"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/jinzhu/gorm"
)

func GetAllActiveCollegeTeams() []structs.Team {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team

	err := db.Where("is_active = ? and is_nba = ?", true, false).
		Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

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

func GetTeamsInConference(db *gorm.DB, conference string) []structs.Team {
	var teams []structs.Team
	err := db.Where("conference = ?", conference).Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}

	return teams
}
