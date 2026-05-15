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
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

func ToggleDraftTime() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	ts.ToggleDraftTime()

	db.Save(&ts)
}

func ConductDraftLottery() {
	db := dbprovider.GetInstance().GetDB()

	// Get the pre-lottery base order from standings/series data.
	// Lottery teams (Chances set) come first, sorted worst→best record.
	// Playoff teams (Chances empty) follow, sorted by elimination round then record.
	baseOrder := GetDraftLotteryOrder(db)
	draftMap := GetDraftPickMap()

	// Split into lottery balls and playoff order based on whether Chances is set.
	lotteryBalls := []structs.DraftLottery{}
	playoffOrder := []structs.DraftLottery{}
	for _, entry := range baseOrder {
		if len(entry.Chances) > 0 {
			lotteryBalls = append(lotteryBalls, entry)
		} else {
			playoffOrder = append(playoffOrder, entry)
		}
	}

	// Run weighted lottery for picks 1-4.
	finalRound1 := []structs.DraftLottery{}
	remaining := make([]structs.DraftLottery, len(lotteryBalls))
	copy(remaining, lotteryBalls)

	for i := 0; i < 4; i++ {
		sum := 0
		for _, l := range remaining {
			l.ApplyCurrentChance(i)
			sum += int(l.CurrentChance)
		}
		chance := util.GenerateIntFromRange(1, sum)
		sum2 := 0
		for _, l := range remaining {
			l.ApplyCurrentChance(i)
			sum2 += int(l.CurrentChance)
			if chance < sum2 {
				finalRound1 = append(finalRound1, l)
				remaining = filterLotteryPicks(remaining, l.ID)
				break
			}
		}
	}
	// Picks 5-N: remaining lottery teams keep their base order.
	finalRound1 = append(finalRound1, remaining...)
	// Picks N+1 through 32: playoff teams in elimination order.
	finalRound1 = append(finalRound1, playoffOrder...)

	// Assign Round 1 draft numbers and collect pick records.
	draftPicks := []structs.DraftPick{}
	round1PickMap := make(map[uint]uint) // teamID → R1 pick number (for R2 tiebreaking)
	for idx, team := range finalRound1 {
		pickNum := uint(idx + 1)
		round1PickMap[team.ID] = pickNum
		key := "1 " + strconv.Itoa(int(team.ID))
		pick := draftMap[key]
		pick.AssignDraftNumber(pickNum)
		draftPicks = append(draftPicks, pick)
	}

	// Build Round 2 order: all teams sorted by record (worst first), no playoff
	// consideration. Ties broken by reverse Round 1 pick (higher R1 pick = earlier R2 pick).
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID - 1))
	allStandings := GetNBAStandingsBySeasonID(seasonID)
	standingsMap := make(map[uint]structs.NBAStandings)
	for _, s := range allStandings {
		if s.TeamID <= 32 {
			standingsMap[s.TeamID] = s
		}
	}

	round2Order := make([]structs.DraftLottery, len(finalRound1))
	copy(round2Order, finalRound1)
	sort.Slice(round2Order, func(i, j int) bool {
		si := standingsMap[round2Order[i].ID]
		sj := standingsMap[round2Order[j].ID]
		if si.TotalWins != sj.TotalWins {
			return si.TotalWins < sj.TotalWins
		}
		return round1PickMap[round2Order[i].ID] > round1PickMap[round2Order[j].ID]
	})

	for idx, team := range round2Order {
		pickNum := uint(idx + 1)
		key := "2 " + strconv.Itoa(int(team.ID))
		pick := draftMap[key]
		pick.AssignDraftNumber(pickNum)
		draftPicks = append(draftPicks, pick)
	}

	sort.Sort(structs.ByDraftNumber(draftPicks))

	for _, pick := range draftPicks {
		fmt.Println("Pick " + strconv.Itoa(int(pick.DraftNumber)) + ": " + pick.OriginalTeam)
		db.Save(&pick)
	}

	// Create a draft lottery forum thread (best-effort).
	season, picks := ts.Season, draftPicks
	go CreateDraftLotteryForumThread(season, picks)
}

func GetDraftLotteryOrder(db *gorm.DB) []structs.DraftLottery {
	lotteryOrder := []structs.DraftLottery{}

	nbaTeams := GetAllActiveNBATeams()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID - 1))
	nbaSeries := GetNBASeriesBySeasonID(seasonID)
	nbaGames := GetNBAMatchesBySeasonID(seasonID)

	// Build full team name map (City + Nickname) for NBA teams only (ID <= 32)
	teamNameMap := make(map[uint]string)
	for _, t := range nbaTeams {
		if t.ID <= 32 {
			teamNameMap[t.ID] = t.Team + " " + t.Nickname
		}
	}

	// Load previous season standings for NBA teams (ID <= 32)
	allStandings := GetNBAStandingsBySeasonID(seasonID)
	standingsMap := make(map[uint]structs.NBAStandings)
	nbaStandings := []structs.NBAStandings{}
	for _, s := range allStandings {
		if s.TeamID <= 32 {
			standingsMap[s.TeamID] = s
			nbaStandings = append(nbaStandings, s)
		}
	}

	// Build head-to-head win map from regular season games only.
	// Exclude playoff games, play-in games (week 18), and international teams.
	headToHead := make(map[uint]map[uint]int)
	for _, g := range nbaGames {
		if g.IsPlayoffGame || g.IsPlayInGame || g.Week >= 18 || !g.GameComplete {
			continue
		}
		if g.HomeTeamID > 32 || g.AwayTeamID > 32 {
			continue
		}
		if headToHead[g.HomeTeamID] == nil {
			headToHead[g.HomeTeamID] = make(map[uint]int)
		}
		if headToHead[g.AwayTeamID] == nil {
			headToHead[g.AwayTeamID] = make(map[uint]int)
		}
		if g.HomeTeamWin {
			headToHead[g.HomeTeamID][g.AwayTeamID]++
		} else if g.AwayTeamWin {
			headToHead[g.AwayTeamID][g.HomeTeamID]++
		}
	}

	// Count playoff series wins per team to determine the round they were eliminated.
	// 0 wins = Round 1 loser, 1 = Round 2 loser, 2 = Conference Finals loser,
	// 3 = Finals loser, 4 = Champion.
	seriesWins := make(map[uint]int)
	for _, s := range nbaSeries {
		if !s.IsPlayoffGame || s.IsInternational || !s.SeriesComplete {
			continue
		}
		if s.HomeTeamWin && s.HomeTeamID <= 32 {
			seriesWins[s.HomeTeamID]++
		} else if !s.HomeTeamWin && s.AwayTeamID <= 32 {
			seriesWins[s.AwayTeamID]++
		}
	}

	// Collect all NBA teams that appeared in a playoff series
	type playoffEntry struct {
		TeamID   uint
		TeamName string
		Wins     int
		Standing structs.NBAStandings
	}
	playoffMap := make(map[uint]playoffEntry)
	for _, s := range nbaSeries {
		if !s.IsPlayoffGame || s.IsInternational || !s.SeriesComplete {
			continue
		}
		for _, id := range []uint{s.HomeTeamID, s.AwayTeamID} {
			if id > 32 {
				continue
			}
			name := teamNameMap[id]
			if name == "" {
				if id == s.HomeTeamID {
					name = s.HomeTeam
				} else {
					name = s.AwayTeam
				}
			}
			playoffMap[id] = playoffEntry{
				TeamID:   id,
				TeamName: name,
				Wins:     seriesWins[id],
				Standing: standingsMap[id],
			}
		}
	}

	playoffTeamIDs := make(map[uint]bool)
	for id := range playoffMap {
		playoffTeamIDs[id] = true
	}

	// --- BOTTOM 16: Lottery teams ---
	// Any NBA team (ID <= 32) not in a playoff series is a lottery team.
	// This naturally includes play-in losers (they never appear in NBASeries).
	lotteryStandings := []structs.NBAStandings{}
	for _, s := range nbaStandings {
		if !playoffTeamIDs[s.TeamID] {
			lotteryStandings = append(lotteryStandings, s)
		}
	}

	// Sort: worst record first (fewest wins).
	// Tiebreaker 1: head-to-head wins between the tied teams.
	// Tiebreaker 2: point differential (lower = picks earlier).
	// Tiebreaker 3: random coin flip.
	sort.Slice(lotteryStandings, func(i, j int) bool {
		si, sj := lotteryStandings[i], lotteryStandings[j]
		if si.TotalWins != sj.TotalWins {
			return si.TotalWins < sj.TotalWins
		}
		h2hI := headToHead[si.TeamID][sj.TeamID]
		h2hJ := headToHead[sj.TeamID][si.TeamID]
		if h2hI != h2hJ {
			return h2hI < h2hJ
		}
		if si.PointsFor-si.PointsAgainst != sj.PointsFor-sj.PointsAgainst {
			return (si.PointsFor - si.PointsAgainst) < (sj.PointsFor - sj.PointsAgainst)
		}
		return rand.Intn(2) == 0
	})

	// Assign lottery ball chances based on position (1 = worst team, 16 = best lottery team)
	for idx, s := range lotteryStandings {
		chances := util.GetLotteryChances(idx + 1)
		name := teamNameMap[s.TeamID]
		if name == "" {
			name = s.TeamAbbr
		}
		lotteryOrder = append(lotteryOrder, structs.DraftLottery{
			ID:      s.TeamID,
			Team:    name,
			Chances: chances,
		})
	}

	// --- TOP 16: Playoff teams ---
	// Grouped by series wins (proxy for elimination round).
	// Within each round group, sort by regular season record (worst first).
	sortPlayoffGroup := func(group []playoffEntry) {
		sort.Slice(group, func(i, j int) bool {
			si, sj := group[i].Standing, group[j].Standing
			if si.TotalWins != sj.TotalWins {
				return si.TotalWins < sj.TotalWins
			}
			h2hI := headToHead[group[i].TeamID][group[j].TeamID]
			h2hJ := headToHead[group[j].TeamID][group[i].TeamID]
			if h2hI != h2hJ {
				return h2hI < h2hJ
			}
			return (si.PointsFor - si.PointsAgainst) < (sj.PointsFor - sj.PointsAgainst)
		})
	}

	roundGroups := make(map[int][]playoffEntry)
	for _, entry := range playoffMap {
		if entry.TeamID > 32 {
			continue
		}
		roundGroups[entry.Wins] = append(roundGroups[entry.Wins], entry)
	}

	// Append in order: R1 losers (0 wins) → Semis losers (1) → CF losers (2) → Finals loser (3) → Champion (4)
	for wins := 0; wins <= 4; wins++ {
		group := roundGroups[wins]
		sortPlayoffGroup(group)
		for _, entry := range group {
			lotteryOrder = append(lotteryOrder, structs.DraftLottery{
				ID:      entry.TeamID,
				Team:    entry.TeamName,
				Chances: []uint{},
			})
		}
	}

	return lotteryOrder
}

func GetDraftPickMap() map[string]structs.DraftPick {
	draftPicks := GetCurrentSeasonDraftPickList()
	draftMap := make(map[string]structs.DraftPick)

	for _, pick := range draftPicks {
		if pick.ID == 0 {
			continue
		}
		keyString := strconv.Itoa(int(pick.DraftRound)) + " " + strconv.Itoa(int(pick.OriginalTeamID))
		draftMap[keyString] = pick
	}
	return draftMap
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
		s2 := util.GenerateIntFromRange(int(d.MidRangeShooting)-3, int(d.MidRangeShooting)+3)
		s2Grade := util.GetDrafteeGrade(uint8(s2))
		s3 := util.GenerateIntFromRange(int(d.ThreePointShooting)-3, int(d.ThreePointShooting)+3)
		s3Grade := util.GetDrafteeGrade(uint8(s3))
		ft := util.GenerateIntFromRange(int(d.FreeThrow)-3, int(d.FreeThrow)+3)
		ftGrade := util.GetDrafteeGrade(uint8(ft))
		fn := util.GenerateIntFromRange(int(d.InsideShooting)-3, int(d.InsideShooting)+3)
		fnGrade := util.GetDrafteeGrade(uint8(fn))
		bw := util.GenerateIntFromRange(int(d.Ballwork)-3, int(d.Ballwork)+3)
		bwGrade := util.GetDrafteeGrade(uint8(bw))
		rb := util.GenerateIntFromRange(int(d.Rebounding)-3, int(d.Rebounding)+3)
		rbGrade := util.GetDrafteeGrade(uint8(rb))
		id := util.GenerateIntFromRange(int(d.InteriorDefense)-3, int(d.InteriorDefense)+3)
		idGrade := util.GetDrafteeGrade(uint8(id))
		pd := util.GenerateIntFromRange(int(d.PerimeterDefense)-3, int(d.PerimeterDefense)+3)
		pdGrade := util.GetDrafteeGrade(uint8(pd))
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
		s2 := util.GenerateIntFromRange(int(d.MidRangeShooting)-3, int(d.MidRangeShooting)+3)
		s3 := util.GenerateIntFromRange(int(d.ThreePointShooting)-3, int(d.ThreePointShooting)+3)
		ft := util.GenerateIntFromRange(int(d.FreeThrow)-3, int(d.FreeThrow)+3)
		fn := util.GenerateIntFromRange(int(d.InsideShooting)-3, int(d.InsideShooting)+3)
		bw := util.GenerateIntFromRange(int(d.Ballwork)-3, int(d.Ballwork)+3)
		rb := util.GenerateIntFromRange(int(d.Rebounding)-3, int(d.Rebounding)+3)
		id := util.GenerateIntFromRange(int(d.InteriorDefense)-3, int(d.InteriorDefense)+3)
		pd := util.GenerateIntFromRange(int(d.PerimeterDefense)-3, int(d.PerimeterDefense)+3)
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

func GetAllRelevantDraftPicks() []structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	seasonID := strconv.Itoa(int(ts.SeasonID))

	seasonIDPlusFive := strconv.Itoa(int(ts.SeasonID + 5))

	draftPicks := []structs.DraftPick{}

	db.Order("season_id asc").Order("draft_number asc").Where("season_id >= ? AND season_id <= ?", seasonID, seasonIDPlusFive).Find(&draftPicks)

	return draftPicks
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

func GetAllNBAWarRooms() []structs.NBAWarRoom {
	db := dbprovider.GetInstance().GetDB()

	warRoom := []structs.NBAWarRoom{}

	err := db.Find(&warRoom).Error
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
			repository.SaveCollegePlayerRecord(c, db)
		}
	}
}

func GetAllScoutingProfiles() []structs.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()

	profiles := []structs.ScoutingProfile{}
	err := db.Where("removed_from_board = ?", false).
		Find(&profiles).Error
	if err != nil {
		return profiles
	}

	return profiles
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
	if ovr < 20 || player.IsRedshirting {
		return false
	}
	odds := util.GenerateIntFromRange(1, 100)
	if ovr > 19 && odds <= 2 {
		return true
	} else if ovr > 21 && odds <= 5 {
		return true
	} else if ovr > 22 && odds <= 8 {
		return true
	} else if ovr > 23 && odds <= 15 {
		return true
	} else if ovr > 24 && odds <= 20 {
		return true
	} else if ovr > 25 && odds <= 50 {
		return true
	} else if ovr > 26 && odds <= 75 {
		return true
	} else if ovr > 27 && odds <= 80 {
		return true
	} else if ovr > 28 && odds <= 85 {
		return true
	} else if ovr > 29 && odds <= 95 {
		return true
	} else if ovr > 30 {
		return true
	}
	return false
}

func InternationalDeclaration(player structs.NBAPlayer, isEligible bool) bool {
	// Redshirt senior or just a senior
	if !isEligible || player.TeamID < 33 || player.Age < 18 || player.IsIntDeclared {
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

	for _, pick := range picks {
		if pick.SelectedPlayerID == 2445 {
			continue
		}
		playerId := strconv.Itoa(int(pick.SelectedPlayerID))
		draftee := GetNBADrafteeByID(playerId)
		nbaPlayer := structs.NBAPlayer{}
		if draftee.College == "DRAFT" {
			// Get International Player Record
			nbaPlayer = GetNBAPlayerByID(playerId)
			nbaPlayer.DraftInternationalPlayer(pick.ID, pick.DraftRound, pick.DraftNumber, pick.TeamID, pick.Team)
			repository.SaveProfessionalPlayerRecord(nbaPlayer, db)
		} else {
			nbaPlayer = structs.NBAPlayer{
				BasePlayer:    draftee.BasePlayer, // Assuming BasePlayer fields are common
				CollegeID:     draftee.CollegeID,
				College:       draftee.College,
				DraftPickID:   draftee.DraftPickID,
				DraftedTeamID: pick.TeamID,
				DraftedTeam:   pick.Team,
				DraftedRound:  pick.DraftRound,
				DraftPick:     pick.DraftNumber,
				IsNBA:         true,
			}
			nbaPlayer.SetID(draftee.PlayerID)
			repository.CreateProfessionalPlayerRecord(nbaPlayer, db)
		}
		draftee.AssignDraftedTeam(strconv.Itoa(int(pick.DraftNumber)), pick.ID, pick.TeamID, pick.Team)
		db.Save(&draftee)
		year1Salary := util.GetDrafteeSalary(pick.DraftNumber, 1)
		year2Salary := util.GetDrafteeSalary(pick.DraftNumber, 2)
		year3Salary := util.GetDrafteeSalary(pick.DraftNumber, 3)
		year4Salary := util.GetDrafteeSalary(pick.DraftNumber, 4)
		yearsRemaining := util.GetYearsRemainingForDrafteeContract(pick.DraftNumber)
		contract := structs.NBAContract{
			PlayerID:       nbaPlayer.PlayerID,
			TeamID:         nbaPlayer.TeamID,
			Team:           nbaPlayer.Team,
			OriginalTeamID: nbaPlayer.TeamID,
			OriginalTeam:   nbaPlayer.Team,
			YearsRemaining: yearsRemaining,
			ContractType:   "Rookie",
			TotalRemaining: year1Salary + year2Salary + year3Salary + year4Salary,
			Year1Total:     year1Salary,
			Year2Total:     year2Salary,
			Year3Total:     year3Salary,
			Year4Total:     year4Salary,
			Year3Opt:       true,
			Year4Opt:       true,
			IsActive:       true,
		}

		repository.CreateProfessionalContractRecord(contract, db)
	}

	draftablePlayers := GetAllNBADraftees()

	for _, draftee := range draftablePlayers {
		if draftee.DraftedTeamID > 0 {
			continue
		}
		playerID := strconv.Itoa(int(draftee.ID))
		nbaPlayer := GetNBAPlayerByID(playerID)
		if nbaPlayer.ID > 0 && nbaPlayer.Team != "DRAFT" {
			continue
		}
		if draftee.College == "DRAFT" {
			nbaPlayer.BecomeUDFA()
		} else {
			nbaPlayer = structs.NBAPlayer{
				BasePlayer:        draftee.BasePlayer, // Assuming BasePlayer fields are common
				CollegeID:         draftee.CollegeID,
				College:           draftee.College,
				DraftPickID:       draftee.DraftPickID,
				DraftedTeamID:     draftee.DraftedTeamID,
				DraftedTeam:       draftee.DraftedTeam,
				IsNBA:             true,
				IsNegotiating:     false,
				IsAcceptingOffers: true,
				IsFreeAgent:       true,
				MinimumValue:      0.7,
			}
			nbaPlayer.PlayerID = draftee.PlayerID
			nbaPlayer.TeamID = 0
			nbaPlayer.Team = "FA"
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
		shooting2 := GetCombineValue(draftee.MidrangeShootingGrade, draftee.MidRangeShooting, false)
		shooting3 := GetCombineValue(draftee.ThreePointShootingGrade, draftee.ThreePointShooting, false)
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
			benchPressDr := util.GenerateIntFromRange(1, 100) + int(strength)
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

func GetCombineModifier(value uint8) float64 {
	return math.Log(float64(value)+1) * 1.7
}

func GetCombineValue(grade string, value uint8, isOverall bool) uint8 {
	gradeVal := GetValueFromGrade(grade)
	if gradeVal == value {
		return value
	}
	if isOverall {
		return uint8(util.GenerateIntFromRange(int(value)-15, int(value)+15))
	}
	min := math.Min(float64(gradeVal), float64(value))
	max := math.Max(float64(gradeVal), float64(value))
	return uint8(util.GenerateIntFromRange(int(min), int(max)))
}

func GetValueFromGrade(grade string) uint8 {
	switch grade {
	case "A+":
		return 24
	case "A":
		return 22
	case "A-":
		return 20
	case "B+":
		return 18
	case "B":
		return 16
	case "B-":
		return 14
	case "C+":
		return 12
	case "C":
		return 10
	case "C-":
		return 8
	case "D":
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
