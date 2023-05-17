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
	WaiverPlayers := GetAllWaiverWirePlayers()
	Offers := GetFreeAgentOffersByTeamID(TeamID)

	return structs.FreeAgencyResponse{
		FreeAgents:    FAs,
		WaiverPlayers: WaiverPlayers,
		TeamOffers:    Offers,
	}
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
	message := "FA " + FreeAgent.Position + " " + FreeAgent.FirstName + " " + FreeAgent.LastName + " has signed with the " + NBATeam.Team + " " + NBATeam.Nickname + " with a contract worth approximately $" + strconv.Itoa(int(Contract.ContractValue)) + " Million Dollars."
	newsLog := structs.NewsLog{
		WeekID:      ts.NBAWeekID,
		SeasonID:    ts.SeasonID,
		MessageType: "Free Agency",
		Message:     message,
		League:      "NBA",
	}

	db.Create(&newsLog)
}

func SyncFreeAgencyOffers() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	FreeAgents := GetAllNBAPlayersByTeamID("0")

	for _, FA := range FreeAgents {

		// Check if still accepting offers
		if ts.IsNBAOffseason && FA.IsAcceptingOffers && ts.FreeAgencyRound < FA.NegotiationRound {
			continue
		}

		if ts.IsNBAOffseason && FA.IsAcceptingOffers && ts.FreeAgencyRound >= FA.NegotiationRound {
			FA.ToggleIsNegotiating()
			db.Save(&FA)
			continue
		}

		// Check if still negotiation
		if ts.IsNBAOffseason && FA.IsNegotiating && ts.FreeAgencyRound < FA.SigningRound {
			continue
		}

		// Is Ready to Sign
		Offers := GetFreeAgentOffersByPlayerID(strconv.Itoa(int(FA.ID)))

		// Sort by highest contract value
		sort.Sort(structs.ByContractValue(Offers))

		WinningOffer := structs.NBAContractOffer{}

		for _, Offer := range Offers {
			// Get the Contract with the best value for the FA
			if Offer.IsActive && WinningOffer.ID == 0 {
				WinningOffer = Offer
			}

			if Offer.IsActive && WinningOffer.ID != 0 && WinningOffer.ID != Offer.ID {
				Offer.CancelOffer()
			}

			db.Save(&Offer)
		}

		SignFreeAgent(WinningOffer, FA, ts)
	}
}

func GetFreeAgentOffersByPlayerID(PlayerID string) []structs.NBAContractOffer {
	db := dbprovider.GetInstance().GetDB()

	offers := []structs.NBAContractOffer{}

	err := db.Where("NBA_player_id = ? AND is_active = ?", PlayerID, true).Find(&offers).Error
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

		id := row[0]
		teamID := row[1]
		playerRecord := GetNBAPlayerRecord(id)
		minimumValue := playerRecord.MinimumValue
		contractLength := util.ConvertStringToInt(row[3])
		totalValue := util.ConvertStringToInt(row[4])
		year1 := util.ConvertStringToInt(row[5])
		year2 := util.ConvertStringToInt(row[6])
		year3 := util.ConvertStringToInt(row[7])
		year4 := util.ConvertStringToInt(row[8])
		year5 := util.ConvertStringToInt(row[9])
		contractStatus := ""
		if playerRecord.MaxRequested {
			contractStatus = "Max"
		}
		if playerRecord.IsSuperMaxQualified {
			contractStatus = "SuperMax"
		}

		pref := playerRecord.FreeAgency

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

func validateFreeAgencyPref(playerRecord structs.NBAPlayer, teamID uint, totalValue int) bool {
	preference := playerRecord.FreeAgency

	if preference == "Average" {
		return true
	}
	if preference == "Drafted team discount" && playerRecord.DraftedTeamID == teamID {
		return true
	}
	if preference == "Loyal" {

	}
	if preference == "Hometown Hero" {

	}
	if preference == "Adversarial" {

	}

	if preference == "I'm the starter" {

	}
	if preference == "Market-driven" {

	}
	if preference == "Money motivated" {

	}
	if preference == "Highest bidder" {

	}
	if preference == "Championship seeking" {

	}
}
