package firebase

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotifyTeamInjury sends a push notification to a team's registered
// coaches / owners when a player is injured during a game.
// The call is idempotent: if a notification with the same SourceEventKey
// already exists for each recipient it will not be written again.
func NotifyTeamInjury(ctx context.Context, in TeamInjuryNotificationInput) error {
	if len(in.RecipientUIDs) == 0 {
		return nil
	}

	weeksStr := fmt.Sprintf("%d week(s)", in.WeeksOfRecovery)
	if in.WeeksOfRecovery == 1 {
		weeksStr = "1 week"
	}

	message := fmt.Sprintf(
		"%s (%s) was injured during a %s game. Injury: %s — estimated recovery %s.",
		in.PlayerName, in.Position, in.League, in.InjuryType, weeksStr,
	)

	notification := ForumNotification{
		Type:           NotificationTypeInjury,
		Domain:         in.Domain,
		Message:        message,
		LinkTo:         BuildTeamRosterRoute(in.League, in.TeamID),
		IsRead:         false,
		CreatedAt:      time.Now(),
		SourceEventKey: in.SourceEventKey,
	}

	return writeNotificationsIfNew(ctx, in.RecipientUIDs, notification)
}

// NotifyRecruitSigned sends a push notification to the winning team's coach
// when a recruit commits.
func NotifyRecruitSigned(ctx context.Context, in RecruitSignedNotificationInput) error {
	if len(in.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf(
		"%s has committed to %s.",
		in.RecruitName, in.TeamName,
	)

	notification := ForumNotification{
		Type:           NotificationTypeRecruiting,
		Domain:         in.Domain,
		Message:        message,
		LinkTo:         BuildTeamRecruitingRoute(in.League, in.TeamID),
		IsRead:         false,
		CreatedAt:      time.Now(),
		SourceEventKey: in.SourceEventKey,
	}

	return writeNotificationsIfNew(ctx, in.RecipientUIDs, notification)
}

// NotifyTransferPortalSigning sends a push notification to the winning team's
// coach when a transfer portal player signs with their team.
func NotifyTransferPortalSigning(ctx context.Context, in TransferPortalSigningNotificationInput) error {
	if len(in.RecipientUIDs) == 0 {
		return nil
	}

	stars := ""
	for i := 0; i < in.Stars; i++ {
		stars += "★"
	}

	message := fmt.Sprintf(
		"%s (%s, %s%s) has transferred from %s to %s.",
		in.PlayerName, in.Position, stars, func() string {
			if in.Stars > 0 {
				return " "
			}
			return ""
		}(),
		in.PreviousTeam, in.TeamName,
	)

	notification := ForumNotification{
		Type:           NotificationTypeTransfer,
		Domain:         DomainCBB,
		Message:        message,
		LinkTo:         BuildTeamRosterRoute("cbb", in.TeamID),
		IsRead:         false,
		CreatedAt:      time.Now(),
		SourceEventKey: in.SourceEventKey,
	}

	return writeNotificationsIfNew(ctx, in.RecipientUIDs, notification)
}

// NotifyNBAFreeAgentSigned sends a push notification to an NBA team's staff
// when a free agent signs with their team.
func NotifyNBAFreeAgentSigned(ctx context.Context, in NBAFreeAgentSignedNotificationInput) error {
	if len(in.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf(
		"%s (%s) has signed with %s on a %d-year deal worth $%.1fM per year.",
		in.PlayerName, in.Position, in.TeamName, in.TotalYears, in.ContractValue,
	)

	notification := ForumNotification{
		Type:           NotificationTypeFreeAgency,
		Domain:         DomainNBA,
		Message:        message,
		LinkTo:         BuildTeamRosterRoute("nba", in.TeamID),
		IsRead:         false,
		CreatedAt:      time.Now(),
		SourceEventKey: in.SourceEventKey,
	}

	return writeNotificationsIfNew(ctx, in.RecipientUIDs, notification)
}

// ─────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────

// writeNotificationsIfNew writes one notification document per UID, skipping
// any recipient that already has a document with the same SourceEventKey in
// their notifications subcollection.
func writeNotificationsIfNew(ctx context.Context, uids []string, base ForumNotification) error {
	client := GetFirestoreClient()

	var lastErr error
	for _, uid := range uids {
		if uid == "" {
			continue
		}

		col := client.Collection("users").Doc(uid).Collection("notifications")

		if base.SourceEventKey != "" {
			exists, err := notificationExists(ctx, col, uid, base.SourceEventKey)
			if err != nil {
				log.Printf("firebase: checking notification existence for uid=%s key=%s: %v", uid, base.SourceEventKey, err)
				continue
			}
			if exists {
				continue
			}
		}

		docRef := col.NewDoc()
		n := base
		n.ID = docRef.ID
		n.UID = uid

		if _, err := docRef.Set(ctx, n); err != nil {
			log.Printf("firebase: writing notification for uid=%s: %v", uid, err)
			lastErr = err
		}
	}

	return lastErr
}

// notificationExists returns true if any document in the collection already
// has the given sourceEventKey value.
func notificationExists(ctx context.Context, col *firestore.CollectionRef, uid, sourceEventKey string) (bool, error) {
	iter := col.Where("sourceEventKey", "==", sourceEventKey).Limit(1).Documents(ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if status.Code(err) == codes.OK || err == nil {
		return true, nil
	}
	return false, err
}

// NotifyScheduleEvent notifies a coach about a game-request lifecycle event
// such as acceptance, rejection, or an admin veto. Idempotent via SourceEventKey.
func NotifyScheduleEvent(ctx context.Context, in ScheduleEventNotificationInput) error {
	if len(in.RecipientUIDs) == 0 {
		return nil
	}
	return writeNotificationsIfNew(ctx, in.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeSystem,
		Domain:         in.Domain,
		LinkTo:         BuildTeamRosterRoute(in.League, in.TeamID),
		Message:        in.Message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		CreatedAt:      time.Now(),
		SourceEventKey: in.SourceEventKey,
	})
}
