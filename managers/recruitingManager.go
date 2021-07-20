package managers

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/jinzhu/gorm"
)

func GetRecruitingProfileByTeamId(teamId string) structs.RecruitingProfile {
	var profile structs.RecruitingProfile
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("Recruits", "removed_from_board = ?", false).Preload("Recruits.Recruit.RecruitingPoints", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points_spent DESC")
	}).Where("id = ?", teamId).Find(&profile).Error
	if err != nil {
		log.Fatal(err)
	}
	return profile
}

func GetOnlyRecruitingProfileByTeamId(teamId string) structs.RecruitingProfile {
	var profile structs.RecruitingProfile
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("id = ?", teamId).Find(&profile).Error
	if err != nil {
		log.Fatal(err)
	}
	return profile
}

func GetAllRecruitsByProfileID(profileID string) []structs.RecruitingPoints {
	db := dbprovider.GetInstance().GetDB()

	var recruitPoints []structs.RecruitingPoints

	db.Preload("Recruit").Where("profile_id = ?", profileID).Find(&recruitPoints)

	return recruitPoints
}

func CreateRecruitingPointsProfileForRecruit(recruitPointsDto structs.CreateRecruitPointsDto) structs.RecruitingPoints {
	db := dbprovider.GetInstance().GetDB()

	recruitProfile := GetRecruitingPointsProfileByPlayerId(strconv.Itoa(recruitPointsDto.PlayerId), strconv.Itoa(recruitPointsDto.ProfileId))

	// If Recruit Already Exists
	if recruitProfile.PlayerID != 0 && recruitProfile.ProfileID != 0 {
		recruitProfile.ReplaceRecruitToBoard()
		db.Save(&recruitProfile)
		return recruitProfile
	}

	recruitingPointProfile := structs.RecruitingPoints{
		SeasonID:               recruitPointsDto.SeasonId,
		PlayerID:               recruitPointsDto.PlayerId,
		ProfileID:              recruitPointsDto.ProfileId,
		Team:                   recruitPointsDto.Team,
		TotalPointsSpent:       0,
		CurrentPointsSpent:     0,
		Scholarship:            false,
		InterestLevel:          "None",
		InterestLevelThreshold: 0,
		Signed:                 false,
		RemovedFromBoard:       false,
	}

	db.Create(&recruitingPointProfile)

	return recruitingPointProfile
}

func GetRecruitingPointsProfileByPlayerId(playerId string, profileId string) structs.RecruitingPoints {
	db := dbprovider.GetInstance().GetDB()

	var recruitingPoints structs.RecruitingPoints
	err := db.Where("player_id = ? AND profile_id = ?", playerId, profileId).Find(&recruitingPoints).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.RecruitingPoints{
				SeasonID:               0,
				PlayerID:               0,
				ProfileID:              0,
				TotalPointsSpent:       0,
				CurrentPointsSpent:     0,
				Scholarship:            false,
				InterestLevel:          "None",
				InterestLevelThreshold: 0,
				Signed:                 false,
				RemovedFromBoard:       false,
			}
		} else {
			log.Fatal(err)
		}

	}
	return recruitingPoints
}

func GetRecruitingPointsByTeamId(id string) []structs.RecruitingPoints {
	db := dbprovider.GetInstance().GetDB()
	var recruits []structs.RecruitingPoints
	db.Where("profile_id = ? AND removed_from_board = ?", id, false).Find(&recruits)

	return recruits
}

func GetRecruitFromRecruitsList(id int, recruits []structs.RecruitingPoints) structs.RecruitingPoints {
	var recruit structs.RecruitingPoints

	for i := 0; i < len(recruits); i++ {
		if recruits[i].PlayerID == id {
			recruit = recruits[i]
			break
		}
	}

	return recruit
}

func AllocateRecruitingPointsForRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	recruitingPointsProfile := GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	recruitingProfile.AllocateSpentPoints(updateRecruitPointsDto.SpentPoints)
	if recruitingProfile.SpentPoints > recruitingProfile.WeeklyPoints {
		fmt.Printf("Recruiting Profile " + strconv.Itoa(updateRecruitPointsDto.ProfileId) + " cannot spend more points than weekly amount")
		return
	}

	recruitingPointsProfile.AllocatePoints(updateRecruitPointsDto.SpentPoints)

	db.Save(&recruitingPointsProfile)

	db.Save(&recruitingProfile)
}

func SendScholarshipToRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.RecruitingPoints, structs.RecruitingProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	if recruitingProfile.ScholarshipsAvailable == 0 {
		log.Fatalf("\nTeamId: " + strconv.Itoa(updateRecruitPointsDto.ProfileId) + " does not have any availabe scholarships")
	}

	recruitingPointsProfile := GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.Scholarship {
		log.Fatalf("\nRecruit " + strconv.Itoa(recruitingPointsProfile.PlayerID) + "already has a scholarship")
	}

	recruitingPointsProfile.AllocateScholarship()
	recruitingProfile.SubtractScholarshipsAvailable()

	db.Save(&recruitingPointsProfile)
	db.Save(&recruitingProfile)

	return recruitingPointsProfile, recruitingProfile
}

func RevokeScholarshipFromRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.RecruitingPoints, structs.RecruitingProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	recruitingPointsProfile := GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.Scholarship {
		fmt.Printf("\nCannot revoke an inexistant scholarship from Recruit " + strconv.Itoa(recruitingPointsProfile.PlayerID))
		return recruitingPointsProfile, recruitingProfile
	}

	recruitingPointsProfile.RevokeScholarship()
	recruitingProfile.ReallocateScholarship()

	db.Save(&recruitingPointsProfile)
	db.Save(&recruitingProfile)

	return recruitingPointsProfile, recruitingProfile
}

func RemoveRecruitFromBoard(updateRecruitPointsDto structs.UpdateRecruitPointsDto) structs.RecruitingPoints {
	db := dbprovider.GetInstance().GetDB()

	recruitingPointsProfile := GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.RemovedFromBoard {
		panic("Recruit already removed from board")
	}

	recruitingPointsProfile.RemoveRecruitFromBoard()
	db.Save(&recruitingPointsProfile)

	return recruitingPointsProfile
}

func UpdateRecruitingProfile(updateRecruitingBoardDto structs.UpdateRecruitingBoardDto) structs.RecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateRecruitingBoardDto.TeamID)

	var profile = GetOnlyRecruitingProfileByTeamId(teamId)

	var recruitingPoints = GetRecruitingPointsByTeamId(teamId)

	var updatedRecruits = updateRecruitingBoardDto.Recruits

	currentPoints := 0

	for i := 0; i < len(recruitingPoints); i++ {
		updatedRecruit := GetRecruitFromRecruitsList(recruitingPoints[i].PlayerID, updatedRecruits)

		if updatedRecruit.CurrentPointsSpent > 0 &&
			recruitingPoints[i].CurrentPointsSpent != updatedRecruit.CurrentPointsSpent {

			// Allocate Points to Profile
			currentPoints += updatedRecruit.CurrentPointsSpent
			profile.AllocateSpentPoints(currentPoints)
			// If total not surpassed, allocate to the recruit and continue
			if profile.SpentPoints <= profile.WeeklyPoints {
				recruitingPoints[i].AllocatePoints(updatedRecruit.CurrentPointsSpent)
				fmt.Println("Saving recruit " + strconv.Itoa(recruitingPoints[i].PlayerID))
				db.Save(&recruitingPoints[i])
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(profile.TeamID) + " than what is allowed.")
			}
		}
	}

	// Save profile
	db.Save(&profile)

	return profile
}
