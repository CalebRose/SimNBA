package managers

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/secrets"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

func GetAllAvailableNBAPlayers(TeamID string) structs.FreeAgencyResponse {
	FAs := GetAllFreeAgentsWithOffers()
	waiverPlayers := GetAllWaiverWirePlayers()
	gLeagePlayer := GetAllGLeaguePlayersForFA()
	islPlayers := GetAllISLPlayersForFA()
	Offers := GetFreeAgentOffersByTeamID(TeamID)
	roster := GetAllNBAPlayersByTeamID(TeamID)
	count := 0
	for _, p := range roster {
		if p.IsGLeague {
			continue
		}
		count += 1
	}

	return structs.FreeAgencyResponse{
		FreeAgents:     FAs,
		WaiverPlayers:  waiverPlayers,
		GLeaguePlayers: gLeagePlayer,
		ISLPlayers:     islPlayers,
		TeamOffers:     Offers,
		RosterCount:    uint(count),
	}
}

func GetAllFreeAgents() []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	fas := []structs.NBAPlayer{}

	db.Where("is_free_agent = ?", true).Find(&fas)

	return fas
}

// GetAllFreeAgentsWithOffers -- For Free Agency UI Page.
func GetAllFreeAgentsWithOffers() []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	fas := []structs.NBAPlayer{}

	db.Preload("Offers", func(db *gorm.DB) *gorm.DB {
		return db.Order("contract_value DESC").Where("is_active = true")
	}).Order("overall desc").Where("is_free_agent = ?", true).Find(&fas)

	return fas
}

func GetAllWaiverWirePlayers() []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	WaivedPlayers := []structs.NBAPlayer{}

	db.Preload("WaiverOffers").Preload("Contract").Where("is_waived = ?", true).Find(&WaivedPlayers)

	return WaivedPlayers
}

func GetAllGLeaguePlayersForFA() []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	gLeaguePlayers := []structs.NBAPlayer{}

	db.Preload("WaiverOffers").Preload("Contract").Where("is_g_league = ?", true).Find(&gLeaguePlayers)

	return gLeaguePlayers
}

func GetAllISLPlayersForFA() []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	islPlayers := []structs.NBAPlayer{}

	db.Preload("WaiverOffers").Preload("Contract").Where("team_id > 32").Find(&islPlayers)

	return islPlayers
}

func GetFreeAgentOffersByTeamID(TeamID string) []structs.NBAContractOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAContractOffer{}

	err := db.Where("team_id = ? AND is_active = ?", TeamID, true).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func CreateFAOffer(offer structs.NBAContractOfferDTO) structs.NBAContractOffer {
	db := dbprovider.GetInstance().GetDB()

	freeAgentOffer := GetFreeAgentOfferByOfferID(strconv.Itoa(int(offer.ID)))

	if freeAgentOffer.ID == 0 {
		id := GetLatestOfferInDB(db)
		freeAgentOffer.AssignID(id)
	}

	freeAgentOffer.CalculateOffer(offer)

	db.Save(&freeAgentOffer)

	fmt.Println("Creating offer!")

	return freeAgentOffer
}

func GetFreeAgentOfferByOfferID(OfferID string) structs.NBAContractOffer {
	db := dbprovider.GetInstance().GetDB()

	offer := structs.NBAContractOffer{}

	err := db.Where("id = ?", OfferID).Find(&offer).Error
	if err != nil {
		return offer
	}

	return offer
}

func GetLatestOfferInDB(db *gorm.DB) uint {
	var latestOffer structs.NBAContractOffer

	err := db.Last(&latestOffer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestOffer.ID + 1
}

func CancelOffer(offer structs.NBAContractOfferDTO) {
	db := dbprovider.GetInstance().GetDB()

	OfferID := strconv.Itoa(int(offer.ID))

	freeAgentOffer := GetFreeAgentOfferByOfferID(OfferID)

	freeAgentOffer.CancelOffer()

	db.Save(&freeAgentOffer)
}

func CreateWaiverOffer(offer structs.NBAWaiverOfferDTO) structs.NBAWaiverOffer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	waiverOffer := GetWaiverOfferByOfferID(strconv.Itoa(int(offer.ID)))

	if waiverOffer.ID == 0 {
		id := GetLatestWaiverOfferInDB(db)
		waiverOffer.AssignID(id)
	}

	if ts.IsFreeAgencyLocked {
		return waiverOffer
	}

	waiverOffer.Map(offer)
	playerIDStr := strconv.Itoa(int(offer.PlayerID))
	nbaPlayer := GetNBAPlayerRecord(playerIDStr)

	if nbaPlayer.IsGLeague && nbaPlayer.TeamID == offer.TeamID {
		// Sign player back to team
		nbaPlayer.ToggleGLeague()
		db.Save(&nbaPlayer)

		otherWaiverOffers := GetWaiverOffersByPlayerID(playerIDStr)

		for _, o := range otherWaiverOffers {
			db.Delete(&o)
		}
		message := "Breaking News! " + nbaPlayer.FirstName + " " + nbaPlayer.LastName + " has been picked up from the GLeague onto his owning team, " + offer.Team + "!"
		CreateNewsLog("NBA", message, "FreeAgency", 0, ts)
		return waiverOffer
	}

	db.Save(&waiverOffer)

	leagueType := ""
	if nbaPlayer.IsGLeague {
		leagueType = "G-League Player "
	} else if nbaPlayer.IsInternational {
		leagueType = "ISL Player "
	}
	message := "Breaking News! " + offer.Team + " have placed a waivering offer on " + leagueType + nbaPlayer.Position + " " + nbaPlayer.FirstName + " " + nbaPlayer.LastName + "!"
	CreateNewsLog("NBA", message, "FreeAgency", int(offer.TeamID), ts)

	fmt.Println("Creating offer!")

	return waiverOffer
}

func CancelWaiverOffer(offer structs.NBAWaiverOfferDTO) {
	db := dbprovider.GetInstance().GetDB()

	OfferID := strconv.Itoa(int(offer.ID))

	waiverOffer := GetWaiverOfferByOfferID(OfferID)

	db.Delete(&waiverOffer)
}

func GetWaiverOfferByOfferID(OfferID string) structs.NBAWaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offer := structs.NBAWaiverOffer{}

	err := db.Where("id = ?", OfferID).Find(&offer).Error
	if err != nil {
		return offer
	}

	return offer
}

func GetLatestWaiverOfferInDB(db *gorm.DB) uint {
	var latestOffer structs.NBAWaiverOffer

	err := db.Last(&latestOffer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestOffer.ID + 1
}

func GetWaiverOffersByPlayerID(playerID string) []structs.NBAWaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAWaiverOffer{}

	err := db.Order("waiver_order asc").Where("player_id = ?", playerID).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func SetWaiverOrder() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	nbaTeams := GetAllActiveNBATeams()

	teamMap := make(map[uint]*structs.NBATeam)

	for i := 0; i < len(nbaTeams); i++ {
		teamMap[nbaTeams[i].ID] = &nbaTeams[i]
	}

	var nbaStandings []structs.NBAStandings

	if ts.IsNBAOffseason || ts.NBAWeek < 3 {
		nbaStandings = GetNBAStandingsBySeasonID(strconv.Itoa(int(ts.SeasonID - 1)))
	} else {
		nbaStandings = GetNBAStandingsBySeasonID(strconv.Itoa(int(ts.SeasonID)))
	}

	for idx, ns := range nbaStandings {
		rank := len(nbaStandings) - idx
		teamMap[ns.TeamID].AssignWaiverOrder(uint(rank))
	}

	for _, t := range nbaTeams {
		db.Save(&t)
	}
}

func SignFreeAgent(offer structs.NBAContractOffer, FreeAgent structs.NBAPlayer, ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	NBATeam := GetNBATeamByTeamID(strconv.Itoa(int(offer.TeamID)))
	FreeAgent.SignWithTeam(NBATeam.ID, NBATeam.Abbr)

	Contract := structs.NBAContract{
		PlayerID:       FreeAgent.PlayerID,
		TeamID:         NBATeam.ID,
		Team:           NBATeam.Abbr,
		OriginalTeamID: NBATeam.ID,
		OriginalTeam:   NBATeam.Abbr,
		YearsRemaining: offer.TotalYears,
		ContractType:   offer.ContractType,
		Year1Total:     offer.Year1Total,
		Year2Total:     offer.Year2Total,
		Year3Total:     offer.Year3Total,
		Year4Total:     offer.Year4Total,
		Year5Total:     offer.Year5Total,
		TotalRemaining: offer.TotalCost,
		ContractValue:  offer.ContractValue,
		Year1Opt:       offer.Year1Opt,
		Year2Opt:       offer.Year2Opt,
		Year3Opt:       offer.Year3Opt,
		Year4Opt:       offer.Year4Opt,
		Year5Opt:       offer.Year5Opt,
		IsActive:       true,
		IsComplete:     false,
		IsExtended:     false,
	}

	db.Create(&Contract)
	db.Save(&FreeAgent)

	// News Log
	message := "FA " + FreeAgent.Position + " " + FreeAgent.FirstName + " " + FreeAgent.LastName + " has signed with the " + NBATeam.Team + " " + NBATeam.Nickname + " with a contract worth approximately $" + strconv.Itoa(int(Contract.TotalRemaining)) + " Million Dollars."
	CreateNewsLog("NBA", message, "Free Agency", 0, ts)
}

func SyncFreeAgencyOffers() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	ts.ToggleFALock()
	db.Save(&ts)

	// Sync Free Agents and their contract offers
	FreeAgents := GetAllFreeAgents()
	faSyncFreeAgents(FreeAgents, ts, db)

	// Iterate through WaiverWire players
	waiverWirePlayers := GetAllWaiverWirePlayers()
	faSyncWaiverWirePlayers(waiverWirePlayers, ts, db)

	// Iterate through GLeague Players
	gLeaguePlayers := GetAllGLeaguePlayersForFA()
	faSyncGLeaguePlayers(gLeaguePlayers, ts, db)

	// Iterate through International Players
	islPlayers := GetAllISLPlayersForFA()
	faSyncISLPlayers(islPlayers, ts, db)

	ts.ToggleFALock()
	ts.ToggleGMActions()
	db.Save(&ts)
}

func GetFreeAgentOffersByPlayerID(PlayerID string) []structs.NBAContractOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAContractOffer{}

	err := db.Where("player_id = ? AND is_active = ?", PlayerID, true).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func GetWaiverOffersByTeamID(teamID string) []structs.NBAWaiverOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAWaiverOffer{}

	err := db.Where("team_id = ?", teamID).Find(&offers).Error
	if err != nil {
		return offers
	}

	return offers
}

func TempExtensionAlgorithm() {
	db := dbprovider.GetInstance().GetDB()
	// DB
	path := secrets.GetPath()["extensions"]
	extensionsCSV := util.ReadCSV(path)
	ts := GetTimestamp()
	// Read CSV
	for idx, row := range extensionsCSV {
		if idx == 0 {
			continue
		}

		id := row[1]
		teamID := row[0]
		playerRecord := GetNBAPlayerRecord(id)
		team := GetNBATeamByTeamID(teamID)
		contractType := row[3]
		contractLength := util.ConvertStringToInt(row[4])
		totalValue := util.ConvertStringToFloat(row[5])
		year1 := util.ConvertStringToFloat(row[6])
		year2 := util.ConvertStringToFloat(row[7])
		year3 := util.ConvertStringToFloat(row[8])
		year4 := util.ConvertStringToFloat(row[9])
		year5 := util.ConvertStringToFloat(row[10])
		year1Opt := util.ConvertStringToBool(row[11])
		year2Opt := util.ConvertStringToBool(row[12])
		year3Opt := util.ConvertStringToBool(row[13])
		year4Opt := util.ConvertStringToBool(row[14])
		year5Opt := util.ConvertStringToBool(row[15])
		nbaContract := structs.NBAContract{
			PlayerID:       playerRecord.ID,
			TeamID:         team.ID,
			ContractType:   contractType,
			YearsRemaining: uint(contractLength),
			TotalRemaining: totalValue,
			Year1Total:     year1,
			Year2Total:     year2,
			Year3Total:     year3,
			Year4Total:     year4,
			Year5Total:     year5,
			Year1Opt:       year1Opt,
			Year2Opt:       year2Opt,
			Year3Opt:       year3Opt,
			Year4Opt:       year4Opt,
			Year5Opt:       year5Opt,
			IsExtended:     true,
			IsActive:       true,
		}
		minimumValue := playerRecord.MinimumValue
		contractStatus := ""
		if playerRecord.MaxRequested {
			contractStatus = "Max"
		}
		if playerRecord.IsSuperMaxQualified {
			contractStatus = "SuperMax"
		}
		multiplier := 1.0
		validation := validateFreeAgencyPref(playerRecord, team, strconv.Itoa(int(ts.SeasonID)), idx)
		if validation && playerRecord.FreeAgency != "Average" {
			multiplier = 0.85
		} else if !validation && playerRecord.FreeAgency != "Average" {
			multiplier = 1.15
		}
		if playerRecord.FreeAgency == "Highest bidder" {
			multiplier = 1
		}
		minimumValue = minimumValue * multiplier
		validOffer := validateContract(nbaContract, contractStatus, minimumValue)

		if !validOffer {
			message := playerRecord.Position + " " + playerRecord.FirstName + " " + playerRecord.LastName + " has rejected an extension offer from " + team.Team + " " + team.Nickname
			CreateNewsLog("NBA", message, "Contract", int(team.ID), ts)
			continue
		}

		message := playerRecord.Position + " " + playerRecord.FirstName + " " + playerRecord.LastName + " has signed an extension with the " + team.Team + " " + team.Nickname + ", worth approximately $" + strconv.Itoa(int(nbaContract.TotalRemaining)) + " Million!"
		CreateNewsLog("NBA", message, "Contract", 0, ts)
		playerRecord.SignWithTeam(team.ID, team.Team)
		db.Save(&playerRecord)
		db.Create(&nbaContract)
	}
	// Iterate through submissions
	// Player Record by ID
	// Get Minimum Value required
	// Check if max/supermax
	// Check FA Preference
	// Compare contract with FA Preference with minimum value
	// If met, player signs
	// If not, continue algorithm
}

func validateFreeAgencyPref(playerRecord structs.NBAPlayer, team structs.NBATeam, seasonID string, idx int) bool {
	preference := playerRecord.FreeAgency

	if preference == "Average" {
		return true
	}
	if preference == "Drafted team discount" && playerRecord.DraftedTeamID == team.ID {
		return true
	}
	if preference == "Loyal" && playerRecord.PreviousTeamID == team.ID {
		return true
	}

	if preference == "Hometown Hero" && playerRecord.State == team.State {
		return true
	}
	if preference == "Adversarial" && playerRecord.PreviousTeamID != team.ID && playerRecord.DraftedTeamID != team.ID {
		return true
	}

	if preference == "I'm the starter" {
		teamRoster := GetAllNBAPlayersByTeamID(strconv.Itoa(int(team.ID)))
		sort.Slice(teamRoster, func(i, j int) bool {
			return teamRoster[i].Overall > teamRoster[j].Overall
		})
		for idx, p := range teamRoster {
			if idx > 4 {
				return false
			}
			if playerRecord.Overall >= p.Overall {
				return true
			}
		}
	}
	if preference == "Market-driven" && checkMarketCity(team.City) {
		return true
	}
	if preference == "Money motivated" {
		return false
	}
	if preference == "Highest bidder" && idx == 0 {
		return true
	}
	if preference == "Championship seeking" {
		standings := GetNBAStandingsRecordByTeamID(strconv.Itoa(int(team.ID)), seasonID)
		if standings.TotalWins > standings.TotalLosses {
			return true
		}
	}
	return false
}

func validateContract(offer structs.NBAContract, status string, minimum float64) bool {
	if status == "Max" || status == "SuperMax" {
		// if offer.YearsRemaining == 5 {
		// 	return minimum < offer.Year1Total && minimum < offer.Year2Total && minimum < offer.Year3Total && minimum < offer.Year4Total && minimum < offer.Year5Total
		// } else if offer.YearsRemaining == 4 {
		// 	return minimum < offer.Year1Total && minimum < offer.Year2Total && minimum < offer.Year3Total && minimum < offer.Year4Total
		// } else if offer.YearsRemaining == 3 {
		// 	return minimum < offer.Year1Total && minimum < offer.Year2Total && minimum < offer.Year3Total
		// } else if offer.YearsRemaining == 2 {
		// 	return minimum < offer.Year1Total && minimum < offer.Year2Total
		// }
		return minimum <= offer.Year1Total
	}
	return minimum <= offer.TotalRemaining
}

func validateOffer(offer structs.NBAContractOffer, status string, minimum float64) bool {
	if status == "Max" || status == "SuperMax" {
		if offer.TotalYears == 5 {
			return minimum < offer.Year1Total && minimum < offer.Year2Total && minimum < offer.Year3Total && minimum < offer.Year4Total && minimum < offer.Year5Total
		} else if offer.TotalYears == 4 {
			return minimum < offer.Year1Total && minimum < offer.Year2Total && minimum < offer.Year3Total && minimum < offer.Year4Total
		} else if offer.TotalYears == 3 {
			return minimum < offer.Year1Total && minimum < offer.Year2Total && minimum < offer.Year3Total
		} else if offer.TotalYears == 2 {
			return minimum < offer.Year1Total && minimum < offer.Year2Total
		}
		return minimum <= offer.Year1Total
	}
	return minimum <= offer.TotalCost
}

func checkMarketCity(city string) bool {
	return city == "Los Angeles" || city == "New York" || city == "Brooklyn" || city == "Chicago" || city == "Philadelphia" || city == "Boston" || city == "Dallas" || city == "Oakland" || city == "Atlanta" || city == "Houston" || city == "Washington"
}

func faSyncFreeAgents(freeAgents []structs.NBAPlayer, ts structs.Timestamp, db *gorm.DB) {
	seasonID := strconv.Itoa(int(ts.SeasonID))
	for _, FA := range freeAgents {

		// Check if still accepting offers
		if ts.IsNBAOffseason && FA.IsAcceptingOffers && ts.FreeAgencyRound < FA.NegotiationRound {
			continue
		}

		// If the Player is ready to negotiate, toggle the Player as IsNegotiating, save the record, and continue to the next player
		if ts.IsNBAOffseason && FA.IsAcceptingOffers && ts.FreeAgencyRound >= FA.NegotiationRound {
			FA.ToggleIsNegotiating()
			db.Save(&FA)
			continue
		}

		// Check if still negotiation
		if ts.IsNBAOffseason && FA.IsNegotiating && ts.FreeAgencyRound < FA.SigningRound {
			continue
		}

		// Is Ready to Sign, Get All Offers on the Free Agent
		Offers := GetFreeAgentOffersByPlayerID(strconv.Itoa(int(FA.ID)))

		// Sort by highest contract value
		sort.Sort(structs.ByContractValue(Offers))

		WinningOffer := &structs.NBAContractOffer{}
		minimumValue := FA.MinimumValue
		// Logic to confirm if the Free Agent is requesting a Max contract or SuperMax contract
		contractStatus := ""
		if FA.MaxRequested {
			contractStatus = "Max"
		}
		if FA.IsSuperMaxQualified {
			contractStatus = "SuperMax"
		}
		for idx, Offer := range Offers {
			minimumValueMultiplier := 1.0
			team := GetNBATeamByTeamID(strconv.Itoa(int(Offer.TeamID)))
			validation := validateFreeAgencyPref(FA, team, seasonID, idx)
			// If the offer is valid and meets the player's free agency bias, reduce the minimum value required by 15%
			if validation && FA.FreeAgency != "Average" && FA.Year > 2 {
				minimumValueMultiplier = 0.85
				// If the offer does not meet the player's free agency bias, increase the minimum value required by 15%
			} else if !validation && FA.FreeAgency != "Average" && FA.Year > 2 {
				minimumValueMultiplier = 1.15
			}
			minimumValue = minimumValue * minimumValueMultiplier
			validOffer := validateOffer(Offer, contractStatus, minimumValue)

			// Get the Contract with the best value for the FA
			if Offer.IsActive && WinningOffer.ID == 0 && validOffer {
				*WinningOffer = Offer
			}

			// If the offer being iterated through ISN'T the winning offer, cancel the offer.
			if Offer.IsActive && WinningOffer.ID != 0 && WinningOffer.ID != Offer.ID {
				Offer.CancelOffer()
				db.Save(&Offer)
			}
		}

		// If there is a winning offer, sign the player
		if WinningOffer.ID > 0 {
			SignFreeAgent(*WinningOffer, FA, ts)
		}
	}
}

func faSyncWaiverWirePlayers(waiverWirePlayers []structs.NBAPlayer, ts structs.Timestamp, db *gorm.DB) {
	for _, w := range waiverWirePlayers {

		waiverWireID := strconv.Itoa(int(w.ID))

		waiverOffers := GetWaiverOffersByPlayerID(waiverWireID)
		if len(waiverOffers) == 0 {
			// Deactivate Contract, convert to Free Agent
			w.ConvertWaivedPlayerToFA()
			contract := GetContractByPlayerID(waiverWireID)
			contract.DeactivateContract()
			db.Save(&contract)
		} else {
			winningOffer := waiverOffers[0]
			winningOfferTeamID := strconv.Itoa(int(winningOffer.TeamID))
			w.SignWithTeam(winningOffer.TeamID, winningOffer.Team)

			contract := GetNBAContractsByPlayerID(waiverWireID)
			contract.TradePlayer(winningOffer.TeamID, winningOffer.Team)
			contract.MakeContractActive()

			db.Save(&contract)

			message := w.Position + " " + w.FirstName + " " + w.LastName + " was picked up on the Waiver Wire by " + winningOffer.Team
			CreateNewsLog("NBA", message, "Free Agency", int(winningOffer.TeamID), ts)

			// Recalibrate winning team's remaining offers
			teamOffers := GetWaiverOffersByTeamID(winningOfferTeamID)
			team := GetNBATeamByTeamID(winningOfferTeamID)

			team.AssignWaiverOrder(team.WaiverOrder + 32)
			db.Save(&team)

			for _, o := range teamOffers {
				o.AssignNewWaiverOrder(team.WaiverOrder + 32)
				db.Save(&o)
			}

			// Delete current waiver offers
			for _, o := range waiverOffers {
				db.Delete(&o)
			}
		}
		db.Save(&w)
	}
}

func faSyncGLeaguePlayers(gLeaguePlayers []structs.NBAPlayer, ts structs.Timestamp, db *gorm.DB) {
	for _, g := range gLeaguePlayers {
		gLeaguePlayerID := strconv.Itoa(int(g.ID))
		Offers := GetWaiverOffersByPlayerID(gLeaguePlayerID)

		if len(Offers) == 0 {
			continue
		}
		ownerTeam := g.TeamID
		ownerOffer := structs.NBAWaiverOffer{}

		for _, o := range Offers {
			if o.TeamID == ownerTeam && o.IsActive {
				ownerOffer = o
				break
			}
		}
		if ownerOffer.ID > 0 {
			g.SignWithTeam(ownerOffer.TeamID, ownerOffer.Team)
			contract := GetNBAContractsByPlayerID(gLeaguePlayerID)
			contract.TradePlayer(ownerOffer.TeamID, ownerOffer.Team)
			db.Save(&contract)
			message := g.Position + " " + g.FirstName + " " + g.LastName + " was picked up from the GLeague by " + ownerOffer.Team
			CreateNewsLog("NBA", message, "Free Agency", int(ownerOffer.TeamID), ts)

			db.Save(&g)
		} else {
			sort.Slice(Offers, func(i, j int) bool {
				return Offers[i].WaiverOrder < Offers[j].WaiverOrder
			})

			WinningOffer := structs.NBAWaiverOffer{}
			for _, Offer := range Offers {
				// Get the Contract with the best value for the FA
				if Offer.IsActive && WinningOffer.ID == 0 {
					WinningOffer = Offer
				}
				if Offer.IsActive {
					Offer.DeactivateWaiverOffer()
				}

				db.Save(&Offer)

				if WinningOffer.ID > 0 {
					g.SignWithTeam(WinningOffer.TeamID, WinningOffer.Team)
					contract := GetNBAContractsByPlayerID(gLeaguePlayerID)
					contract.TradePlayer(WinningOffer.TeamID, WinningOffer.Team)
					db.Save(&contract)
					message := g.Position + " " + g.FirstName + " " + g.LastName + " was picked up from the GLeague by " + WinningOffer.Team
					CreateNewsLog("NBA", message, "Free Agency", int(WinningOffer.TeamID), ts)
					db.Save(&g)
				} else if ts.IsNBAOffseason {
					g.WaitUntilStartOfSeason()
					db.Save(&g)
				}
			}
		}

		for _, o := range Offers {
			db.Delete(&o)
		}
	}
}

func faSyncISLPlayers(islPlayers []structs.NBAPlayer, ts structs.Timestamp, db *gorm.DB) {
	for _, i := range islPlayers {
		islPlayerID := strconv.Itoa(int(i.ID))
		Offers := GetWaiverOffersByPlayerID(islPlayerID)

		if len(Offers) == 0 {
			continue
		}
		ownerTeam := i.TeamID
		ownerOffer := structs.NBAWaiverOffer{}

		for _, o := range Offers {
			if o.TeamID == ownerTeam && o.IsActive {
				ownerOffer = o
				break
			}
		}

		contract := GetNBAContractsByPlayerID(islPlayerID)
		contract.TradePlayer(ownerOffer.TeamID, ownerOffer.Team)
		db.Save(&contract)
		i.SignWithTeam(ownerOffer.TeamID, ownerOffer.Team)
		message := i.Position + " " + i.FirstName + " " + i.LastName + " was picked up from the GLeague by " + ownerOffer.Team
		CreateNewsLog("NBA", message, "Free Agency", int(ownerOffer.TeamID), ts)

		db.Save(&i)

		if ownerOffer.ID > 0 {
			contract := GetNBAContractsByPlayerID(islPlayerID)
			contract.TradePlayer(ownerOffer.TeamID, ownerOffer.Team)
			db.Save(&contract)
			i.SignWithTeam(ownerOffer.TeamID, ownerOffer.Team)
			message := i.Position + " " + i.FirstName + " " + i.LastName + " was picked up from the GLeague by " + ownerOffer.Team
			CreateNewsLog("NBA", message, "Free Agency", int(ownerOffer.TeamID), ts)
		} else {
			sort.Slice(Offers, func(i, j int) bool {
				return Offers[i].WaiverOrder < Offers[j].WaiverOrder
			})

			WinningOffer := structs.NBAWaiverOffer{}
			for _, Offer := range Offers {
				// Get the Contract with the best value for the FA
				if Offer.IsActive && WinningOffer.ID == 0 {
					WinningOffer = Offer
				}
				if Offer.IsActive {
					Offer.DeactivateWaiverOffer()
				}

				db.Save(&Offer)

				if WinningOffer.ID > 0 {
					contract := GetNBAContractsByPlayerID(islPlayerID)
					contract.TradePlayer(WinningOffer.TeamID, WinningOffer.Team)
					db.Save(&contract)
					i.SignWithTeam(WinningOffer.TeamID, WinningOffer.Team)
					message := i.Position + " " + i.FirstName + " " + i.LastName + " was picked up from the GLeague by " + WinningOffer.Team
					CreateNewsLog("NBA", message, "Free Agency", int(WinningOffer.TeamID), ts)
				} else if ts.IsNBAOffseason {
					i.WaitUntilStartOfSeason()
					db.Save(&i)
				}
			}
		}

		for _, o := range Offers {
			db.Delete(&o)
		}
	}
}
