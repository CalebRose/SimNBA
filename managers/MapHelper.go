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
