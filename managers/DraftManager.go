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

	err := db.Preload("DraftPicks").
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
	if scoutProfile.PlayerID != 0 && scoutProfile.TeamID != 0 {
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

	warRoom := GetNBAWarRoomByTeamID(strconv.Itoa(int(dto.TeamID)))

	if warRoom.ID == 0 || warRoom.SpentPoints >= 100 || warRoom.SpentPoints+dto.Points > 100 {
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
	if ovr < 60 || player.IsRedshirting {
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
