package managers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"gorm.io/gorm"
)

// GetRecruitingProfileForDashboardByTeamID -- Dashboard
func GetRecruitingProfileForDashboardByTeamID(TeamID string) structs.TeamRecruitingProfile {
	return repository.FindTeamRecruitingProfileRecord(repository.TeamRecruitingProfileClauses{TeamID: TeamID, IncludeRecruitProfiles: true, OrderByPoints: true})
}

func GetRecruitingProfileForTeamBoardByTeamID(TeamID string) structs.SimTeamBoardResponse {
	db := dbprovider.GetInstance().GetDB()

	var profile structs.TeamRecruitingProfile

	err := db.Preload("Recruits.Recruit.RecruitProfiles", func(db *gorm.DB) *gorm.DB {
		return db.Order("total_points DESC")
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

		// croot.Map(profile.Recruits[i].Recruit)

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

func GetAllRecruitsByProfileID(profileID string) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var recruitPoints []structs.RecruitPlayerProfile

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

	return croots
}

func GetAllRecruitPlayerProfiles() []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile
	err := db.Where("removed_from_board = false").Order("total_points desc").Find(&croots).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.RecruitPlayerProfile{}
		} else {
			log.Fatal(err)
		}
	}
	return croots
}

func GetRecruitPlayerProfilesByRecruitId(recruitID string) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var croots []structs.RecruitPlayerProfile
	err := db.Where("recruit_id = ? AND removed_from_board = false", recruitID).Order("total_points desc").Find(&croots).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []structs.RecruitPlayerProfile{}
		} else {
			log.Fatal(err)
		}
	}
	return croots
}

func AddRecruitToTeamBoard(recruitProfileDto structs.CreateRecruitProfileDto) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitProfile := GetPlayerRecruitProfileByPlayerId(strconv.Itoa(recruitProfileDto.PlayerId), strconv.Itoa(recruitProfileDto.ProfileId))

	// If Recruit Already Exists
	if recruitProfile.RecruitID != 0 && recruitProfile.ProfileID != 0 {
		recruitProfile.ReplaceRecruitToBoard()
		repository.SaveRecruitProfileRecord(recruitProfile, db)
		return recruitProfile
	}

	newProfileForRecruit := structs.RecruitPlayerProfile{
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

func AddRecruitToTeamBoardV2(recruitProfileDto structs.CreateRecruitProfileDtoV2) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitProfile := GetPlayerRecruitProfileByPlayerId(strconv.Itoa(int(recruitProfileDto.PlayerID)), strconv.Itoa(int(recruitProfileDto.ProfileID)))

	teamRecruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(strconv.Itoa(int(recruitProfileDto.ProfileID)))

	// If Recruit Already Exists
	if recruitProfile.RecruitID != 0 && recruitProfile.ProfileID != 0 {
		recruitProfile.ReplaceRecruitToBoard()
		repository.SaveRecruitProfileRecord(recruitProfile, db)
		return recruitProfile
	}

	modifier := CalculateModifierTowardsRecruit(recruitProfileDto.Recruit.PlayerPreferences, teamRecruitingProfile)

	state := recruitProfileDto.HasStateBonus
	region := recruitProfileDto.HasRegionBonus

	newProfileForRecruit := structs.RecruitPlayerProfile{
		SeasonID:           recruitProfileDto.SeasonID,
		RecruitID:          recruitProfileDto.PlayerID,
		ProfileID:          recruitProfileDto.ProfileID,
		TeamAbbreviation:   recruitProfileDto.Team,
		TotalPoints:        0,
		Modifier:           modifier,
		CurrentWeeksPoints: 0,
		Scholarship:        false,
		InterestLevel:      "None",
		IsSigned:           false,
		RemovedFromBoard:   false,
		HasStateBonus:      state,
		HasRegionBonus:     region,
	}

	repository.CreateRecruitPlayerProfileRecord(newProfileForRecruit, db)

	return newProfileForRecruit
}

func GetPlayerRecruitProfileByPlayerId(playerId string, profileId string) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	var recruitingPoints structs.RecruitPlayerProfile
	err := db.Where("recruit_id = ? AND profile_id = ?", playerId, profileId).Find(&recruitingPoints).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return structs.RecruitPlayerProfile{
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

func FindRecruitPlayerProfileByProfileID(profiles []structs.RecruitPlayerProfile, profileID uint) structs.RecruitPlayerProfile {
	if len(profiles) == 0 {
		return structs.RecruitPlayerProfile{}
	}

	for _, p := range profiles {
		if p.ProfileID == profileID {
			return p
		}
	}
	return structs.RecruitPlayerProfile{}
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

func GetRecruitingProfilesByTeamId(id string) []structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()
	var recruits []structs.RecruitPlayerProfile
	db.Where("profile_id = ? AND removed_from_board = ?", id, false).Find(&recruits)

	return recruits
}

func GetRecruitFromRecruitsList(id int, recruits []structs.RecruitPlayerProfile) structs.RecruitPlayerProfile {
	var recruit structs.RecruitPlayerProfile

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
		fmt.Printf("%s", "Recruiting Profile "+strconv.Itoa(updateRecruitPointsDto.ProfileId)+" cannot spend more points than weekly amount")
		return
	}

	recruitingPointsProfile.AllocatePoints(updateRecruitPointsDto.SpentPoints)
	repository.SaveRecruitProfileRecord(recruitingPointsProfile, db)
	repository.SaveTeamRecruitingProfileRecord(recruitingProfile, db)
}

func SendScholarshipToRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.RecruitPlayerProfile, structs.TeamRecruitingProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	if recruitingProfile.ScholarshipsAvailable == 0 {
		log.Fatalf("%s", "\nTeamId: "+strconv.Itoa(updateRecruitPointsDto.ProfileId)+" does not have any availabe scholarships")
	}

	crootProfile := GetPlayerRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	crootProfile.ToggleScholarship(updateRecruitPointsDto.RewardScholarship, updateRecruitPointsDto.RevokeScholarship)
	if crootProfile.Scholarship {
		recruitingProfile.SubtractScholarshipsAvailable()
	} else {
		recruitingProfile.ReallocateScholarship()
	}

	repository.SaveRecruitProfileRecord(crootProfile, db)
	repository.SaveTeamRecruitingProfileRecord(recruitingProfile, db)
	return crootProfile, recruitingProfile
}

func SendScholarshipToRecruitV2(updateRecruitPointsDto structs.UpdateRecruitPointsDtoV2) (structs.RecruitPlayerProfile, structs.TeamRecruitingProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(strconv.Itoa(updateRecruitPointsDto.ProfileID))

	if recruitingProfile.ScholarshipsAvailable == 0 {
		log.Fatalf("%s", "\nTeamId: "+strconv.Itoa(updateRecruitPointsDto.ProfileID)+" does not have any availabe scholarships")
	}

	crootProfile := GetPlayerRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.RecruitID),
		strconv.Itoa(updateRecruitPointsDto.ProfileID),
	)

	crootProfile.ToggleScholarship(updateRecruitPointsDto.RewardScholarship, updateRecruitPointsDto.RevokeScholarship)
	if crootProfile.Scholarship {
		recruitingProfile.SubtractScholarshipsAvailable()
	} else {
		recruitingProfile.ReallocateScholarship()
	}

	repository.SaveRecruitProfileRecord(crootProfile, db)
	repository.SaveTeamRecruitingProfileRecord(recruitingProfile, db)
	return crootProfile, recruitingProfile
}

func RevokeScholarshipFromRecruit(updateRecruitPointsDto structs.UpdateRecruitPointsDto) (structs.RecruitPlayerProfile, structs.TeamRecruitingProfile) {
	db := dbprovider.GetInstance().GetDB()

	recruitingProfile := GetRecruitingProfileByTeamId(strconv.Itoa(updateRecruitPointsDto.ProfileId))

	recruitingPointsProfile := GetPlayerRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.Scholarship {
		fmt.Printf("%s", "\nCannot revoke an inexistant scholarship from Recruit "+strconv.Itoa(int(recruitingPointsProfile.RecruitID)))
		return recruitingPointsProfile, recruitingProfile
	}

	recruitingPointsProfile.RevokeScholarship()
	recruitingProfile.ReallocateScholarship()

	repository.SaveRecruitProfileRecord(recruitingPointsProfile, db)
	repository.SaveTeamRecruitingProfileRecord(recruitingProfile, db)

	return recruitingPointsProfile, recruitingProfile
}

func RemoveRecruitFromBoard(updateRecruitPointsDto structs.UpdateRecruitPointsDto) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitingPointsProfile := GetPlayerRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.PlayerId),
		strconv.Itoa(updateRecruitPointsDto.ProfileId),
	)

	if recruitingPointsProfile.RemovedFromBoard {
		panic("Recruit already removed from board")
	}

	recruitingPointsProfile.RemoveRecruitFromBoard()
	repository.SaveRecruitProfileRecord(recruitingPointsProfile, db)

	return recruitingPointsProfile
}

func RemoveRecruitFromBoardV2(updateRecruitPointsDto structs.UpdateRecruitPointsDtoV2) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitingPointsProfile := GetPlayerRecruitProfileByPlayerId(
		strconv.Itoa(updateRecruitPointsDto.RecruitID),
		strconv.Itoa(updateRecruitPointsDto.ProfileID),
	)

	if recruitingPointsProfile.RemovedFromBoard {
		panic("Recruit already removed from board")
	}

	recruitingPointsProfile.RemoveRecruitFromBoard()
	repository.SaveRecruitProfileRecord(recruitingPointsProfile, db)

	return recruitingPointsProfile
}

func UpdateRecruitingProfile(updateRecruitingBoardDto structs.UpdateRecruitingBoardDto) structs.TeamRecruitingProfile {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateRecruitingBoardDto.TeamID)

	var profile = GetOnlyTeamRecruitingProfileByTeamID(teamId)

	var recruitingPoints = GetRecruitingProfilesByTeamId(teamId)

	var updatedRecruits = updateRecruitingBoardDto.Recruits

	currentPoints := 0

	for i := 0; i < len(recruitingPoints); i++ {
		updatedRecruit := GetRecruitFromRecruitsList(int(recruitingPoints[i].RecruitID), updatedRecruits)

		if recruitingPoints[i].CurrentWeeksPoints != uint8(updatedRecruit.CurrentWeeksPoints) {

			// Allocate Points to Profile
			currentPoints += int(updatedRecruit.CurrentWeeksPoints)
			profile.AllocateSpentPoints(currentPoints)
			// If total not surpassed, allocate to the recruit and continue
			if profile.SpentPoints <= profile.WeeklyPoints {
				recruitingPoints[i].AllocatePoints(int(updatedRecruit.CurrentWeeksPoints))
				fmt.Println("Saving recruit " + strconv.Itoa(int(recruitingPoints[i].RecruitID)))
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(int(profile.TeamID)) + " than what is allowed.")
			}
			repository.SaveRecruitProfileRecord(recruitingPoints[i], db)
		} else {
			currentPoints += int(recruitingPoints[i].CurrentWeeksPoints)
			profile.AllocateSpentPoints(currentPoints)
		}
	}

	// Save profile
	repository.SaveTeamRecruitingProfileRecord(profile, db)

	return profile
}

func ScoutAttribute(dto structs.ScoutAttributeDTO) structs.RecruitPlayerProfile {
	db := dbprovider.GetInstance().GetDB()

	recruitID := strconv.Itoa(int(dto.RecruitID))
	profileID := strconv.Itoa(int(dto.ProfileID))

	teamProfile := repository.FindTeamRecruitingProfileRecord(repository.TeamRecruitingProfileClauses{TeamID: profileID})

	recruitProfile := repository.FindRecruitPlayerProfileRecord(repository.RecruitProfileClauses{RecruitID: recruitID, ProfileID: profileID})

	if teamProfile.ID == 0 || recruitProfile.ID == 0 {
		log.Panic("ERROR: IDs PROVIDED DON'T LINE UP")
	}

	if teamProfile.WeeklyScoutingPoints == 0 {
		return recruitProfile
	}

	recruitProfile.ApplyScoutingAttribute(dto.Attribute)

	teamProfile.SubtractScoutingPoints(dto.Attribute)

	repository.SaveTeamRecruitingProfileRecord(teamProfile, db)
	repository.SaveRecruitProfileRecord(recruitProfile, db)

	return recruitProfile
}

// ToggleAIBehavior -- Toggle whether a Team will use AI recruiting or not
func ToggleAIBehavior(teamID string) {
	db := dbprovider.GetInstance().GetDB()
	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(teamID)
	recruitingProfile.ToggleAIBehavior()
	repository.SaveTeamRecruitingProfileRecord(recruitingProfile, db)
}

// SaveAIBehavior -- Toggle whether a Team will use AI recruiting or not
func SaveAIBehavior(profile structs.TeamRecruitingProfile) {
	db := dbprovider.GetInstance().GetDB()
	TeamID := strconv.Itoa(int(profile.TeamID))
	recruitingProfile := GetOnlyTeamRecruitingProfileByTeamID(TeamID)
	recruitingProfile.UpdateAIBehavior(profile.IsAI, profile.AIAutoOfferscholarships, profile.AIStarMax, profile.AIStarMin, profile.AIMinThreshold, profile.AIMaxThreshold, profile.AIValue, profile.OffensiveScheme, profile.DefensiveScheme)
	repository.SaveTeamRecruitingProfileRecord(recruitingProfile, db)
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
	expectations := util.GetPlaytimeExpectations(int(dto.Stars), 1, 0)
	rankMod := 0.95 + rand.Float64()*(1.05-0.95)

	collegeRecruit := &structs.Recruit{
		PlayerID:        newID,
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

func DetermineRecruitingClassSize() {
	db := dbprovider.GetInstance().GetDB()
	recruitingProfiles := GetRecruitingProfileForRecruitSync()
	limit := 15
	for _, rp := range recruitingProfiles {
		existingRoster := GetCollegePlayersByTeamId(strconv.Itoa(int(rp.ID)))
		count := 0
		for _, p := range existingRoster {
			if !p.WillDeclare {
				count++
			}
		}
		nonLeaving := limit - count
		classSize := limit - nonLeaving
		if rp.RecruitClassSize != classSize {
			rp.SetClassSize(classSize)
			repository.SaveTeamRecruitingProfileRecord(rp, db)
		}
	}
}

func GetTeamProfileMap() map[string]*structs.TeamRecruitingProfile {
	teamRecruitingProfiles := GetRecruitingProfileForRecruitSync()

	teamMap := make(map[string]*structs.TeamRecruitingProfile)
	for i := 0; i < len(teamRecruitingProfiles); i++ {
		teamMap[strconv.Itoa(int(teamRecruitingProfiles[i].ID))] = &teamRecruitingProfiles[i]
	}

	return teamMap
}

func CalculateModifierTowardsRecruit(prefs structs.PlayerPreferences, team structs.TeamRecruitingProfile) float32 {
	programMod := calculateMultiplier(uint(team.ProgramPrestige), uint(prefs.ProgramPref))
	professionalDevMod := calculateMultiplier(uint(team.ProfessionalPrestige), uint(prefs.ProfDevPref))
	traditionsMod := calculateMultiplier(uint(team.Traditions), uint(prefs.TraditionsPref))
	facilitiesMod := calculateMultiplier(uint(team.Facilities), uint(prefs.FacilitiesPref))
	atmosphereMod := calculateMultiplier(uint(team.Atmosphere), uint(prefs.AtmospherePref))
	academicsMod := calculateMultiplier(uint(team.Academics), uint(prefs.AcademicsPref))
	conferenceMod := calculateMultiplier(uint(team.ConferencePrestige), uint(prefs.ConferencePref))
	coachMod := calculateMultiplier(uint(team.CoachRating), uint(prefs.CoachPref))
	seasonMod := calculateMultiplier(uint(team.SeasonMomentum), uint(prefs.SeasonMomentumPref))
	collegeLifeMod := calculateMultiplier(uint(team.CampusLife), uint(prefs.CampusLifePref))

	return (programMod + professionalDevMod + traditionsMod + facilitiesMod + atmosphereMod + academicsMod + conferenceMod + coachMod + seasonMod + collegeLifeMod) / 10
}

func calculateBaseModifier(attr int) float32 {
	return 1 + float32(attr-5)/5
}

func calculateAdjustmentFactor(teamAttr, playerPref int) float32 {
	return 1 + float32((teamAttr-playerPref)/10)
}

func calculateMultiplier(teamAttr uint, playerPref uint) float32 {
	baseMod := calculateBaseModifier(int(teamAttr))
	adjFactor := calculateAdjustmentFactor(int(teamAttr), int(playerPref))
	return baseMod * adjFactor
}
