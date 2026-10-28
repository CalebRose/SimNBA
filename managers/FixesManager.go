package managers

import (
	"strconv"

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

func FixTeamSeasonStats() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	collegeTeams := GetAllActiveCollegeTeams()
	nbaTeams := GetAllActiveNBATeams()
	teamStats := repository.FindCollegeTeamGameStatsRecords(seasonID, "", "", "")
	teamStatsMap := MakeCBBTeamStatsMap(teamStats)
	teamSeasonStats := repository.FindCollegeTeamSeasonStatsRecords(seasonID, "")
	preseasonNBAStats := repository.FindProTeamGameStatsRecords(seasonID, "", "1", "")
	regularSeasonNBAStats := repository.FindProTeamGameStatsRecords(seasonID, "", "2", "")
	preseasonNBAStatsMap := MakeNBATeamStatsMap(preseasonNBAStats)
	regularSeasonNBAStatsMap := MakeNBATeamStatsMap(regularSeasonNBAStats)
	preseasonNBASeasonStats := repository.FindProTeamSeasonStatsRecords(seasonID, "1")
	regularSeasonNBASeasonStats := repository.FindProTeamSeasonStatsRecords(seasonID, "2")
	preseasonNBASeasonStatMap := MakeNBATeamSeasonStatsMap(preseasonNBASeasonStats)
	regularSeasonNBASeasonStatMap := MakeNBATeamSeasonStatsMap(regularSeasonNBASeasonStats)
	teamSeasonStatMap := MakeCBBTeamSeasonStatsMap(teamSeasonStats)

	for _, team := range collegeTeams {
		stats := teamStatsMap[team.ID]
		seasonStat := teamSeasonStatMap[team.ID]
		if seasonStat.ID == 0 {
			seasonStat.TeamID = team.ID
			seasonStat.SeasonID = ts.SeasonID
			for _, stat := range stats {
				seasonStat.AddStatsToSeasonRecord(stat)
			}
			seasonStat.GameType = 2
			repository.SaveCollegeTeamSeasonStatRecord(seasonStat, db)
		} else {
			seasonStat.ResetSeasonsRecord()
			for _, stat := range stats {
				seasonStat.AddStatsToSeasonRecord(stat)
			}
			seasonStat.GameType = 2
			repository.SaveCollegeTeamSeasonStatRecord(seasonStat, db)
		}
	}

	for _, team := range nbaTeams {
		preseasonStats := preseasonNBAStatsMap[team.ID]
		regularSeasonStasts := regularSeasonNBAStatsMap[team.ID]
		preseasonSeasonStat := preseasonNBASeasonStatMap[team.ID]
		regularSeasonSeasonStat := regularSeasonNBASeasonStatMap[team.ID]
		if preseasonSeasonStat.ID == 0 {
			preseasonSeasonStat.TeamID = team.ID
			preseasonSeasonStat.SeasonID = ts.SeasonID
			for _, stat := range preseasonStats {
				preseasonSeasonStat.AddStatsToSeasonRecord(stat)
			}
			preseasonSeasonStat.GameType = 1
			repository.SaveNBATeamSeasonStatRecord(preseasonSeasonStat, db)
		} else {
			preseasonSeasonStat.ResetSeasonsRecord()
			for _, stat := range preseasonStats {
				preseasonSeasonStat.AddStatsToSeasonRecord(stat)
			}
			preseasonSeasonStat.GameType = 1
			repository.SaveNBATeamSeasonStatRecord(preseasonSeasonStat, db)
		}

		if regularSeasonSeasonStat.ID == 0 {
			regularSeasonSeasonStat.TeamID = team.ID
			regularSeasonSeasonStat.SeasonID = ts.SeasonID
			for _, stat := range regularSeasonStasts {
				regularSeasonSeasonStat.AddStatsToSeasonRecord(stat)
			}
			regularSeasonSeasonStat.GameType = 2
			repository.SaveNBATeamSeasonStatRecord(regularSeasonSeasonStat, db)
		} else {
			regularSeasonSeasonStat.ResetSeasonsRecord()
			for _, stat := range regularSeasonStasts {
				regularSeasonSeasonStat.AddStatsToSeasonRecord(stat)
			}
			regularSeasonSeasonStat.GameType = 2
			repository.SaveNBATeamSeasonStatRecord(regularSeasonSeasonStat, db)
		}
	}
}
