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
	UpdateStandings(ts)
	UpdateSeasonStats(ts)
	ts.SyncToNextWeek()
	err := db.Save(&ts).Error
	if err != nil {
		log.Fatalln("Could not save timestamp and sync week")
	}
}
