package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func GetCBBStatsPageData() structs.StatsPageResponse {
	db := dbprovider.GetInstance().GetDB()

	var teamList []structs.CollegeTeamResponse
	var playerList []structs.CollegePlayerResponse
	var conferences []structs.CollegeConference

	db.Find(&conferences)

	// Teams
	teams := GetAllActiveCollegeTeamsWithSeasonStats()

	for _, team := range teams {
		seasonsResponse := structs.TeamSeasonStatsResponse{
			ID:                        team.TeamSeasonStats.ID,
			TeamID:                    team.ID,
			SeasonID:                  team.TeamSeasonStats.SeasonID,
			GamesPlayed:               team.TeamSeasonStats.GamesPlayed,
			Points:                    team.TeamSeasonStats.Points,
			PointsAgainst:             team.TeamSeasonStats.PointsAgainst,
			PPG:                       team.TeamSeasonStats.PPG,
			PAPG:                      team.TeamSeasonStats.PAPG,
			PointsDiff:                team.TeamSeasonStats.PPG - team.TeamSeasonStats.PAPG,
			Possessions:               team.TeamSeasonStats.Possessions,
			PossessionsPerGame:        team.TeamSeasonStats.PossessionsPerGame,
			FGM:                       team.TeamSeasonStats.FGM,
			FGA:                       team.TeamSeasonStats.FGA,
			FGPercent:                 team.TeamSeasonStats.FGPercent,
			FGMPG:                     team.TeamSeasonStats.FGMPG,
			FGAPG:                     team.TeamSeasonStats.FGAPG,
			FGMAgainst:                team.TeamSeasonStats.FGMAgainst,
			FGAAgainst:                team.TeamSeasonStats.FGAAgainst,
			FGPercentAgainst:          team.TeamSeasonStats.FGPercentAgainst,
			FGMAPG:                    team.TeamSeasonStats.FGMAPG,
			FGAAPG:                    team.TeamSeasonStats.FGAAPG,
			FGMDiff:                   team.TeamSeasonStats.FGMPG - team.TeamSeasonStats.FGMAPG,
			FGADiff:                   team.TeamSeasonStats.FGAPG - team.TeamSeasonStats.FGAAPG,
			FGPercentDiff:             team.TeamSeasonStats.FGPercent - team.TeamSeasonStats.FGPercentAgainst,
			ThreePointsMade:           team.TeamSeasonStats.ThreePointsMade,
			ThreePointAttempts:        team.TeamSeasonStats.ThreePointAttempts,
			ThreePointPercent:         team.TeamSeasonStats.ThreePointPercent,
			ThreePointsMadeAgainst:    team.TeamSeasonStats.ThreePointsMadeAgainst,
			ThreePointAttemptsAgainst: team.TeamSeasonStats.ThreePointAttemptsAgainst,
			ThreePointPercentAgainst:  team.TeamSeasonStats.ThreePointPercentAgainst,
			TPMPG:                     team.TeamSeasonStats.TPMPG,
			TPAPG:                     team.TeamSeasonStats.TPAPG,
			TPMAPG:                    team.TeamSeasonStats.TPMAPG,
			TPAAPG:                    team.TeamSeasonStats.TPAAPG,
			TPMDiff:                   team.TeamSeasonStats.TPMPG - team.TeamSeasonStats.TPMAPG,
			TPADiff:                   team.TeamSeasonStats.TPAPG - team.TeamSeasonStats.TPAAPG,
			TPPercentDiff:             team.TeamSeasonStats.ThreePointPercent - team.TeamSeasonStats.ThreePointPercentAgainst,
			FTM:                       team.TeamSeasonStats.FTM,
			FTA:                       team.TeamSeasonStats.FTA,
			FTPercent:                 team.TeamSeasonStats.FTPercent,
			FTMAgainst:                team.TeamSeasonStats.FTMAgainst,
			FTAAgainst:                team.TeamSeasonStats.FTAAgainst,
			FTPercentAgainst:          team.TeamSeasonStats.FTPercentAgainst,
			FTMPG:                     team.TeamSeasonStats.FTMPG,
			FTAPG:                     team.TeamSeasonStats.FTAPG,
			FTMAPG:                    team.TeamSeasonStats.FTMAPG,
			FTAAPG:                    team.TeamSeasonStats.FTAAPG,
			FTMDiff:                   team.TeamSeasonStats.FTMPG - team.TeamSeasonStats.FTMAPG,
			FTADiff:                   team.TeamSeasonStats.FTAPG - team.TeamSeasonStats.FTAAPG,
			FTPercentDiff:             team.TeamSeasonStats.FTPercent - team.TeamSeasonStats.FTPercentAgainst,
			Rebounds:                  team.TeamSeasonStats.Rebounds,
			OffRebounds:               team.TeamSeasonStats.OffRebounds,
			DefRebounds:               team.TeamSeasonStats.DefRebounds,
			ReboundsPerGame:           team.TeamSeasonStats.ReboundsPerGame,
			OffReboundsPerGame:        team.TeamSeasonStats.OffReboundsPerGame,
			DefReboundsPerGame:        team.TeamSeasonStats.DefReboundsPerGame,
			ReboundsAllowed:           team.TeamSeasonStats.ReboundsAllowed,
			ReboundsAllowedPerGame:    team.TeamSeasonStats.ReboundsAllowedPerGame,
			OffReboundsAllowed:        team.TeamSeasonStats.OffReboundsAllowed,
			OffReboundsAllowedPerGame: team.TeamSeasonStats.OffReboundsAllowedPerGame,
			DefReboundsAllowed:        team.TeamSeasonStats.DefReboundsAllowed,
			DefReboundsAllowedPerGame: team.TeamSeasonStats.DefReboundsAllowedPerGame,
			ReboundsDiff:              team.TeamSeasonStats.ReboundsPerGame - team.TeamSeasonStats.ReboundsAllowedPerGame,
			OReboundsDiff:             team.TeamSeasonStats.OffReboundsPerGame - team.TeamSeasonStats.OffReboundsAllowedPerGame,
			DReboundsDiff:             team.TeamSeasonStats.DefReboundsPerGame - team.TeamSeasonStats.DefReboundsAllowedPerGame,
			Assists:                   team.TeamSeasonStats.Assists,
			AssistsAllowed:            team.TeamSeasonStats.AssistsAllowed,
			AssistsPerGame:            team.TeamSeasonStats.AssistsPerGame,
			AssistsAllowedPerGame:     team.TeamSeasonStats.AssistsAllowedPerGame,
			AssistsDiff:               team.TeamSeasonStats.AssistsPerGame - team.TeamSeasonStats.AssistsAllowedPerGame,
			Steals:                    team.TeamSeasonStats.Steals,
			StealsAllowed:             team.TeamSeasonStats.StealsAllowed,
			StealsPerGame:             team.TeamSeasonStats.StealsPerGame,
			StealsAllowedPerGame:      team.TeamSeasonStats.StealsAllowedPerGame,
			StealsDiff:                team.TeamSeasonStats.StealsPerGame - team.TeamSeasonStats.StealsAllowedPerGame,
			Blocks:                    team.TeamSeasonStats.Blocks,
			BlocksAllowed:             team.TeamSeasonStats.BlocksAllowed,
			BlocksPerGame:             team.TeamSeasonStats.BlocksPerGame,
			BlocksAllowedPerGame:      team.TeamSeasonStats.BlocksAllowedPerGame,
			BlocksDiff:                team.TeamSeasonStats.BlocksPerGame - team.TeamSeasonStats.BlocksAllowedPerGame,
			TotalTurnovers:            team.TeamSeasonStats.TotalTurnovers,
			TurnoversAllowed:          team.TeamSeasonStats.TurnoversAllowed,
			TurnoversPerGame:          team.TeamSeasonStats.TurnoversPerGame,
			TurnoversAllowedPerGame:   team.TeamSeasonStats.TurnoversAllowedPerGame,
			TODiff:                    team.TeamSeasonStats.TurnoversPerGame - team.TeamSeasonStats.TurnoversAllowedPerGame,
			Fouls:                     team.TeamSeasonStats.Fouls,
			FoulsPerGame:              team.TeamSeasonStats.FoulsPerGame,
		}

		var teamRes = structs.CollegeTeamResponse{
			ID:           team.ID,
			Team:         team.Team,
			Nickname:     team.Nickname,
			Abbr:         team.Abbr,
			ConferenceID: team.ConferenceID,
			Conference:   team.Conference,
			Coach:        team.Coach,
			OverallGrade: team.OverallGrade,
			OffenseGrade: team.OffenseGrade,
			DefenseGrade: team.DefenseGrade,
			IsNBA:        team.IsNBA,
			IsActive:     team.IsActive,
			SeasonStats:  seasonsResponse,
		}

		teamList = append(teamList, teamRes)
	}

	players := GetAllCollegePlayersWithSeasonStats()

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
		}

		playerList = append(playerList, playerRes)
	}

	return structs.StatsPageResponse{
		CollegeConferences: conferences,
		CollegeTeams:       teamList,
		CollegePlayers:     playerList,
	}
}

func GetNBAStatsPageData() structs.NBAStatsPageResponse {
	db := dbprovider.GetInstance().GetDB()

	var teamList []structs.NBATeamResponse
	var playerList []structs.NBAPlayerResponse
	var conferences []structs.NBAConference

	db.Find(&conferences)

	// Teams
	teams := GetAllActiveNBATeams()

	for _, team := range teams {
		seasonsResponse := structs.TeamSeasonStatsResponse{
			ID:                        team.TeamSeasonStats.ID,
			TeamID:                    team.ID,
			SeasonID:                  team.TeamSeasonStats.SeasonID,
			GamesPlayed:               team.TeamSeasonStats.GamesPlayed,
			Points:                    team.TeamSeasonStats.Points,
			PointsAgainst:             team.TeamSeasonStats.PointsAgainst,
			PPG:                       team.TeamSeasonStats.PPG,
			PAPG:                      team.TeamSeasonStats.PAPG,
			PointsDiff:                team.TeamSeasonStats.PPG - team.TeamSeasonStats.PAPG,
			Possessions:               team.TeamSeasonStats.Possessions,
			PossessionsPerGame:        team.TeamSeasonStats.PossessionsPerGame,
			FGM:                       team.TeamSeasonStats.FGM,
			FGA:                       team.TeamSeasonStats.FGA,
			FGPercent:                 team.TeamSeasonStats.FGPercent,
			FGMPG:                     team.TeamSeasonStats.FGMPG,
			FGAPG:                     team.TeamSeasonStats.FGAPG,
			FGMAgainst:                team.TeamSeasonStats.FGMAgainst,
			FGAAgainst:                team.TeamSeasonStats.FGAAgainst,
			FGPercentAgainst:          team.TeamSeasonStats.FGPercentAgainst,
			FGMAPG:                    team.TeamSeasonStats.FGMAPG,
			FGAAPG:                    team.TeamSeasonStats.FGAAPG,
			FGMDiff:                   team.TeamSeasonStats.FGMPG - team.TeamSeasonStats.FGMAPG,
			FGADiff:                   team.TeamSeasonStats.FGAPG - team.TeamSeasonStats.FGAAPG,
			FGPercentDiff:             team.TeamSeasonStats.FGPercent - team.TeamSeasonStats.FGPercentAgainst,
			ThreePointsMade:           team.TeamSeasonStats.ThreePointsMade,
			ThreePointAttempts:        team.TeamSeasonStats.ThreePointAttempts,
			ThreePointPercent:         team.TeamSeasonStats.ThreePointPercent,
			ThreePointsMadeAgainst:    team.TeamSeasonStats.ThreePointsMadeAgainst,
			ThreePointAttemptsAgainst: team.TeamSeasonStats.ThreePointAttemptsAgainst,
			ThreePointPercentAgainst:  team.TeamSeasonStats.ThreePointPercentAgainst,
			TPMPG:                     team.TeamSeasonStats.TPMPG,
			TPAPG:                     team.TeamSeasonStats.TPAPG,
			TPMAPG:                    team.TeamSeasonStats.TPMAPG,
			TPAAPG:                    team.TeamSeasonStats.TPAAPG,
			TPMDiff:                   team.TeamSeasonStats.TPMPG - team.TeamSeasonStats.TPMAPG,
			TPADiff:                   team.TeamSeasonStats.TPAPG - team.TeamSeasonStats.TPAAPG,
			TPPercentDiff:             team.TeamSeasonStats.ThreePointPercent - team.TeamSeasonStats.ThreePointPercentAgainst,
			FTM:                       team.TeamSeasonStats.FTM,
			FTA:                       team.TeamSeasonStats.FTA,
			FTPercent:                 team.TeamSeasonStats.FTPercent,
			FTMAgainst:                team.TeamSeasonStats.FTMAgainst,
			FTAAgainst:                team.TeamSeasonStats.FTAAgainst,
			FTPercentAgainst:          team.TeamSeasonStats.FTPercentAgainst,
			FTMPG:                     team.TeamSeasonStats.FTMPG,
			FTAPG:                     team.TeamSeasonStats.FTAPG,
			FTMAPG:                    team.TeamSeasonStats.FTMAPG,
			FTAAPG:                    team.TeamSeasonStats.FTAAPG,
			FTMDiff:                   team.TeamSeasonStats.FTMPG - team.TeamSeasonStats.FTMAPG,
			FTADiff:                   team.TeamSeasonStats.FTAPG - team.TeamSeasonStats.FTAAPG,
			FTPercentDiff:             team.TeamSeasonStats.FTPercent - team.TeamSeasonStats.FTPercentAgainst,
			Rebounds:                  team.TeamSeasonStats.Rebounds,
			OffRebounds:               team.TeamSeasonStats.OffRebounds,
			DefRebounds:               team.TeamSeasonStats.DefRebounds,
			ReboundsPerGame:           team.TeamSeasonStats.ReboundsPerGame,
			OffReboundsPerGame:        team.TeamSeasonStats.OffReboundsPerGame,
			DefReboundsPerGame:        team.TeamSeasonStats.DefReboundsPerGame,
			ReboundsAllowed:           team.TeamSeasonStats.ReboundsAllowed,
			ReboundsAllowedPerGame:    team.TeamSeasonStats.ReboundsAllowedPerGame,
			OffReboundsAllowed:        team.TeamSeasonStats.OffReboundsAllowed,
			OffReboundsAllowedPerGame: team.TeamSeasonStats.OffReboundsAllowedPerGame,
			DefReboundsAllowed:        team.TeamSeasonStats.DefReboundsAllowed,
			DefReboundsAllowedPerGame: team.TeamSeasonStats.DefReboundsAllowedPerGame,
			ReboundsDiff:              team.TeamSeasonStats.ReboundsPerGame - team.TeamSeasonStats.ReboundsAllowedPerGame,
			OReboundsDiff:             team.TeamSeasonStats.OffReboundsPerGame - team.TeamSeasonStats.OffReboundsAllowedPerGame,
			DReboundsDiff:             team.TeamSeasonStats.DefReboundsPerGame - team.TeamSeasonStats.DefReboundsAllowedPerGame,
			Assists:                   team.TeamSeasonStats.Assists,
			AssistsAllowed:            team.TeamSeasonStats.AssistsAllowed,
			AssistsPerGame:            team.TeamSeasonStats.AssistsPerGame,
			AssistsAllowedPerGame:     team.TeamSeasonStats.AssistsAllowedPerGame,
			AssistsDiff:               team.TeamSeasonStats.AssistsPerGame - team.TeamSeasonStats.AssistsAllowedPerGame,
			Steals:                    team.TeamSeasonStats.Steals,
			StealsAllowed:             team.TeamSeasonStats.StealsAllowed,
			StealsPerGame:             team.TeamSeasonStats.StealsPerGame,
			StealsAllowedPerGame:      team.TeamSeasonStats.StealsAllowedPerGame,
			StealsDiff:                team.TeamSeasonStats.StealsPerGame - team.TeamSeasonStats.StealsAllowedPerGame,
			Blocks:                    team.TeamSeasonStats.Blocks,
			BlocksAllowed:             team.TeamSeasonStats.BlocksAllowed,
			BlocksPerGame:             team.TeamSeasonStats.BlocksPerGame,
			BlocksAllowedPerGame:      team.TeamSeasonStats.BlocksAllowedPerGame,
			BlocksDiff:                team.TeamSeasonStats.BlocksPerGame - team.TeamSeasonStats.BlocksAllowedPerGame,
			TotalTurnovers:            team.TeamSeasonStats.TotalTurnovers,
			TurnoversAllowed:          team.TeamSeasonStats.TurnoversAllowed,
			TurnoversPerGame:          team.TeamSeasonStats.TurnoversPerGame,
			TurnoversAllowedPerGame:   team.TeamSeasonStats.TurnoversAllowedPerGame,
			TODiff:                    team.TeamSeasonStats.TurnoversPerGame - team.TeamSeasonStats.TurnoversAllowedPerGame,
			Fouls:                     team.TeamSeasonStats.Fouls,
			FoulsPerGame:              team.TeamSeasonStats.FoulsPerGame,
		}

		var teamRes = structs.NBATeamResponse{
			ID:           team.ID,
			Team:         team.Team,
			Nickname:     team.Nickname,
			Abbr:         team.Abbr,
			ConferenceID: team.ConferenceID,
			Conference:   team.Conference,
			Coach:        team.NBACoachName,
			OverallGrade: team.OverallGrade,
			OffenseGrade: team.OffenseGrade,
			DefenseGrade: team.DefenseGrade,
			IsActive:     team.IsActive,
			SeasonStats:  seasonsResponse,
		}

		teamList = append(teamList, teamRes)
	}

	players := GetAllNBAPlayersWithSeasonStats()

	for _, p := range players {
		shooting2Grade := util.GetAttributeGrade(p.Shooting2)
		shooting3Grade := util.GetAttributeGrade(p.Shooting3)
		freeThrowGrade := util.GetAttributeGrade(p.FreeThrow)
		finishingGrade := util.GetAttributeGrade(p.Finishing)
		reboundingGrade := util.GetAttributeGrade(p.Rebounding)
		ballworkGrade := util.GetAttributeGrade(p.Ballwork)
		interiorDefense := util.GetAttributeGrade(p.InteriorDefense)
		perimeterDefense := util.GetAttributeGrade(p.PerimeterDefense)
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
			InteriorDefenseGrade:  interiorDefense,
			PerimeterDefenseGrade: perimeterDefense,
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
		}

		playerList = append(playerList, playerRes)
	}

	return structs.NBAStatsPageResponse{
		NBAConferences: conferences,
		NBATeams:       teamList,
		NBAPlayers:     playerList,
	}
}

func GetPlayerStatsByPlayerId(playerId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("player_id = ?", playerId).Find(&playerStats)

	return playerStats
}

func GetPlayerStatsBySeason(playerId string, seasonId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("player_id = ? AND season_id = ?", playerId, seasonId).Find(&playerStats)

	return playerStats
}

func GetPlayerStatsInConferenceBySeason(seasonId string, conference string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("season_id = ? AND conference = ?", seasonId, conference).Find(&playerStats)

	return playerStats
}

func GetPlayerStatsByMatch(matchId string) []structs.CollegePlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerStats

	db.Where("match_id = ?", matchId).Find(&playerStats)

	return playerStats
}

func GetNBAPlayerStatsByMatch(matchId string) []structs.NBAPlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NBAPlayerStats

	db.Where("match_id = ?", matchId).Find(&playerStats)

	return playerStats
}

func GetTeamStatsBySeason(teamId string, seasonId string) []structs.PlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.PlayerStats

	db.Where("team_id = ? AND season_id = ?", teamId, seasonId).Find(&playerStats)

	return playerStats
}

func GetCBBTeamStatsByMatch(teamId string, matchId string) structs.TeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.TeamStats

	db.Where("team_id = ? AND match_id = ?", teamId, matchId).Find(&teamStats)

	return teamStats
}

func GetNBATeamStatsByMatch(teamId string, matchId string) structs.NBATeamStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.NBATeamStats

	db.Where("team_id = ? AND match_id = ?", teamId, matchId).Find(&teamStats)

	return teamStats
}

func GetPlayerSeasonStatsByPlayerID(playerID string, seasonID string) structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var seasonStats structs.CollegePlayerSeasonStats

	err := db.Where("college_player_id = ? AND season_id = ?", playerID, seasonID).Find(&seasonStats)
	if err != nil {
		fmt.Println("Could not find existing record for player... generating new one.")
	}

	return seasonStats
}

func GetTeamSeasonStatsByTeamID(teamID string, seasonID string) structs.TeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var seasonStats structs.TeamSeasonStats

	err := db.Where("team_id = ? AND season_id = ?", teamID, seasonID).Find(&seasonStats)
	if err != nil {
		fmt.Println("Could not find existing record for team... generating new one.")
	}

	return seasonStats
}

func UpdateSeasonStats(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	weekId := strconv.Itoa(int(ts.CollegeWeekID))
	seasonId := strconv.Itoa(int(ts.SeasonID))

	matches := GetMatchesByWeekId(weekId, seasonId, MatchType)

	for _, match := range matches {
		homeTeamStats := GetCBBTeamStatsByMatch(strconv.Itoa(int(match.HomeTeamID)), strconv.Itoa(int(match.ID)))

		homeSeasonStats := GetTeamSeasonStatsByTeamID(strconv.Itoa(int(match.HomeTeamID)), seasonId)

		homeSeasonStats.AddStatsToSeasonRecord(homeTeamStats)

		err := db.Save(&homeSeasonStats).Error
		if err != nil {
			log.Fatalln("Could not save season stats for " + strconv.Itoa(int(match.HomeTeamID)))
		}

		awayTeamStats := GetCBBTeamStatsByMatch(strconv.Itoa(int(match.AwayTeamID)), strconv.Itoa(int(match.ID)))

		awaySeasonStats := GetTeamSeasonStatsByTeamID(strconv.Itoa(int(match.AwayTeamID)), seasonId)

		awaySeasonStats.AddStatsToSeasonRecord(awayTeamStats)

		err = db.Save(&awaySeasonStats).Error
		if err != nil {
			log.Fatalln("Could not save season stats for " + strconv.Itoa(int(match.AwayTeamID)))
		}

		playerStats := GetPlayerStatsByMatch(strconv.Itoa(int(match.ID)))

		for _, stat := range playerStats {
			if stat.Minutes <= 0 {
				continue
			}
			playerSeasonStats := GetPlayerSeasonStatsByPlayerID(strconv.Itoa(int(stat.CollegePlayerID)), seasonId)
			playerSeasonStats.AddStatsToSeasonRecord(stat)

			err = db.Save(&playerSeasonStats).Error
			if err != nil {
				log.Fatalln("Could not save season stats for " + strconv.Itoa(int(playerSeasonStats.CollegePlayerID)))
			}
		}
	}
}

func RegressSeasonStats(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	weekId := strconv.Itoa(int(ts.CollegeWeekID))
	seasonId := strconv.Itoa(int(ts.SeasonID))

	matches := GetMatchesByWeekId(weekId, seasonId, MatchType)

	for _, match := range matches {
		homeTeamStats := GetCBBTeamStatsByMatch(strconv.Itoa(int(match.HomeTeamID)), strconv.Itoa(int(match.ID)))

		homeSeasonStats := GetTeamSeasonStatsByTeamID(strconv.Itoa(int(match.HomeTeamID)), seasonId)

		homeSeasonStats.RemoveStatsToSeasonRecord(homeTeamStats)

		err := db.Save(&homeSeasonStats).Error
		if err != nil {
			log.Fatalln("Could not save season stats for " + strconv.Itoa(int(match.HomeTeamID)))
		}

		awayTeamStats := GetCBBTeamStatsByMatch(strconv.Itoa(int(match.AwayTeamID)), strconv.Itoa(int(match.ID)))

		awaySeasonStats := GetTeamSeasonStatsByTeamID(strconv.Itoa(int(match.AwayTeamID)), seasonId)

		awaySeasonStats.RemoveStatsToSeasonRecord(awayTeamStats)

		err = db.Save(&awaySeasonStats).Error
		if err != nil {
			log.Fatalln("Could not save season stats for " + strconv.Itoa(int(match.AwayTeamID)))
		}

		playerStats := GetPlayerStatsByMatch(strconv.Itoa(int(match.ID)))

		for _, stat := range playerStats {
			if stat.Minutes <= 0 {
				continue
			}
			playerSeasonStats := GetPlayerSeasonStatsByPlayerID(strconv.Itoa(int(stat.CollegePlayerID)), seasonId)
			playerSeasonStats.RemoveStatsToSeasonRecord(stat)

			err = db.Save(&playerSeasonStats).Error
			if err != nil {
				log.Fatalln("Could not save season stats for " + strconv.Itoa(int(playerSeasonStats.CollegePlayerID)))
			}
		}
	}
}
