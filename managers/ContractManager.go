package managers

import (
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetNBACapsheetByTeamID(teamID string) structs.NBACapSheet {
	db := dbprovider.GetInstance().GetDB()

	var capsheet structs.NBACapSheet

	db.Where("team_id = ?", teamID).Find(&capsheet)
	return capsheet
}

func GetNBAContractsByPlayerID(playerID string) structs.NBAContract {
	db := dbprovider.GetInstance().GetDB()

	var contracts structs.NBAContract

	db.Where("player_id = ? AND is_active = true", playerID).Find(&contracts)
	return contracts
}

func GetNBAContractOffersByPlayerID(playerID string, seasonID string, weekID string) []structs.NBAContractOffer {
	db := dbprovider.GetInstance().GetDB()

	var offers []structs.NBAContractOffer

	db.Order("total_cost DESC").Where("player_id = ? AND season_id = ? AND week_id = ?", playerID, seasonID, weekID).Find(&offers)
	return offers
}

func SyncFreeAgentOffers(playerID string, seasonID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	weekID := strconv.Itoa(int(ts.NBAWeekID))

	freeAgents := GetAllNBAPlayersByTeamID("0")

	for _, player := range freeAgents {
		if !player.IsFreeAgent {
			continue
		}

		playerID := strconv.Itoa(int(player.ID))

		offers := GetNBAContractOffersByPlayerID(playerID, seasonID, weekID)

		for _, offer := range offers {
			if !player.IsFreeAgent {
				offer.RejectOffer()
				db.Save(&offer)
				continue
			}
			// Check for cap space
			teamID := strconv.Itoa(int(offer.TeamID))
			// Get Cap Sheet
			capsheet := GetNBACapsheetByTeamID(teamID)
			year1Diff := capsheet.Year1Cap - (capsheet.Year1Total + offer.Year1Total)
			if year1Diff < 0 {
				offer.RejectOffer()
				db.Save(&offer)
				continue
			}
			year2Diff := capsheet.Year2Cap - (capsheet.Year2Total + offer.Year2Total)
			if year2Diff < 0 {
				offer.RejectOffer()
				db.Save(&offer)
				continue
			}
			year3Diff := capsheet.Year3Cap - (capsheet.Year3Total + offer.Year3Total)
			if year3Diff < 0 {
				offer.RejectOffer()
				db.Save(&offer)
				continue
			}
			year4Diff := capsheet.Year4Cap - (capsheet.Year4Total + offer.Year4Total)
			if year4Diff < 0 {
				offer.RejectOffer()
				db.Save(&offer)
				continue
			}
			year5Diff := capsheet.Year5Cap - (capsheet.Year5Total + offer.Year5Total)
			if year5Diff < 0 {
				offer.RejectOffer()
				db.Save(&offer)
				continue
			}
			// If below the cap, accept offer
			offer.AcceptOffer()
			player.SignWithTeam(offer.TeamID, offer.Team)
			contract := GetNBAContractsByPlayerID(playerID)
			if contract.ID == 0 {
				contract.MapFromOffer(offer)
				db.Create(&contract)
			} else {
				contract.MapFromOffer(offer)
				db.Save(&contract)
			}
			db.Save(&offer)
			db.Save(&player)
		}
	}

}
