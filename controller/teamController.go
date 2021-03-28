package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/structs"
	"github.com/jinzhu/gorm"
)

// var db *gorm.DB
// var c = config.Config()

func AllTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllActiveTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Where("first_season is not null").Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllAvailableTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Where("first_season is not null AND coach is null OR coach = ?", "AI").Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func AllCoachedTeams(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var teams []structs.Team
	db.Where("coach is not null").Order("team asc").Find(&teams)
	json.NewEncoder(w).Encode(teams)
}

func GetTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Team Endpoint Hit")
}

func NewTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "New Team Endpoint Hit")
}

func RemoveTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete Team Endpoint Hit")
}

func UpdateTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update Team Endpoint Hit")
}
