package managers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetCBBStatsPageData(seasonID, weekID, matchType, viewType string) structs.StatsPageResponse {
	db := dbprovider.GetInstance().GetDB()

	var teamList []structs.CollegeTeamResponse
	var playerList []structs.CollegePlayerResponse
	var conferences []structs.CollegeConference

	db.Find(&conferences)

	teamsChan := make(chan []structs.CollegeTeamResponse)
	playersChan := make(chan []structs.CollegePlayerResponse)

	go func() {
		ct := GetAllActiveCollegeTeamsWithSeasonStats(seasonID, weekID, matchType, viewType)
		teamsChan <- ct
	}()

	go func() {
		cp := GetAllCollegePlayersWithSeasonStats(seasonID, weekID, matchType, viewType)
		playersChan <- cp
	}()

	// Teams
	teamList = <-teamsChan
	close(teamsChan)

	playerList = <-playersChan
	close(playersChan)

	return structs.StatsPageResponse{
		CollegeConferences: conferences,
		CollegeTeams:       teamList,
		CollegePlayers:     playerList,
	}
}

func GetNBAStatsPageData(seasonID, weekID, matchType, viewType string) structs.NBAStatsPageResponse {
	db := dbprovider.GetInstance().GetDB()

	var teamList []structs.NBATeamResponse
	var playerList []structs.NBAPlayerResponse
	var conferences []structs.NBAConference

	db.Find(&conferences)

	teamsChan := make(chan []structs.NBATeamResponse)
	playersChan := make(chan []structs.NBAPlayerResponse)

	go func() {
		ct := GetAllActiveNBATeamsWithSeasonStats(seasonID, weekID, matchType, viewType)
		teamsChan <- ct
	}()

	go func() {
		cp := GetAllNBAPlayersWithSeasonStats(seasonID, weekID, matchType, viewType)
		playersChan <- cp
	}()

	// Teams
	teamList = <-teamsChan
	close(teamsChan)

	playerList = <-playersChan
	close(playersChan)

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

func GetNBAPlayerStatsBySeason(playerId string, seasonId string) []structs.NBAPlayerStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.NBAPlayerStats

	db.Where("nba_player_id = ? AND season_id = ?", playerId, seasonId).Find(&playerStats)

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

func GetCBBTeamResultsByMatch(teamId string, matchId string) structs.MatchResultsTeam {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.TeamStats

	db.Where("team_id = ? AND match_id = ?", teamId, matchId).Find(&teamStats)

	return structs.MatchResultsTeam{
		FirstHalfScore:  teamStats.FirstHalfScore,
		SecondHalfScore: teamStats.SecondHalfScore,
		OvertimeScore:   teamStats.OvertimeScore,
		Points:          teamStats.Points,
		Possessions:     teamStats.Possessions,
	}
}

func GetNBATeamResultsByMatch(teamId string, matchId string) structs.MatchResultsTeam {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.NBATeamStats

	db.Where("team_id = ? AND match_id = ?", teamId, matchId).Find(&teamStats)

	return structs.MatchResultsTeam{
		FirstHalfScore:  teamStats.FirstHalfScore,
		SecondHalfScore: teamStats.SecondHalfScore,
		OvertimeScore:   teamStats.OvertimeScore,
		Points:          teamStats.Points,
		Possessions:     teamStats.Possessions,
	}
}

func GetPlayerSeasonStatsByPlayerID(playerID string, seasonID string) structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var seasonStats structs.CollegePlayerSeasonStats

	err := db.Where("college_player_id = ? AND season_id = ?", playerID, seasonID).Find(&seasonStats).Error
	if err != nil {
		fmt.Println("Could not find existing record for player... generating new one.")
	}

	return seasonStats
}

func GetNBAPlayerSeasonStatsByPlayerID(playerID string, seasonID string) structs.NBAPlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var seasonStats structs.NBAPlayerSeasonStats

	err := db.Where("nba_player_id = ? AND season_id = ?", playerID, seasonID).Find(&seasonStats).Error
	if err != nil {
		fmt.Println("Could not find existing record for player... generating new one.")
	}

	return seasonStats
}

func GetTeamSeasonStatsByTeamID(teamID string, seasonID string) structs.TeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var seasonStats structs.TeamSeasonStats

	err := db.Where("team_id = ? AND season_id = ?", teamID, seasonID).Find(&seasonStats).Error
	if err != nil {
		fmt.Println("Could not find existing record for team... generating new one.")
	}

	return seasonStats
}

func GetNBATeamSeasonStatsByTeamID(teamID string, seasonID string) structs.NBATeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var seasonStats structs.NBATeamSeasonStats

	err := db.Where("team_id = ? AND season_id = ?", teamID, seasonID).Find(&seasonStats).Error
	if err != nil {
		fmt.Println("Could not find existing record for team... generating new one.")
	}

	return seasonStats
}

func UpdateSeasonStats(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	weekId := strconv.Itoa(int(ts.CollegeWeekID))
	seasonId := strconv.Itoa(int(ts.SeasonID))

	matches := GetMatchesByWeekIdAndMatchType(weekId, seasonId, MatchType)

	for _, match := range matches {
		if !match.GameComplete {
			continue
		}
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

	nbaGames := GetNBATeamMatchesByMatchType(strconv.Itoa(int(ts.NBAWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)

	for _, match := range nbaGames {
		if !match.GameComplete {
			continue
		}
		homeTeamStats := GetNBATeamStatsByMatch(strconv.Itoa(int(match.HomeTeamID)), strconv.Itoa(int(match.ID)))

		homeSeasonStats := GetNBATeamSeasonStatsByTeamID(strconv.Itoa(int(match.HomeTeamID)), seasonId)

		homeSeasonStats.AddStatsToSeasonRecord(homeTeamStats)

		err := db.Save(&homeSeasonStats).Error
		if err != nil {
			log.Fatalln("Could not save season stats for " + strconv.Itoa(int(match.HomeTeamID)))
		}

		awayTeamStats := GetNBATeamStatsByMatch(strconv.Itoa(int(match.AwayTeamID)), strconv.Itoa(int(match.ID)))

		awaySeasonStats := GetNBATeamSeasonStatsByTeamID(strconv.Itoa(int(match.AwayTeamID)), seasonId)

		awaySeasonStats.AddStatsToSeasonRecord(awayTeamStats)

		err = db.Save(&awaySeasonStats).Error
		if err != nil {
			log.Fatalln("Could not save season stats for " + strconv.Itoa(int(match.AwayTeamID)))
		}

		playerStats := GetNBAPlayerStatsByMatch(strconv.Itoa(int(match.ID)))

		for _, stat := range playerStats {
			if stat.Minutes <= 0 {
				continue
			}
			playerSeasonStats := GetNBAPlayerSeasonStatsByPlayerID(strconv.Itoa(int(stat.NBAPlayerID)), seasonId)
			playerSeasonStats.AddStatsToSeasonRecord(stat)

			err = db.Save(&playerSeasonStats).Error
			if err != nil {
				log.Fatalln("Could not save season stats for " + strconv.Itoa(int(playerSeasonStats.NBAPlayerID)))
			}
		}
	}
}

func RegressSeasonStats(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	weekId := strconv.Itoa(int(ts.CollegeWeekID))
	seasonId := strconv.Itoa(int(ts.SeasonID))

	matches := GetMatchesByWeekIdAndMatchType(weekId, seasonId, MatchType)

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

func FixNBASeasonTables() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	nbaPlayers := GetAllNBAPlayers()

	for _, p := range nbaPlayers {
		id := strconv.Itoa(int(p.ID))
		stats := GetNBAPlayerStatsBySeason(id, seasonID)
		seasonStats := GetNBAPlayerSeasonStatsByPlayerID(id, seasonID)
		for _, s := range stats {
			seasonStats.AddStatsToSeasonRecord(s)
		}

		db.Create(&seasonStats)
	}
}
