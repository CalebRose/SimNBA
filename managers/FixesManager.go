package managers

import (
	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
)

func FixExistingModifiersForRecruits() {
	db := dbprovider.GetInstance().GetDB()
	recruits := GetAllRecruitRecords()
	recruitMap := MakeRecruitMap(recruits)
	recruitProfiles := repository.FindRecruitPlayerProfileRecords(repository.RecruitProfileClauses{})
	teamProfiles := GetRecruitingProfileForRecruitSync()
	teamProfileMap := MakeTeamRecruitingProfileMapByTeamID(teamProfiles)

	for _, profile := range recruitProfiles {
		recruit := recruitMap[uint(profile.RecruitID)]
		teamProfile := teamProfileMap[uint(profile.ProfileID)]

		stateBonus := CheckStateBonus(recruit.State, teamProfile.State)
		regionBonus := CheckRegionBonus(recruit.Country, recruit.State, teamProfile.Region)

		modifier := CalculateModifierTowardsRecruit(recruit.PlayerPreferences, teamProfile, stateBonus, regionBonus)

		profile.Modifier = modifier
		totalPoints := profile.Modifier * float32(profile.PreviouslySpentPoints)
		profile.TotalPoints = float64(totalPoints)
		repository.SaveRecruitProfileRecord(profile, db)
	}
}
