package managers

import "github.com/CalebRose/SimNBA/structs"

type StatsObject interface {
	MapToStats(teamID, matchID, weekID, seasonID uint, TeamOne, TeamTwo structs.TeamResultsDTO) interface{}
	MapToPlayerStats(player structs.PlayerDTO, id, matchID int, seasonID uint) interface{}
}
