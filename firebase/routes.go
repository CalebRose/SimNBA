package firebase

import (
	"fmt"
	"strconv"
	"strings"
)

// ─────────────────────────────────────────────
// App routes (used as linkTo values in notifications)
// ─────────────────────────────────────────────

// BuildTeamRosterRoute returns the frontend route for a team roster page.
// league should be "cbb" or "nba".
func BuildTeamRosterRoute(league string, teamID uint) string {
	return fmt.Sprintf("/%s/team/%s", league, UintToString(teamID))
}

// BuildTeamRecruitingRoute returns the frontend route for a team's recruiting page.
func BuildTeamRecruitingRoute(league string, teamID uint) string {
	return fmt.Sprintf("/%s/recruiting", league)
}

// BuildTeamFreeAgencyRoute returns the frontend route for the NBA free agency page.
func BuildTeamFreeAgencyRoute(teamID uint) string {
	return "/nba/free-agency"
}

// BuildForumThreadRoute returns the frontend route for a specific forum thread.
func BuildForumThreadRoute(threadID string) string {
	return fmt.Sprintf("/forums/thread/%s", threadID)
}

// BuildForumPostRoute returns the frontend route for a specific post within a thread.
func BuildForumPostRoute(threadID, postID string) string {
	return fmt.Sprintf("/forums/thread/%s#post-%s", threadID, postID)
}

// BuildSourceEventKey returns a colon-separated idempotency key from the given
// parts (e.g. "injury", "cbb", gameID, playerID).
func BuildSourceEventKey(parts ...string) string {
	return strings.Join(parts, ":")
}

// UintToString converts a uint to its decimal string representation.
func UintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
