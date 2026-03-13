package repository

import (
	"log"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

type RecruitProfileClauses struct {
	ProfileID       string
	RecruitID       string
	IncludeRecruit  bool
	OrderByPoints   bool
	RemoveFromBoard bool
}

type TeamRecruitingProfileClauses struct {
	TeamID                 string
	IncludeRecruits        bool
	IncludeRecruitProfiles bool
	OrderByPoints          bool
}

func FindRecruitPlayerProfileRecord(clauses RecruitProfileClauses) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots structs.RecruitPlayerProfile

	query := db.Model(&croots)

	if clauses.IncludeRecruit {
		query = query.Preload("Recruit")
	}

	if len(clauses.ProfileID) > 0 {
		query = query.Where("profile_id = ?", clauses.ProfileID)
	}

	if len(clauses.RecruitID) > 0 {
		query = query.Where("recruit_id = ?", clauses.RecruitID)
	}

	if clauses.RemoveFromBoard {
		query = query.Where("removed_from_board = ?", false)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&croots).Error; err != nil {
		return structs.RecruitPlayerProfile{}
	}

	return croots
}

func FindRecruitPlayerProfileRecords(clauses RecruitProfileClauses) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile

	query := db.Model(&croots)

	if clauses.IncludeRecruit {
		query = query.Preload("Recruit")
	}

	if len(clauses.ProfileID) > 0 {
		query = query.Where("profile_id = ?", clauses.ProfileID)
	}

	if len(clauses.RecruitID) > 0 {
		query = query.Where("recruit_id = ?", clauses.RecruitID)
	}

	if clauses.RemoveFromBoard {
		query = query.Where("removed_from_board = ?", false)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&croots).Error; err != nil {
		return []structs.RecruitPlayerProfile{}
	}

	return croots
}

func FindTeamRecruitingProfileRecord(clauses TeamRecruitingProfileClauses) structs.TeamRecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TeamRecruitingProfile

	query := db.Model(&profile)

	if clauses.IncludeRecruits {
		query = query.Preload("Recruits")
	}

	if clauses.IncludeRecruitProfiles {
		query = query.Preload("Recruits.Recruit.RecruitProfiles", func(db *gorm.DB) *gorm.DB {
			return db.Order("total_points DESC").Where("total_points > 0")
		})
	}

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points DESC")
	}

	if err := query.First(&profile).Error; err != nil {
		return structs.TeamRecruitingProfile{}
	}

	return profile
}

// Create
func CreateRecruitPlayerProfileRecord(profile structs.RecruitPlayerProfile, db *gorm.DB) {
	err := db.Create(&profile).Error
	if err != nil {
		log.Panicln("Could not create recruit profile record!")
	}
}

func CreateRecruitPointAllocationRecord(allocation structs.RecruitPointAllocation, db *gorm.DB) {
	err := db.Create(&allocation).Error
	if err != nil {
		log.Panicln("Could not create recruit point allocation record!")
	}
}

// Saves
func SaveRecruitRecord(recruit structs.Recruit, db *gorm.DB) {
	recruit.RecruitProfiles = nil
	err := db.Save(&recruit).Error
	if err != nil {
		log.Panicln("Could not save recruit record!")
	}
}

func SaveRecruitProfileRecord(profile structs.RecruitPlayerProfile, db *gorm.DB) {
	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save recruit profile record!")
	}
}

func SaveTeamRecruitingProfileRecord(tp structs.TeamRecruitingProfile, db *gorm.DB) {
	tp.Recruits = nil
	err := db.Save(&tp).Error
	if err != nil {
		log.Panicln("Could not save team recruiting profile record!")
	}
}
