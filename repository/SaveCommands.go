package repository

import (
	"log"

	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func SaveCollegePlayerRecord(player structs.CollegePlayer, db *gorm.DB) {
	player.Stats = nil
	player.SeasonStats = structs.CollegePlayerSeasonStats{}
	player.Profiles = nil
	err := db.Save(&player).Error
	if err != nil {
		log.Panicln("Could not save player record")
	}
}
