package managers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

// GetRecruitingProfileForDashboardByTeamID -- Dashboard
func GetRecruitingProfileForDashboardByTeamID(TeamID string) structs.TeamRecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TeamRecruitingProfile

	err := db.Preload("Recruits.Recruit.RecruitProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC").Where("total_points > 0")
	}).Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Panicln(err)
	}

	return profile
}

func GetRecruitingProfileForTeamBoardByTeamID(TeamID string) structs.SimTeamBoardResponse {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TeamRecruitingProfile

	err := db.Preload("Recruits.Recruit.RecruitProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC").Where("total_points > 0")
	}).Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Panicln(err)
	}

	var teamProfileResponse structs.SimTeamBoardResponse
	var crootProfiles []structs.CrootProfile

	// iterate through player recruit profiles --> get recruit with preload to player profiles
	for i := 0; i < len(profile.Recruits); i++ {
		var crootProfile structs.CrootProfile
		var croot structs.Croot

		croot.Map(profile.Recruits[i].Recruit)

		crootProfile.Map(profile.Recruits[i], croot)

		crootProfiles = append(crootProfiles, crootProfile)
	}

	sort.Sort(structs.ByCrootProfileTotal(crootProfiles))

	teamProfileResponse.Map(profile, crootProfiles)

	return teamProfileResponse
}

func GetTeamRecruitingProfilesForRecruitSync() []structs.TeamRecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TeamRecruitingProfile

	err := db.Find(&profiles).Error
	if err != nil {
		log.Panicln(err)
	}

	return profiles
}

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

// GetRecruitingProfileByTeamID
func GetOnlyTeamRecruitingProfileByTeamID(TeamID string) structs.TeamRecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TeamRecruitingProfile

	err := db.Where("id = ?", TeamID).Find(&profile).Error
	if err != nil {
		log.Fatal(err)
	}

	return profile
}

func GetOnlyAITeamRecruitingProfiles() []structs.TeamRecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var AIProfiles []structs.TeamRecruitingProfile

	err := db.Where("is_ai = ?", true).Find(&AIProfiles).Error
	if err != nil {
		log.Fatal(err)
	}

	return AIProfiles
}

func GetRecruitByRecruitID(RecruitID string) structs.Recruit {
	db := dbprovider.GetInstance().GetDB()
	var recruit structs.Recruit

	db.Where("id = ?", RecruitID).Find(&recruit)

	return recruit
}

func GetAllRecruitsByProfileID(profileID string) []structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()

	var recruitPoints []structs.PlayerRecruitProfile

	db.Preload("Recruit", func(db *gorm.DB) *gorm.DB {
		return db.Order("stars DESC")
	}).Where("profile_id = ? AND removed_from_board = ?", profileID, false).Order("total_points DESC").Order("").Find(&recruitPoints)

	return recruitPoints
}

func GetAllUnsignedRecruits() []structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.Recruit

	db.Order("stars DESC").Order("overall DESC").Where("is_signed = ?", false).Find(&croots)

	return croots
}

func GetSignedRecruitsByTeamProfileID(ProfileID string) []structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.Recruit

	err := db.Order("overall DESC").Where("team_id = ? AND is_signed = ?", ProfileID, true).Find(&croots).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.Recruit{}
		} else {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	return croots
}

func GetRecruitPlayerProfilesByRecruitId(recruitID string) []structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.PlayerRecruitProfile
	err := db.Where("recruit_id = ? AND removed_from_board = false", recruitID).Order("total_points desc").Find(&croots).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.PlayerRecruitProfile{}
		} else {
			log.Fatal(err)
		}
	}
	return croots
}

func AddRecruitToTeamBoard(recruitProfileDto structs.CreateRecruitProfileDto) structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitProfile := GetPlayerRecruitProfileByPlayerId(strconv.Itoa(recruitProfileDto.PlayerId), strconv.Itoa(recruitProfileDto.ProfileId))

	// If Recruit Already Exists
	if recruitProfile.RecruitID != 0 && recruitProfile.ProfileID != 0 {
		recruitProfile.ReplaceRecruitToBoard()
		db.Save(&recruitProfile)
		return recruitProfile
	}

	newProfileForRecruit := structs.PlayerRecruitProfile{
		SeasonID:           uint(recruitProfileDto.SeasonId),
		RecruitID:          uint(recruitProfileDto.PlayerId),
		ProfileID:          uint(recruitProfileDto.ProfileId),
		TeamAbbreviation:   recruitProfileDto.Team,
		TotalPoints:        0,
		CurrentWeeksPoints: 0,
		Scholarship:        false,
		InterestLevel:      "None",
		IsSigned:           false,
		RemovedFromBoard:   false,
		HasStateBonus:      recruitProfileDto.HasStateBonus,
		HasRegionBonus:     recruitProfileDto.HasRegionBonus,
	}

	db.Create(&newProfileForRecruit)

	return newProfileForRecruit
}

func GetPlayerRecruitProfileByPlayerId(playerId string, profileId string) structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()

	var recruitingPoints structs.PlayerRecruitProfile
	err := db.Where("recruit_id = ? AND profile_id = ?", playerId, profileId).Find(&recruitingPoints).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.PlayerRecruitProfile{
				SeasonID:           0,
				RecruitID:          0,
				ProfileID:          0,
				TotalPoints:        0,
				CurrentWeeksPoints: 0,
				Scholarship:        false,
				InterestLevel:      "None",
				IsSigned:           false,
				RemovedFromBoard:   false,
			}
		} else {
			log.Fatal(err)
		}

	}
	return recruitingPoints
}

func GetRecruitingProfileForRecruitSync() []structs.TeamRecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TeamRecruitingProfile

	err := db.Find(&profiles).Error
	if err != nil {
		log.Panicln(err)
	}

	return profiles
}

func GetRecruitWithPlayerProfilesByRecruitID(recruitID string) structs.Recruit {
	db := dbprovider.GetInstance().GetDB()

	var croot structs.Recruit

	db.Preload("PlayerRecruitProfiles").Where("id = ?", recruitID).Find(&croot)

	return croot
}

func GetRecruitingPointsByTeamId(id string) []structs.PlayerRecruitProfile {
	db := dbprovider.GetInstance().GetDB()
	var recruits []structs.PlayerRecruitProfile
	db.Where("profile_id = ? AND removed_from_board = ?", id, false).Find(&recruits)

	return recruits
}

func GetRecruitFromRecruitsList(id int, recruits []structs.CrootProfile) structs.CrootProfile {
	var recruit structs.CrootProfile

	for i := 0; i < len(recruits); i++ {
		if recruits[i].RecruitID == uint(id) {
			recruit = recruits[i]
			break
		}
	}

	return recruit
}

func AllocateRecruitingPointsForRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	recruitingPointsProfile := GetPlayerRecruitProfileByPlayerId(
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

	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	if recruitingProfile.ScholarshipsAvailable == 0 {
		log.Fatalf("\nTeamId: " + strconv.Itoa(updateRecruitPointsDto.ProfileId) + " does not have any availabe scholarships")
	}

	crootProfile := GetPlayerRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	crootProfile.ToggleScholarship(updateRecruitPointsDto.RewardScholarship, updateRecruitPointsDto.RevokeScholarship)
	if crootProfile.Scholarship {
		recruitingProfile.SubtractScholarshipsAvailable()

		ts := GetTimestamp()
		recruit := GetRecruitByRecruitID(strconv.Itoa(int(updateRecruitPointsDto.PlayerId)))

		stars := recruit.Stars

		location := ""

		if len(recruit.State) == 0 {
			location = recruit.Country
		} else {
			location = recruit.State + ", " + recruit.Country
		}

		if stars >= 4 {
			message := recruit.FirstName + " " + recruit.LastName + ", " + strconv.Itoa(stars) + " star " + recruit.Position + " from " + location + " has received an offer from " + updateRecruitPointsDto.Team
			newLog := structs.NewsLog{
				WeekID:      ts.CollegeWeekID,
				Week:        uint(ts.CollegeWeek),
				SeasonID:    ts.SeasonID,
				MessageType: "Recruiting",
				Message:     message,
			}

			err := db.Create(&newLog).Error
			if err != nil {
				fmt.Println(err.Error())
				log.Fatalln("ERROR! Could not save news log for scholarship on " + recruit.FirstName + " " + recruit.LastName + " " + strconv.Itoa(int(recruit.ID)))
			}
		}
	} else {
		recruitingProfile.ReallocateScholarship()
	}

	db.Save(&crootProfile)
	db.Save(&recruitingProfile)

	return crootProfile, recruitingProfile
}

func RevokeScholarshipFromRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.PlayerRecruitProfile, structs.TeamRecruitingProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	recruitingPointsProfile := GetPlayerRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.Scholarship {
		fmt.Printf("\nCannot revoke an inexistant scholarship from Recruit " + strconv.Itoa(int(recruitingPointsProfile.RecruitID)))
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

	recruitingPointsProfile := GetPlayerRecruitProfileByPlayerId(
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

	var profile = GetOnlyTeamRecruitingProfileByTeamID(teamId)

	var recruitingPoints = GetRecruitingPointsByTeamId(teamId)

	var updatedRecruits = updateRecruitingBoardDto.Recruits

	currentPoints := 0

	for i := 0; i < len(recruitingPoints); i++ {
		updatedRecruit := GetRecruitFromRecruitsList(int(recruitingPoints[i].RecruitID), updatedRecruits)

		if recruitingPoints[i].CurrentWeeksPoints != updatedRecruit.CurrentWeeksPoints {

			// Allocate Points to Profile
			currentPoints += updatedRecruit.CurrentWeeksPoints
			profile.AllocateSpentPoints(currentPoints)
			// If total not surpassed, allocate to the recruit and continue
			if profile.SpentPoints <= profile.WeeklyPoints {
				recruitingPoints[i].AllocatePoints(updatedRecruit.CurrentWeeksPoints)
				fmt.Println("Saving recruit " + strconv.Itoa(int(recruitingPoints[i].RecruitID)))
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(int(profile.TeamID)) + " than what is allowed.")
			}
			db.Save(&recruitingPoints[i])
		} else {
			currentPoints += recruitingPoints[i].CurrentWeeksPoints
			profile.AllocateSpentPoints(currentPoints)
		}
	}

	// Save profile
	db.Save(&profile)

	return profile
}

func CreateRecruit(dto structs.CreateRecruitDTO) {
	db := dbprovider.GetInstance().GetDB()

	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	newID := lastPlayerRecord.ID + 1
	threshold := GetRecruitModifier(dto.Stars)
	expectations := util.GetPlaytimeExpectations(dto.Stars, 1, 0)
	rankMod := 0.95 + rand.Float64()*(1.05-0.95)

	collegeRecruit := &structs.Recruit{
		RecruitModifier: threshold,
		TopRankModifier: rankMod,
	}
	collegeRecruit.Map(dto, newID, expectations)

	playerRecord := structs.GlobalPlayer{
		RecruitID:       newID,
		CollegePlayerID: newID,
		NBAPlayerID:     newID,
	}
	playerRecord.SetID(newID)
	// Create Player Record
	db.Create(&playerRecord)
	// Save Recruit
	db.Create(&collegeRecruit)
}

func GetRecruitingClassByTeamID(id string) []structs.Croot {
	db := dbprovider.GetInstance().GetDB()

	var class []structs.Recruit
	var recruitingClass []structs.Croot

	db.Where("is_signed = true AND team_id = ?", id).Find(&class)

	for i := 0; i < len(class); i++ {
		var croot structs.Croot

		croot.Map(class[i])

		recruitingClass = append(recruitingClass, croot)
	}

	return recruitingClass
}
