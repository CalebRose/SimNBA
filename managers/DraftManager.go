package managers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/secrets"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func ToggleDraftTime() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	ts.ToggleDraftTime()

	db.Save(&ts)
}

func ConductDraftLottery() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["draftlottery"]
	lotteryCSV := util.ReadCSV(path)
	lotteryBalls := []structs.DraftLottery{}
	draftPicks := []structs.DraftPick{}
	draftMap := GetDraftPickMap()

	for idx, row := range lotteryCSV {
		if idx == 0 {
			continue
		}
		teamID := util.ConvertStringToInt(row[0])
		team := row[1]

		// Rows 1-16 of the CSV, the 16 teams for the draft lottery
		if idx < 17 {
			chances := util.GetLotteryChances(idx)
			lottery := structs.DraftLottery{
				ID:      uint(teamID),
				Team:    team,
				Chances: chances,
			}
			lotteryBalls = append(lotteryBalls, lottery)
		} else {
			break
		}
	}
	lotteryPicks := 16
	draftOrder := []structs.DraftLottery{}
	for i := range lotteryPicks {
		if i <= 3 {
			sum := 0
			for _, l := range lotteryBalls {
				l.ApplyCurrentChance(i)
				sum += int(l.CurrentChance)
			}

			chance := util.GenerateIntFromRange(1, sum)
			sum2 := 0
			for _, l := range lotteryBalls {
				l.ApplyCurrentChance(i)
				sum2 += int(l.CurrentChance)
				if chance < sum2 {
					draftOrder = append(draftOrder, l)
					lotteryBalls = filterLotteryPicks(lotteryBalls, l.ID)
					break
				}
			}
		} else {
			draftOrder = append(draftOrder, lotteryBalls...)
			break
		}
	}

	for idx, do := range draftOrder {
		key := "1 " + do.Team
		pick := idx + 1
		draftPick := draftMap[key]
		draftPick.AssignDraftNumber(uint(pick))
		draftPicks = append(draftPicks, draftPick)
	}

	sort.Sort(structs.ByDraftNumber(draftPicks))

	for idx, row := range lotteryCSV {
		if idx < 17 {
			continue
		}
		pickNumber := idx
		team := row[1]
		roundStr := "1"
		if pickNumber > 32 {
			roundStr = "2"
		}
		key := roundStr + " " + team
		draftpick := draftMap[key]
		draftpick.AssignDraftNumber(uint(pickNumber))
		draftPicks = append(draftPicks, draftpick)
	}

	for _, pick := range draftPicks {
		fmt.Println("Pick " + strconv.Itoa(int(pick.DraftNumber)) + ": " + pick.OriginalTeam)
		db.Save(&pick)
	}
}

// Gets all Current Season and Beyond Draft Picks
func GetDraftPicksByTeamID(TeamID string) []structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(int(ts.SeasonID))
	var picks []structs.DraftPick

	db.Where("team_id = ? AND season_id >= ?", TeamID, seasonID).Find(&picks)

	return picks
}

// Gets all Current Season and Beyond Draft Picks
func GetDraftPickByDraftPickID(DraftPickID string) structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()

	var pick structs.DraftPick

	db.Where("id = ?", DraftPickID).Find(&pick)

	return pick
}

func GenerateDraftLetterGrades() {
	db := dbprovider.GetInstance().GetDB()

	draftees := GetAllNBADraftees()

	for _, d := range draftees {
		s2 := util.GenerateIntFromRange(d.Shooting2-3, d.Shooting2+3)
		s2Grade := util.GetDrafteeGrade(s2)
		s3 := util.GenerateIntFromRange(d.Shooting3-3, d.Shooting3+3)
		s3Grade := util.GetDrafteeGrade(s3)
		ft := util.GenerateIntFromRange(d.FreeThrow-3, d.FreeThrow+3)
		ftGrade := util.GetDrafteeGrade(ft)
		fn := util.GenerateIntFromRange(d.Finishing-3, d.Finishing+3)
		fnGrade := util.GetDrafteeGrade(fn)
		bw := util.GenerateIntFromRange(d.Ballwork-3, d.Ballwork+3)
		bwGrade := util.GetDrafteeGrade(bw)
		rb := util.GenerateIntFromRange(d.Rebounding-3, d.Rebounding+3)
		rbGrade := util.GetDrafteeGrade(rb)
		id := util.GenerateIntFromRange(d.InteriorDefense-3, d.InteriorDefense+3)
		idGrade := util.GetDrafteeGrade(id)
		pd := util.GenerateIntFromRange(d.PerimeterDefense-3, d.PerimeterDefense+3)
		pdGrade := util.GetDrafteeGrade(pd)
		ovrVal := ((s2 + s3 + ft) / 3) + fn + bw + rb + ((id + pd) / 2)
		ovr := util.GetOverallDraftGrade(ovrVal)

		d.ApplyGrades(s2Grade, s3Grade, ftGrade, fnGrade, bwGrade, rbGrade, idGrade, pdGrade, ovr)

		if d.ProPotentialGrade == 0 {
			pot := util.GeneratePotential()
			d.AssignProPotentialGrade(pot)
		}

		d.GetNBAPotentialGrade()

		db.Save(&d)
	}
}

func DraftPredictionRound() {
	db := dbprovider.GetInstance().GetDB()

	draftees := GetAllNBADraftees()

	for _, d := range draftees {
		s2 := util.GenerateIntFromRange(d.Shooting2-3, d.Shooting2+3)
		s3 := util.GenerateIntFromRange(d.Shooting3-3, d.Shooting3+3)
		ft := util.GenerateIntFromRange(d.FreeThrow-3, d.FreeThrow+3)
		fn := util.GenerateIntFromRange(d.Finishing-3, d.Finishing+3)
		bw := util.GenerateIntFromRange(d.Ballwork-3, d.Ballwork+3)
		rb := util.GenerateIntFromRange(d.Rebounding-3, d.Rebounding+3)
		id := util.GenerateIntFromRange(d.InteriorDefense-3, d.InteriorDefense+3)
		pd := util.GenerateIntFromRange(d.PerimeterDefense-3, d.PerimeterDefense+3)
		ovrVal := ((s2 + s3 + ft) / 3) + fn + bw + rb + ((id + pd) / 2)
		round := 0
		if ovrVal > 88 {
			round = 1
		} else if ovrVal > 85 {
			round = 2
		} else if ovrVal > 82 {
			round = 3
		} else if ovrVal > 79 {
			round = 4
		} else if ovrVal > 76 {
			round = 5
		} else if ovrVal > 73 {
			round = 6
		} else {
			round = 7
		}

		d.PredictRound(round)

		db.Save(&d)
	}
}

func GetDraftPickMap() map[string]structs.DraftPick {
	draftPicks := GetCurrentSeasonDraftPickList()
	draftMap := make(map[string]structs.DraftPick)

	for _, pick := range draftPicks {
		if pick.ID == 0 {
			continue
		}
		keyString := strconv.Itoa(int(pick.DraftRound)) + " " + pick.OriginalTeam
		draftMap[keyString] = pick
	}
	return draftMap
}

func GetCurrentSeasonDraftPickList() []structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	draftPicks := []structs.DraftPick{}

	db.Order("draft_number asc").Where("season_id = ?", strconv.Itoa(int(ts.SeasonID))).Find(&draftPicks)

	return draftPicks
}

func GetAllCurrentSeasonDraftPicks() [2][]structs.DraftPick {
	draftPicks := GetCurrentSeasonDraftPickList()

	draftList := [2][]structs.DraftPick{}
	for _, pick := range draftPicks {
		roundIdx := int(pick.DraftRound) - 1
		if roundIdx >= 0 && roundIdx < len(draftList) {
			draftList[roundIdx] = append(draftList[roundIdx], pick)
		} else {
			log.Panicln("Invalid round to insert pick!")
		}

	}

	return draftList
}

func GetOnlyNBAWarRoomByTeamID(TeamID string) structs.NBAWarRoom {
	db := dbprovider.GetInstance().GetDB()

	warRoom := structs.NBAWarRoom{}

	err := db.
		Where("team_id = ?", TeamID).Find(&warRoom).Error
	if err != nil {
		return warRoom
	}

	return warRoom
}

func GetNBAWarRoomByTeamID(TeamID string) structs.NBAWarRoom {
	db := dbprovider.GetInstance().GetDB()

	warRoom := structs.NBAWarRoom{}
	ts := GetTimestamp()

	err := db.Preload("DraftPicks", "season_id = ?", strconv.Itoa(int(ts.SeasonID))).
		Preload("ScoutProfiles.Draftee").
		Preload("ScoutProfiles", "removed_from_board = ?", false).
		Where("team_id = ?", TeamID).Find(&warRoom).Error
	if err != nil {
		return warRoom
	}

	return warRoom
}

func GetNBADrafteesForDraftPage() []structs.NBADraftee {
	db := dbprovider.GetInstance().GetDB()
	draftees := []structs.NBADraftee{}

	db.Find(&draftees)

	sort.Slice(draftees, func(i, j int) bool {
		iVal := util.GetNumericalSortValueByLetterGrade(draftees[i].OverallGrade)
		jVal := util.GetNumericalSortValueByLetterGrade(draftees[j].OverallGrade)
		return iVal < jVal
	})

	return draftees
}

func RunDeclarationsAlgorithm() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()

	for _, c := range collegePlayers {
		if c.IsRedshirting {
			continue
		}
		willDeclare := DetermineIfDeclaring(c)
		if willDeclare {
			c.SetDeclarationStatus()
			db.Save(&c)
		}
	}
}

func GetScoutProfileByScoutProfileID(profileID string) structs.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()

	var scoutProfile structs.ScoutingProfile

	err := db.Where("id = ?", profileID).Find(&scoutProfile).Error
	if err != nil {
		return structs.ScoutingProfile{}
	}

	return scoutProfile
}

func GetOnlyScoutProfileByPlayerIDandTeamID(playerID, teamID string) structs.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()

	var scoutProfile structs.ScoutingProfile

	err := db.Where("player_id = ? AND team_id = ?", playerID, teamID).Error
	if err != nil {
		return structs.ScoutingProfile{}
	}

	return scoutProfile
}

func CreateScoutingProfile(dto structs.ScoutingProfileDTO) structs.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := GetOnlyScoutProfileByPlayerIDandTeamID(strconv.Itoa(int(dto.PlayerID)), strconv.Itoa(int(dto.TeamID)))

	// If Recruit Already Exists
	if scoutProfile.PlayerID > 0 && scoutProfile.TeamID > 0 {
		scoutProfile.ReplaceOnBoard()
		db.Save(&scoutProfile)
		return scoutProfile
	}

	newScoutingProfile := structs.ScoutingProfile{
		PlayerID:         dto.PlayerID,
		TeamID:           dto.TeamID,
		ShowCount:        0,
		RemovedFromBoard: false,
	}

	db.Create(&newScoutingProfile)

	return newScoutingProfile
}

func RemovePlayerFromScoutBoard(id string) {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := GetScoutProfileByScoutProfileID(id)

	scoutProfile.RemoveFromBoard()

	db.Save(&scoutProfile)
}

func GetScoutingDataByPlayerID(id string) structs.ScoutingDataResponse {
	ts := GetTimestamp()

	lastSeasonID := ts.SeasonID - 1
	lastSeasonIDSTR := strconv.Itoa(int(lastSeasonID))

	draftee := GetHistoricCollegePlayerByID(id)

	seasonStats := GetPlayerSeasonStatsByPlayerID(id, lastSeasonIDSTR)
	teamID := strconv.Itoa(int(draftee.TeamID))
	collegeStandings := GetStandingsRecordByTeamID(teamID, lastSeasonIDSTR)

	return structs.ScoutingDataResponse{
		DrafteeSeasonStats: seasonStats,
		TeamStandings:      collegeStandings,
	}
}

func RevealScoutingAttribute(dto structs.RevealAttributeDTO) bool {
	db := dbprovider.GetInstance().GetDB()

	scoutProfile := GetScoutProfileByScoutProfileID(strconv.Itoa(int(dto.ScoutProfileID)))

	if scoutProfile.ID == 0 {
		return false
	}

	scoutProfile.RevealAttribute(dto.Attribute)

	warRoom := GetOnlyNBAWarRoomByTeamID(strconv.Itoa(int(dto.TeamID)))

	if warRoom.ID == 0 || warRoom.SpentPoints >= warRoom.ScoutingPoints || warRoom.SpentPoints+dto.Points > warRoom.ScoutingPoints {
		return false
	}

	warRoom.AddToSpentPoints(dto.Points)

	err := db.Save(&scoutProfile).Error
	if err != nil {
		return false
	}
	err = db.Save(&warRoom).Error
	return err == nil
}

func DetermineIfDeclaring(player structs.CollegePlayer) bool {
	// Redshirt senior or just a senior
	if (player.IsRedshirt && player.Year == 5) || (!player.IsRedshirt && player.Year == 4 && !player.IsRedshirting) {
		return true
	}
	ovr := player.Overall
	if ovr < 64 || player.IsRedshirting {
		return false
	}
	odds := util.GenerateIntFromRange(1, 100)
	if ovr > 64 && odds <= 20 {
		return true
	} else if ovr > 68 && odds <= 35 {
		return true
	} else if ovr > 70 && odds <= 40 {
		return true
	} else if ovr > 74 && odds <= 50 {
		return true
	} else if ovr > 76 && odds <= 75 {
		return true
	} else if ovr > 80 && odds <= 80 {
		return true
	} else if ovr > 82 && odds <= 85 {
		return true
	} else if ovr > 84 && odds <= 95 {
		return true
	} else if ovr > 89 {
		return true
	}
	return false
}

func InternationalDeclaration(player structs.NBAPlayer, isEligible bool) bool {
	// Redshirt senior or just a senior
	if !isEligible || player.TeamID < 33 || player.Age < 18 {
		return false
	}
	ovr := player.Overall
	if ovr < 60 {
		return false
	}
	odds := util.GenerateIntFromRange(1, 100)
	if ovr > 60 && odds <= 20 {
		return true
	} else if ovr > 64 && odds <= 35 {
		return true
	} else if ovr > 67 && odds <= 40 {
		return true
	} else if ovr > 69 && odds <= 50 {
		return true
	} else if ovr > 72 && odds <= 75 {
		return true
	} else if ovr > 74 && odds <= 80 {
		return true
	} else if ovr > 76 && odds <= 85 {
		return true
	} else if ovr > 79 && odds <= 95 {
		return true
	} else if ovr > 84 {
		return true
	}
	return false
}

func ExportDraftedPlayers(picks []structs.DraftPick) bool {
	db := dbprovider.GetInstance().GetDB()

	// for _, pick := range picks {
	// 	playerId := strconv.Itoa(int(pick.SelectedPlayerID))
	// 	draftee := GetNBADrafteeByID(playerId)
	// 	nbaPlayer := structs.NBAPlayer{}
	// 	if draftee.College == "DRAFT" {
	// 		// Get International Player Record
	// 		nbaPlayer = GetNBAPlayerByID(playerId)
	// 		nbaPlayer.DraftInternationalPlayer(pick.ID, pick.DraftRound, pick.DraftNumber, pick.TeamID, pick.Team)
	// 		repository.SaveProfessionalPlayerRecord(nbaPlayer, db)
	// 	} else {
	// 		nbaPlayer = structs.NBAPlayer{
	// 			BasePlayer:      draftee.BasePlayer, // Assuming BasePlayer fields are common
	// 			PlayerID:        draftee.PlayerID,
	// 			TeamID:          pick.TeamID,
	// 			TeamAbbr:        pick.Team,
	// 			CollegeID:       draftee.CollegeID,
	// 			College:         draftee.College,
	// 			DraftPickID:     draftee.DraftPickID,
	// 			DraftedTeamID:   pick.TeamID,
	// 			DraftedTeamAbbr: pick.Team,
	// 			DraftedRound:    pick.DraftRound,
	// 			DraftPick:       pick.DraftNumber,
	// 			PrimeAge:        uint(draftee.PrimeAge),
	// 			IsNBA:           true,
	// 		}
	// 		nbaPlayer.SetID(draftee.PlayerID)
	// 		repository.CreateProfessionalPlayerRecord(nbaPlayer, db)
	// 	}
	// 	draftee.AssignDraftedTeam(strconv.Itoa(int(pick.DraftNumber)), pick.ID, pick.TeamID, pick.Team)
	// 	db.Save(&draftee)
	// 	year1Salary := util.GetDrafteeSalary(pick.DraftNumber, 1)
	// 	year2Salary := util.GetDrafteeSalary(pick.DraftNumber, 2)
	// 	year3Salary := util.GetDrafteeSalary(pick.DraftNumber, 3)
	// 	year4Salary := util.GetDrafteeSalary(pick.DraftNumber, 4)
	// 	yearsRemaining := util.GetYearsRemainingForDrafteeContract(pick.DraftNumber)
	// 	contract := structs.NBAContract{
	// 		PlayerID:       nbaPlayer.PlayerID,
	// 		TeamID:         nbaPlayer.TeamID,
	// 		Team:           nbaPlayer.TeamAbbr,
	// 		OriginalTeamID: nbaPlayer.TeamID,
	// 		OriginalTeam:   nbaPlayer.TeamAbbr,
	// 		YearsRemaining: yearsRemaining,
	// 		ContractType:   "Rookie",
	// 		TotalRemaining: year1Salary + year2Salary + year3Salary + year4Salary,
	// 		Year1Total:     year1Salary,
	// 		Year2Total:     year2Salary,
	// 		Year3Total:     year3Salary,
	// 		Year4Total:     year4Salary,
	// 		Year3Opt:       true,
	// 		Year4Opt:       true,
	// 		IsActive:       true,
	// 	}

	// 	repository.CreateProfessionalContractRecord(contract, db)
	// }

	draftablePlayers := GetAllNBADraftees()

	for _, draftee := range draftablePlayers {
		if draftee.DraftedTeamID > 0 || draftee.ID < 154 {
			continue
		}
		playerID := strconv.Itoa(int(draftee.ID))
		nbaPlayer := GetNBAPlayerByID(playerID)
		if nbaPlayer.ID > 0 && nbaPlayer.TeamAbbr != "DRAFT" {
			continue
		}
		if draftee.College == "DRAFT" {
			nbaPlayer.BecomeUDFA()
		} else {
			nbaPlayer = structs.NBAPlayer{
				BasePlayer:        draftee.BasePlayer, // Assuming BasePlayer fields are common
				PlayerID:          draftee.PlayerID,
				TeamID:            0,
				TeamAbbr:          "FA",
				CollegeID:         draftee.CollegeID,
				College:           draftee.College,
				DraftPickID:       draftee.DraftPickID,
				DraftedTeamID:     draftee.DraftedTeamID,
				DraftedTeamAbbr:   draftee.DraftedTeamAbbr,
				PrimeAge:          uint(draftee.PrimeAge),
				IsNBA:             true,
				IsNegotiating:     false,
				IsAcceptingOffers: true,
				IsFreeAgent:       true,
				MinimumValue:      0.7,
			}
			nbaPlayer.SetID(draftee.PlayerID)
		}

		NegotiationRound := 0
		if draftee.Overall < 80 {
			NegotiationRound = util.GenerateIntFromRange(2, 4)
		} else {
			NegotiationRound = util.GenerateIntFromRange(3, 6)
		}

		SigningRound := 7
		nbaPlayer.AssignFAPreferences(uint(NegotiationRound), uint(SigningRound))
		if draftee.College == "DRAFT" {
			repository.SaveProfessionalPlayerRecord(nbaPlayer, db)
		} else {
			repository.CreateProfessionalPlayerRecord(nbaPlayer, db)
		}
	}

	return true
}

func NBACombineForDraft() {
	db := dbprovider.GetInstance().GetDB()

	draftees := GetAllNBADraftees()
	combineGrades := []structs.NBACombineResults{}

	for _, draftee := range draftees {
		// Disregard all candidates under 55 overall
		if draftee.Overall < 55 {
			continue
		}
		strength := GetCombineValue(draftee.OverallGrade, draftee.Overall, true)
		agility := GetCombineValue(draftee.OverallGrade, draftee.Overall, true)
		shooting2 := GetCombineValue(draftee.Shooting2Grade, draftee.Shooting2, false)
		shooting3 := GetCombineValue(draftee.Shooting3Grade, draftee.Shooting3, false)
		passing := GetCombineValue(draftee.BallworkGrade, draftee.Ballwork, false)
		blocking := GetCombineValue(draftee.InteriorDefenseGrade, draftee.PerimeterDefense, false)
		stealing := GetCombineValue(draftee.BallworkGrade, draftee.Ballwork, false)
		rebound := GetCombineValue(draftee.ReboundingGrade, draftee.Rebounding, false)

		// Shooting Drills
		twoPointShootingCount := 0
		threePointShootingCount := 0
		maxReps := 30
		successfulAssists := 0
		successfulBlocks := 0
		successfulSteals := 0
		for range 30 {
			twoPointDR := util.GenerateIntFromRange(1, 20) + int(GetCombineModifier(shooting2))
			threePointDR := util.GenerateIntFromRange(1, 20) + int(GetCombineModifier(shooting3))
			if twoPointDR > 15 {
				twoPointShootingCount++
			}
			if threePointDR > 15 {
				threePointShootingCount++
			}
			assistsDr := util.GenerateIntFromRange(1, 20) + int(GetCombineModifier(passing))
			if assistsDr > 15 {
				successfulAssists++
			}
			stealsDr := util.GenerateIntFromRange(1, 20) + int(GetCombineModifier(stealing))
			if stealsDr > 15 {
				successfulSteals++
			}
			blocksDr := util.GenerateIntFromRange(1, 20) + int(GetCombineModifier(blocking))
			if blocksDr > 15 {
				successfulBlocks++
			}
			benchPressDr := util.GenerateIntFromRange(1, 100) + strength
			if benchPressDr < 120 {
				maxReps--
			}
		}

		standingVerticalLeapMod := GetCombineModifier(rebound)
		laneAgilityMod := GetCombineModifier(agility)
		shuttleRunMod := GetCombineModifier(agility)
		standingVerticalLeap := CalculateStandingVerticalLeap(standingVerticalLeapMod)
		maxVerticalLeap := CalculateMaxVerticalLeap(standingVerticalLeapMod)
		laneAgilityTime := CalculateLaneAgility(laneAgilityMod)
		shuttleRunTime := CalculateShuttleRun(shuttleRunMod)

		combineGrade := structs.NBACombineResults{
			PlayerID:             draftee.ID,
			TwoPointShooting:     uint8(twoPointShootingCount),
			ThreePointShooting:   uint8(threePointShootingCount),
			PassingDrills:        uint8(successfulAssists),
			BlockingDrills:       uint8(successfulBlocks),
			StealDrills:          uint8(successfulSteals),
			BenchPress:           uint8(maxReps),
			StandingVerticalLeap: float32(standingVerticalLeap),
			MaxVerticalLeap:      float32(maxVerticalLeap),
			LaneAgility:          float32(laneAgilityTime),
			ShuttleRun:           float32(shuttleRunTime),
		}
		combineGrades = append(combineGrades, combineGrade)
	}

	repository.CreateNBACombineRecordsBatch(db, combineGrades, 250)
}

func GetCombineModifier(value int) float64 {
	return math.Log(float64(value)+1) * 1.7
}

func GetCombineValue(grade string, value int, isOverall bool) int {
	gradeVal := GetValueFromGrade(grade)
	if gradeVal == value {
		return value
	}
	if isOverall {
		return util.GenerateIntFromRange(value-15, value+15)
	}
	min := math.Min(float64(gradeVal), float64(value))
	max := math.Max(float64(gradeVal), float64(value))
	return util.GenerateIntFromRange(int(min), int(max))
}

func GetValueFromGrade(grade string) int {
	if grade == "A+" {
		return 24
	} else if grade == "A" {
		return 22
	} else if grade == "A-" {
		return 20
	} else if grade == "B+" {
		return 18
	} else if grade == "B" {
		return 16
	} else if grade == "B-" {
		return 14
	} else if grade == "C+" {
		return 12
	} else if grade == "C" {
		return 10
	} else if grade == "C-" {
		return 8
	} else if grade == "D" {
		return 5
	}
	return 1
}

func CalculateStandingVerticalLeap(modifier float64) float64 {
	base := 24.0
	max := 42.0
	factor := 0.1

	// mean = base + modifier*factor
	mean := base + float64(modifier)*factor
	// choose a standard deviation (e.g. 4 inches)
	sd := 4.0

	leap := mean + rand.NormFloat64()*sd
	// clamp into [base, max]
	return math.Min(math.Max(leap, base), max)
}

func CalculateMaxVerticalLeap(modifier float64) float64 {
	base := 30.0
	max := 42.0
	factor := 0.12

	mean := base + float64(modifier)*factor
	sd := 5.0 // a bit wider spread
	leap := mean + rand.NormFloat64()*sd
	return math.Min(math.Max(leap, base), max)
}

func CalculateLaneAgility(modifier float64) float64 {
	base := 12.0
	best := 10.0
	factor := 0.02

	mean := base - float64(modifier)*factor
	sd := 0.5 // half-second std dev

	t := mean + rand.NormFloat64()*sd
	return math.Max(t, best)
}

func CalculateShuttleRun(modifier float64) float64 {
	base := 3.7
	best := 2.9
	factor := 0.015

	centre := base - float64(modifier)*factor
	spread := 0.7 // ±0.4 seconds

	min := math.Max(best, centre-spread)
	maxVal := centre + spread
	return min + rand.Float64()*(maxVal-min)
}
