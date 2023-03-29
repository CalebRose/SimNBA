package managers

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

func GetTradeBlockDataByTeamID(TeamID string) structs.NBATradeBlockResponse {
	var waitgroup sync.WaitGroup

	waitgroup.Add(5)
	NBATeamChan := make(chan structs.NBATeam)
	playersChan := make(chan []structs.NBAPlayer)
	picksChan := make(chan []structs.DraftPick)
	proposalsChan := make(chan structs.NBATeamProposals)
	preferencesChan := make(chan structs.NBATradePreferences)

	go func() {
		waitgroup.Wait()
		close(NBATeamChan)
		close(playersChan)
		close(picksChan)
		close(proposalsChan)
		close(preferencesChan)
	}()

	go func() {
		defer waitgroup.Done()
		team := GetNBATeamWithCapsheetByTeamID(TeamID)
		NBATeamChan <- team
	}()

	go func() {
		defer waitgroup.Done()
		players := GetTradableNBAPlayersByTeamID(TeamID)
		playersChan <- players
	}()

	go func() {
		defer waitgroup.Done()
		picks := GetDraftPicksByTeamID(TeamID)
		picksChan <- picks
	}()

	go func() {
		defer waitgroup.Done()
		proposals := GetTradeProposalsByNBAID(TeamID)
		proposalsChan <- proposals
	}()

	go func() {
		defer waitgroup.Done()
		pref := GetTradePreferencesByTeamID(TeamID)
		preferencesChan <- pref
	}()

	NBATeam := <-NBATeamChan
	tradablePlayers := <-playersChan
	draftPicks := <-picksChan
	teamProposals := <-proposalsChan
	tradePreferences := <-preferencesChan

	// close(NBATeamChan)
	// close(playersChan)
	// close(picksChan)
	// close(proposalsChan)

	return structs.NBATradeBlockResponse{
		Team:                   NBATeam,
		TradablePlayers:        tradablePlayers,
		DraftPicks:             draftPicks,
		SentTradeProposals:     teamProposals.SentTradeProposals,
		ReceivedTradeProposals: teamProposals.ReceivedTradeProposals,
		TradePreferences:       tradePreferences,
	}
}

// GetTradeProposalsByNBAID -- Returns all trade proposals that were either sent or received
func GetTradeProposalsByNBAID(TeamID string) structs.NBATeamProposals {
	db := dbprovider.GetInstance().GetDB()

	proposals := []structs.NBATradeProposal{}

	db.Preload("NBATeamTradeOptions").Where("NBA_team_id = ? OR recepient_team_id = ?", TeamID, TeamID).Find(&proposals)

	SentProposals := []structs.NBATradeProposalDTO{}
	ReceivedProposals := []structs.NBATradeProposalDTO{}

	id := uint(util.ConvertStringToInt(TeamID))

	for _, proposal := range proposals {
		if proposal.IsTradeAccepted || proposal.IsTradeRejected {
			continue
		}
		sentOptions := []structs.NBATradeOptionObj{}
		receivedOptions := []structs.NBATradeOptionObj{}
		for _, option := range proposal.NBATeamTradeOptions {
			opt := structs.NBATradeOptionObj{
				ID:               option.Model.ID,
				TradeProposalID:  option.TradeProposalID,
				NBATeamID:        option.NBATeamID,
				SalaryPercentage: option.SalaryPercentage,
				OptionType:       option.OptionType,
			}
			if option.NBAPlayerID > 0 {
				player := GetNBAPlayerRecord(strconv.Itoa(int(option.NBAPlayerID)))
				opt.AssignPlayer(player)
			} else if option.NBADraftPickID > 0 {
				draftPick := GetDraftPickByDraftPickID(strconv.Itoa((int(option.NBADraftPickID))))
				opt.AssignPick(draftPick)
			}
			if option.NBATeamID == proposal.NBATeamID {
				sentOptions = append(sentOptions, opt)
			} else {
				receivedOptions = append(receivedOptions, opt)
			}
		}

		proposalResponse := structs.NBATradeProposalDTO{
			ID:                        proposal.Model.ID,
			NBATeamID:                 proposal.NBATeamID,
			NBATeam:                   proposal.NBATeam,
			RecepientTeamID:           proposal.RecepientTeamID,
			RecepientTeam:             proposal.RecepientTeam,
			IsTradeAccepted:           proposal.IsTradeAccepted,
			IsTradeRejected:           proposal.IsTradeRejected,
			NBATeamTradeOptions:       sentOptions,
			RecepientTeamTradeOptions: receivedOptions,
		}

		if proposal.NBATeamID == id {
			SentProposals = append(SentProposals, proposalResponse)
		} else {
			ReceivedProposals = append(ReceivedProposals, proposalResponse)
		}
	}
	return structs.NBATeamProposals{
		SentTradeProposals:     SentProposals,
		ReceivedTradeProposals: ReceivedProposals,
	}
}

func GetTradePreferencesByTeamID(TeamID string) structs.NBATradePreferences {
	db := dbprovider.GetInstance().GetDB()

	preferences := structs.NBATradePreferences{}

	db.Where("id = ?", TeamID).Find(&preferences)

	return preferences
}

func UpdateTradePreferences(pref structs.NBATradePreferencesDTO) {
	db := dbprovider.GetInstance().GetDB()

	preferences := GetTradePreferencesByTeamID(strconv.Itoa(int(pref.NBATeamID)))

	preferences.UpdatePreferences(pref)

	db.Save(&preferences)
}

func GetAcceptedTradeProposals() []structs.NBATradeProposalDTO {
	db := dbprovider.GetInstance().GetDB()

	proposals := []structs.NBATradeProposal{}

	db.Preload("NBATeamTradeOptions").Where("is_trade_accepted = ? AND is_synced = ?", true, false).Find(&proposals)

	acceptedProposals := []structs.NBATradeProposalDTO{}

	for _, proposal := range proposals {
		sentOptions := []structs.NBATradeOptionObj{}
		receivedOptions := []structs.NBATradeOptionObj{}
		for _, option := range proposal.NBATeamTradeOptions {
			opt := structs.NBATradeOptionObj{
				ID:               option.Model.ID,
				TradeProposalID:  option.TradeProposalID,
				NBATeamID:        option.NBATeamID,
				SalaryPercentage: option.SalaryPercentage,
				OptionType:       option.OptionType,
			}
			if option.NBAPlayerID > 0 {
				player := GetNBAPlayerRecord(strconv.Itoa(int(option.NBAPlayerID)))
				opt.AssignPlayer(player)
			} else if option.NBADraftPickID > 0 {
				draftPick := GetDraftPickByDraftPickID(strconv.Itoa((int(option.NBADraftPickID))))
				opt.AssignPick(draftPick)
			}
			if option.NBATeamID == proposal.NBATeamID {
				sentOptions = append(sentOptions, opt)
			} else {
				receivedOptions = append(receivedOptions, opt)
			}
		}

		proposalResponse := structs.NBATradeProposalDTO{
			ID:                        proposal.Model.ID,
			NBATeamID:                 proposal.NBATeamID,
			NBATeam:                   proposal.NBATeam,
			RecepientTeamID:           proposal.RecepientTeamID,
			RecepientTeam:             proposal.RecepientTeam,
			IsTradeAccepted:           proposal.IsTradeAccepted,
			IsTradeRejected:           proposal.IsTradeRejected,
			NBATeamTradeOptions:       sentOptions,
			RecepientTeamTradeOptions: receivedOptions,
		}

		acceptedProposals = append(acceptedProposals, proposalResponse)
	}

	return acceptedProposals
}

func GetRejectedTradeProposals() []structs.NBATradeProposal {
	db := dbprovider.GetInstance().GetDB()

	proposals := []structs.NBATradeProposal{}

	db.Preload("NBATeamTradeOptions").Where("is_trade_rejected = ?", true).Find(&proposals)

	return proposals
}

func PlaceNBAPlayerOnTradeBlock(playerID string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetNBAPlayerRecord(playerID)

	player.ToggleTradeBlock()

	db.Save(&player)
}

func CreateTradeProposal(TradeProposal structs.NBATradeProposalDTO) {
	db := dbprovider.GetInstance().GetDB()
	latestID := GetLatestProposalInDB(db)

	// Create Trade Proposal Object
	proposal := structs.NBATradeProposal{
		NBATeamID:       TradeProposal.NBATeamID,
		NBATeam:         TradeProposal.NBATeam,
		RecepientTeamID: TradeProposal.RecepientTeamID,
		RecepientTeam:   TradeProposal.RecepientTeam,
		IsTradeAccepted: false,
		IsTradeRejected: false,
	}
	proposal.AssignID(latestID)

	db.Create(&proposal)

	// Create Trade Options
	SentTradeOptions := TradeProposal.NBATeamTradeOptions
	ReceivedTradeOptions := TradeProposal.RecepientTeamTradeOptions

	for _, sentOption := range SentTradeOptions {
		tradeOption := structs.NBATradeOption{
			TradeProposalID:  latestID,
			NBATeamID:        TradeProposal.NBATeamID,
			NBAPlayerID:      sentOption.NBAPlayerID,
			NBADraftPickID:   sentOption.NBADraftPickID,
			SalaryPercentage: sentOption.SalaryPercentage,
			OptionType:       sentOption.OptionType,
		}
		db.Create(&tradeOption)
	}

	for _, recepientOption := range ReceivedTradeOptions {
		tradeOption := structs.NBATradeOption{
			TradeProposalID:  latestID,
			NBATeamID:        TradeProposal.RecepientTeamID,
			NBAPlayerID:      recepientOption.NBAPlayerID,
			NBADraftPickID:   recepientOption.NBADraftPickID,
			SalaryPercentage: recepientOption.SalaryPercentage,
			OptionType:       recepientOption.OptionType,
		}
		db.Create(&tradeOption)
	}
}

func GetOnlyTradeProposalByProposalID(proposalID string) structs.NBATradeProposal {
	db := dbprovider.GetInstance().GetDB()

	proposal := structs.NBATradeProposal{}

	db.Preload("NBATeamTradeOptions").Where("id = ?", proposalID).Find(&proposal)

	return proposal
}

func AcceptTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)

	proposal.AcceptTrade()

	// Create News Log
	newsLogMessage := proposal.RecepientTeam + " has accepted a trade offer from " + proposal.NBATeam + " for trade the following players:\n\n"

	for _, options := range proposal.NBATeamTradeOptions {
		if options.NBATeamID == proposal.NBATeamID {
			if options.NBAPlayerID > 0 {
				playerRecord := GetNBAPlayerRecord(strconv.Itoa(int(options.NBAPlayerID)))
				ovrGrade := util.GetOverallGrade(playerRecord.Overall)
				ovr := playerRecord.Overall
				if playerRecord.Year > 1 {
					newsLogMessage += playerRecord.Position + " " + strconv.Itoa(ovr) + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + proposal.RecepientTeam + "\n"
				} else {
					newsLogMessage += playerRecord.Position + " " + ovrGrade + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + proposal.RecepientTeam + "\n"
				}
			} else if options.NBADraftPickID > 0 {
				draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(options.NBADraftPickID)))
				pickRound := strconv.Itoa(int(draftPick.DraftRound))
				roundAbbreviation := util.GetRoundAbbreviation(pickRound)
				season := strconv.Itoa(int(draftPick.Season))
				newsLogMessage += season + " " + roundAbbreviation + " pick to " + proposal.RecepientTeam + "\n"
			}
		} else {
			if options.NBAPlayerID > 0 {
				playerRecord := GetNBAPlayerRecord(strconv.Itoa(int(options.NBAPlayerID)))
				newsLogMessage += playerRecord.Position + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + proposal.NBATeam + "\n"
			} else if options.NBADraftPickID > 0 {
				draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(options.NBADraftPickID)))
				pickRound := strconv.Itoa(int(draftPick.DraftRound))
				roundAbbreviation := util.GetRoundAbbreviation(pickRound)
				season := strconv.Itoa(int(draftPick.Season))
				newsLogMessage += season + " " + roundAbbreviation + " pick to " + proposal.NBATeam + "\n"
			}
		}
	}
	newsLogMessage += "\n"

	newsLog := structs.NewsLog{
		WeekID:      ts.NBAWeekID,
		Week:        uint(ts.NBAWeek),
		SeasonID:    ts.SeasonID,
		League:      "NBA",
		MessageType: "Trade",
		Message:     newsLogMessage,
	}

	db.Create(&newsLog)
	db.Save(&proposal)
}

func RejectTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)

	proposal.RejectTrade()
	newsLog := structs.NewsLog{
		WeekID:      ts.NBAWeekID,
		Week:        uint(ts.NBAWeek),
		SeasonID:    ts.SeasonID,
		League:      "NBA",
		MessageType: "Trade",
		Message:     proposal.RecepientTeam + " has rejected a trade from " + proposal.NBATeam,
	}

	db.Create(&newsLog)
	db.Save(&proposal)
}

func CancelTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)
	options := proposal.NBATeamTradeOptions

	for _, option := range options {
		db.Delete(&option)
	}

	db.Delete(&proposal)
}

func GetLatestProposalInDB(db *gorm.DB) uint {
	var latestProposal structs.NBATradeProposal

	err := db.Last(&latestProposal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestProposal.ID + 1
}

func RemoveRejectedTrades() {
	db := dbprovider.GetInstance().GetDB()

	rejectedProposals := GetRejectedTradeProposals()

	for _, proposal := range rejectedProposals {
		sentOptions := proposal.NBATeamTradeOptions
		deleteOptions(db, sentOptions)

		// Delete Proposal
		db.Delete(&proposal)
	}
}

func SyncAcceptedTrade(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)
	SentOptions := proposal.NBATeamTradeOptions

	syncAcceptedOptions(db, SentOptions, proposal.NBATeamID, proposal.NBATeam, proposal.RecepientTeamID, proposal.RecepientTeam)

	proposal.ToggleSyncStatus()

	db.Save(&proposal)
}

func syncAcceptedOptions(db *gorm.DB, options []structs.NBATradeOption, senderID uint, senderTeam string, recepientID uint, recepientTeam string) {
	sendingTeam := GetNBATeamByTeamID(strconv.Itoa(int(senderID)))
	receivingTeam := GetNBATeamByTeamID(strconv.Itoa(int(recepientID)))
	SendersCapsheet := GetCapsheetByTeamID(strconv.Itoa(int(senderID)))
	recepientCapsheet := GetCapsheetByTeamID(strconv.Itoa(int(recepientID)))
	for _, option := range options {
		// Contract
		percentage := option.SalaryPercentage
		if option.NBAPlayerID > 0 {
			playerRecord := GetNBAPlayerRecord(strconv.Itoa(int(option.NBAPlayerID)))
			contract := playerRecord.Contract
			if playerRecord.TeamID == senderID {
				sendersPercentage := percentage * 0.01
				receiversPercentage := (100 - percentage) * 0.01
				SendersCapsheet.SubtractFromCapsheetViaTrade(contract)
				SendersCapsheet.NegotiateSalaryDifference(contract.Year1Total, contract.Year1Total*sendersPercentage)
				recepientCapsheet.AddContractViaTrade(contract, contract.Year1Total*receiversPercentage)
				playerRecord.TradePlayer(recepientID, receivingTeam.Abbr)
				contract.TradePlayer(recepientID, receivingTeam.Abbr, receiversPercentage)
			} else {
				receiversPercentage := percentage * 0.01
				sendersPercentage := (100 - percentage) * 0.01
				recepientCapsheet.SubtractFromCapsheetViaTrade(contract)
				recepientCapsheet.NegotiateSalaryDifference(contract.Year1Total, contract.Year1Total*receiversPercentage)
				SendersCapsheet.AddContractViaTrade(contract, contract.Year1Total*sendersPercentage)
				playerRecord.TradePlayer(senderID, sendingTeam.Abbr)
				contract.TradePlayer(senderID, sendingTeam.Abbr, sendersPercentage)
			}

			db.Save(&playerRecord)
			db.Save(&contract)

		} else if option.NBADraftPickID > 0 {
			draftPick := GetDraftPickByDraftPickID(strconv.Itoa(int(option.NBADraftPickID)))
			if draftPick.TeamID == senderID {
				draftPick.TradePick(recepientID, recepientTeam)
			} else {
				draftPick.TradePick(senderID, senderTeam)
			}

			db.Save(&draftPick)
		}

		db.Delete(&option)
	}
	db.Save(&SendersCapsheet)
	db.Save(&recepientCapsheet)
}

func VetoTrade(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := GetOnlyTradeProposalByProposalID(proposalID)
	SentOptions := proposal.NBATeamTradeOptions

	deleteOptions(db, SentOptions)

	db.Delete(&proposal)
}

func deleteOptions(db *gorm.DB, options []structs.NBATradeOption) {
	// Delete Recepient Trade Options
	for _, option := range options {
		// Deletes the option
		db.Delete(&option)
	}
}
