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

func GetRecruitingProfileByTeamId(teamId string) structs.TeamRecruitingProfile {
	var profile structs.TeamRecruitingProfile
	db := dbprovider.GetInstance().GetDB()
	err := db.Preload("Recruits", "removed_from_board = ?", false).Preload("Recruits.Recruit.RecruitingPoints", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points_spent DESC")
	}).Where("id = ?", teamId).Find(&profile).Error
	if err != nil {
		log.Fatal(err)
	}
	return profile
}

func GetOnlyRecruitingProfileByTeamId(teamId string) structs.TeamRecruitingProfile {
	var profile structs.TeamRecruitingProfile
	db := dbprovider.GetInstance().GetDB()
	err := db.Where("id = ?", teamId).Find(&profile).Error
	if err != nil {
		log.Fatal(err)
	}
	return profile
}

func GetAllRecruitsByProfileID(profileID string) []structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()

	var recruitPoints []structs.PlayerRecruitProfile

	db.Preload("Recruit").Where("profile_id = ?", profileID).Find(&recruitPoints)

	return recruitPoints
}

func CreateRecruitingPointsProfileForRecruit(recruitPointsDto structs.CreateRecruitPointsDto) structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitProfile := GetRecruitingPointsProfileByPlayerId(strconv.Itoa(recruitPointsDto.PlayerId), strconv.Itoa(recruitPointsDto.ProfileId))

	// If Recruit Already Exists
	if recruitProfile.RecruitID != 0 && recruitProfile.ProfileID != 0 {
		recruitProfile.ReplaceRecruitToBoard()
		db.Save(&recruitProfile)
		return recruitProfile
	}

	recruitingPointProfile := structs.PlayerRecruitProfile{
		SeasonID:               recruitPointsDto.SeasonId,
		RecruitID:              recruitPointsDto.PlayerId,
		ProfileID:              recruitPointsDto.ProfileId,
		TeamAbbreviation:       recruitPointsDto.Team,
		TotalPoints:            0,
		CurrentWeeksPoints:     0,
		Scholarship:            false,
		InterestLevel:          "None",
		InterestLevelThreshold: 0,
		IsSigned:               false,
		RemovedFromBoard:       false,
	}

	db.Create(&recruitingPointProfile)

	return recruitingPointProfile
}

func GetRecruitingPointsProfileByPlayerId(playerId string, profileId string) structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()

	var recruitingPoints structs.PlayerRecruitProfile
	err := db.Where("player_id = ? AND profile_id = ?", playerId, profileId).Find(&recruitingPoints).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.PlayerRecruitProfile{
				SeasonID:               0,
				RecruitID:              0,
				ProfileID:              0,
				TotalPoints:            0,
				CurrentWeeksPoints:     0,
				Scholarship:            false,
				InterestLevel:          "None",
				InterestLevelThreshold: 0,
				IsSigned:               false,
				RemovedFromBoard:       false,
			}
		} else {
			log.Fatal(err)
		}

	}
	return recruitingPoints
}

func GetRecruitingPointsByTeamId(id string) []structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()
	var recruits []structs.PlayerRecruitProfile
	db.Where("profile_id = ? AND removed_from_board = ?", id, false).Find(&recruits)

	return recruits
}

func GetRecruitFromRecruitsList(id int, recruits []structs.PlayerRecruitProfile) structs.PlayerRecruitProfile {
	var recruit structs.PlayerRecruitProfile

	for i := 0; i < len(recruits); i++ {
		if recruits[i].RecruitID == id {
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

func SendScholarshipToRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.PlayerRecruitProfile, structs.TeamRecruitingProfile) {
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
		log.Fatalf("\nRecruit " + strconv.Itoa(recruitingPointsProfile.RecruitID) + "already has a scholarship")
	}

	recruitingPointsProfile.AllocateScholarship()
	recruitingProfile.SubtractScholarshipsAvailable()

	db.Save(&recruitingPointsProfile)
	db.Save(&recruitingProfile)

	return recruitingPointsProfile, recruitingProfile
}

func RevokeScholarshipFromRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.PlayerRecruitProfile, structs.TeamRecruitingProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	recruitingPointsProfile := GetRecruitingPointsProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.Scholarship {
		fmt.Printf("\nCannot revoke an inexistant scholarship from Recruit " + strconv.Itoa(recruitingPointsProfile.RecruitID))
		return recruitingPointsProfile, recruitingProfile
	}

	recruitingPointsProfile.RevokeScholarship()
	recruitingProfile.ReallocateScholarship()

	db.Save(&recruitingPointsProfile)
	db.Save(&recruitingProfile)

	return recruitingPointsProfile, recruitingProfile
}

func RemoveRecruitFromBoard(updateRecruitPointsDto structs.UpdateRecruitPointsDto) structs.PlayerRecruitProfile {
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

func UpdateRecruitingProfile(updateRecruitingBoardDto structs.UpdateRecruitingBoardDto) structs.TeamRecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateRecruitingBoardDto.TeamID)

	var profile = GetOnlyRecruitingProfileByTeamId(teamId)

	var recruitingPoints = GetRecruitingPointsByTeamId(teamId)

	var updatedRecruits = updateRecruitingBoardDto.Recruits

	currentPoints := 0

	for i := 0; i < len(recruitingPoints); i++ {
		updatedRecruit := GetRecruitFromRecruitsList(recruitingPoints[i].RecruitID, updatedRecruits)

		if updatedRecruit.CurrentWeeksPoints > 0 &&
			recruitingPoints[i].CurrentWeeksPoints != updatedRecruit.CurrentWeeksPoints {

			// Allocate Points to Profile
			currentPoints += updatedRecruit.CurrentWeeksPoints
			profile.AllocateSpentPoints(currentPoints)
			// If total not surpassed, allocate to the recruit and continue
			if profile.SpentPoints <= profile.WeeklyPoints {
				recruitingPoints[i].AllocatePoints(updatedRecruit.CurrentWeeksPoints)
				fmt.Println("Saving recruit " + strconv.Itoa(recruitingPoints[i].RecruitID))
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
