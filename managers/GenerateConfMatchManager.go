package managers

import (
	"math/rand"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
)

func GenerateSimCBBConferenceSchedules() {
	db := dbprovider.GetInstance().GetDB()
	collegeTeams := GetAllActiveCollegeTeams()

	ts := GetTimestamp()

	finalMatchesUpload := []structs.Match{}

	// Map everything into individual conferences
	conferenceMap := make(map[uint][]structs.Team)
	for _, team := range collegeTeams {
		conferenceMap[team.ConferenceID] = append(conferenceMap[team.ConferenceID], team)
	}

	ConferenceIDList := make([]uint, 0, len(conferenceMap))
	for conferenceID := range conferenceMap {
		ConferenceIDList = append(ConferenceIDList, conferenceID)
	}

	// Ideally should have list from 1-31.

	for _, conferenceID := range ConferenceIDList {
		teams := conferenceMap[conferenceID]
		if len(teams) == 0 {
			continue
		}

		scheduleTemplate := getScheduleTemplate(len(teams), conferenceID)

		matches := GenerateConferenceSchedule(teams, scheduleTemplate, ts, conferenceID)
		finalMatchesUpload = append(finalMatchesUpload, matches...)
	}

	repository.CreateCollegeMatchesRecordsBatch(db, finalMatchesUpload, 200)
}

// GenerateConferenceSchedule builds all 56 conference games (14 per team,
// home-and-away vs. every opponent) for an 8-team conference spanning weeks 9-15.
// The teams slice must contain exactly 8 teams in the desired seeding order.
func GenerateConferenceSchedule(teams []structs.Team, entries []scheduleEntry, ts structs.Timestamp, conferenceID uint) []structs.Match {

	var shuffledTeams []structs.Team

	switch conferenceID {
	case 1: // ACC – fixed protected-pair ordering
		shuffledTeams = sortTeamsForConferencePairs(teams, accPairs)
	case 2: // Big Ten – fixed protected-pair ordering
		shuffledTeams = sortTeamsForConferencePairs(teams, bigTenPairs)
	default:
		shuffledTeams = make([]structs.Team, len(teams))
		copy(shuffledTeams, teams)
		rand.Shuffle(len(shuffledTeams), func(i, j int) {
			shuffledTeams[i], shuffledTeams[j] = shuffledTeams[j], shuffledTeams[i]
		})
	}

	matches := make([]structs.Match, 0, len(entries))

	for _, e := range entries {
		home := shuffledTeams[e.homeIdx]
		away := shuffledTeams[e.awayIdx]

		match := structs.Match{
			SeasonID:      ts.SeasonID,
			Week:          e.week,
			WeekID:        ts.CollegeWeekID + e.week,
			MatchOfWeek:   e.slot,
			HomeTeamID:    home.ID,
			HomeTeam:      home.Abbr,
			HomeTeamCoach: home.Coach,
			Arena:         home.Arena,
			City:          home.City,
			State:         home.State,
			AwayTeamID:    away.ID,
			AwayTeam:      away.Abbr,
			AwayTeamCoach: away.Coach,
			IsConference:  true,
		}

		matches = append(matches, match)
	}

	return matches
}

func getScheduleTemplate(numOfTeams int, conferenceID uint) []scheduleEntry {
	switch numOfTeams {
	case 8:
		return eightTeamSchedule
	case 9:
		return nineTeamSchedule
	case 10:
		return tenTeamSchedule
	case 11:
		if conferenceID == 1 {
			return elevenTeam18GameSchedule
		} else {
			return elevenTeam20GameSchedule
		}
	case 12:
		return twelveTeamSchedule
	case 13:
		return thirteenTeamSchedule
	case 14:
		return fourteenTeamSchedule
	case 15:
		return fifteenTeamSchedule
	case 16:
		return sixteenTeamSchedule
	case 18:
		return eighteenTeamSchedule
	default:
		return []scheduleEntry{}
	}
}
