package managers

import (
	"strconv"
	"sync"

	"github.com/CalebRose/SimNBA/structs"
)

func MakeNBAPlayerMap(nbaPlayers []structs.NBAPlayer) map[uint]structs.NBAPlayer {
	playerMap := make(map[uint]structs.NBAPlayer)

	for _, p := range nbaPlayers {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakeCollegePlayerMap(players []structs.CollegePlayer) map[uint]structs.CollegePlayer {
	playerMap := make(map[uint]structs.CollegePlayer)

	for _, p := range players {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakeCollegePlayerMapByTeamID(players []structs.CollegePlayer, excludeUnsigned bool) map[uint][]structs.CollegePlayer {
	playerMap := make(map[uint][]structs.CollegePlayer)

	for _, p := range players {
		if p.TeamID == 0 && excludeUnsigned {
			continue
		}
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.CollegePlayer{p}
		}
	}

	return playerMap
}

func MakeNBAPlayerMapByTeamID(players []structs.NBAPlayer, excludeFAs bool) map[uint][]structs.NBAPlayer {
	playerMap := make(map[uint][]structs.NBAPlayer)

	for _, p := range players {
		if p.TeamID == 0 && excludeFAs {
			continue
		}
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.NBAPlayer{p}
		}
	}

	return playerMap
}

func MakeFullTransferPortalProfileMap(players []structs.CollegePlayer) map[uint][]structs.TransferPortalProfile {
	playerIDs := []string{}
	for _, p := range players {
		playerID := strconv.Itoa(int(p.ID))
		playerIDs = append(playerIDs, playerID)
	}
	portalProfiles := GetTransferPortalProfilesByPlayerIDs(playerIDs)
	portalMap := make(map[uint][]structs.TransferPortalProfile)
	var mu sync.Mutex     // to safely update the map
	var wg sync.WaitGroup // to wait for all goroutines to finish
	semaphore := make(chan struct{}, 10)
	for _, p := range portalProfiles {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(c structs.TransferPortalProfile) {
			defer wg.Done()
			mu.Lock()
			if len(portalMap[c.ID]) == 0 {
				portalMap[c.ID] = []structs.TransferPortalProfile{c}
			} else {
				portalMap[c.ID] = append(portalMap[c.ID], c)
			}
			mu.Unlock()

			<-semaphore
		}(p)
	}
	wg.Wait()
	close(semaphore)
	return portalMap
}

func MakeContractMap(contracts []structs.NBAContract) map[uint]structs.NBAContract {
	contractMap := make(map[uint]structs.NBAContract)

	for _, c := range contracts {
		contractMap[uint(c.PlayerID)] = c
	}

	return contractMap
}

func MakeExtensionMap(extensions []structs.NBAExtensionOffer) map[uint]structs.NBAExtensionOffer {
	contractMap := make(map[uint]structs.NBAExtensionOffer)

	for _, c := range extensions {
		contractMap[uint(c.NBAPlayerID)] = c
	}

	return contractMap
}

func MakeFreeAgencyOfferMap(offers []structs.NBAContractOffer) map[uint][]structs.NBAContractOffer {
	offerMap := make(map[uint][]structs.NBAContractOffer)

	for _, offer := range offers {
		if len(offerMap[offer.PlayerID]) > 0 {
			offerMap[offer.PlayerID] = append(offerMap[uint(offer.PlayerID)], offer)
		} else {
			offerMap[offer.PlayerID] = []structs.NBAContractOffer{offer}
		}
	}

	return offerMap
}

func MakeFreeAgencyOfferMapByTeamID(offers []structs.NBAContractOffer) map[uint][]structs.NBAContractOffer {
	offerMap := make(map[uint][]structs.NBAContractOffer)

	for _, offer := range offers {
		if len(offerMap[offer.TeamID]) > 0 {
			offerMap[offer.TeamID] = append(offerMap[uint(offer.TeamID)], offer)
		} else {
			offerMap[offer.TeamID] = []structs.NBAContractOffer{offer}
		}
	}

	return offerMap
}

func MakeCollegeGameplanMap(gamePlans []structs.Gameplan) map[uint]structs.Gameplan {
	gameplanMap := make(map[uint]structs.Gameplan)

	for _, gp := range gamePlans {
		gameplanMap[gp.TeamID] = gp
	}

	return gameplanMap
}

func MakeCollegeTeamMap(teams []structs.Team) map[uint]structs.Team {
	gameplanMap := make(map[uint]structs.Team)

	for _, team := range teams {
		gameplanMap[team.ID] = team
	}

	return gameplanMap
}

func MakeNBAGameplanMap(gamePlans []structs.NBAGameplan) map[uint]structs.NBAGameplan {
	gameplanMap := make(map[uint]structs.NBAGameplan)

	for _, gp := range gamePlans {
		gameplanMap[gp.TeamID] = gp
	}

	return gameplanMap
}

func MakeNBATeamMap(teams []structs.NBATeam) map[uint]structs.NBATeam {
	gameplanMap := make(map[uint]structs.NBATeam)

	for _, team := range teams {
		gameplanMap[team.ID] = team
	}

	return gameplanMap
}

func MakeNBATradePreferencesMap(tradePreferences []structs.NBATradePreferences) map[uint]structs.NBATradePreferences {
	preferencesMap := make(map[uint]structs.NBATradePreferences)

	for _, c := range tradePreferences {
		preferencesMap[uint(c.NBATeamID)] = c
	}

	return preferencesMap
}

func MakeNBAWarRoomMap(warRooms []structs.NBAWarRoom) map[uint]structs.NBAWarRoom {
	warRoomMap := make(map[uint]structs.NBAWarRoom)

	for _, t := range warRooms {
		warRoomMap[t.TeamID] = t
	}

	return warRoomMap
}

func MakeScoutingProfileMapByTeam(profiles []structs.ScoutingProfile) map[uint]structs.ScoutingProfile {
	profileMap := make(map[uint]structs.ScoutingProfile)

	for _, t := range profiles {
		profileMap[t.TeamID] = t
	}

	return profileMap
}
