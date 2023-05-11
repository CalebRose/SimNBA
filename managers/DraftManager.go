package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

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
