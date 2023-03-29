package managers

import (
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

// Gets all Current Season and Beyond Draft Picks
func GetDraftPicksByTeamID(TeamID string) []structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(int(ts.SeasonID))
	var picks []structs.DraftPick

	db.Where("team_id = ? AND season_id >= ?", TeamID, seasonID).Find(&picks)

	return picks
}

// Gets all Current Season and Beyond Draft Picks
func GetDraftPickByDraftPickID(DraftPickID string) structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()

	var pick structs.DraftPick

	db.Where("id = ?", DraftPickID).Find(&pick)

	return pick
}
