package managers

import (
	"fmt"
	"log"
	"sort"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
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
		freeThrowGrade := util.GetAttributeGrade(p.FreeThrow)
		finishingGrade := util.GetAttributeGrade(p.Finishing)
		reboundingGrade := util.GetAttributeGrade(p.Rebounding)
		ballworkGrade := util.GetAttributeGrade(p.Ballwork)
		interiorDefenseGrade := util.GetAttributeGrade(p.InteriorDefense)
		perimeterDefenseGrade := util.GetAttributeGrade(p.PerimeterDefense)
		potentialGrade := util.GetPotentialGrade(p.Potential)
		overallGrade := util.GetPlayerOverallGrade(p.Overall)

		res := structs.CollegePlayerResponse{
			FirstName:             p.FirstName,
			LastName:              p.LastName,
			Position:              p.Position,
			Age:                   p.Age,
			Year:                  p.Year,
			State:                 p.State,
			Country:               p.Country,
			Stars:                 p.Stars,
			Height:                p.Height,
			PotentialGrade:        potentialGrade,
			Shooting2Grade:        shooting2Grade,
			Shooting3Grade:        shooting3Grade,
			FreeThrowGrade:        freeThrowGrade,
			FinishingGrade:        finishingGrade,
			BallworkGrade:         ballworkGrade,
			ReboundingGrade:       reboundingGrade,
			InteriorDefenseGrade:  interiorDefenseGrade,
			PerimeterDefenseGrade: perimeterDefenseGrade,
			OverallGrade:          overallGrade,
			Stamina:               p.Stamina,
			PlaytimeExpectations:  p.PlaytimeExpectations,
			Minutes:               p.Minutes,
			Potential:             p.Potential,
			Personality:           p.Personality,
			RecruitingBias:        p.RecruitingBias,
			WorkEthic:             p.WorkEthic,
			AcademicBias:          p.AcademicBias,
			PlayerID:              p.PlayerID,
			TeamID:                p.TeamID,
			TeamAbbr:              p.TeamAbbr,
			IsRedshirting:         p.IsRedshirting,
			IsRedshirt:            p.IsRedshirt,
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

func GetCollegePlayersWithMatchStatsByTeamId(teamId string, matchID string) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer
	err := db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("match_id = ?", matchID)
	}).Where("team_id = ?", teamId).Find(&players).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}
	sort.Sort(structs.ByPlayedMinutes(players))
	return players
}

func GetCollegePlayerByPlayerID(playerID string) structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.CollegePlayer
	err := db.Where("id = ?", playerID).Find(&player).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	return player
}

func GetAllCollegePlayers() []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer

	db.Find(&players)

	return players
}

func GetAllCollegePlayersWithSeasonStats() []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer

	db.Preload("SeasonStats").Find(&players)

	return players
}

func GetAllNBAPlayersWithSeasonStats() []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Preload("SeasonStats").Find(&players)

	return players
}

func GetAllCollegePlayersFromOldTable() []structs.Player {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.Player

	db.Where("is_nba = ?", false).Find(&players)

	return players
}

func GetAllRecruitRecords() []structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Recruit
	db.Find(&recruits)

	return recruits
}

func GetAllCollegeRecruits() []structs.Croot {
	db := dbprovider.GetInstance().GetDB()

	var recruits []structs.Recruit
	db.Preload("RecruitProfiles", "total_points > ?", 0, func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC")
	}).Find(&recruits)

	var croots []structs.Croot
	for _, recruit := range recruits {
		var croot structs.Croot
		croot.Map(recruit)

		overallGrade := util.GetOverallGrade(recruit.Overall)

		croot.SetOverallGrade(overallGrade)

		croots = append(croots, croot)
	}

	sort.Sort(structs.ByCrootRank(croots))

	return croots
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

func GetAllNBAPlayers() []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Find(&players)

	return players
}

func GetCollegePlayerByPlayerId(playerId string) structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.CollegePlayer

	err := db.Where("id = ?", playerId).Find(&player).Error
	if err != nil {
		log.Fatal(err)
	}

	return player
}

func SetRedshirtStatusForPlayer(playerId string) structs.CollegePlayer {
	player := GetCollegePlayerByPlayerId(playerId)

	player.SetRedshirtingStatus()

	UpdatePlayer(player)

	return player
}

func UpdatePlayer(p structs.CollegePlayer) {
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

func GetNBADrafteeByNameAndCollege(firstName string, lastName string, college string) structs.HistoricCollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.HistoricCollegePlayer

	err := db.Where("first_name = ? and last_name = ? and team_abbr = ?", firstName, lastName, college).Find(&player)
	if err != nil {
		fmt.Println("Could not find player in historics DB")
	}

	return player
}

func GetAllNBAPlayersByTeamID(teamID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Where("team_id = ?", teamID).Find(&players)
	return players
}

func GetNBAPlayersWithContractsByTeamID(TeamID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Preload("Contract").Where("team_id = ?", TeamID).Find(&players)

	return players
}

func GetNBAPlayerRecord(playerID string) structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.NBAPlayer

	db.Preload("Contract").Where("id = ?", playerID).Find(&player)

	return player
}

func GetTradableNBAPlayersByTeamID(TeamID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Preload("Contract").Where("team_id = ? AND is_on_trade_block = ?", TeamID, true).Find(&players)

	return players
}
