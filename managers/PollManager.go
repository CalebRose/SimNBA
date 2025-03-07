package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
)

func GetAllCollegePollsByWeekIDAndSeasonID(weekID, seasonID string) []structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()

	submissions := []structs.CollegePollSubmission{}

	err := db.Where("week_id = ? AND season_id = ?", weekID, seasonID).Find(&submissions).Error
	if err != nil {
		return []structs.CollegePollSubmission{}
	}

	return submissions
}

func GetPollSubmissionBySubmissionID(id string) structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()

	submission := structs.CollegePollSubmission{}

	err := db.Where("id = ?", id).Find(&submission).Error
	if err != nil {
		return structs.CollegePollSubmission{}
	}

	return submission
}

func GetPollSubmissionByUsernameWeekAndSeason(username string) structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	weekID := strconv.Itoa(int(ts.CollegeWeekID + 1))
	seasonID := strconv.Itoa(int(ts.SeasonID))

	submission := structs.CollegePollSubmission{}

	err := db.Where("username = ? AND week_id = ? AND season_id = ?", username, weekID, seasonID).Find(&submission).Error
	if err != nil {
		return structs.CollegePollSubmission{}
	}

	return submission
}

func SyncCollegePollSubmissionForCurrentWeek(week, weekID, seasonID uint) {
	db := dbprovider.GetInstance().GetDB()

	if week > 15 && week < 20 {
		return
	}

	weekIDStr := strconv.Itoa(int(weekID))
	seasonIDStr := strconv.Itoa(int(seasonID))

	submissions := GetAllCollegePollsByWeekIDAndSeasonID(weekIDStr, seasonIDStr)

	allCollegeTeams := GetAllActiveCollegeTeams()

	voteMap := make(map[uint]*structs.TeamVote)

	for _, t := range allCollegeTeams {
		voteMap[t.ID] = &structs.TeamVote{TeamID: t.ID, Team: t.Abbr}
	}

	for _, s := range submissions {
		voteMap[s.Rank1ID].AddVotes(1)
		voteMap[s.Rank2ID].AddVotes(2)
		voteMap[s.Rank3ID].AddVotes(3)
		voteMap[s.Rank4ID].AddVotes(4)
		voteMap[s.Rank5ID].AddVotes(5)
		voteMap[s.Rank6ID].AddVotes(6)
		voteMap[s.Rank7ID].AddVotes(7)
		voteMap[s.Rank8ID].AddVotes(8)
		voteMap[s.Rank9ID].AddVotes(9)
		voteMap[s.Rank10ID].AddVotes(10)
		voteMap[s.Rank11ID].AddVotes(11)
		voteMap[s.Rank12ID].AddVotes(12)
		voteMap[s.Rank13ID].AddVotes(13)
		voteMap[s.Rank14ID].AddVotes(14)
		voteMap[s.Rank15ID].AddVotes(15)
		voteMap[s.Rank16ID].AddVotes(16)
		voteMap[s.Rank17ID].AddVotes(17)
		voteMap[s.Rank18ID].AddVotes(18)
		voteMap[s.Rank19ID].AddVotes(19)
		voteMap[s.Rank20ID].AddVotes(20)
		voteMap[s.Rank21ID].AddVotes(21)
		voteMap[s.Rank22ID].AddVotes(22)
		voteMap[s.Rank23ID].AddVotes(23)
		voteMap[s.Rank24ID].AddVotes(24)
		voteMap[s.Rank25ID].AddVotes(25)
	}

	allVotes := []structs.TeamVote{}

	for _, t := range allCollegeTeams {
		v := voteMap[t.ID]
		if v.TotalVotes == 0 {
			continue
		}
		newVoteObj := structs.TeamVote{TeamID: v.TeamID, Team: v.Team, TotalVotes: v.TotalVotes, Number1Votes: v.Number1Votes}

		allVotes = append(allVotes, newVoteObj)
	}

	sort.Slice(allVotes, func(i, j int) bool {
		return allVotes[i].TotalVotes > allVotes[j].TotalVotes
	})

	officialPoll := structs.CollegePollOfficial{
		WeekID:   weekID,
		Week:     week,
		SeasonID: seasonID,
	}
	for idx, v := range allVotes {
		if idx > 24 {
			break
		}
		officialPoll.AssignRank(idx, v)
		// Get Standings
		teamID := strconv.Itoa(int(v.TeamID))
		teamStandings := GetStandingsRecordByTeamID(teamID, seasonIDStr)
		rank := idx + 1
		teamStandings.AssignRank(rank)
		db.Save(&teamStandings)

		matches := GetMatchesByTeamIdAndSeasonId(teamID, seasonIDStr)

		for _, m := range matches {
			if m.Week < week {
				continue
			}
			if m.Week > week {
				break
			}
			m.AssignRank(v.TeamID, uint(rank))
			db.Save(&m)
		}
	}

	db.Create(&officialPoll)
}

func CreatePoll(dto structs.CollegePollSubmission) structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()
	existingPoll := GetPollSubmissionBySubmissionID(strconv.Itoa(int(dto.ID)))

	if existingPoll.ID > 0 {
		dto.AssignID(existingPoll.ID)
		db.Save(&dto)
	} else {
		db.Create(&dto)
	}

	return dto
}

func GetOfficialPollBySeasonID(seasonID string) []structs.CollegePollOfficial {
	db := dbprovider.GetInstance().GetDB()
	officialPoll := []structs.CollegePollOfficial{}

	err := db.Where("season_id = ?", seasonID).Find(&officialPoll).Error
	if err != nil {
		return []structs.CollegePollOfficial{}
	}

	return officialPoll
}
