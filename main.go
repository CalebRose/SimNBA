package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CalebRose/SimNBA/config"
	"github.com/CalebRose/SimNBA/controller"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
)

var db *gorm.DB
var err error
var c = config.Config()

func InitialMigration() {
	// 23.252.52.222
	// 68.66.216.54
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	db.AutoMigrate(&structs.Player{})
	db.AutoMigrate(&structs.Team{})

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", helloWorld).Methods("GET")
	myRouter.HandleFunc("/players", controller.AllPlayers).Methods("GET")
	myRouter.HandleFunc("/player/add/{firstname}/{lastname}", controller.NewPlayer).Methods("POST")
	myRouter.HandleFunc("/player/remove/{ID}", controller.RemovePlayer).Methods("DELETE")
	myRouter.HandleFunc("/player/update/{ID}", controller.UpdatePlayer).Methods("PUT")
	myRouter.HandleFunc("/teams", controller.AllTeams).Methods("GET")
	myRouter.HandleFunc("/teams/active", controller.AllActiveTeams).Methods("GET")
	myRouter.HandleFunc("/teams/available", controller.AllAvailableTeams).Methods("GET")
	myRouter.HandleFunc("/teams/coached", controller.AllCoachedTeams).Methods("GET")

	handler := cors.AllowAll().Handler(myRouter)

	log.Fatal(http.ListenAndServe(":8081", handler))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World.")
}

func main() {
	fmt.Println("GORM initiation")

	InitialMigration()

	handleRequests()
}
