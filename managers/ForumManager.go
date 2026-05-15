package managers

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	fbsvc "github.com/CalebRose/SimNBA/firebase"
	"github.com/CalebRose/SimNBA/structs"
)

// ─────────────────────────────────────────────
// Forum IDs & paths
// ─────────────────────────────────────────────

const (
	cbbPostGameForumID = "postgame-discussions-simcbb"
	cbbMediaForumID    = "media-simcbb"
	nbaMediaForumID    = "media-simnba"
)

// ─────────────────────────────────────────────
// CBB Post-game discussion thread
// ─────────────────────────────────────────────

// CreatePostGameDiscussionThreadForCBBGame creates a system-generated post-game
// discussion thread for a completed (and user-coached) college basketball game.
// Intended to be called in a goroutine after the game record is saved.
// The operation is idempotent: calling it twice for the same game has no effect.
func CreatePostGameDiscussionThreadForCBBGame(
	game structs.Match,
	seasonID uint,
	homeTeamStats structs.TeamStats,
	awayTeamStats structs.TeamStats,
	homePlayerStats []structs.CollegePlayerStats,
	awayPlayerStats []structs.CollegePlayerStats,
	collegePlayersMap map[uint]structs.CollegePlayer,
) {
	ctx := context.Background()

	gameID := strconv.Itoa(int(game.ID))
	eventKey := fmt.Sprintf("postgame_thread:cbb:season%d:game%s", seasonID, gameID)

	title := buildCBBPostGameThreadTitle(game)
	nodes := buildCBBPostGameNodes(game, homeTeamStats, awayTeamStats, homePlayerStats, awayPlayerStats, collegePlayersMap)
	bodyText := nodesToPlainText(nodes)
	richBody := buildRichDoc(nodes)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           cbbPostGameForumID,
		ForumPath:         []string{"postgame-discussions", "simcbb"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeGameReference,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedGameID:  gameID,
		ReferencedLeague:  "cbb",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create CBB post-game thread for game %s: %v", gameID, err)
		return
	}

	log.Printf("ForumManager: created CBB post-game thread %s for game %s (%s)", thread.ID, gameID, title)
}

// buildCBBPostGameThreadTitle returns a forum-ready title for a CBB game.
func buildCBBPostGameThreadTitle(game structs.Match) string {
	season := game.SeasonID + 2024
	if game.MatchName != "" {
		return fmt.Sprintf("[%d] Week %d: %s: %s vs %s", season, game.Week, game.MatchName, game.AwayTeam, game.HomeTeam)
	}
	return fmt.Sprintf("[%d] Week %d: %s at %s", season, game.Week, game.AwayTeam, game.HomeTeam)
}

// buildCBBPostGameNodes builds the rich-text node list for a CBB post-game thread.
func buildCBBPostGameNodes(
	game structs.Match,
	homeTeamStats structs.TeamStats,
	awayTeamStats structs.TeamStats,
	homePlayerStats []structs.CollegePlayerStats,
	awayPlayerStats []structs.CollegePlayerStats,
	collegePlayersMap map[uint]structs.CollegePlayer,
) []map[string]interface{} {
	nodes := []map[string]interface{}{}
	homeTeam := game.HomeTeam
	awayTeam := game.AwayTeam

	// ── Final score ──────────────────────────────────────────────────────────
	nodes = append(nodes, rtBoldParagraph(fmt.Sprintf(
		"FINAL: %s %d, %s %d",
		awayTeam, game.AwayTeamScore, homeTeam, game.HomeTeamScore,
	)))

	// ── Scoring by half ───────────────────────────────────────────────────────
	nodes = append(nodes, rtHeading(3, "Scoring by Half"))
	hasOT := awayTeamStats.OvertimeScore > 0 || homeTeamStats.OvertimeScore > 0
	if hasOT {
		nodes = append(nodes, rtTableNode(
			[]string{"Team", "H1", "H2", "OT", "Total"},
			[][]string{
				{awayTeam,
					fmt.Sprintf("%d", awayTeamStats.FirstHalfScore),
					fmt.Sprintf("%d", awayTeamStats.SecondHalfScore),
					fmt.Sprintf("%d", awayTeamStats.OvertimeScore),
					fmt.Sprintf("%d", game.AwayTeamScore)},
				{homeTeam,
					fmt.Sprintf("%d", homeTeamStats.FirstHalfScore),
					fmt.Sprintf("%d", homeTeamStats.SecondHalfScore),
					fmt.Sprintf("%d", homeTeamStats.OvertimeScore),
					fmt.Sprintf("%d", game.HomeTeamScore)},
			},
		))
	} else {
		nodes = append(nodes, rtTableNode(
			[]string{"Team", "H1", "H2", "Total"},
			[][]string{
				{awayTeam,
					fmt.Sprintf("%d", awayTeamStats.FirstHalfScore),
					fmt.Sprintf("%d", awayTeamStats.SecondHalfScore),
					fmt.Sprintf("%d", game.AwayTeamScore)},
				{homeTeam,
					fmt.Sprintf("%d", homeTeamStats.FirstHalfScore),
					fmt.Sprintf("%d", homeTeamStats.SecondHalfScore),
					fmt.Sprintf("%d", game.HomeTeamScore)},
			},
		))
	}

	// ── Venue ─────────────────────────────────────────────────────────────────
	venueStr := game.Arena
	if game.City != "" || game.State != "" {
		loc := game.City
		if game.State != "" {
			if loc != "" {
				loc += ", " + game.State
			} else {
				loc = game.State
			}
		}
		if loc != "" {
			venueStr += " — " + loc
		}
	}
	if venueStr != "" {
		nodes = append(nodes, rtParagraph("Venue: "+venueStr))
	}

	// ── Team stats ────────────────────────────────────────────────────────────
	nodes = append(nodes, rtHeading(3, "Team Stats"))
	nodes = append(nodes, rtTableNode(
		[]string{"Team", "FG", "FG%", "3PT", "3PT%", "FT", "FT%", "REB", "AST", "STL", "BLK", "TO"},
		[][]string{
			{awayTeam,
				fmt.Sprintf("%d/%d", awayTeamStats.FGM, awayTeamStats.FGA),
				fmt.Sprintf("%.1f%%", awayTeamStats.FGPercent),
				fmt.Sprintf("%d/%d", awayTeamStats.ThreePointsMade, awayTeamStats.ThreePointAttempts),
				fmt.Sprintf("%.1f%%", awayTeamStats.ThreePointPercent),
				fmt.Sprintf("%d/%d", awayTeamStats.FTM, awayTeamStats.FTA),
				fmt.Sprintf("%.1f%%", awayTeamStats.FTPercent),
				fmt.Sprintf("%d", awayTeamStats.Rebounds),
				fmt.Sprintf("%d", awayTeamStats.Assists),
				fmt.Sprintf("%d", awayTeamStats.Steals),
				fmt.Sprintf("%d", awayTeamStats.Blocks),
				fmt.Sprintf("%d", awayTeamStats.TotalTurnovers)},
			{homeTeam,
				fmt.Sprintf("%d/%d", homeTeamStats.FGM, homeTeamStats.FGA),
				fmt.Sprintf("%.1f%%", homeTeamStats.FGPercent),
				fmt.Sprintf("%d/%d", homeTeamStats.ThreePointsMade, homeTeamStats.ThreePointAttempts),
				fmt.Sprintf("%.1f%%", homeTeamStats.ThreePointPercent),
				fmt.Sprintf("%d/%d", homeTeamStats.FTM, homeTeamStats.FTA),
				fmt.Sprintf("%.1f%%", homeTeamStats.FTPercent),
				fmt.Sprintf("%d", homeTeamStats.Rebounds),
				fmt.Sprintf("%d", homeTeamStats.Assists),
				fmt.Sprintf("%d", homeTeamStats.Steals),
				fmt.Sprintf("%d", homeTeamStats.Blocks),
				fmt.Sprintf("%d", homeTeamStats.TotalTurnovers)},
		},
	))

	// ── Player stats ──────────────────────────────────────────────────────────
	nodes = appendCBBPlayerStatTable(nodes, awayTeam, awayPlayerStats, collegePlayersMap)
	nodes = appendCBBPlayerStatTable(nodes, homeTeam, homePlayerStats, collegePlayersMap)

	nodes = append(nodes, rtParagraph("Postgame discussion is open. Share your reactions below."))
	return nodes
}

// appendCBBPlayerStatTable appends a per-player stats section for one team.
func appendCBBPlayerStatTable(
	nodes []map[string]interface{},
	teamName string,
	playerStats []structs.CollegePlayerStats,
	playerMap map[uint]structs.CollegePlayer,
) []map[string]interface{} {
	type row struct {
		label string
		stats structs.CollegePlayerStats
	}

	rows := make([]row, 0, len(playerStats))
	for _, s := range playerStats {
		p, ok := playerMap[uint(s.CollegePlayerID)]
		if !ok {
			continue
		}
		label := fmt.Sprintf("%s %s (%s)", p.FirstName, p.LastName, p.Position)
		rows = append(rows, row{label: label, stats: s})
	}

	// Sort by points descending.
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].stats.Points > rows[j].stats.Points
	})

	if len(rows) == 0 {
		return nodes
	}

	nodes = append(nodes, rtHeading(3, teamName+" — Player Stats"))
	tableRows := make([][]string, 0, len(rows))
	for _, r := range rows {
		s := r.stats
		tableRows = append(tableRows, []string{
			r.label,
			fmt.Sprintf("%d", s.Minutes),
			fmt.Sprintf("%d", s.Points),
			fmt.Sprintf("%d", s.TotalRebounds),
			fmt.Sprintf("%d", s.Assists),
			fmt.Sprintf("%d", s.Steals),
			fmt.Sprintf("%d", s.Blocks),
			fmt.Sprintf("%d", s.Turnovers),
			fmt.Sprintf("%d/%d", s.FGM, s.FGA),
			fmt.Sprintf("%d/%d", s.ThreePointsMade, s.ThreePointAttempts),
			fmt.Sprintf("%d/%d", s.FTM, s.FTA),
		})
	}

	nodes = append(nodes, rtTableNode(
		[]string{"Player", "MIN", "PTS", "REB", "AST", "STL", "BLK", "TO", "FG", "3PT", "FT"},
		tableRows,
	))
	return nodes
}

// ─────────────────────────────────────────────
// NBA Draft Lottery forum thread
// ─────────────────────────────────────────────

// CreateDraftLotteryForumThread creates a system-generated forum thread
// summarising the NBA draft lottery results for the given season.
// draftPicks should include all lottery picks (typically the first 14 R1 picks).
// The operation is idempotent: calling it twice for the same season has no effect.
func CreateDraftLotteryForumThread(season int, draftPicks []structs.DraftPick) {
	ctx := context.Background()

	title := fmt.Sprintf("SimNBA: Season %d Draft Lottery Results", season)
	eventKey := fmt.Sprintf("draft_lottery:nba:season%d", season)

	nodes := buildDraftLotteryNodes(season, draftPicks)
	bodyText := nodesToPlainText(nodes)
	richBody := buildRichDoc(nodes)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           nbaMediaForumID,
		ForumPath:         []string{"media", "simnba"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "nba",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create draft lottery thread for season %d: %v", season, err)
		return
	}

	log.Printf("ForumManager: created draft lottery thread %s for season %d", thread.ID, season)
}

// buildDraftLotteryNodes builds the rich-text node list for a draft lottery thread.
func buildDraftLotteryNodes(season int, draftPicks []structs.DraftPick) []map[string]interface{} {
	nodes := []map[string]interface{}{}

	nodes = append(nodes, rtBoldParagraph(fmt.Sprintf("Season %d NBA Draft Lottery", season)))
	nodes = append(nodes, rtParagraph("The draft lottery has been conducted. Here are the results:"))

	// Separate and sort picks by round then pick number.
	var r1, r2 []structs.DraftPick
	for _, p := range draftPicks {
		if p.DraftRound == 1 {
			r1 = append(r1, p)
		} else if p.DraftRound == 2 {
			r2 = append(r2, p)
		}
	}
	sort.Slice(r1, func(i, j int) bool { return r1[i].DraftNumber < r1[j].DraftNumber })
	sort.Slice(r2, func(i, j int) bool { return r2[i].DraftNumber < r2[j].DraftNumber })

	if len(r1) > 0 {
		nodes = append(nodes, rtHeading(3, "Round 1 Pick Order"))
		pickRows := make([][]string, 0, len(r1))
		for _, p := range r1 {
			notes := p.Notes
			if notes == "" {
				notes = "—"
			}
			pickRows = append(pickRows, []string{
				fmt.Sprintf("#%d", p.DraftNumber),
				p.Team,
				notes,
			})
		}
		nodes = append(nodes, rtTableNode([]string{"Pick", "Team", "Via"}, pickRows))
	}

	if len(r2) > 0 {
		nodes = append(nodes, rtHeading(3, "Round 2 Pick Order"))
		pickRows := make([][]string, 0, len(r2))
		for _, p := range r2 {
			notes := p.Notes
			if notes == "" {
				notes = "—"
			}
			pickRows = append(pickRows, []string{
				fmt.Sprintf("#%d", p.DraftNumber),
				p.Team,
				notes,
			})
		}
		nodes = append(nodes, rtTableNode([]string{"Pick", "Team", "Via"}, pickRows))
	}

	nodes = append(nodes, rtParagraph("React to the draft lottery results below!"))
	return nodes
}

// ─────────────────────────────────────────────
// CBB Recruiting sync forum thread
// ─────────────────────────────────────────────

// CreateRecruitingSyncForumThread creates a system-generated weekly thread
// listing every recruit that committed to a CBB program during the sync.
// signings is a list of human-readable labels built at commit time.
// The operation is idempotent: calling it twice for the same season/week has no effect.
func CreateRecruitingSyncForumThread(season, week int, signings []string) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCBB: Season %d Week %d Recruiting Commitments", season, week)
	eventKey := fmt.Sprintf("recruiting_sync:cbb:season%d:week%d", season, week)

	paragraphs := buildRecruitingSyncParagraphs(season, week, signings)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           cbbMediaForumID,
		ForumPath:         []string{"media", "simcbb"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "cbb",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create recruiting sync thread for season %d week %d: %v", season, week, err)
		return
	}

	log.Printf("ForumManager: created recruiting sync thread %s for season %d week %d", thread.ID, season, week)
}

func buildRecruitingSyncParagraphs(season, week int, signings []string) []string {
	var paras []string
	count := len(signings)

	if count == 0 {
		paras = append(paras, fmt.Sprintf(
			"Week %d recruiting is complete for Season %d. No recruits committed to a program this week.",
			week, season,
		))
	} else {
		paras = append(paras, fmt.Sprintf(
			"Week %d recruiting results are in for Season %d. A total of %d recruit(s) have committed to a program this week.",
			week, season, count,
		))
		for _, label := range signings {
			paras = append(paras, label)
		}
	}

	paras = append(paras, "React to the latest commitments and discuss your team's recruiting class below!")
	return paras
}

// ─────────────────────────────────────────────
// CBB Transfer Portal sync forum thread
// ─────────────────────────────────────────────

// CreateTransferPortalSyncForumThread creates a system-generated thread
// summarising the signings from a single transfer portal sync round.
// signings is a list of human-readable labels for every player that signed.
// The operation is idempotent: calling it twice for the same season/round has no effect.
func CreateTransferPortalSyncForumThread(season, round int, signings []string) {
	ctx := context.Background()

	title := fmt.Sprintf("SimCBB: Season %d Transfer Portal — Round %d Results", season, round)
	eventKey := fmt.Sprintf("transfer_portal_sync:cbb:season%d:round%d", season, round)

	paragraphs := buildTransferPortalSyncParagraphs(season, round, signings)
	bodyText := strings.Join(paragraphs, "\n\n")
	richBody := buildRichPostBody(paragraphs)

	input := fbsvc.CreateForumThreadInput{
		ForumID:           cbbMediaForumID,
		ForumPath:         []string{"media", "simcbb"},
		Title:             title,
		AuthorUID:         "system",
		AuthorUsername:    "SimSN",
		AuthorDisplayName: "SimSN System",
		CreatedByType:     fbsvc.CreatedBySystem,
		ThreadType:        fbsvc.ThreadTypeStandard,
		FirstPostBodyText: bodyText,
		FirstPostBody:     richBody,
		ReferencedLeague:  "cbb",
		ExternalEventKey:  eventKey,
	}

	thread, err := fbsvc.CreateThread(ctx, input)
	if err != nil {
		log.Printf("ForumManager: failed to create transfer portal sync thread for season %d round %d: %v", season, round, err)
		return
	}

	log.Printf("ForumManager: created transfer portal sync thread %s for season %d round %d", thread.ID, season, round)
}

func buildTransferPortalSyncParagraphs(season, round int, signings []string) []string {
	var paras []string
	count := len(signings)

	if count == 0 {
		paras = append(paras, fmt.Sprintf(
			"Transfer Portal Round %d is complete for Season %d. No players signed with new programs this round.",
			round, season,
		))
	} else {
		paras = append(paras, fmt.Sprintf(
			"Transfer Portal Round %d results are in for Season %d. A total of %d player(s) have signed with new programs this round.",
			round, season, count,
		))
		for _, label := range signings {
			paras = append(paras, label)
		}
	}

	paras = append(paras, "Discuss the latest transfer portal news below!")
	return paras
}

// ─────────────────────────────────────────────
// Rich text helpers
// ─────────────────────────────────────────────

// buildRichPostBody converts a slice of paragraph strings into a ProseMirror document.
func buildRichPostBody(paragraphs []string) map[string]interface{} {
	content := make([]map[string]interface{}, 0, len(paragraphs))
	for _, p := range paragraphs {
		content = append(content, map[string]interface{}{
			"type": "paragraph",
			"content": []map[string]interface{}{
				{"type": "text", "text": p},
			},
		})
	}
	return map[string]interface{}{
		"type":    "doc",
		"content": content,
	}
}

// buildRichDoc wraps a slice of content nodes into a top-level ProseMirror doc.
func buildRichDoc(nodes []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":    "doc",
		"content": nodes,
	}
}

// rtParagraph creates a plain paragraph node.
func rtParagraph(text string) map[string]interface{} {
	return map[string]interface{}{
		"type": "paragraph",
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
	}
}

// rtBoldParagraph creates a paragraph with bold-marked text.
func rtBoldParagraph(text string) map[string]interface{} {
	return map[string]interface{}{
		"type": "paragraph",
		"content": []map[string]interface{}{
			{
				"type":  "text",
				"text":  text,
				"marks": []map[string]interface{}{{"type": "bold"}},
			},
		},
	}
}

// rtHeading creates a heading node at the given level (1–6).
func rtHeading(level int, text string) map[string]interface{} {
	return map[string]interface{}{
		"type":  "heading",
		"attrs": map[string]interface{}{"level": level, "textAlign": "left"},
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
	}
}

// rtTableCell creates a single table header or data cell.
func rtTableCell(text string, isHeader bool) map[string]interface{} {
	cellType := "tableCell"
	if isHeader {
		cellType = "tableHeader"
	}
	return map[string]interface{}{
		"type":  cellType,
		"attrs": map[string]interface{}{"colspan": 1, "rowspan": 1, "colwidth": nil},
		"content": []map[string]interface{}{
			{
				"type":  "paragraph",
				"attrs": map[string]interface{}{"textAlign": nil},
				"content": []map[string]interface{}{
					{"type": "text", "text": text},
				},
			},
		},
	}
}

// rtTableNode builds a TipTap-compatible table node from header strings and row data.
func rtTableNode(headers []string, rows [][]string) map[string]interface{} {
	tableRows := []map[string]interface{}{}

	headerCells := make([]map[string]interface{}, len(headers))
	for i, h := range headers {
		headerCells[i] = rtTableCell(h, true)
	}
	tableRows = append(tableRows, map[string]interface{}{
		"type":    "tableRow",
		"content": headerCells,
	})

	for _, row := range rows {
		cells := make([]map[string]interface{}, len(row))
		for i, cell := range row {
			cells[i] = rtTableCell(cell, false)
		}
		tableRows = append(tableRows, map[string]interface{}{
			"type":    "tableRow",
			"content": cells,
		})
	}

	return map[string]interface{}{
		"type":    "table",
		"content": tableRows,
	}
}

// nodesToPlainText extracts readable plain text from rich nodes for the bodyText field.
func nodesToPlainText(nodes []map[string]interface{}) string {
	var lines []string
	for _, node := range nodes {
		switch node["type"] {
		case "paragraph", "heading":
			if text := extractInlineText(node); text != "" {
				lines = append(lines, text)
			}
		case "table":
			if rows, ok := node["content"].([]map[string]interface{}); ok {
				for _, row := range rows {
					if cells, ok := row["content"].([]map[string]interface{}); ok {
						var cellTexts []string
						for _, cell := range cells {
							cellTexts = append(cellTexts, extractCellPlainText(cell))
						}
						lines = append(lines, strings.Join(cellTexts, "  |  "))
					}
				}
			}
		}
	}
	return strings.Join(lines, "\n\n")
}

func extractInlineText(node map[string]interface{}) string {
	if content, ok := node["content"].([]map[string]interface{}); ok {
		var texts []string
		for _, child := range content {
			if t, ok := child["text"].(string); ok {
				texts = append(texts, t)
			}
		}
		return strings.Join(texts, "")
	}
	return ""
}

func extractCellPlainText(cell map[string]interface{}) string {
	if content, ok := cell["content"].([]map[string]interface{}); ok {
		for _, para := range content {
			return extractInlineText(para)
		}
	}
	return ""
}
