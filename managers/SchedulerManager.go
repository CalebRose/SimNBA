package managers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	fbsvc "github.com/CalebRose/SimNBA/firebase"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
)

// ─────────────────────────────────────────────
// CBB Game Request
// ─────────────────────────────────────────────

// CreateCBBGameRequest saves a new CBBGameRequest record to the database.
func CreateCBBGameRequest(request structs.CBBGameRequest) {
	db := dbprovider.GetInstance().GetDB()
	repository.CreateCBBGameRequest(request, db)
}

// AcceptCBBGameRequest marks the request as accepted and notifies the sending
// team's coach if they are a user-managed team.
func AcceptCBBGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCBBGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	request.Accepted()
	repository.SaveCBBGameRequest(request, db)

	sendingTeam := GetTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	if isCBBUserTeam(sendingTeam) {
		receivingTeam := GetTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
		ctx := context.Background()
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "cbb",
			Domain:         fbsvc.DomainCBB,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("%s has accepted your game request for Week %d.", receivingTeam.Team, request.Week),
			SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "cbb", "accept", requestID),
		})
	}
}

// RejectCBBGameRequest deletes the request and notifies the sending team's coach
// if they are a user-managed team.
func RejectCBBGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCBBGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	sendingTeam := GetTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))

	repository.DeleteCBBGameRequest(request, db)

	if isCBBUserTeam(sendingTeam) {
		receivingTeam := GetTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))
		ctx := context.Background()
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "cbb",
			Domain:         fbsvc.DomainCBB,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        fmt.Sprintf("%s has rejected your game request for Week %d.", receivingTeam.Team, request.Week),
			SourceEventKey: fbsvc.BuildSourceEventKey("gamerequest", "cbb", "reject", requestID),
		})
	}
}

// ProcessCBBGameRequest creates a Match record from the accepted game request
// and marks the request as approved.
func ProcessCBBGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCBBGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	homeTeam := GetTeamByTeamID(strconv.Itoa(int(request.HomeTeamID)))
	awayTeam := GetTeamByTeamID(strconv.Itoa(int(request.AwayTeamID)))
	arena := repository.FindArenaByID(request.ArenaID)

	match := structs.Match{
		HomeTeamID:    request.HomeTeamID,
		HomeTeam:      homeTeam.Team,
		AwayTeamID:    request.AwayTeamID,
		AwayTeam:      awayTeam.Team,
		Week:          request.Week,
		WeekID:        request.WeekID,
		SeasonID:      request.SeasonID,
		Arena:         arena.ArenaName,
		City:          arena.City,
		State:         arena.State,
		TimeSlot:      request.Timeslot,
		IsNeutralSite: request.IsNeutralSite,
		IsConference:  homeTeam.ConferenceID == awayTeam.ConferenceID,
	}

	_ = repository.CreateCollegeMatchesRecordsBatch(db, []structs.Match{match}, 1)

	request.Approved()
	repository.SaveCBBGameRequest(request, db)
}

// VetoCBBGameRequest deletes the request and notifies both teams' coaches
// if either is a user-managed team.
func VetoCBBGameRequest(requestID string) {
	db := dbprovider.GetInstance().GetDB()

	request := repository.FindCBBGameRequestRecord(repository.SchedulerQuery{ID: requestID})

	sendingTeam := GetTeamByTeamID(strconv.Itoa(int(request.SendingTeamID)))
	receivingTeam := GetTeamByTeamID(strconv.Itoa(int(request.RequestingTeamID)))

	repository.DeleteCBBGameRequest(request, db)

	ctx := context.Background()
	msg := fmt.Sprintf("The game request between %s and %s for Week %d has been vetoed.", sendingTeam.Team, receivingTeam.Team, request.Week)
	vetoKey := fbsvc.BuildSourceEventKey("gamerequest", "cbb", "veto", requestID)

	if isCBBUserTeam(sendingTeam) {
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{sendingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "cbb",
			Domain:         fbsvc.DomainCBB,
			TeamID:         sendingTeam.ID,
			RecipientUIDs:  uids,
			Message:        msg,
			SourceEventKey: vetoKey + ":sending",
		})
	}
	if isCBBUserTeam(receivingTeam) {
		uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{receivingTeam.Coach})
		_ = fbsvc.NotifyScheduleEvent(ctx, fbsvc.ScheduleEventNotificationInput{
			League:         "cbb",
			Domain:         fbsvc.DomainCBB,
			TeamID:         receivingTeam.ID,
			RecipientUIDs:  uids,
			Message:        msg,
			SourceEventKey: vetoKey + ":receiving",
		})
	}
}

// isCBBUserTeam returns true if the given Team is managed by a human coach.
func isCBBUserTeam(team structs.Team) bool {
	return team.IsUserCoached && team.Coach != ""
}
