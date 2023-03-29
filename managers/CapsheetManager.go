package managers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetCapsheetByTeamID(TeamID string) structs.NBACapsheet {
	db := dbprovider.GetInstance().GetDB()

	capSheet := structs.NBACapsheet{}

	err := db.Where("NBA_team_id = ?", TeamID).Find(&capSheet).Error
	if err != nil {
		fmt.Println("Could not find capsheet, returning new one")
		return structs.NBACapsheet{}
	}

	return capSheet
}

func AllocateCapsheets() {
	db := dbprovider.GetInstance().GetDB()

	teams := GetAllActiveNBATeams()

	for _, team := range teams {
		TeamID := strconv.Itoa(int(team.ID))

		players := GetNBAPlayersWithContractsByTeamID(TeamID)

		Capsheet := GetCapsheetByTeamID(TeamID)

		if Capsheet.ID == 0 {
			Capsheet.AssignID(team.ID)
		}

		Capsheet.ResetCapsheet()

		sort.Sort(structs.ByTotalContract(players))

		y1 := 0.0
		y2 := 0.0
		y3 := 0.0
		y4 := 0.0
		y5 := 0.0

		for idx, player := range players {
			if idx > 50 {
				break
			}
			contract := player.Contract
			y1 += contract.Year1Total
			y2 += contract.Year2Total
			y3 += contract.Year3Total
			y4 += contract.Year4Total
			y5 += contract.Year5Total
		}

		Capsheet.SyncTotals(y1, y2, y3, y4, y5)

		db.Save(&Capsheet)
	}
}

func GetContractByPlayerID(PlayerID string) structs.NBAContract {
	db := dbprovider.GetInstance().GetDB()

	contract := structs.NBAContract{}

	err := db.Where("NBA_player_id = ? AND is_active = ?", PlayerID, true).Find(&contract).Error
	if err != nil {
		log.Fatalln("Could not find active contract for player" + PlayerID)
	}

	return contract
}

func GetAllContracts() []structs.NBAContract {
	db := dbprovider.GetInstance().GetDB()

	contracts := []structs.NBAContract{}

	err := db.Where("is_active = ?", true).Find(&contracts).Error
	if err != nil {
		log.Fatalln("Could not find all active contracts")
	}

	return contracts
}

func CalculateContractValues() {
	db := dbprovider.GetInstance().GetDB()

	contracts := GetAllContracts()

	for _, c := range contracts {
		c.CalculateContract()
		db.Save(&c)
	}
}
