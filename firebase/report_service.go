package firebase

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/CalebRose/SimNBA/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MissedRecruitingReports(profiles []structs.TeamRecruitingProfile, ts structs.Timestamp) []AdminReport {
	reports := make([]AdminReport, 0)
	if uint64(ts.CollegeWeekID) > math.MaxInt64 || uint64(ts.Phase) > math.MaxInt64 {
		log.Printf("RecruitingSyncManager: cannot report missed recruiting: week ID or phase exceeds Firestore integer range")
		return reports
	}
	for _, profile := range profiles {
		if !profile.IsUserTeam || profile.SpentPoints != 0 {
			continue
		}
		if uint64(profile.TeamID) > math.MaxInt64 {
			log.Printf("RecruitingSyncManager: cannot report missed recruiting for team %d: ID exceeds Firestore integer range", profile.TeamID)
			continue
		}
		subject := profile.TeamAbbr
		if profile.Recruiter != "" {
			subject = fmt.Sprintf("%s (%s)", profile.Recruiter, profile.TeamAbbr)
		}
		reports = append(reports, AdminReport{
			League:       "chl",
			TeamID:       int64(profile.TeamID),
			ReportType:   "missed_recruiting",
			Message:      fmt.Sprintf("%s spent no recruiting points in week %d.", subject, ts.CollegeWeek),
			ResponseType: AdminReportMissedRecruiting,
			Status:       "pending",
			Phase:        int64(ts.Phase),
			WeekID:       int64(ts.CollegeWeekID),
			ReportedBy:   "system",
		})
	}
	return reports
}

func TradeTransactionReport(teamID uint, message string, ts structs.Timestamp) (AdminReport, error) {
	if uint64(teamID) > math.MaxInt64 || uint64(ts.NBAWeekID) > math.MaxInt64 || uint64(ts.Phase) > math.MaxInt64 {
		return AdminReport{}, fmt.Errorf("firebase: trade report identifiers exceed Firestore integer range")
	}
	return AdminReport{
		League:       "nba",
		TeamID:       int64(teamID),
		ReportType:   "trade_transaction",
		Message:      message,
		ResponseType: AdminReportTradeTransaction,
		Status:       "pending",
		Phase:        int64(ts.Phase),
		WeekID:       int64(ts.NBAWeekID),
		ReportedBy:   "system",
	}, nil
}

func BrokenLineupReport(league string, teamID uint, message string, weekID uint, phase uint8) (AdminReport, error) {
	if uint64(teamID) > math.MaxInt64 || uint64(weekID) > math.MaxInt64 {
		return AdminReport{}, fmt.Errorf("firebase: lineup report identifiers exceed Firestore integer range")
	}
	return AdminReport{
		League:       league,
		TeamID:       int64(teamID),
		ReportType:   "broken_gameplan",
		Message:      message,
		ResponseType: AdminReportBrokenGameplan,
		Status:       "pending",
		Phase:        int64(phase),
		WeekID:       int64(weekID),
		ReportedBy:   "system",
	}, nil
}

// CreateAdminReport creates a report without overwriting an existing one.
func CreateAdminReport(ctx context.Context, id string, report AdminReport) error {
	if id == "" || strings.Contains(id, "/") {
		return fmt.Errorf("firebase: invalid admin report document ID %q", id)
	}
	client := GetFirestoreClient()
	if client == nil {
		return fmt.Errorf("firebase: failed to get Firestore client")
	}
	return createAdminReport(ctx, client, id, report)
}

func createAdminReport(ctx context.Context, client *firestore.Client, id string, report AdminReport) error {
	now := time.Now().UTC()
	report.CreatedAt = now
	report.UpdatedAt = now

	_, err := client.Collection("adminReports").Doc(id).Create(ctx, report)
	if status.Code(err) == codes.AlreadyExists {
		return nil
	}
	if err != nil {
		return fmt.Errorf("firebase: failed to create admin report %s: %w", id, err)
	}
	return nil
}
