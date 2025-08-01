package managers

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

func GetAllAvailableNBAPlayers(TeamID string) structs.FreeAgencyResponse {
	ts := GetTimestamp()
	seasonID := ts.SeasonID

	if ts.IsNBAOffseason {
		seasonID = ts.SeasonID - 1
	}
	seasonIDStr := strconv.Itoa(int(seasonID))
	FAs := GetAllFreeAgentsWithOffers(seasonIDStr)
	waiverPlayers := GetAllWaiverWirePlayers(seasonIDStr)
	gLeagePlayer := GetAllGLeaguePlayersForFA(seasonIDStr)
	islPlayers := GetAllISLPlayersForFA(seasonIDStr)
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

func GetAllAvailableNBAPlayersViaChan(TeamID string, ch chan<- structs.FreeAgencyResponse) {
	ts := GetTimestamp()
	seasonID := ts.SeasonID

	if ts.IsNBAOffseason {
		seasonID = ts.SeasonID - 1
	}
	seasonIDStr := strconv.Itoa(int(seasonID))
	FAs := GetAllFreeAgentsWithOffers(seasonIDStr)
	waiverPlayers := GetAllWaiverWirePlayers(seasonIDStr)
	gLeagePlayer := GetAllGLeaguePlayersForFA(seasonIDStr)
	islPlayers := GetAllISLPlayersForFA(seasonIDStr)
	Offers := GetFreeAgentOffersByTeamID(TeamID)
	roster := GetAllNBAPlayersByTeamID(TeamID)
	count := 0
	for _, p := range roster {
		if p.IsGLeague {
			continue
		}
		count += 1
	}

	ch <- structs.FreeAgencyResponse{
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

	db.Order("overall desc").Where("is_free_agent = ?", true).Find(&fas)

	return fas
}

// GetAllFreeAgentsWithOffers -- For Free Agency UI Page.
func GetAllFreeAgentsWithOffers(seasonID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	fas := []structs.NBAPlayer{}

	db.Preload("Offers", func(db *gorm.DB) *gorm.DB {
		return db.Order("contract_value DESC").Where("is_active = true")
	}).Preload("SeasonStats", "season_id = ?", seasonID).Order("overall desc").Where("is_free_agent = ?", true).Find(&fas)

	return fas
}

func GetAllWaiverWirePlayers(seasonID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	WaivedPlayers := []structs.NBAPlayer{}

	db.Preload("WaiverOffers").Preload("Contract").Preload("SeasonStats", "season_id = ?", seasonID).Where("is_waived = ?", true).Find(&WaivedPlayers)

	return WaivedPlayers
}

func GetAllGLeaguePlayersForFA(seasonID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	gLeaguePlayers := []structs.NBAPlayer{}

	db.Preload("WaiverOffers").Preload("Contract").Preload("SeasonStats", "season_id = ?", seasonID).Where("is_g_league = ?", true).Find(&gLeaguePlayers)

	return gLeaguePlayers
}

func GetAllISLPlayersForFA(seasonID string) []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	islPlayers := []structs.NBAPlayer{}

	db.Preload("WaiverOffers").Preload("Contract").Preload("SeasonStats", "season_id = ?", seasonID).Where("team_id > 32").Find(&islPlayers)

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

	if nbaPlayer.IsGLeague && nbaPlayer.TeamID != offer.TeamID {
		message := "Breaking News! " + nbaPlayer.FirstName + " " + nbaPlayer.LastName + " has received an offer from the " + offer.Team + "!"
		CreateNotification("CBB", message, "Offer", nbaPlayer.TeamID)
	} else if nbaPlayer.IsGLeague && nbaPlayer.TeamID == offer.TeamID {
		// Sign player back to team
		nbaPlayer.ToggleGLeague()
		repository.SaveNBAPlayerRecord(nbaPlayer, db)

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

func CreateExtensionOffer(offer structs.NBAContractOfferDTO) structs.NBAExtensionOffer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	extensionOffer := GetExtensionOfferByOfferID(strconv.Itoa(int(offer.ID)))
	player := GetNBAPlayerRecord(strconv.Itoa(int(offer.PlayerID)))

	extensionOffer.CalculateOffer(offer)

	// If the owning team is sending an offer to a player
	if extensionOffer.ID == 0 {
		id := GetLatestExtensionOfferInDB(db)
		extensionOffer.AssignID(id)
		db.Create(&extensionOffer)
		fmt.Println("Creating Extension Offer!")

		message := offer.Team + " have offered a " + strconv.Itoa(int(offer.TotalYears)) + " year contract extension for " + player.Position + " " + player.FirstName + " " + player.LastName + "."
		CreateNewsLog("NFL", message, "Free Agency", int(player.TeamID), ts)
	} else if extensionOffer.IsActive {
		fmt.Println("Updating Extension Offer!")
		db.Save(&extensionOffer)
	}

	return extensionOffer
}

func CancelExtensionOffer(offer structs.NBAContractOfferDTO) {
	db := dbprovider.GetInstance().GetDB()

	OfferID := strconv.Itoa(int(offer.ID))

	freeAgentOffer := GetExtensionOfferByOfferID(OfferID)

	freeAgentOffer.CancelOffer()

	db.Save(&freeAgentOffer)
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

func GetExtensionOffersByPlayerID(PlayerID string) []structs.NBAExtensionOffer {
	db := dbprovider.GetInstance().GetDB()

	offer := []structs.NBAExtensionOffer{}

	err := db.Where("nba_player_id = ?", PlayerID).Find(&offer).Error
	if err != nil {
		return offer
	}

	return offer
}

func GetExtensionOfferByOfferID(OfferID string) structs.NBAExtensionOffer {
	db := dbprovider.GetInstance().GetDB()

	offer := structs.NBAExtensionOffer{}

	err := db.Where("id = ?", OfferID).Find(&offer).Error
	if err != nil {
		return offer
	}

	return offer
}

func GetLatestExtensionOfferInDB(db *gorm.DB) uint {
	var latestOffer structs.NBAExtensionOffer

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
	newMinimumValue := offer.ContractValue * (float64(FreeAgent.Age) / 30)
	FreeAgent.SignWithTeam(NBATeam.ID, NBATeam.Abbr, true, newMinimumValue)

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
	repository.SaveProfessionalPlayerRecord(FreeAgent, db)

	// News Log
	message := "FA " + FreeAgent.Position + " " + FreeAgent.FirstName + " " + FreeAgent.LastName + " has signed with the " + NBATeam.Team + " " + NBATeam.Nickname + " with a contract worth approximately $" + strconv.Itoa(int(Contract.TotalRemaining)) + " Million Dollars."
	CreateNewsLog("NBA", message, "Free Agency", 0, ts)
}

func SyncFreeAgencyOffers() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	// NO FA BEFORE AND DURING DRAFT
	if ts.IsDraftTime {
		return
	}

	ts.ToggleFALock()
	seasonID := ts.SeasonID
	if ts.IsNBAOffseason {
		seasonID = seasonID - 1
	}
	repository.SaveTimeStamp(ts, db)
	seasonIDStr := strconv.Itoa(int(seasonID))
	// Sync Free Agents and their contract offers
	FreeAgents := GetAllFreeAgents()
	faSyncFreeAgents(FreeAgents, ts, db)

	// Iterate through WaiverWire players
	waiverWirePlayers := GetAllWaiverWirePlayers(seasonIDStr)
	faSyncWaiverWirePlayers(waiverWirePlayers, ts, db)

	// Iterate through GLeague Players
	gLeaguePlayers := GetAllGLeaguePlayersForFA(seasonIDStr)
	faSyncGLeaguePlayers(gLeaguePlayers, ts, db)

	// Iterate through International Players
	islPlayers := GetAllISLPlayersForFA(seasonIDStr)
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

func RunExtensionsAlgorithm() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	nbaTeams := GetAllActiveNBATeams()

	for _, team := range nbaTeams {
		if team.ID > 32 {
			break
		}
		teamID := strconv.Itoa(int(team.ID))
		roster := GetNBAPlayersWithContractsAndExtensionsByTeamID(teamID)

		for _, player := range roster {
			min := player.MinimumValue
			contract := player.Contract
			if contract.ID == 0 {
				// Yeah this is an error
				continue
			}
			if contract.YearsRemaining == 1 && len(player.Extensions) > 0 {
				for idx, e := range player.Extensions {
					if e.IsRejected || !e.IsActive || e.IsAccepted {
						continue
					}
					minimumValueMultiplier := 1.0
					validation := validateFreeAgencyPref(player, roster, team, seasonID, int(e.TotalYears), idx)
					// If the offer is valid and meets the player's free agency bias, reduce the minimum value required by 15%
					if validation && player.FreeAgency != "Average" {
						minimumValueMultiplier = 0.85
						// If the offer does not meet the player's free agency bias, increase the minimum value required by 15%
					} else if !validation && player.FreeAgency != "Average" {
						minimumValueMultiplier = 1.15
					}
					percentage := (e.ContractValue / (min * minimumValueMultiplier) * 100)
					odds := getExtensionPercentageOdds(percentage)
					// Run Check on the Extension

					roll := util.GenerateIntFromRange(1, 100)
					message := ""
					if odds == 0 || float64(roll) > odds {
						// Rejects offer
						e.DeclineOffer(ts.NBAWeek)
						player.DeclineOffer(ts.NBAWeek)
						if e.IsRejected || player.Rejections > 2 {
							message = player.Position + " " + player.FirstName + " " + player.LastName + " has rejected an extension offer from " + e.Team + " worth approximately $" + strconv.Itoa(int(e.ContractValue)) + " Million Dollars and will enter Free Agency."
						} else {
							message = player.Position + " " + player.FirstName + " " + player.LastName + " has declined an extension offer from " + e.Team + " with an extension worth approximately $" + strconv.Itoa(int(e.ContractValue)) + " Million Dollars, and is still negotiating."
						}
						CreateNewsLog("NBA", message, "Free Agency", int(e.TeamID), ts)
						db.Save(&player)
					} else {
						e.AcceptOffer()
						message = player.Position + " " + player.FirstName + " " + player.LastName + " has accepted an extension offer from " + e.Team + " worth approximately $" + strconv.Itoa(int(e.ContractValue)) + " Million Dollars."
						CreateNewsLog("NBA", message, "Free Agency", int(e.TeamID), ts)
						db.Save(&team)
					}
					db.Save(&e)
				}
			}
		}
	}

	ts.MoveUpFreeAgencyRound()
	repository.SaveTimeStamp(ts, db)
}

func GetContractMap() map[uint]structs.NBAContract {
	contracts := repository.FindAllProContracts(true)
	return MakeContractMap(contracts)
}

func GetExtensionMap() map[uint]structs.NBAExtensionOffer {
	extensions := repository.FindAllProExtensions(true)
	return MakeExtensionMap(extensions)
}

func validateFreeAgencyPref(playerRecord structs.NBAPlayer, roster []structs.NBAPlayer, team structs.NBATeam, seasonID string, offerLength, idx int) bool {
	preference := playerRecord.FreeAgency

	if preference == "Average" || preference == "Will eventually play for LeBron" {
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
		teamRoster := roster
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
	if preference == "Market-driven" && offerLength < 3 {
		return true
	}
	if preference == "Wants Extension" && offerLength > 2 {
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
	rosterMap := GetFullRosterNBAMap()
	capsheetMap := GetPointedCapsheetMap()
	nbaTeams := GetAllActiveNBATeams()
	nbaTeamMap := MakeNBATeamMap(nbaTeams)
	for _, FA := range freeAgents {
		// Check if still accepting offers
		if ts.IsNBAOffseason && !FA.IsAcceptingOffers {
			continue
		}

		// Is Ready to Sign, Get All Offers on the Free Agent
		Offers := GetFreeAgentOffersByPlayerID(strconv.Itoa(int(FA.ID)))

		if len(Offers) == 0 {
			continue
		}
		maxDay := 1000
		for _, offer := range Offers {
			if maxDay > int(offer.Syncs) {
				maxDay = int(offer.Syncs)
			}
		}
		if maxDay < 3 {
			for _, offer := range Offers {
				offer.IncrementSyncs()
				repository.SaveContractOfferRecord(offer, db)
			}
		} else {
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
				roster := rosterMap[Offer.TeamID]
				validation := validateFreeAgencyPref(FA, roster, team, seasonID, int(Offer.TotalYears), idx)
				// If the offer is valid and meets the player's free agency bias, reduce the minimum value required by 15%
				if validation && FA.FreeAgency != "Average" && FA.Year > 2 {
					minimumValueMultiplier = 0.85
					// If the offer does not meet the player's free agency bias, increase the minimum value required by 15%
				} else if !validation && FA.FreeAgency != "Average" && FA.Year > 2 {
					minimumValueMultiplier = 1.15
				}
				minimumValue = minimumValue * minimumValueMultiplier
				validOffer := validateOffer(Offer, contractStatus, minimumValue)
				teamCapsheet := capsheetMap[Offer.TeamID]
				// Check to make sure that
				belowCap := Offer.Year1Total+teamCapsheet.Year1Total+teamCapsheet.Year1Cap+teamCapsheet.Year1CashTransferred-teamCapsheet.Year1CashReceived <= ts.Y1Capspace
				nbaTeam := nbaTeamMap[Offer.TeamID]
				isAi := nbaTeam.NBAOwnerName == "" || nbaTeam.NBAOwnerName == "AI"
				withinRosterLimit := len(roster) < 18
				// Get the Contract with the best value for the FA
				if Offer.IsActive && WinningOffer.ID == 0 && validOffer && belowCap {
					if isAi && !withinRosterLimit {
						Offer.CancelOffer()
					} else {
						*WinningOffer = Offer
					}
				}

				// If the offer being iterated through ISN'T the winning offer, cancel the offer.
				if Offer.IsActive && WinningOffer.ID != 0 && WinningOffer.ID != Offer.ID {
					Offer.CancelOffer()
				}
				repository.SaveContractOfferRecord(Offer, db)
			}

			// If there is a winning offer, sign the player
			if WinningOffer.ID > 0 {
				SignFreeAgent(*WinningOffer, FA, ts)
				teamCapsheet := capsheetMap[WinningOffer.TeamID]
				newTotal1 := teamCapsheet.Year1Total + WinningOffer.Year1Total
				newTotal2 := teamCapsheet.Year2Total + WinningOffer.Year2Total
				newTotal3 := teamCapsheet.Year3Total + WinningOffer.Year3Total
				newTotal4 := teamCapsheet.Year4Total + WinningOffer.Year4Total
				newTotal5 := teamCapsheet.Year5Total + WinningOffer.Year5Total
				teamCapsheet.SyncTotals(newTotal1, newTotal2, newTotal3, newTotal4, newTotal5)
			}
		}
	}
}

func faSyncWaiverWirePlayers(waiverWirePlayers []structs.NBAPlayer, ts structs.Timestamp, db *gorm.DB) {
	capsheets := GetCapsheetMap()
	for _, w := range waiverWirePlayers {

		waiverWireID := strconv.Itoa(int(w.ID))

		waiverOffers := GetWaiverOffersByPlayerID(waiverWireID)
		if len(waiverOffers) == 0 {
			// Deactivate Contract, convert to Free Agent
			contract := GetContractByPlayerID(waiverWireID)
			capsheet := capsheets[contract.TeamID]
			capsheet.CutPlayerFromCapsheet(contract)
			db.Save(&capsheet)
			contract.DeactivateContract()
			repository.SaveProfessionalContractRecord(contract, db)
			w.ConvertWaivedPlayerToFA()
		} else {
			winningOffer := waiverOffers[0]
			winningOfferTeamID := strconv.Itoa(int(winningOffer.TeamID))
			w.SignWithTeam(winningOffer.TeamID, winningOffer.Team, false, 0)

			contract := GetNBAContractsByPlayerID(waiverWireID)
			contract.TradePlayer(winningOffer.TeamID, winningOffer.Team)
			contract.MakeContractActive()

			repository.SaveProfessionalContractRecord(contract, db)

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
		repository.SaveProfessionalPlayerRecord(w, db)
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
			g.SignWithTeam(ownerOffer.TeamID, ownerOffer.Team, false, 0)
			contract := GetNBAContractsByPlayerID(gLeaguePlayerID)
			contract.TradePlayer(ownerOffer.TeamID, ownerOffer.Team)
			repository.SaveProfessionalContractRecord(contract, db)
			message := g.Position + " " + g.FirstName + " " + g.LastName + " was picked up from the GLeague by " + ownerOffer.Team
			CreateNewsLog("NBA", message, "Free Agency", int(ownerOffer.TeamID), ts)

			repository.SaveProfessionalPlayerRecord(g, db)
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
					g.SignWithTeam(WinningOffer.TeamID, WinningOffer.Team, false, 0)
					contract := GetNBAContractsByPlayerID(gLeaguePlayerID)
					contract.TradePlayer(WinningOffer.TeamID, WinningOffer.Team)
					repository.SaveProfessionalContractRecord(contract, db)
					message := g.Position + " " + g.FirstName + " " + g.LastName + " was picked up from the GLeague by " + WinningOffer.Team
					CreateNewsLog("NBA", message, "Free Agency", int(WinningOffer.TeamID), ts)
					repository.SaveProfessionalPlayerRecord(g, db)
				} else if ts.IsNBAOffseason {
					g.WaitUntilStartOfSeason()
					repository.SaveProfessionalPlayerRecord(g, db)
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
		repository.SaveProfessionalContractRecord(contract, db)
		i.SignWithTeam(contract.TeamID, contract.Team, false, 0)
		message := i.Position + " " + i.FirstName + " " + i.LastName + " was picked up from the GLeague by " + ownerOffer.Team
		CreateNewsLog("NBA", message, "Free Agency", int(ownerOffer.TeamID), ts)

		repository.SaveProfessionalPlayerRecord(i, db)

		if ownerOffer.ID > 0 {
			contract := GetNBAContractsByPlayerID(islPlayerID)
			contract.TradePlayer(ownerOffer.TeamID, ownerOffer.Team)
			repository.SaveProfessionalContractRecord(contract, db)
			i.SignWithTeam(ownerOffer.TeamID, ownerOffer.Team, false, 0)
			message := i.Position + " " + i.FirstName + " " + i.LastName + " was picked up from the ISL by " + ownerOffer.Team
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
					repository.SaveProfessionalContractRecord(contract, db)
					i.SignWithTeam(WinningOffer.TeamID, WinningOffer.Team, false, 0)
					message := i.Position + " " + i.FirstName + " " + i.LastName + " was picked up from the GLeague by " + WinningOffer.Team
					CreateNewsLog("NBA", message, "Free Agency", int(WinningOffer.TeamID), ts)
				} else if ts.IsNBAOffseason {
					i.WaitUntilStartOfSeason()
					repository.SaveProfessionalPlayerRecord(i, db)
				}
			}
		}

		for _, o := range Offers {
			db.Delete(&o)
		}
	}
}

func getExtensionPercentageOdds(percentage float64) float64 {
	if percentage >= 100 {
		return 100
	} else if percentage >= 90 {
		return 75
	} else if percentage >= 80 {
		return 50
	} else if percentage >= 70 {
		return 25
	}
	return 0
}

func GetAllFreeAgencyOffers() []structs.NBAContractOffer {
	return repository.FindAllFreeAgentOffers(repository.FreeAgencyQuery{
		IsActive: true,
	})
}

func SyncAIOffers() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	teams := GetAllActiveNBATeams()

	offers := GetAllFreeAgencyOffers()
	offerMapByTeamID := MakeFreeAgencyOfferMapByTeamID(offers)
	offerMapByPlayerID := MakeFreeAgencyOfferMap(offers)
	freeAgents := GetAllFreeAgents()
	players := GetAllNBAPlayers()
	playerMap := MakeNBAPlayerMapByTeamID(players, true)
	capsheetMap := GetCapsheetMap()

	for _, team := range teams {
		if team.ID > 32 {
			continue
		}
		if len(team.NBAOwnerName) > 0 && team.NBAOwnerName != "AI" && team.NBAOwnerName != "" {
			continue
		}
		if team.NBAOwnerName == "" && len(team.NBAGMName) > 0 && team.NBAGMName != "AI" {
			continue
		}
		capsheet := capsheetMap[team.ID]
		offersByTeam := offerMapByTeamID[team.ID]
		if len(offersByTeam) > 7 {
			continue
		}
		freeAgentOfferMap := MakeFreeAgencyOfferMap(offersByTeam)
		roster := playerMap[team.ID]
		if len(roster) > 17 {
			continue
		}
		cCount := 0
		pfCount := 0
		sfCount := 0
		sgCount := 0
		pgCount := 0
		cBids := 0
		pfBids := 0
		sfBids := 0
		sgBids := 0
		pgBids := 0
		for _, p := range roster {
			if p.Position == "C" {
				cCount++
			} else if p.Position == "PF" {
				pfCount++
			} else if p.Position == "SF" {
				sfCount++
			} else if p.Position == "SG" {
				sgCount++
			} else {
				pgCount++
			}
		}

		// Iterate through FA list to get bids
		for _, fa := range freeAgents {
			existingOffers := freeAgentOfferMap[fa.ID]
			if len(existingOffers) > 0 {
				if fa.Position == "C" {
					cBids++
				} else if fa.Position == "PF" {
					pfBids++
				} else if fa.Position == "SF" {
					sfBids++
				} else if fa.Position == "SG" {
					sgBids++
				} else {
					pgBids++
				}
			}
		}

		for _, fa := range freeAgents {
			existingOffers := freeAgentOfferMap[fa.ID]
			if len(existingOffers) > 0 {
				continue
			}
			if fa.Position == "C" && (cCount > 3 || cBids > 2) {
				continue
			}
			if fa.Position == "PF" && (pfCount > 4 || pfBids > 3) {
				continue
			}
			if fa.Position == "SF" && (sfCount > 4 || sfBids > 3) {
				continue
			}
			if fa.Position == "SG" && (sgCount > 4 || sgBids > 3) {
				continue
			}
			if fa.Position == "PG" && (pgCount > 3 || pgBids > 2) {
				continue
			}
			coinFlip := util.GenerateIntFromRange(1, 2)
			if coinFlip == 2 {
				continue
			}

			existingCompetition := offerMapByPlayerID[fa.ID]
			if len(existingCompetition) > 4 {
				continue
			}
			maxPercentage := GetMaxPercentage(fa.Year, fa.MaxRequested, fa.IsSuperMaxQualified)
			minRequired := maxPercentage * ts.Y1Capspace
			if minRequired == 0 {
				minRequired = fa.MinimumValue
			}
			// Okay, now we found an open player. Send a bid.
			basePay := 1.0
			if fa.Age < 25 || fa.Overall < 80 {
				basePay = 0.7
			} else if fa.Overall > 79 {
				rangedPay := util.GenerateFloatFromRange(1, 3.5)
				if rangedPay < minRequired {
					rangedPay = util.GenerateFloatFromRange(minRequired, minRequired+3.5)
				}
				basePay = RoundToFixedDecimalPlace(rangedPay, 2)
			}

			yearsOnContract := 2
			if fa.Overall > 79 {
				yearsOnContract = 3
			} else {
				yearsOnContract = 1
			}
			y1 := basePay
			if capsheet.Year1Total+capsheet.Year1Cap+capsheet.Year1CashTransferred-capsheet.Year1CashReceived+y1 > ts.Y1Capspace {
				continue
			}
			y2 := 0.0
			y3 := 0.0
			if yearsOnContract > 2 {
				y3 = basePay
			}
			if yearsOnContract > 1 {
				y2 = basePay
			}
			if fa.Position == "C" {
				cBids++
			} else if fa.Position == "PF" {
				pfBids++
			} else if fa.Position == "SF" {
				sfBids++
			} else if fa.Position == "SG" {
				sgBids++
			} else {
				pgBids++
			}
			offer := structs.NBAContractOffer{
				Year1Total:    y1,
				Year2Total:    y2,
				Year3Total:    y3,
				TotalCost:     basePay * float64(yearsOnContract),
				ContractValue: basePay,
				IsActive:      true,
				PlayerID:      fa.ID,
				TeamID:        team.ID,
				TotalYears:    uint(yearsOnContract),
			}
			offerMapByPlayerID[fa.ID] = append(offerMapByPlayerID[fa.ID], offer)
			offerMapByTeamID[team.ID] = append(offerMapByTeamID[team.ID], offer)

			repository.SaveContractOfferRecord(offer, db)
		}
	}
}

func GetMaxPercentage(year int, maxRequested, isSuperMax bool) float64 {
	if isSuperMax {
		if year > 9 {
			return 0.35
		} else if year > 6 {
			return 0.3
		} else {
			return 0.25
		}
	} else if maxRequested {
		if year > 9 {
			return 0.3
		} else if year > 6 {
			return 0.25
		} else {
			return 0.2
		}
	}
	return 0
}
