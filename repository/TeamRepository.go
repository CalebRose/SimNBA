package repository

import (
	"log"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func FindAllArenaRecords() []structs.Arena {
	db := dbprovider.GetInstance().GetDB()
	arenas := []structs.Arena{}
	err := db.Find(&arenas).Error
	if err != nil {
		log.Fatal(err)
	}
	return arenas
}
