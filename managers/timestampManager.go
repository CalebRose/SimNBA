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

func SyncToNextWeek() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	ts.SyncToNextWeek()
	err := db.Save(&ts).Error
	if err != nil {
		log.Fatalln("Could not save timestamp and sync week")
	}
}

func ShowAGames() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	UpdateStandings(ts, "A")
	UpdateSeasonStats(ts, "A")
	ts.ToggleGamesARan()
	err := db.Save(&ts).Error
	if err != nil {
		log.Fatalln("Could not save timestamp and sync week")
	}
}

func ShowBGames() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	UpdateStandings(ts, "B")
	UpdateSeasonStats(ts, "B")
	ts.ToggleGamesBRan()
	err := db.Save(&ts).Error
	if err != nil {
		log.Fatalln("Could not save timestamp and sync week")
	}
}
