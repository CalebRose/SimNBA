package managers

import (
	"fmt"
	"log"
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

	err := db.Save(&request).Error
	if err != nil {
		log.Fatalln("Could not create record to DB:" + err.Error())
	}
}

func ApproveTeamRequest(request structs.Request) {
	db := dbprovider.GetInstance().GetDB()
	request.ApproveTeamRequest()

	fmt.Println("Team Approved...")

	db.Save(&request)

	fmt.Println("Assigning team...")

	// Assign Team
	team := GetTeamByTeamID(strconv.Itoa(int(request.TeamID)))

	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(strconv.Itoa(int(request.TeamID)))

	recruitingProfile.ToggleAIBehavior(false)

	db.Save(&recruitingProfile)

	standing := GetStandingsRecordByTeamID(strconv.Itoa(int(request.TeamID)))

	standing.UpdateCoach(request.Username)

	team.AssignUserToTeam(request.Username)

	db.Save(&team)

	db.Save(&standing)
}

func RejectTeamRequest(request structs.Request) {
	db := dbprovider.GetInstance().GetDB()

	request.RejectTeamRequest()

	db.Delete(&request)
}
