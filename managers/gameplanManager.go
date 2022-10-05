package managers

import (
	"fmt"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func UpdateGameplan(updateGameplanDto structs.UpdateGameplanDto) {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateGameplanDto.TeamID)

	// Get Gameplans
	var gameplans = GetGameplansByTeam(teamId)

	for i := 0; i < len(gameplans); i++ {
		updatedGameplan := updateGameplanDto.Gameplans[i]

		// If no changes made to gameplan
		if gameplans[i].Pace == updatedGameplan.Pace &&
			gameplans[i].ThreePointProportion == updatedGameplan.ThreePointProportion &&
			gameplans[i].JumperProportion == updatedGameplan.JumperProportion &&
			gameplans[i].PaintProportion == updatedGameplan.PaintProportion {
			continue
		}

		// Otherwise, update the gameplan
		gameplans[i].UpdatePace(updateGameplanDto.Gameplans[i].Pace)
		gameplans[i].Update3PtProportion(updatedGameplan.ThreePointProportion)
		gameplans[i].UpdateJumperProportion(updatedGameplan.JumperProportion)
		gameplans[i].UpdatePaintProportion(updatedGameplan.PaintProportion)
		fmt.Printf("Saving Gameplan for Team " + teamId + "\n")
		db.Save(&gameplans[i])
	}

	// Get Players
	var players = GetCollegePlayersByTeamId(teamId)

	for i := 0; i < len(players); i++ {
		updatedPlayer := updateGameplanDto.Players[i]
		if players[i].Minutes == updatedPlayer.MinutesA &&
			players[i].Minutes == updatedPlayer.MinutesB &&
			players[i].Minutes == updatedPlayer.MinutesC {
			continue
		}
		players[i].UpdateMinutes(updatedPlayer.MinutesA)

		// If player is an NBA player, update Minutes for C Game
		// if players[i].IsNBA {
		// 	players[i].UpdateMinutesC(updateGameplanDto.Players[i].MinutesC)
		// }
		fmt.Printf("Saving Player " + players[i].FirstName + " " + players[i].LastName + "\n")
		db.Save(&players[i])
	}
}

func GetGameplansByTeam(teamId string) []structs.Gameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans []structs.Gameplan

	db.Where("team_id = ?", teamId).Order("game asc").Find(&gameplans)

	return gameplans
}
