package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/structs"
	"github.com/jinzhu/gorm"
)

func GetTeamRequests(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	defer db.Close()

	var requests []structs.Request
	db.Where("deleted_date is null AND is_approved = 0").Find(&requests)
	json.NewEncoder(w).Encode(requests)
}

func CreateTeamRequest(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	fmt.Println("Booting Up DB")

	defer db.Close()

	var request structs.Request
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(request)

	db.Create(&request)

	fmt.Fprintf(w, "Request Successfully Created")
}

func ApproveTeamRequest(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	fmt.Println("Booting Up DB")

	defer db.Close()

	var request structs.Request
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request.ApproveTeamRequest()

	db.Model(&request).Update("is_approved", request.IsApproved)

	fmt.Fprintf(w, "Request: %+v", request)
}

func RejectTeamRequest(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(c["db"], c["cs"])
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}

	fmt.Println("Booting Up DB")

	defer db.Close()

	var request structs.Request

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request.RejectTeamRequest()

	db.Delete(&request)

	fmt.Fprintf(w, "Request: %+v", request)
}
