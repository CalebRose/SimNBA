package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/config"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// var connectionString = "simfbaah_davidross10:bestpunterev3r!@tcp(68.66.216.54)/simfbaah_simnba?parseTime=true"

var db *gorm.DB
var c = config.Config()

func AllPlayers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var players []structs.Player
	db.Find(&players)
	json.NewEncoder(w).Encode(players)
}

func NewPlayer(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	vars := mux.Vars(r)
	firstName := vars["firstname"]
	lastName := vars["lastname"]
	fmt.Println(firstName)
	fmt.Println(lastName)

	db.Create(&structs.Player{FirstName: firstName, LastName: lastName,
		Position: "C", Year: 4, State: "WA", Country: "USA",
		Stars: 3, Height: "7'0", TeamID: 10, Shooting: 14,
		Finishing: 20, Ballwork: 18, Rebounding: 20, Defense: 19,
		PotentialGrade: 20, Stamina: 36, PlaytimeExpectations: 25,
		Minutes: 35, Overall: 20, IsNBA: false,
		IsRedshirt: false, IsRedshirting: false})

	fmt.Fprintf(w, "New Player Successfully Created")
}

func RemovePlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete User Endpoint Hit")
}

func UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update User Endpoint Hit")
}
