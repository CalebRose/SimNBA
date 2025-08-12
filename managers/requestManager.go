package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func GetAllTeamRequests() []structs.RequestDTO {
	db := dbprovider.GetInstance().GetDB()
	var requests []structs.RequestDTO
	db.Raw("SELECT requests.id, requests.team_id, teams.team, teams.abbr, requests.username, teams.conference, teams.is_nba, requests.is_approved FROM simfbaah_simnba.requests INNER JOIN simfbaah_simnba.teams on teams.id = requests.team_id WHERE requests.deleted_at is null AND requests.is_approved = 0").
		Scan(&requests)

	return requests
}

func GetAllNBATeamRequests() []structs.NBARequest {
	db := dbprovider.GetInstance().GetDB()
	var NBATeamRequests []structs.NBARequest

	//NBA Team Requests
	db.Where("is_approved = false").Find(&NBATeamRequests)

	return NBATeamRequests
}

func CreateTeamRequest(request structs.Request) {
	db := dbprovider.GetInstance().GetDB()

	err := db.Save(&request).Error
	if err != nil {
		log.Fatalln("Could not create record to DB:" + err.Error())
	}
}

func CreateNBATeamRequest(request structs.NBARequest) {
	db := dbprovider.GetInstance().GetDB()

	var existingRequest structs.NBARequest
	err := db.Where("username = ? AND nba_team_id = ? AND is_owner = ? AND is_manager = ? AND is_coach = ? AND is_assistant = ? AND is_approved = false AND deleted_at is null", request.Username, request.NBATeamID, request.IsOwner, request.IsManager, request.IsCoach, request.IsAssistant).Find(&existingRequest).Error
	if err != nil {
		// Then there's no existing record, I guess? Which is fine.
		fmt.Println("Creating Team Request for TEAM " + strconv.Itoa(int(request.NBATeamID)))
	}
	if existingRequest.ID != 0 {
		// There is already an existing record.
		log.Fatalln("There is already an existing request in place for the user. Please be patient while admin approves your formal request. If there is an issue, please reach out to TuscanSota.")
	}

	db.Create(&request)
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

	recruitingProfile.ActivateUserTeam()

	repository.SaveCBBTeamRecruitingProfile(recruitingProfile, db)

	ts := GetTimestamp()

	standing := GetStandingsRecordByTeamID(strconv.Itoa(int(request.TeamID)), strconv.Itoa(int(ts.SeasonID)))

	standing.UpdateCoach(request.Username)

	matches := GetMatchesByTeamIdAndSeasonId(strconv.Itoa(int(request.TeamID)), strconv.Itoa(int(ts.SeasonID)))

	for _, match := range matches {
		match.UpdateCoach(int(request.TeamID), request.Username)
		db.Save(&match)
	}

	team.AssignUserToTeam(request.Username)

	db.Save(&team)

	db.Save(&standing)
}

func RejectTeamRequest(request structs.Request) {
	db := dbprovider.GetInstance().GetDB()

	request.RejectTeamRequest()

	db.Delete(&request)
}

func ApproveNBATeamRequest(request structs.NBARequest) structs.NBARequest {
	db := dbprovider.GetInstance().GetDB()

	timestamp := GetTimestamp()

	// Approve Request
	request.ApproveTeamRequest()

	fmt.Println("Team Approved...")

	db.Save(&request)

	// Assign Team
	fmt.Println("Assigning team...")

	team := GetNBATeamByTeamID(strconv.Itoa(int(request.NBATeamID)))

	coach := GetNBAUserByUsername(request.Username)

	coach.SetTeam(request)

	team.AssignNBAUserToTeam(request, coach)

	// seasonalGames := GetCollegeGamesByTeamIdAndSeasonId(strconv.Itoa(request.TeamID), strconv.Itoa(timestamp.CollegeSeasonID))

	// for _, game := range seasonalGames {
	// 	if game.Week >= timestamp.CollegeWeek {
	// 		game.UpdateCoach(int(request.NBATeamID), request.Username)
	// 		db.Save(&game)
	// 	}
	// }

	db.Save(&team)

	db.Save(&coach)

	message := "Breaking News! The " + team.Team + " " + team.Nickname + " have hired " + coach.Username + " to their staff for the " + strconv.Itoa(timestamp.Season) + " season!"
	CreateNewsLog("CBB", message, "CoachJob", int(team.ID), timestamp)
	return request
}

func RejectNBATeamRequest(request structs.NBARequest) {
	db := dbprovider.GetInstance().GetDB()

	request.RejectTeamRequest()

	err := db.Delete(&request).Error
	if err != nil {
		log.Fatalln("Could not delete request: " + err.Error())
	}
}

func RemoveUserFromTeam(teamId string) structs.Team {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	team := GetTeamByTeamID(teamId)

	team.RemoveUser()
	team.AssignDiscordID("")

	standings := GetStandingsRecordByTeamID(teamId, strconv.Itoa(int(ts.SeasonID)))

	standings.UpdateCoach("AI")

	matches := GetMatchesByTeamIdAndSeasonId(teamId, strconv.Itoa(int(ts.SeasonID)))
	for _, match := range matches {
		match.UpdateCoach(int(team.ID), "AI")
		db.Save(&match)
	}

	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(teamId)

	recruitingProfile.DeactivateUserTeam()

	db.Save(&team)
	repository.SaveCBBTeamRecruitingProfile(recruitingProfile, db)

	db.Save(&standings)

	return team
}

func RemoveUserFromNBATeam(request structs.NBARequest) {
	db := dbprovider.GetInstance().GetDB()

	teamID := strconv.Itoa(int(request.NBATeamID))

	team := GetNBATeamByTeamID(teamID)

	user := GetNBAUserByUsername(request.Username)

	message := ""

	team.AssignDiscordID("", request.Username)

	if request.Username == team.NBAOwnerName {
		user.RemoveOwnership()
		message = request.Username + " has decided to step down as Owner of the " + team.Team + " " + team.Nickname + "!"
	}

	if request.Username == team.NBAGMName {
		user.RemoveManagerPosition()
		message = request.Username + " has decided to step down as Manager of the " + team.Team + " " + team.Nickname + "!"
	}

	if request.Username == team.NBACoachName {
		user.RemoveCoachPosition()
		message = request.Username + " has decided to step down as Head Coach of the " + team.Team + " " + team.Nickname + "!"
	}

	if request.Username == team.NBAAssistantName {
		user.RemoveAssistantPosition()
		message = request.Username + " has decided to step down as an Assistant of the " + team.Team + " " + team.Nickname + "!"
	}

	team.RemoveUser(user.Username)

	db.Save(&team)

	db.Save(&user)

	timestamp := GetTimestamp()
	CreateNewsLog("NBA", message, "CoachJob", int(team.ID), timestamp)
}

func GetCBBTeamForAvailableTeamsPage(teamID string) structs.TeamRecordResponse {
	historicalDataResponse := GetHistoricalCBBRecordsByTeamID(teamID)

	// Get top 3 players on roster
	roster := GetCollegePlayersByTeamId(teamID)
	sort.Slice(roster, func(i, j int) bool {
		return roster[i].Overall > roster[j].Overall
	})

	topPlayers := []structs.TopPlayer{}

	for i := range roster {
		if i > 4 {
			break
		}
		tp := structs.TopPlayer{}
		grade := util.GetPlayerOverallGrade(roster[i].Overall)
		tp.MapCollegePlayer(roster[i], grade)
		topPlayers = append(topPlayers, tp)
	}

	historicalDataResponse.AddTopPlayers(topPlayers)

	return historicalDataResponse
}

func GetNBATeamForAvailableTeamsPage(teamID string) structs.TeamRecordResponse {
	historicalDataResponse := GetHistoricalNBARecordsByTeamID(teamID)

	// Get top 3 players on roster
	roster := GetAllNBAPlayersByTeamID(teamID)
	sort.Slice(roster, func(i, j int) bool {
		return roster[i].Overall > roster[j].Overall
	})

	topPlayers := []structs.TopPlayer{}

	for i := range roster {
		if i > 4 {
			break
		}
		tp := structs.TopPlayer{}
		tp.MapNBAPlayer(roster[i])
		topPlayers = append(topPlayers, tp)
	}

	historicalDataResponse.AddTopPlayers(topPlayers)

	return historicalDataResponse
}
