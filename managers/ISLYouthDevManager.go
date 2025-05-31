package managers

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func GetAllYouthDevelopmentPlayers() []structs.NBAPlayer {
	db := dbprovider.GetInstance().GetDB()

	var players []structs.NBAPlayer

	db.Where("is_int_generated = ? AND age < ? and country != ? AND team_id = ? AND team_abbr != ?", true, "23", "USA", "0", "DRAFT").Find(&players)

	return players
}

func GetAllYouthDevelopmentPlayerMap() map[uint]structs.NBAPlayer {
	players := GetAllYouthDevelopmentPlayers()
	playerMap := make(map[uint]structs.NBAPlayer)
	for _, p := range players {
		playerMap[p.ID] = p
	}
	return playerMap
}

func GetYouthDevelopmentPlayerCount() int {
	players := GetAllYouthDevelopmentPlayers()
	return len(players)
}

func GetScoutingDeptByTeamID(id string) structs.ISLScoutingDept {
	db := dbprovider.GetInstance().GetDB()

	dept := structs.ISLScoutingDept{}

	err := db.Where("team_id = ?", id).Find(&dept).Error
	if err != nil {
		return dept
	}

	return dept
}

func GetAllScoutingDepts() []structs.ISLScoutingDept {
	db := dbprovider.GetInstance().GetDB()

	dept := []structs.ISLScoutingDept{}

	err := db.Find(&dept).Error
	if err != nil {
		return dept
	}

	return dept
}

func GetScoutingReportByPlayerAndTeam(pid, tid string) structs.ISLScoutingReport {
	db := dbprovider.GetInstance().GetDB()

	report := structs.ISLScoutingReport{}

	err := db.Where("player_id = ? and team_id = ?", pid, tid).Find(&report).Error
	if err != nil {
		return report
	}

	return report
}

func GetScoutingReportsByPlayerID(pid string) []structs.ISLScoutingReport {
	db := dbprovider.GetInstance().GetDB()

	report := []structs.ISLScoutingReport{}

	err := db.Where("player_id = ? AND removed_from_board = ?", pid, false).Find(&report).Error
	if err != nil {
		return report
	}

	return report
}

func GetScoutingReportByPlayerIDMAP(youthPlayers []structs.NBAPlayer) map[uint][]structs.ISLScoutingReport {
	newMap := make(map[uint][]structs.ISLScoutingReport)

	for _, p := range youthPlayers {
		pid := strconv.Itoa(int(p.ID))
		reports := GetScoutingReportsByPlayerID(pid)
		newMap[p.ID] = reports
	}

	return newMap
}

func GetScoutingReportsByTeamID(tid string) []structs.ISLScoutingReport {
	db := dbprovider.GetInstance().GetDB()

	report := []structs.ISLScoutingReport{}

	err := db.Where("team_id = ?", tid).Find(&report).Error
	if err != nil {
		return report
	}

	return report
}

func ISLIdentityPhase() {
	db := dbprovider.GetInstance().GetDB()

	depts := GetAllScoutingDepts()
	fmt.Println("Loading available players...")
	youthPlayers := GetAllYouthDevelopmentPlayers()
	fmt.Println("Loading existing reports...")
	reportMap := GetScoutingReportByPlayerIDMAP(youthPlayers)

	// To help teams build up, set player list by OVR.
	fmt.Println("Ordering players...")
	sort.Slice(youthPlayers, func(i, j int) bool {
		return youthPlayers[i].Overall > youthPlayers[j].Overall
	})

	fmt.Println("Shuffling list of departments...")
	// Shuffle the departments so that teams don't pick based on order of creation
	rand.Shuffle(len(depts), func(i, j int) {
		depts[i], depts[j] = depts[j], depts[i]
	})

	teamMap := GetProfessionalTeamMap()
	adjCountryMap := util.GetAdjacentCountryMap()
	for _, d := range depts {
		teamId := strconv.Itoa(int(d.TeamID))
		currentRoster := GetAllNBAPlayersByTeamID(teamId)
		if len(currentRoster) > 12 {
			continue
		}
		team := teamMap[d.TeamID]
		pointCost := 8 - d.IdentityMod
		if pointCost > 100 {
			pointCost = 1
		}
		for _, p := range youthPlayers {
			if d.Resources <= 20 || d.Resources-pointCost <= 0 {
				break
			}
			playerID := strconv.Itoa(int(p.ID))
			existingReport := GetScoutingReportByPlayerAndTeam(playerID, teamId)
			if existingReport.ID > 0 {
				continue
			}
			base := 80
			// Determine if player is in team's country
			adjCountries := adjCountryMap[p.Country]
			if p.Country == team.Country {
				base += 95
			} else if util.CheckIfStringInList(team.Country, adjCountries) {
				base += 45
			}
			roll := util.GenerateIntFromRange(1, 100)
			proposedPoints := int(d.Resources - pointCost)
			reports := reportMap[p.ID]
			if roll <= base && proposedPoints >= 0 && len(reports) < 2 {
				d.IncrementPool(1, pointCost)
				// Add player to board
				report := structs.ISLScoutingReport{
					TeamID:   d.TeamID,
					PlayerID: p.ID,
				}

				reportMap[p.ID] = append(reportMap[p.ID], report)

				repository.CreateISLScoutingReportRecord(report, db)
			}
		}
		repository.SaveISLScoutingDeptRecord(d, db)
	}
}

func ISLScoutingPhase() {
	db := dbprovider.GetInstance().GetDB()
	playerMap := GetAllYouthDevelopmentPlayerMap()
	depts := GetAllScoutingDepts()

	for _, d := range depts {
		teamId := strconv.Itoa(int(d.TeamID))
		selectionList := []string{"fn", "sh2", "sh3", "ft", "bw", "rb", "ind", "prd", "pot", "idn"}
		rand.Shuffle(len(selectionList), func(i, j int) {
			selectionList[i], selectionList[j] = selectionList[j], selectionList[i]
		})
		scoutingReports := GetScoutingReportsByTeamID(teamId)
		pointsRemaining := int(d.Resources)
		if pointsRemaining < 0 {
			pointsRemaining = 0
		}
		pointsToSpend := 200
		if d.BehaviorBias == 1 {
			pointsToSpend += 10
		} else if d.BehaviorBias == 3 {
			pointsToSpend -= 10
		}
		if pointsToSpend > pointsRemaining {
			pointsToSpend = pointsRemaining
		}
		if pointsToSpend < 0 {
			pointsToSpend = 0
		}
		teamNeedsMap := make(map[string]bool)
		positionCount := make(map[string]int)
		currentRoster := GetAllNBAPlayersByTeamID(teamId)
		if len(currentRoster) > 12 {
			continue
		}
		if _, ok := teamNeedsMap["PG"]; !ok {
			teamNeedsMap["PG"] = true
		}
		if _, ok := teamNeedsMap["SG"]; !ok {
			teamNeedsMap["SG"] = true
		}
		if _, ok := teamNeedsMap["SF"]; !ok {
			teamNeedsMap["SF"] = true
		}
		if _, ok := teamNeedsMap["PF"]; !ok {
			teamNeedsMap["PF"] = true
		}
		if _, ok := teamNeedsMap["C"]; !ok {
			teamNeedsMap["C"] = true
		}

		if _, ok := positionCount["PG"]; !ok {
			positionCount["PG"] = 0
		}
		if _, ok := positionCount["SG"]; !ok {
			positionCount["SG"] = 0
		}
		if _, ok := positionCount["SF"]; !ok {
			positionCount["SF"] = 0
		}
		if _, ok := positionCount["PF"]; !ok {
			positionCount["PF"] = 0
		}
		if _, ok := positionCount["C"]; !ok {
			positionCount["C"] = 0
		}

		for _, r := range currentRoster {
			positionCount[r.Position] += 1
		}

		if positionCount["PG"] >= 3 {
			teamNeedsMap["PG"] = false
		} else if positionCount["SG"] >= 4 {
			teamNeedsMap["SG"] = false
		} else if positionCount["SF"] >= 4 {
			teamNeedsMap["SF"] = false
		} else if positionCount["PF"] >= 4 {
			teamNeedsMap["PF"] = false
		} else if positionCount["C"] >= 3 {
			teamNeedsMap["C"] = false
		}

		for _, s := range scoutingReports {
			if s.RemovedFromBoard || s.IsLocked {
				continue
			}
			change := false
			if pointsToSpend <= 0 {
				break
			}
			// Skip over if overall is already revealed
			if s.Overall || s.RemovedFromBoard {
				continue
			}
			player := playerMap[s.PlayerID]
			if !teamNeedsMap[player.Position] {
				s.RemovePlayerFromBoard()
				repository.SaveISLScoutingReportRecord(s, db)
				continue
			}
			baseCost := 10
			for _, attr := range selectionList {
				if s.Overall {
					break
				}
				if attr == "fn" && !s.Finishing {
					pointCost := baseCost - int(d.Finishing)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("fn")
						change = true
					}
				} else if attr == "sh2" && !s.Shooting2 {
					pointCost := baseCost - int(d.Shooting2)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("sh2")
						change = true
					}
				} else if attr == "sh3" && !s.Shooting3 {
					pointCost := baseCost - int(d.Shooting3)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						pointsToSpend -= int(pointCost)
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("sh3")
						change = true
					}
				} else if attr == "ft" && !s.FreeThrow {
					pointCost := baseCost - int(d.FreeThrow)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						pointsToSpend -= int(pointCost)
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("ft")
						change = true
					}
				} else if attr == "bw" && !s.Ballwork {
					pointCost := baseCost - int(d.Ballwork)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						pointsToSpend -= int(pointCost)
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("bw")
						change = true
					}
				} else if attr == "rb" && !s.Rebounding {
					pointCost := baseCost - int(d.Rebounding)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						pointsToSpend -= int(pointCost)
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("rb")
						change = true
					}
				} else if attr == "ind" && !s.IntDefense {
					pointCost := baseCost - int(d.IntDefense)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						pointsToSpend -= int(pointCost)
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("ind")
						change = true
					}
				} else if attr == "prd" && !s.PerDefense {
					pointCost := baseCost - int(d.PerDefense)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						pointsToSpend -= int(pointCost)
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("prd")
						change = true
					}
				} else if attr == "pot" && !s.Potential {
					pointCost := baseCost - int(d.Potential)
					if pointCost < 0 {
						pointCost = 0
					}
					if pointsToSpend >= pointCost {
						pointsToSpend -= pointCost
						pointsToSpend -= int(pointCost)
						d.IncrementPool(2, uint8(pointCost))
						s.RevealAttribute("pot")
						change = true
					}
				}
			}
			// if Overall was just revealed
			if s.Overall {
				pointRequirement := 30
				coinFlip := util.GenerateIntFromRange(1, 2)
				if d.Prestige > 3 && player.Overall <= 40 && coinFlip == 2 {
					s.RemovePlayerFromBoard()
				} else if player.Overall > 65 {
					pointRequirement = player.Overall
				}
				s.SetPointRequirement(uint(pointRequirement))
			}
			if change {
				repository.SaveISLScoutingReportRecord(s, db)
			}
		}
		repository.SaveISLScoutingDeptRecord(d, db)
	}
}

func ISLInvestingPhase() {
	db := dbprovider.GetInstance().GetDB()

	depts := GetAllScoutingDepts()
	playerMap := GetAllYouthDevelopmentPlayerMap()

	for _, d := range depts {
		teamId := strconv.Itoa(int(d.TeamID))
		scoutingReports := GetScoutingReportsByTeamID(teamId)
		pointsRemaining := int(d.Resources)
		if pointsRemaining < 0 {
			pointsRemaining = 0
		}
		for _, s := range scoutingReports {
			if pointsRemaining <= 0 {
				break
			}
			// Skip over if overall is already revealed
			if !s.Overall || s.RemovedFromBoard || s.IsLocked || s.CurrentPoints > 0 {
				continue
			}

			player := playerMap[s.PlayerID]
			if player.ID == 0 {
				continue
			}

			demand := 30
			if player.Overall > 58 && d.Prestige > 2 {
				demand *= int(d.Prestige)
			} else if player.Overall > 58 && d.Prestige <= 2 {
				demand = util.GenerateIntFromRange(30, 45)
			}

			if demand > pointsRemaining {
				demand = pointsRemaining
			}

			pointsRemaining -= demand
			d.IncrementPool(3, uint8(demand))
			s.AllocatePoints(uint8(demand))

			repository.SaveISLScoutingReportRecord(s, db)

		}
		repository.SaveISLScoutingDeptRecord(d, db)
	}
}

func SyncISLYouthDevelopment() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	players := GetAllYouthDevelopmentPlayers()
	teamMap := GetProfessionalTeamMap()

	for _, p := range players {
		if p.ID == 0 || p.TeamID > 0 || p.TeamAbbr == "DRAFT" {
			continue
		}
		playerID := strconv.Itoa(int(p.ID))
		reports := GetScoutingReportsByPlayerID(playerID)
		qualifyingReports := []structs.ISLScoutingReport{}
		qualifyingPool := 0
		for _, r := range reports {
			if r.RemovedFromBoard || r.CurrentPoints == 0 || !r.Overall {
				continue
			}

			r.IncrementTotalPoints()
			if r.TotalPoints >= uint8(r.PointRequirement) {
				qualifyingReports = append(qualifyingReports, r)
				qualifyingPool += int(r.TotalPoints)
			}
			repository.SaveISLScoutingReportRecord(r, db)
		}

		// Make a decision
		if len(qualifyingReports) > 0 {
			winningTeamID := 0
			sort.Slice(qualifyingReports, func(i, j int) bool {
				return qualifyingReports[i].TotalPoints > qualifyingReports[j].TotalPoints
			})
			roll := util.GenerateIntFromRange(1, qualifyingPool)
			weight := 0
			for _, r := range qualifyingReports {
				weight += int(r.TotalPoints)
				teamID := strconv.Itoa(int(r.TeamID))
				currentRoster := GetAllNBAPlayersByTeamID(teamID)
				if winningTeamID == 0 && roll <= weight && len(currentRoster) < 15 {
					// Player signs with THIS team
					winningTeamID = int(r.TeamID)
				}
			}
			if winningTeamID > 0 {
				team := teamMap[uint(winningTeamID)]
				label := strings.TrimSpace(team.Team + " " + team.Nickname)
				p.SignWithTeam(team.ID, label, false, 0)
				playerLabel := strconv.Itoa(p.Age) + " year old " + p.Position + " " + p.FirstName + " " + p.LastName
				message := "Breaking News! " + playerLabel + " has signed with ISL Team " + label + " in " + team.Country + "!"
				CreateNewsLog("NBA", message, "FreeAgency", 0, ts)

				repository.SaveProfessionalPlayerRecord(p, db)
				yearsRemaining := 2
				if p.Age < 22 {
					yearsRemaining = 22 - p.Age
					if yearsRemaining > 5 {
						yearsRemaining = 5
					}
				}
				year1Salary := 0.0
				year2Salary := 0.0
				year3Salary := 0.0
				year4Salary := 0.0
				year5Salary := 0.0
				for i := 1; i <= yearsRemaining; i++ {
					if i == 1 {
						year1Salary = 0.5
					}
					if i == 2 {
						year2Salary = 0.5
					}
					if i == 3 {
						year3Salary = 0.5
					}
					if i == 4 {
						year4Salary = 0.5
					}
					if i == 5 {
						year5Salary = 0.5
					}
				}

				c := structs.NBAContract{
					PlayerID:       p.ID,
					TeamID:         uint(winningTeamID),
					Team:           label,
					OriginalTeamID: uint(winningTeamID),
					OriginalTeam:   label,
					YearsRemaining: uint(yearsRemaining),
					ContractType:   "ISL",
					TotalRemaining: year1Salary + year2Salary + year3Salary + year4Salary + year5Salary,
					IsActive:       true,
				}
				repository.CreateProfessionalContractRecord(c, db)

				for _, r := range reports {
					r.LockBoard(r.TeamID == uint(winningTeamID))
					repository.SaveISLScoutingReportRecord(r, db)
				}
			}
		}
	}
}

func ISLResetAllPoints() {
	db := dbprovider.GetInstance().GetDB()

	db.Model(&structs.ISLScoutingDept{}).Updates(map[string]interface{}{"resources": 100, "identity_pool": 0, "scouting_pool": 0, "investing_pool": 0})
}

func PickUpISLPlayers() {
	db := dbprovider.GetInstance().GetDB()

	islPlayers := GetAllYouthDevelopmentPlayers()
	internationalTeams := GetInternationalTeams()
	contractUpload := []structs.NBAContract{}

	for _, t := range internationalTeams {
		teamID := strconv.Itoa(int(t.ID))
		roster := GetAllNBAPlayersByTeamID(teamID)

		if len(roster) > 14 {
			continue
		}

		teamNeedsMap := make(map[string]bool)
		positionCount := make(map[string]int)

		if _, ok := teamNeedsMap["PG"]; !ok {
			teamNeedsMap["PG"] = true
		}
		if _, ok := teamNeedsMap["SG"]; !ok {
			teamNeedsMap["SG"] = true
		}
		if _, ok := teamNeedsMap["SF"]; !ok {
			teamNeedsMap["SF"] = true
		}
		if _, ok := teamNeedsMap["PF"]; !ok {
			teamNeedsMap["PF"] = true
		}
		if _, ok := teamNeedsMap["C"]; !ok {
			teamNeedsMap["C"] = true
		}

		if _, ok := positionCount["PG"]; !ok {
			positionCount["PG"] = 0
		}
		if _, ok := positionCount["SG"]; !ok {
			positionCount["SG"] = 0
		}
		if _, ok := positionCount["SF"]; !ok {
			positionCount["SF"] = 0
		}
		if _, ok := positionCount["PF"]; !ok {
			positionCount["PF"] = 0
		}
		if _, ok := positionCount["C"]; !ok {
			positionCount["C"] = 0
		}

		for _, r := range roster {
			positionCount[r.Position] += 1
		}

		if positionCount["PG"] >= 3 {
			teamNeedsMap["PG"] = false
		} else if positionCount["SG"] >= 4 {
			teamNeedsMap["SG"] = false
		} else if positionCount["SF"] >= 4 {
			teamNeedsMap["SF"] = false
		} else if positionCount["PF"] >= 4 {
			teamNeedsMap["PF"] = false
		} else if positionCount["C"] >= 3 {
			teamNeedsMap["C"] = false
		}

		// Team Needs Map Acquired
		// Position Counts acquired
		// Loop through

		for _, p := range islPlayers {
			if !teamNeedsMap[p.Position] {
				continue
			}

			if p.Age > 22 || p.Overall > 71 {
				continue
			}

			odds := 25

			if p.Country == t.Country {
				odds += 25
			}

			dr := util.GenerateIntFromRange(1, 100)

			if odds > dr {
				// Sign Player
				c := structs.NBAContract{
					PlayerID:       p.ID,
					TeamID:         t.ID,
					Team:           t.Team,
					OriginalTeamID: t.ID,
					OriginalTeam:   t.Team,
					YearsRemaining: 3,
					ContractType:   "ISL",
					TotalRemaining: 3,
					Year1Total:     1,
					Year2Total:     1,
					Year3Total:     1,
					IsActive:       true,
				}
				contractUpload = append(contractUpload, c)
				p.SignWithTeam(t.ID, t.Team, false, 1)
				repository.SaveNBAPlayerRecord(p, db)
				positionCount[p.Position] += 1
				if positionCount[p.Position] >= 3 && (p.Position == "PG" || p.Position == "C") {
					teamNeedsMap[p.Position] = false
				} else if positionCount[p.Position] >= 4 && (p.Position == "SG" || p.Position == "SF" || p.Position == "PF") {
					teamNeedsMap[p.Position] = false
				}
			}
		}
	}

	repository.CreateProContractRecordsBatch(db, contractUpload, 100)
}
