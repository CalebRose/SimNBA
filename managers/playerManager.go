package managers

import (
	"log"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"github.com/jinzhu/gorm"
)

func GetAllPlayers() []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player

	db.Find(&players)

	return players
}

func GetTeamRosterForRosterPage(teamId string) []structs.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer

	var responseList []structs.CollegePlayerResponse
	err := db.Order("team_id asc").Where("team_id = ?", teamId).Find(&players).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	for _, p := range players {
		shooting2Grade := util.GetAttributeGrade(p.Shooting2)
		shooting3Grade := util.GetAttributeGrade(p.Shooting3)
		finishingGrade := util.GetAttributeGrade(p.Finishing)
		reboundingGrade := util.GetAttributeGrade(p.Rebounding)
		ballworkGrade := util.GetAttributeGrade(p.Ballwork)
		defenseGrade := util.GetAttributeGrade(p.Defense)
		potentialGrade := util.GetPotentialGrade(p.Potential)
		overallGrade := util.GetPlayerOverallGrade(p.Overall)

		res := structs.CollegePlayerResponse{
			FirstName:            p.FirstName,
			LastName:             p.LastName,
			Position:             p.Position,
			Age:                  p.Age,
			Year:                 p.Year,
			State:                p.State,
			Country:              p.Country,
			Stars:                p.Stars,
			Height:               p.Height,
			PotentialGrade:       potentialGrade,
			Shooting2Grade:       shooting2Grade,
			Shooting3Grade:       shooting3Grade,
			FinishingGrade:       finishingGrade,
			BallworkGrade:        ballworkGrade,
			ReboundingGrade:      reboundingGrade,
			DefenseGrade:         defenseGrade,
			OverallGrade:         overallGrade,
			Stamina:              p.Stamina,
			PlaytimeExpectations: p.PlaytimeExpectations,
			Minutes:              p.Minutes,
			Potential:            p.Potential,
			Personality:          p.Personality,
			RecruitingBias:       p.RecruitingBias,
			WorkEthic:            p.WorkEthic,
			AcademicBias:         p.AcademicBias,
			PlayerID:             p.PlayerID,
			TeamID:               p.TeamID,
			TeamAbbr:             p.TeamAbbr,
		}

		responseList = append(responseList, res)
	}

	return responseList
}

func GetCollegePlayersByTeamId(teamId string) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer
	err := db.Order("team_id asc").Where("team_id = ?", teamId).Find(&players).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	return players
}

func GetAllCollegePlayers() []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer

	db.Find(&players)

	return players
}

func GetAllCollegePlayersFromOldTable() []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player

	db.Where("is_nba = ?", false).Find(&players)

	return players
}

func GetAllCollegeRecruits() []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Player
	db.Preload("RecruitingPoints", "total_points_spent > ?", 0, func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points_spent DESC")
	}).Where("is_nba = ? AND team_id = 0", false).Find(&recruits)

	return recruits
}

func GetAllJUCOCollegeRecruits() []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Player

	db.Where("is_nba = ? AND team_id = 0 AND year > 0", false).Find(&recruits)

	return recruits
}

func GetPlayersByConference(seasonId string, conference string) []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player

	db.Preload("PlayerStats", "season_id = ?", seasonId).Joins("Team").Where("Team.Conference = ?", conference).Find(&players)

	return players
}

func GetAllNBAPlayers() []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player

	db.Where("is_nba = ?", true).Find(&players)

	return players
}

func GetAllNBAFreeAgents() []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player
	db.Where("is_nba = ? AND team_id is null", true).Find(&players)

	return players
}

func GetPlayerByPlayerId(playerId string) structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var player structs.Player

	err := db.Where("id = ?", playerId).Find(&player).Error
	if err != nil {
		log.Fatal(err)
	}

	return player
}

func SetRedshirtStatusForPlayer(playerId string) structs.Player {
	player := GetPlayerByPlayerId(playerId)

	player.SetRedshirtingStatus()

	UpdatePlayer(player)

	return player
}

func UpdatePlayer(p structs.Player) {
	db := dbprovider.GetInstance().GetDB()
	err := db.Save(&p).Error
	if err != nil {
		log.Fatal(err)
	}
}

func CreateNewPlayer(firstName string, lastName string) {
	db := dbprovider.GetInstance().GetDB()

	player := &structs.Player{FirstName: firstName, LastName: lastName,
		Position: "C", Year: 4, State: "WA", Country: "USA",
		Stars: 3, Height: "7'0", TeamID: 10, Shooting: 14,
		Finishing: 20, Ballwork: 18, Rebounding: 20, Defense: 19,
		PotentialGrade: 20, Stamina: 36, PlaytimeExpectations: 25,
		MinutesA: 35, Overall: 20, IsNBA: false,
		IsRedshirt: false, IsRedshirting: false}

	db.Create(&player)
}
