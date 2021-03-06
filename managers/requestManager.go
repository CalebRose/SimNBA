package managers

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetAllTeamRequests() []structs.RequestDTO {
	db := dbprovider.GetInstance().GetDB()
	var requests []structs.RequestDTO
	db.Raw("SELECT requests.id, requests.team_id, teams.team, teams.abbr, requests.username, teams.conference, teams.is_nba, requests.is_approved FROM simfbaah_simnba.requests INNER JOIN simfbaah_simnba.teams on teams.id = requests.team_id WHERE requests.deleted_at is null AND requests.is_approved = 0").
		Scan(&requests)

	return requests
}

func CreateTeamRequest(request structs.Request) {
	db := dbprovider.GetInstance().GetDB()

	db.Create(&request)
}

func ApproveTeamRequest(request structs.Request) {
	db := dbprovider.GetInstance().GetDB()
	request.ApproveTeamRequest()

	fmt.Println("Team Approved...")

	db.Save(&request)

	fmt.Println("Assigning team...")

	// Assign Team
	team := GetTeamByTeamID(strconv.Itoa(request.TeamID))

	team.AssignUserToTeam(request.Username)

	db.Save(&team)
}

func RejectTeamRequest(request structs.Request) {
	db := dbprovider.GetInstance().GetDB()

	request.RejectTeamRequest()

	db.Delete(&request)
}
