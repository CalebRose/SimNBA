package managers

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

// UpdateGameplan -- Need to update
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

// GetGameplansByTeam
func GetGameplansByTeam(teamId string) structs.Gameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans structs.Gameplan

	db.Where("team_id = ?", teamId).Order("game asc").Find(&gameplans)

	return gameplans
}

func SetAIGameplans() bool {
	db := dbprovider.GetInstance().GetDB()

	teams := GetAllActiveCollegeTeams()

	for _, team := range teams {
		if !team.IsActive {
			continue
		}

		roster := GetCollegePlayersByTeamId(strconv.Itoa(int(team.ID)))
		sort.Sort(structs.ByPlayerOverall(roster))
		totalMinutes := 0
		for _, player := range roster {
			totalMinutes += player.Minutes
		}

		if totalMinutes == 200 {
			continue
		}

		for idx, player := range roster {
			if idx < 5 {
				player.SetMinutes(25)
			} else if idx < 7 {
				player.SetMinutes(20)
			} else if idx == 7 {
				player.SetMinutes(8)
			} else if idx == 8 {
				player.SetMinutes(7)
			} else if idx < 13 {
				player.SetMinutes(5)
			} else {
				continue
			}
			db.Save(&player)
		}

		gameplan := GetGameplansByTeam(strconv.Itoa(int(team.ID)))
		gameplan.Update3PtProportion(20)
		gameplan.UpdateJumperProportion(40)
		gameplan.UpdatePaintProportion(40)

		db.Save(&gameplan)
	}

	return true
}
