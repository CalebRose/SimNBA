package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
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
	// Preload("RecruitingProfile").
	err := db.Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func RemoveUserFromTeam(teamId string) structs.Team {
	db := dbprovider.GetInstance().GetDB()

	team := GetTeamByTeamID(teamId)

	team.RemoveUser()

	standings := GetStandingsRecordByTeamID(teamId)

	standings.UpdateCoach("AI")

	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(teamId)

	recruitingProfile.ToggleAIBehavior(true)

	db.Save(&team)

	db.Save(&standings)

	db.Save(&recruitingProfile)

	return team
}

func GetTeamsInConference(db *gorm.DB, conference string) []structs.Team {
	var teams []structs.Team
	err := db.Where("conference = ?", conference).Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}

	return teams
}

func GetTeamRatings(t structs.Team) {
	db := dbprovider.GetInstance().GetDB()
	teamIDINT := int(t.ID)

	players := GetCollegePlayersByTeamId(strconv.Itoa(teamIDINT))

	offenseRating := 0
	defenseRating := 0
	overallRating := 0

	offenseSum := 0
	defenseSum := 0

	for _, player := range players {
		offenseSum += player.Shooting2 + player.Shooting3 + player.Finishing
		defenseSum += player.Ballwork + player.Rebounding + player.Defense
	}

	offenseRating = offenseSum / len(players)
	defenseRating = defenseSum / len(players)
	overallRating = (offenseRating + defenseRating) / 2

	offLetterGrade := util.GetOffenseGrade(offenseRating)
	defLetterGrade := util.GetDefenseGrade(defenseRating)
	ovrLetterGrade := util.GetOverallGrade(overallRating)

	t.AssignRatings(offLetterGrade, defLetterGrade, ovrLetterGrade)

	err := db.Save(&t).Error
	if err != nil {
		log.Fatalln("Could not save team rating for " + t.Abbr)
	}
}
