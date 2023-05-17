package managers

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/secrets"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

func ConductDraftLottery() {
	// db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	path := secrets.GetPath()["draftlottery"]
	lotteryCSV := util.ReadCSV(path)
	ts := GetTimestamp()
	lotteryBalls := []structs.DraftLottery{}
	draftPicks := []structs.DraftPick{}

	for idx, row := range lotteryCSV {
		if idx == 0 {
			continue
		}
		teamID := util.ConvertStringToInt(row[0])
		team := row[1]

		// Rows 2-17 of the CSV, the 16 teams for the draft lottery
		if idx < 18 {
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
	for i := 0; i < lotteryPicks; i++ {
		if i <= 3 {
			sum := 0
			for _, l := range lotteryBalls {
				sum += int(l.Chances)
			}

			chance := util.GenerateIntFromRange(1, sum)
			sum2 := 0
			for _, l := range lotteryBalls {
				sum2 += int(l.Chances)
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
		pick := idx + 1
		draftPick := structs.DraftPick{
			SeasonID:       ts.SeasonID,
			Season:         uint(ts.Season),
			DraftRound:     1,
			DraftNumber:    uint(pick),
			TeamID:         do.ID,
			Team:           do.Team,
			OriginalTeamID: do.ID,
			OriginalTeam:   do.Team,
			DraftValue:     0,
		}
		draftPicks = append(draftPicks, draftPick)
	}

	sort.Sort(structs.ByDraftNumber(draftPicks))

	for idx, row := range lotteryCSV {
		if idx < 18 {
			continue
		}
		pickNumber := idx
		teamID := util.ConvertStringToInt(row[0])
		team := row[1]
		round := 1
		if pickNumber > 32 {
			round = 2
		}
		pick := structs.DraftPick{
			SeasonID:       ts.SeasonID,
			Season:         uint(ts.Season),
			DraftRound:     uint(round),
			DraftNumber:    uint(pickNumber),
			TeamID:         uint(teamID),
			Team:           team,
			OriginalTeamID: uint(teamID),
			OriginalTeam:   team,
			DraftValue:     0,
		}
		draftPicks = append(draftPicks, pick)
	}

	for _, pick := range draftPicks {
		fmt.Println(pick)
		// db.Create(&pick)
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

func GetAllCurrentSeasonDraftPicks() []structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	draftPicks := []structs.DraftPick{}

	db.Where("season_id = ?", strconv.Itoa(int(ts.SeasonID))).Find(&draftPicks)

	return draftPicks
}

func GetNBAWarRoomByTeamID(TeamID string) structs.NBAWarRoom {
	db := dbprovider.GetInstance().GetDB()

	warRoom := structs.NBAWarRoom{}

	db.Preload("DraftPicks").Preload("ScoutProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Where("removed_from_board = false")
	}).Where("team_id = ?", TeamID).Find(&warRoom)

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
