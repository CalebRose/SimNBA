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

	db.Where("conference_id = ? AND season_id = ?", id, seasonID).Order("conference_losses asc").Order("conference_wins desc").
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

	db.Where("season_id = ?", seasonID).Order("conference_id asc").Order("conference_losses asc").Order("conference_wins desc").
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

func UpdateStandings(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	games := GetMatchesByWeekId(strconv.Itoa(int(ts.CollegeWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)

	for i := 0; i < len(games); i++ {
		HomeID := games[i].HomeTeamID
		AwayID := games[i].AwayTeamID

		homeStandings := GetStandingsRecordByTeamID(strconv.Itoa(int(HomeID)), strconv.Itoa(int(ts.SeasonID)))
		awayStandings := GetStandingsRecordByTeamID(strconv.Itoa(int(AwayID)), strconv.Itoa(int(ts.SeasonID)))

		homeStandings.UpdateCollegeStandings(games[i])
		awayStandings.UpdateCollegeStandings(games[i])

		err := db.Save(&homeStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(HomeID)))
		}

		err = db.Save(&awayStandings).Error
		if err != nil {
			log.Panicln("Could not save standings for team " + strconv.Itoa(int(AwayID)))
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
}

func RegressStandings(ts structs.Timestamp, MatchType string) {
	db := dbprovider.GetInstance().GetDB()

	games := GetMatchesByWeekId(strconv.Itoa(int(ts.CollegeWeekID)), strconv.Itoa(int(ts.SeasonID)), MatchType)

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
