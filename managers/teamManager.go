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

func GetAllActiveCollegeTeamsWithSeasonStats() []structs.Team {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.Team

	err := db.Preload("TeamSeasonStats").Where("is_active = ? and is_nba = ?", true, false).
		Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

func GetAllActiveNBATeamsWithSeasonStats() []structs.NBATeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam

	err := db.Preload("TeamSeasonStats").Where("is_active = ?", true).
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

func GetNBATeamByTeamID(teamId string) structs.NBATeam {
	var team structs.NBATeam
	db := dbprovider.GetInstance().GetDB()
	// Preload("RecruitingProfile").
	err := db.Preload("Capsheet").Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func RemoveUserFromTeam(teamId string) structs.Team {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	team := GetTeamByTeamID(teamId)

	team.RemoveUser()

	standings := GetStandingsRecordByTeamID(teamId, strconv.Itoa(int(ts.SeasonID)))

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

	for idx, player := range players {
		if idx > 9 {
			break
		}
		offenseSum += player.Shooting2 + player.Shooting3 + player.Finishing
		defenseSum += player.Ballwork + player.Rebounding + player.Defense
	}

	offenseRating = offenseSum / 9
	defenseRating = defenseSum / 9
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

func GetNBATeamRatings(t structs.NBATeam) {
	db := dbprovider.GetInstance().GetDB()
	teamIDINT := int(t.ID)

	players := GetNBAPlayersWithContractsByTeamID(strconv.Itoa(teamIDINT))

	offenseRating := 0
	defenseRating := 0
	overallRating := 0

	offenseSum := 0
	defenseSum := 0

	for idx, player := range players {
		if idx > 9 {
			break
		}
		offenseSum += player.Shooting2 + player.Shooting3 + player.Finishing + player.FreeThrow
		defenseSum += player.Ballwork + player.Rebounding + player.InteriorDefense + player.PerimeterDefense
	}

	offenseRating = offenseSum / 9
	defenseRating = defenseSum / 9
	overallRating = (offenseRating + defenseRating) / 2

	offLetterGrade := util.GetNBATeamGrade(offenseRating)
	defLetterGrade := util.GetNBATeamGrade(defenseRating)
	ovrLetterGrade := util.GetNBATeamGrade(overallRating)

	t.AssignRatings(offLetterGrade, defLetterGrade, ovrLetterGrade)

	err := db.Save(&t).Error
	if err != nil {
		log.Fatalln("Could not save team rating for " + t.Abbr)
	}
}

func GetCBBTeamByAbbreviation(abbr string) structs.Team {
	var team structs.Team
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("abbr = ?", abbr).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func GetOnlyNBATeams() []structs.NBATeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam

	err := db.Where("league_id = 1").Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

func GetAllActiveNBATeams() []structs.NBATeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.NBATeam

	err := db.Find(&teams).Error
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

// GetTeamByTeamID - straightforward
func GetNBATeamWithCapsheetByTeamID(teamId string) structs.NBATeam {
	var team structs.NBATeam
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("Capsheet").Where("id = ?", teamId).Find(&team).Error
	if err != nil {
		log.Fatal(err)
	}
	return team
}
