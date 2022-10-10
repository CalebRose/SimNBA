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
	var gameplan = GetGameplansByTeam(teamId)

	updatedGameplan := updateGameplanDto.Gameplan

	// If no changes made to gameplan
	if gameplan.Pace == updatedGameplan.Pace &&
		gameplan.ThreePointProportion == updatedGameplan.ThreePointProportion &&
		gameplan.JumperProportion == updatedGameplan.JumperProportion &&
		gameplan.PaintProportion == updatedGameplan.PaintProportion {
		// Do Nothing
	} else {
		// Otherwise, update the gameplan
		gameplan.UpdatePace(updateGameplanDto.Gameplan.Pace)
		gameplan.Update3PtProportion(updatedGameplan.ThreePointProportion)
		gameplan.UpdateJumperProportion(updatedGameplan.JumperProportion)
		gameplan.UpdatePaintProportion(updatedGameplan.PaintProportion)
		fmt.Printf("Saving Gameplan for Team " + teamId + "\n")
		db.Save(&gameplan)
	}

	// Get Players
	var players = GetCollegePlayersByTeamId(teamId)

	for i := 0; i < len(players); i++ {
		updatedPlayer := updateGameplanDto.Players[i]
		if players[i].Minutes == updatedPlayer.Minutes {
			continue
		}
		players[i].UpdateMinutes(updatedPlayer.Minutes)

		// If player is an NBA player, update Minutes for C Game
		// if players[i].IsNBA {
		// 	players[i].UpdateMinutesC(updateGameplanDto.Players[i].MinutesC)
		// }
		fmt.Printf("Saving Player " + players[i].FirstName + " " + players[i].LastName + "\n")
		db.Save(&players[i])
	}
}

func GetGameplansByTeam(teamId string) structs.Gameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans structs.Gameplan

	db.Where("team_id = ?", teamId).Order("game asc").Find(&gameplans)

	return gameplans
}
