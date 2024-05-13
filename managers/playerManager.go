package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"

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
			PositionOne:           p.PositionOne,
			PositionTwo:           p.PositionTwo,
			PositionThree:         p.PositionThree,
			P1Minutes:             p.P1Minutes,
			P2Minutes:             p.P2Minutes,
			P3Minutes:             p.P3Minutes,
			InsideProportion:      p.InsideProportion,
			MidRangeProportion:    p.MidRangeProportion,
			ThreePointProportion:  p.ThreePointProportion,
			TransferStatus:        p.TransferStatus,
			TransferLikeliness:    p.TransferLikeliness,
		}

		responseList = append(responseList, res)
	}

	return responseList
}

func GetCollegePlayersByTeamId(teamId string) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer
	err := db.Order("overall desc").Order("team_id asc").Where("team_id = ?", teamId).Find(&players).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	return players
}

func GetCollegePlayersByTeamIdForProgression(teamId string, ts structs.Timestamp) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	seasonID := strconv.Itoa(int(ts.SeasonID))

	var players []structs.CollegePlayer
	err := db.Preload("Stats", "season_id = ?", seasonID).
		Order("overall desc").
		Order("team_id asc").
		Where("team_id = ?", teamId).
		Find(&players).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	return players
}

func GetCollegePlayersWithMatchStatsByTeamId(teamId string, matchID string) []structs.MatchResultsPlayer {
	db := dbprovider.GetInstance().GetDB()

	var matchRows []structs.MatchResultsPlayer

	var players []structs.CollegePlayer
	err := db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("match_id = ?", matchID)
	}).Where("team_id = ?", teamId).Find(&players).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	for _, p := range players {
		if len(p.Stats) == 0 {
			continue
		}
		s := p.Stats[0]
		if s.Minutes == 0 {
			continue
		}
		row := structs.MatchResultsPlayer{
			FirstName:          p.FirstName,
			LastName:           p.LastName,
			Position:           p.Position,
			Archetype:          p.Archetype,
			Year:               s.Year,
			League:             "CBB",
			Minutes:            s.Minutes,
			Possessions:        s.Possessions,
			FGM:                s.FGM,
			FGA:                s.FGA,
			FGPercent:          s.FGPercent,
			ThreePointsMade:    s.ThreePointsMade,
			ThreePointAttempts: s.ThreePointAttempts,
			ThreePointPercent:  s.ThreePointPercent,
			FTM:                s.FTM,
			FTA:                s.FTA,
			FTPercent:          s.FTPercent,
			Points:             s.Points,
			TotalRebounds:      s.TotalRebounds,
			OffRebounds:        s.OffRebounds,
			DefRebounds:        s.DefRebounds,
			Assists:            s.Assists,
			Steals:             s.Steals,
			Blocks:             s.Blocks,
			Turnovers:          s.Turnovers,
			Fouls:              s.Fouls,
		}

		matchRows = append(matchRows, row)
	}

	var historicPlayers []structs.HistoricCollegePlayer
	err = db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("match_id = ?", matchID)
	}).Where("team_id = ?", teamId).Find(&historicPlayers).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	for _, p := range historicPlayers {
		if len(p.Stats) == 0 {
			continue
		}
		s := p.Stats[0]
		if s.Minutes == 0 {
			continue
		}
		row := structs.MatchResultsPlayer{
			FirstName:          p.FirstName,
			LastName:           p.LastName,
			Position:           p.Position,
			Archetype:          p.Archetype,
			League:             "CFB",
			Year:               s.Year,
			Minutes:            s.Minutes,
			Possessions:        s.Possessions,
			FGM:                s.FGM,
			FGA:                s.FGA,
			FGPercent:          s.FGPercent,
			ThreePointsMade:    s.ThreePointsMade,
			ThreePointAttempts: s.ThreePointAttempts,
			ThreePointPercent:  s.ThreePointPercent,
			FTM:                s.FTM,
			FTA:                s.FTA,
			FTPercent:          s.FTPercent,
			Points:             s.Points,
			TotalRebounds:      s.TotalRebounds,
			OffRebounds:        s.OffRebounds,
			DefRebounds:        s.DefRebounds,
			Assists:            s.Assists,
			Steals:             s.Steals,
			Blocks:             s.Blocks,
			Turnovers:          s.Turnovers,
			Fouls:              s.Fouls,
		}

		matchRows = append(matchRows, row)
	}

	// Merge both sets of players into one -- new struct: GameResultRow struct

	sort.Sort(structs.ByPlayedMinutes(matchRows))
	return matchRows
}

func GetNBAPlayersWithMatchStatsByTeamId(teamId string, matchID string) []structs.MatchResultsPlayer {
	db := dbprovider.GetInstance().GetDB()

	var matchRows []structs.MatchResultsPlayer

	var players []structs.NBAPlayer
	err := db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("match_id = ?", matchID)
	}).Where("team_id = ?", teamId).Find(&players).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	for _, p := range players {
		if len(p.Stats) == 0 {
			continue
		}
		s := p.Stats[0]
		if s.Minutes == 0 {
			continue
		}
		row := structs.MatchResultsPlayer{
			FirstName:          p.FirstName,
			LastName:           p.LastName,
			Position:           p.Position,
			Archetype:          p.Archetype,
			Year:               s.Year,
			League:             "Pro",
			Minutes:            s.Minutes,
			Possessions:        s.Possessions,
			FGM:                s.FGM,
			FGA:                s.FGA,
			FGPercent:          s.FGPercent,
			ThreePointsMade:    s.ThreePointsMade,
			ThreePointAttempts: s.ThreePointAttempts,
			ThreePointPercent:  s.ThreePointPercent,
			FTM:                s.FTM,
			FTA:                s.FTA,
			FTPercent:          s.FTPercent,
			Points:             s.Points,
			TotalRebounds:      s.TotalRebounds,
			OffRebounds:        s.OffRebounds,
			DefRebounds:        s.DefRebounds,
			Assists:            s.Assists,
			Steals:             s.Steals,
			Blocks:             s.Blocks,
			Turnovers:          s.Turnovers,
			Fouls:              s.Fouls,
		}

		matchRows = append(matchRows, row)
	}

	var historicPlayers []structs.RetiredPlayer
	err = db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
		return db.Where("match_id = ?", matchID)
	}).Where("team_id = ?", teamId).Find(&historicPlayers).Error
	if err != nil {
		log.Fatalln("Could not retrieve players from CollegePlayer Table")
	}

	for _, p := range historicPlayers {
		if len(p.Stats) == 0 {
			continue
		}
		s := p.Stats[0]
		if s.Minutes == 0 {
			continue
		}
		row := structs.MatchResultsPlayer{
			FirstName:          p.FirstName,
			LastName:           p.LastName,
			Position:           p.Position,
			Archetype:          p.Archetype,
			League:             "Pro",
			Year:               s.Year,
			Minutes:            s.Minutes,
			Possessions:        s.Possessions,
			FGM:                s.FGM,
			FGA:                s.FGA,
			FGPercent:          s.FGPercent,
			ThreePointsMade:    s.ThreePointsMade,
			ThreePointAttempts: s.ThreePointAttempts,
			ThreePointPercent:  s.ThreePointPercent,
			FTM:                s.FTM,
			FTA:                s.FTA,
			FTPercent:          s.FTPercent,
			Points:             s.Points,
			TotalRebounds:      s.TotalRebounds,
			OffRebounds:        s.OffRebounds,
			DefRebounds:        s.DefRebounds,
			Assists:            s.Assists,
			Steals:             s.Steals,
			Blocks:             s.Blocks,
			Turnovers:          s.Turnovers,
			Fouls:              s.Fouls,
		}

		matchRows = append(matchRows, row)
	}

	// Merge both sets of players into one -- new struct: GameResultRow struct

	sort.Sort(structs.ByPlayedMinutes(matchRows))
	return matchRows
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

func GetTransferPortalPlayers() []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer

	db.Where("transfer_status > 0").Find(&players)

	return players
}

func GetTransferPortalPlayersForPage() []structs.TransferPlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.CollegePlayer

	db.Preload("Profiles").Where("transfer_status = 2").Find(&players)

	playerList := []structs.TransferPlayerResponse{}

	for _, p := range players {
		res := structs.TransferPlayerResponse{}
		ovr := util.GetPlayerOverallGrade(p.Overall)
		res.Map(p, ovr)

		playerList = append(playerList, res)
	}

	return playerList
}

func GetCollegePlayerMap() map[uint]structs.CollegePlayer {

	portalMap := make(map[uint]structs.CollegePlayer)

	players := GetAllCollegePlayers()

	for _, p := range players {
		portalMap[p.ID] = p
	}

	return portalMap
}

func GetAllCollegePlayersWithSeasonStats(seasonID, weekID, matchType, viewType string) []structs.CollegePlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	seasonIDVal := util.ConvertStringToInt(seasonID)

	var players []structs.CollegePlayer
	var distinctCollegeStats []structs.CollegePlayerSeasonStats
	db.Distinct("college_player_id").Where("minutes > 0 AND season_id = ?", seasonID).Find(&distinctCollegeStats)
	distinctCollegePlayerIDs := util.GetCollegePlayerIDsBySeasonStats(distinctCollegeStats)

	if viewType == "SEASON" {
		db.Preload("SeasonStats", "season_id = ?", seasonID).
			Where("id in ?", distinctCollegePlayerIDs).Find(&players)
	} else {
		db.Preload("Stats", "season_id = ? AND week_id = ? AND match_type = ?", seasonID, weekID, matchType).
			Where("id in ?", distinctCollegePlayerIDs).Find(&players)
	}

	playerList := []structs.CollegePlayerResponse{}

	for _, p := range players {
		if len(p.Stats) == 0 && viewType == "WEEK" {
			continue
		}
		var stat structs.CollegePlayerStats
		if viewType == "WEEK" {
			stat = p.Stats[0]
		}
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
		var playerRes = structs.CollegePlayerResponse{
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
			SeasonStats:           p.SeasonStats,
			Stats:                 stat,
		}

		playerList = append(playerList, playerRes)
	}

	if seasonIDVal < int(ts.SeasonID) {
		var historicCollegePlayers []structs.HistoricCollegePlayer
		if viewType == "SEASON" {
			db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ?", seasonID)
			}).Where("id in ?", distinctCollegePlayerIDs).Find(&historicCollegePlayers)
		} else {
			db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ? AND week_id = ?", seasonID, weekID)
			}).Where("id in ?", distinctCollegePlayerIDs).Find(&historicCollegePlayers)
		}

		for _, p := range historicCollegePlayers {
			if len(p.Stats) == 0 && viewType == "WEEK" {
				continue
			}
			var stat structs.CollegePlayerStats
			if viewType == "WEEK" {
				stat = p.Stats[0]
			}
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
			var playerRes = structs.CollegePlayerResponse{
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
				SeasonStats:           p.SeasonStats,
				Stats:                 stat,
			}

			playerList = append(playerList, playerRes)
		}
	}

	return playerList
}

func GetAllNBAPlayersWithSeasonStats(seasonID, weekID, matchType, viewType string) []structs.NBAPlayerResponse {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	seasonIDVal := util.ConvertStringToInt(seasonID)

	var players []structs.NBAPlayer
	var distinctNBAStats []structs.NBAPlayerSeasonStats
	db.Distinct("nba_player_id").Where("minutes > 0 AND season_id = ?", seasonID).Find(&distinctNBAStats)
	distinctNBAPlayerIDs := util.GetNBAPlayerIDsBySeasonStats(distinctNBAStats)

	if viewType == "SEASON" {
		db.Preload("SeasonStats", "season_id = ?", seasonID).
			Where("id in ?", distinctNBAPlayerIDs).Find(&players)
	} else {
		db.Preload("Stats", "season_id = ? AND week_id = ? AND match_type = ? AND minutes > 0", seasonID, weekID, matchType).
			Where("id in ?", distinctNBAPlayerIDs).Find(&players)
	}

	playerList := []structs.NBAPlayerResponse{}

	for _, p := range players {
		if len(p.Stats) == 0 && viewType == "WEEK" {
			continue
		}
		var stat structs.NBAPlayerStats
		if viewType == "WEEK" {
			stat = p.Stats[0]
		}
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
		var playerRes = structs.NBAPlayerResponse{
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
			SeasonStats:           p.SeasonStats,
			Stats:                 stat,
		}

		playerList = append(playerList, playerRes)
	}

	if seasonIDVal < int(ts.SeasonID) {
		var historicNBAPlayers []structs.RetiredPlayer
		if viewType == "SEASON" {
			db.Preload("SeasonStats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ?", seasonID)
			}).Where("id in ?", distinctNBAPlayerIDs).Find(&historicNBAPlayers)
		} else {
			db.Preload("Stats", func(db *gorm.DB) *gorm.DB {
				return db.Where("season_id = ? AND week_id = ?", seasonID, weekID)
			}).Where("id in ?", distinctNBAPlayerIDs).Find(&historicNBAPlayers)
		}

		for _, p := range historicNBAPlayers {
			if len(p.Stats) == 0 && viewType == "WEEK" {
				continue
			}
			var stat structs.NBAPlayerStats
			if viewType == "WEEK" {
				stat = p.Stats[0]
			}
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
			var playerRes = structs.NBAPlayerResponse{
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
				SeasonStats:           p.SeasonStats,
				Stats:                 stat,
			}

			playerList = append(playerList, playerRes)
		}
	}

	return playerList
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

func GetAllRetiredPlayers() []structs.RetiredPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.RetiredPlayer

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

func GetHistoricCollegePlayerByID(id string) structs.HistoricCollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var player structs.HistoricCollegePlayer

	err := db.Where("id = ?", id).Find(&player).Error
	if err != nil {
		fmt.Println("Could not find player in historics DB")
	}

	return player
}

func GetNBADrafteeByID(id string) structs.NBADraftee {
	db := dbprovider.GetInstance().GetDB()

	var player structs.NBADraftee

	err := db.Where("id = ?", id).Find(&player).Error
	if err != nil {
		fmt.Println("Could not find player in historics DB")
	}

	return player
}

func GetAllNBADraftees() []structs.NBADraftee {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBADraftee

	db.Find(&players)

	return players
}

func GetOnlyNBAPlayersByTeamID(teamID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Where("team_id = ?", teamID).Find(&players)
	return players
}

func GetAllNBAPlayersByTeamID(teamID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Preload("Contract", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	}).Preload("Extensions", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	}).Where("team_id = ?", teamID).Find(&players)
	return players
}

func GetNBAPlayersWithContractsByTeamID(TeamID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Preload("Contract").Where("team_id = ?", TeamID).Order("overall desc").Find(&players)

	return players
}

func GetNBAPlayersWithContractsAndExtensionsByTeamID(TeamID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Preload("Contract").Preload("Extensions").Where("team_id = ?", TeamID).Order("overall desc").Find(&players)

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

func PlaceNBAPlayerInGLeague(playerID string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetNBAPlayerRecord(playerID)

	player.ToggleGLeague()

	db.Save(&player)
}

func AssignPlayerAsTwoWay(playerID string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetNBAPlayerRecord(playerID)

	player.ToggleTwoWay()

	db.Save(&player)
}

func ActivateNextYearOption(contractID string) {
	db := dbprovider.GetInstance().GetDB()

	contract := GetNBAContractByID(contractID)
	if contract.Year2Opt {
		contract.ActivateOption()
		db.Save(&contract)
	}
}

func CutCBBPlayer(playerID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	player := GetCollegePlayerByPlayerID(playerID)

	player.DismissFromTeam()

	if !ts.CollegeSeasonOver && !ts.IsOffSeason && ts.CollegeWeek == 0 {
		teamID := strconv.Itoa(int(player.PreviousTeamID))
		teamProfile := GetRecruitingProfileByTeamId(teamID)
		teamProfile.IncreaseClassSize()
		db.Save(&teamProfile)
	}

	message := "Breaking News! " + strconv.Itoa(player.Stars) + " star " + player.Position + " " + player.FirstName + " " + player.LastName + " has been dismissed from the " + player.PreviousTeam + " basketball team. They will immediately enter the transfer portal."
	CreateNewsLog("CBB", message, "Transfer Portal", int(player.PreviousTeamID), ts)

	db.Save(&player)
}

func CutNBAPlayer(playerID string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetNBAPlayerRecord(playerID)

	player.WaivePlayer()

	db.Save(&player)
}

func GetFullTeamRosterWithCrootsMap() map[uint][]structs.CollegePlayer {
	m := &sync.Mutex{}
	var wg sync.WaitGroup
	collegeTeams := GetAllActiveCollegeTeams()
	fullMap := make(map[uint][]structs.CollegePlayer)
	wg.Add(len(collegeTeams))
	semaphore := make(chan struct{}, 10)
	for _, team := range collegeTeams {
		semaphore <- struct{}{}
		go func(t structs.Team) {
			defer wg.Done()
			id := strconv.Itoa(int(t.ID))
			collegePlayers := GetCollegePlayersByTeamId(id)
			croots := GetSignedRecruitsByTeamProfileID(id)
			fullList := collegePlayers
			for _, croot := range croots {
				p := structs.CollegePlayer{}
				p.MapFromRecruit(croot)

				fullList = append(fullList, p)
			}

			m.Lock()
			fullMap[t.ID] = fullList
			m.Unlock()
			<-semaphore
		}(team)
	}

	wg.Wait()
	close(semaphore)
	return fullMap
}

func GetFullRosterNBAMap() map[uint][]structs.NBAPlayer {
	m := &sync.Mutex{}
	var wg sync.WaitGroup
	collegeTeams := GetAllActiveCollegeTeams()
	fullMap := make(map[uint][]structs.NBAPlayer)
	wg.Add(len(collegeTeams))
	semaphore := make(chan struct{}, 10)
	for _, team := range collegeTeams {
		semaphore <- struct{}{}
		go func(t structs.Team) {
			defer wg.Done()
			id := strconv.Itoa(int(t.ID))
			nbaPlayers := GetOnlyNBAPlayersByTeamID(id)

			m.Lock()
			fullMap[t.ID] = nbaPlayers
			m.Unlock()
			<-semaphore
		}(team)
	}

	wg.Wait()
	close(semaphore)
	return fullMap
}

func ProcessEarlyDeclareeAnnouncements() {
	collegePlayers := GetAllCollegePlayers()
	ts := GetTimestamp()
	for _, c := range collegePlayers {
		if (!c.WillDeclare) ||
			(c.WillDeclare && c.Year == 4 && !c.IsRedshirt) ||
			(c.WillDeclare && c.Year == 5 && c.IsRedshirt) {
			continue
		}
		playerLabel := c.TeamAbbr + " " + strconv.Itoa(c.Stars) + " star " + c.Position + " " + c.FirstName + " " + c.LastName
		message := "Breaking News! " + playerLabel + " has announced their early declaration for the upcoming SimNBA Draft!"
		CreateNewsLog("CBB", message, "Graduation", int(c.TeamID), ts)
	}
}
