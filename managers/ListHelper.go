package managers

import (
	"sort"

	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func MakeCollegeInjuryList(players []structs.CollegePlayer) []structs.CollegePlayer {
	injuryList := []structs.CollegePlayer{}

	for _, p := range players {
		if p.IsInjured {
			injuryList = append(injuryList, p)
		}
	}
	return injuryList
}

func MakeCollegePortalList(players []structs.CollegePlayer) []structs.CollegePlayer {
	portalList := []structs.CollegePlayer{}

	for _, p := range players {
		if p.TransferStatus > 0 {
			portalList = append(portalList, p)
		}
	}
	return portalList
}

func MakeProInjuryList(players []structs.NBAPlayer) []structs.NBAPlayer {
	injuryList := []structs.NBAPlayer{}

	for _, p := range players {
		if p.IsInjured {
			injuryList = append(injuryList, p)
		}
	}
	return injuryList
}

func MakeTransferPortalPlayerResponseList(players []structs.CollegePlayer, profileMap map[uint][]structs.TransferPortalProfile) []structs.TransferPlayerResponse {
	responseList := []structs.TransferPlayerResponse{}

	for _, p := range players {
		profiles := profileMap[p.ID]
		p.AssignTransferProfiles(profiles)
		ovr := util.GetPlayerOverallGrade(p.Overall)

		tp := structs.TransferPlayerResponse{}
		tp.Map(p, ovr)

		responseList = append(responseList, tp)
	}

	return responseList
}

func GetCBBOrderedListByStatType(statType string, teamID uint, CollegeStats []structs.CollegePlayerSeasonStats, collegePlayerMap map[uint]structs.CollegePlayer) []structs.CollegePlayer {
	orderedStats := CollegeStats
	resultList := []structs.CollegePlayer{}
	if statType == "POINTS" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].PPG > orderedStats[j].PPG
		})
	} else if statType == "ASSISTS" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].AssistsPerGame > orderedStats[j].AssistsPerGame
		})
	} else if statType == "REBOUNDS" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].ReboundsPerGame > orderedStats[j].ReboundsPerGame
		})
	}

	teamLeaderInTopStats := false
	for idx, stat := range orderedStats {
		if idx > 4 {
			break
		}
		player := collegePlayerMap[stat.CollegePlayerID]
		if stat.TeamID == teamID {
			teamLeaderInTopStats = true
		}
		player.AddSeasonStats(stat)
		resultList = append(resultList, player)
	}

	if !teamLeaderInTopStats {
		for _, stat := range orderedStats {
			if stat.TeamID == teamID {
				player := collegePlayerMap[stat.CollegePlayerID]
				player.AddSeasonStats(stat)
				resultList = append(resultList, player)
				break
			}
		}
	}
	return resultList
}

func getNFLOrderedListByStatType(statType string, teamID uint, CollegeStats []structs.NBAPlayerSeasonStats, proPlayerMap map[uint]structs.NBAPlayer) []structs.NBAPlayer {
	orderedStats := CollegeStats
	resultList := []structs.NBAPlayer{}
	if statType == "POINTS" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].PPG > orderedStats[j].PPG
		})
	} else if statType == "ASSISTS" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].AssistsPerGame > orderedStats[j].AssistsPerGame
		})
	} else if statType == "REBOUNDS" {
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].ReboundsPerGame > orderedStats[j].ReboundsPerGame
		})
	}

	teamLeaderInTopStats := false
	for idx, stat := range orderedStats {
		if idx > 4 {
			break
		}
		player := proPlayerMap[stat.NBAPlayerID]
		if stat.TeamID == teamID {
			teamLeaderInTopStats = true
		}
		player.AddSeasonStats(stat)
		resultList = append(resultList, player)
	}

	if !teamLeaderInTopStats {
		for _, stat := range orderedStats {
			if stat.TeamID == teamID {
				player := proPlayerMap[stat.NBAPlayerID]
				player.AddSeasonStats(stat)
				resultList = append(resultList, player)
				break
			}
		}
	}
	return resultList
}
