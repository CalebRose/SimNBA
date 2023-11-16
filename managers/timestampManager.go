package managers

import (
	"log"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetTimestamp() structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()

	var timeStamp structs.Timestamp

	err := db.Find(&timeStamp).Error
	if err != nil {
		log.Fatal(err)
	}

	return timeStamp
}

func LockRecruiting() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	ts.ToggleLockRecruiting()

	err := db.Save(&ts).Error
	if err != nil {
		log.Fatal(err)
	}
}

func MoveUpInOffseasonFreeAgency() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	if ts.IsNBAOffseason {
		ts.MoveUpFreeAgencyRound()
	}
	db.Save(&ts)
}

func SyncToNextWeek() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	ts.SyncToNextWeek()
	err := db.Save(&ts).Error
	if err != nil {
		log.Fatalln("Could not save timestamp and sync week")
	}
}

func ShowGames() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	matchType := ""
	if !ts.GamesARan {
		matchType = "A"
	} else if !ts.GamesBRan {
		matchType = "B"
	} else if !ts.GamesCRan {
		matchType = "C"
	} else if !ts.GamesDRan {
		matchType = "D"
	}
	if matchType == "" {
		log.Fatalln("Cannot sync results!")
	}
	UpdateStandings(ts, matchType)
	UpdateSeasonStats(ts, matchType)
	ts.ToggleGames(matchType)
	err := db.Save(&ts).Error
	if err != nil {
		log.Fatalln("Could not save timestamp and sync week")
	}
}

func RegressGames(match string) {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	RegressStandings(ts, match)
	RegressSeasonStats(ts, match)
	ts.ToggleGamesARan()
	err := db.Save(&ts).Error
	if err != nil {
		log.Fatalln("Could not save timestamp and sync week")
	}
}
