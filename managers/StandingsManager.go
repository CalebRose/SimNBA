package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetConferenceStandingsByConferenceID(id string, seasonID string) []structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.CollegeStandings

	db.Where("conference_id = ? AND season_id = ?", id, seasonID).Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetNBAConferenceStandingsByConferenceID(id string, seasonID string) []structs.NBAStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.NBAStandings

	db.Where("conference_id = ? AND season_id = ?", id, seasonID).
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetAllConferenceStandingsBySeasonID(seasonID string) []structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.CollegeStandings

	db.Where("season_id = ?", seasonID).Order("conference_id asc").Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetAllNBAConferenceStandingsBySeasonID(seasonID string) []structs.NBAStandings {
	db := dbprovider.GetInstance().GetDB()

	var standings []structs.NBAStandings

	db.Where("season_id = ?", seasonID).Order("conference_id asc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings)

	return standings
}

func GetStandingsRecordByTeamID(id string, seasonID string) structs.CollegeStandings {
	db := dbprovider.GetInstance().GetDB()

	var standing structs.CollegeStandings

	db.Where("team_id = ? AND season_id = ?", id, seasonID).Find(&standing)

	return standing
}

func GetNBAStandingsRecordByTeamID(id string, seasonID string) structs.NBAStandings {
	db := dbprovider.GetInstance().GetDB()

	var standing structs.NBAStandings

	db.Where("team_id = ? AND season_id = ?", id, seasonID).Find(&standing)

	return standing
}

func GetNBAStandingsBySeasonID(seasonID string) []structs.NBAStandings {
	var standings []structs.NBAStandings
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("season_id = ?", seasonID).Order("total_losses asc").Order("total_wins desc").
		Find(&standings).Error
	if err != nil {
		log.Fatal(err)
	}
	return standings
}

func UpdateStandings(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	games := GetMatchesByWeekIdAndMatchType(strconv.Itoa(int(ts.CollegeWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)

	for i := 0; i < len(games); i++ {
		game := games[i]
		if !game.GameComplete {
			continue
		}
		HomeID := game.HomeTeamID
		AwayID := game.AwayTeamID

		homeStandings := GetStandingsRecordByTeamID(strconv.Itoa(int(HomeID)), strconv.Itoa(int(ts.SeasonID)))
		awayStandings := GetStandingsRecordByTeamID(strconv.Itoa(int(AwayID)), strconv.Itoa(int(ts.SeasonID)))

		homeStandings.UpdateCollegeStandings(game)
		awayStandings.UpdateCollegeStandings(game)

		err := db.Save(&homeStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(HomeID)))
		}

		err = db.Save(&awayStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(AwayID)))
		}

		if game.NextGameID > 0 {

			nextGameID := strconv.Itoa(int(game.NextGameID))
			winningTeamID := 0
			winningTeam := ""
			winningCoach := ""
			winningTeamRank := 0
			if game.HomeTeamWin {
				winningTeamID = int(game.HomeTeamID)
				winningTeam = game.HomeTeam
				winningTeamRank = int(game.HomeTeamRank)
				winningCoach = game.HomeTeamCoach
			} else {
				winningTeamID = int(game.AwayTeamID)
				winningTeam = game.AwayTeam
				winningTeamRank = int(game.AwayTeamRank)
				winningCoach = game.AwayTeamCoach
			}

			nextGame := GetMatchByMatchId(nextGameID)

			nextGame.AddTeam(game.NextGameHOA == "H", uint(winningTeamID), uint(winningTeamRank), winningTeam, winningCoach)

			db.Save(&nextGame)
		}

		// if games[i].HomeTeamCoach != "AI" {
		// 	homeCoach := GetCollegeCoachByCoachName(games[i].HomeTeamCoach)
		// 	homeCoach.UpdateCoachRecord(games[i])

		// 	err = db.Save(&homeCoach).Error
		// 	if err != nil {
		// 		log.Panicln("Could not save coach record for team " + strconv.Itoa(HomeID))
		// 	}
		// }

		// if games[i].AwayTeamCoach != "AI" {
		// 	awayCoach := GetCollegeCoachByCoachName(games[i].AwayTeamCoach)
		// 	awayCoach.UpdateCoachRecord(games[i])
		// 	err = db.Save(&awayCoach).Error
		// 	if err != nil {
		// 		log.Panicln("Could not save coach record for team " + strconv.Itoa(AwayID))
		// 	}
		// }
	}

	nbaGames := GetNBATeamMatchesByMatchType(strconv.Itoa(int(ts.NBAWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)

	for _, game := range nbaGames {
		if !game.GameComplete {
			continue
		}
		HomeID := game.HomeTeamID
		AwayID := game.AwayTeamID

		homeStandings := GetNBAStandingsRecordByTeamID(strconv.Itoa(int(HomeID)), strconv.Itoa(int(ts.SeasonID)))
		awayStandings := GetNBAStandingsRecordByTeamID(strconv.Itoa(int(AwayID)), strconv.Itoa(int(ts.SeasonID)))

		homeStandings.UpdateNBAStandings(game)
		awayStandings.UpdateNBAStandings(game)

		err := db.Save(&homeStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(HomeID)))
		}

		err = db.Save(&awayStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(AwayID)))
		}

		if game.NextGameID > 0 {

			nextGameID := strconv.Itoa(int(game.NextGameID))
			winningTeamID := 0
			winningTeam := ""
			winningCoach := ""
			winningTeamRank := 0
			if game.HomeTeamWin {
				winningTeamID = int(game.HomeTeamID)
				winningTeam = game.HomeTeam
				winningCoach = game.HomeTeamCoach
			} else {
				winningTeamID = int(game.AwayTeamID)
				winningTeam = game.AwayTeam
				winningCoach = game.AwayTeamCoach
			}
			nextGame := GetMatchByMatchId(nextGameID)

			nextGame.AddTeam(game.NextGameHOA == "H", uint(winningTeamID), uint(winningTeamRank), winningTeam, winningCoach)

			db.Save(&nextGame)
		}
	}
}

func RegressStandings(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	games := GetMatchesByWeekIdAndMatchType(strconv.Itoa(int(ts.CollegeWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)

	for i := 0; i < len(games); i++ {
		HomeID := games[i].HomeTeamID
		AwayID := games[i].AwayTeamID

		homeStandings := GetStandingsRecordByTeamID(strconv.Itoa(int(HomeID)), strconv.Itoa(int(ts.SeasonID)))
		awayStandings := GetStandingsRecordByTeamID(strconv.Itoa(int(AwayID)), strconv.Itoa(int(ts.SeasonID)))

		homeStandings.RegressCollegeStandings(games[i])
		awayStandings.RegressCollegeStandings(games[i])

		err := db.Save(&homeStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(HomeID)))
		}

		err = db.Save(&awayStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(AwayID)))
		}
	}
}

func ResetCollegeStandingsRanks() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))

	db.Model(&structs.CollegeStandings{}).Where("season_id = ?", seasonID).Updates(structs.CollegeStandings{Rank: 0})
}
