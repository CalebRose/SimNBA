package managers

import (
	"encoding/csv"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func ExportCroots(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition", "attachment;filename=ezacos_secret_croot_list.csv")
	w.Header().Set("Transfer-Encoding", "chunked")

	writer := csv.NewWriter(w)

	croots := GetAllCollegeRecruits()

	HeaderRow := []string{
		"First Name", "Last Name", "Position",
		"Stars", "College", "State", "Country", "Height",
		"Overall", "Shooting 2s", "Shooting 3s", "Finishing",
		"Ballwork", "Rebounding", "InteriorDefense", "PerimeterDefense", "Potential Grade",
		"Personality", "Recruiting Bias", "Academic Bias", "Work Ethic",
		"ESPN Rank", "Rivals Rank", "247 Rank", "LeadingTeams",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, croot := range croots {
		var leadingAbbr []string

		for _, lt := range croot.LeadingTeams {
			leadingAbbr = append(leadingAbbr, lt.TeamAbbr)
		}

		crootRow := []string{
			croot.FirstName, croot.LastName, croot.Position, strconv.Itoa(croot.Stars),
			croot.College, croot.State, croot.Country, croot.Height,
			croot.OverallGrade, croot.Shooting2, croot.Shooting3, croot.Finishing,
			croot.Ballwork, croot.Rebounding, croot.InteriorDefense, croot.PerimeterDefense, croot.PotentialGrade,
			croot.Personality, croot.RecruitingBias, croot.AcademicBias, croot.WorkEthic,
			strconv.Itoa(int(croot.ESPNRank)), strconv.Itoa(int(croot.RivalsRank)), strconv.Itoa(int(croot.Rank247)), strings.Join(leadingAbbr, ", "),
		}

		err = writer.Write(crootRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportCollegePlayers(w http.ResponseWriter) {
	ts := GetTimestamp()
	season := strconv.Itoa(ts.Season)
	filename := "sagebows_" + season + "_secret_player_list.csv"
	w.Header().Set("Content-Disposition", "attachment;filename="+filename)
	w.Header().Set("Transfer-Encoding", "chunked")

	writer := csv.NewWriter(w)

	players := GetAllCollegePlayers()

	sort.Slice(players, func(i, j int) bool {
		return players[i].TeamID < players[j].TeamID
	})

	HeaderRow := []string{
		"College", "First Name", "Last Name", "Position", "Year", "Is_Redshirt", "Age",
		"Stars", "State", "Country", "Height",
		"Overall", "Shooting 2s", "Shooting 3s", "Free Throwing", "Finishing",
		"Ballwork", "Rebounding", "InteriorDefense", "PerimeterDefense", "Stamina", "Potential Grade",
		"Personality", "RecruitingBias", "Work Ethic", "Previous Team",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {

		shooting2Grade := util.GetAttributeGrade(player.Shooting2)
		shooting3Grade := util.GetAttributeGrade(player.Shooting3)
		freeThrowGrade := util.GetAttributeGrade(player.FreeThrow)
		finishingGrade := util.GetAttributeGrade(player.Finishing)
		reboundingGrade := util.GetAttributeGrade(player.Rebounding)
		ballworkGrade := util.GetAttributeGrade(player.Ballwork)
		interiorDefenseGrade := util.GetAttributeGrade(player.InteriorDefense)
		perimeterDefenseGrade := util.GetAttributeGrade(player.PerimeterDefense)
		potentialGrade := util.GetPotentialGrade(player.Potential)
		overallGrade := util.GetPlayerOverallGrade(player.Overall)
		sta := strconv.Itoa(player.Stamina)

		playerRow := []string{
			player.TeamAbbr, player.FirstName, player.LastName, player.Position, strconv.Itoa(player.Year), strconv.FormatBool(player.IsRedshirt), strconv.Itoa(player.Age),
			strconv.Itoa(player.Stars), player.State, player.Country, player.Height,
			overallGrade, shooting2Grade, shooting3Grade, freeThrowGrade, finishingGrade,
			ballworkGrade, reboundingGrade, interiorDefenseGrade, perimeterDefenseGrade, sta, potentialGrade,
			player.Personality, player.RecruitingBias, player.WorkEthic, player.PreviousTeam,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportCBBPreseasonRanks(w http.ResponseWriter) {
	db := dbprovider.GetInstance().GetDB()
	w.Header().Set("Content-Disposition", "attachment;filename=toucan_preseason_list.csv")
	w.Header().Set("Transfer-Encoding", "chunked")

	writer := csv.NewWriter(w)

	teams := GetAllActiveCollegeTeams()
	HeaderRow := []string{
		"College", "Conference", "Off. Rating", "Def. Rating",
		"Overall", "Star Power",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, t := range teams {
		teamID := strconv.Itoa(int(t.ID))
		var players []structs.CollegePlayer

		db.Order("overall desc").Where("team_id = ?", teamID).Find(&players)

		offenseRank := 0.0
		defenseRank := 0.0
		overall := 0.0
		starPower := 0.0
		count := 0

		for idx, player := range players {
			if idx > 9 {
				break
			}
			o := ((float64(player.Shooting2) + float64(player.Shooting3) + float64(player.Finishing) + float64(player.FreeThrow)) / 4)
			d := ((float64(player.Ballwork) + float64(player.Rebounding) + float64(player.InteriorDefense) + float64(player.PerimeterDefense)) / 4)
			offenseRank += o
			defenseRank += d
			starPower += float64(player.Stars)
			count++
		}

		offenseRank = offenseRank / float64(count)
		defenseRank = defenseRank / float64(count)

		starPower = starPower / float64(count)
		overall = offenseRank + defenseRank

		teamRow := []string{t.Abbr, t.Conference, strconv.FormatFloat(offenseRank, 'E', -1, 64),
			strconv.FormatFloat(defenseRank, 'E', -1, 64),
			strconv.FormatFloat(overall, 'E', -1, 64),
			strconv.FormatFloat(starPower, 'E', -1, 64)}

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportCBBRosterToCSV(TeamID string, w http.ResponseWriter) {
	// Get Team Data
	team := GetTeamByTeamID(TeamID)
	w.Header().Set("Content-Disposition", "attachment;filename="+team.Team+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetCollegePlayersByTeamId(TeamID)

	csvRoster := []structs.CollegePlayerResponse{}

	for _, player := range players {
		shooting2Grade := util.GetAttributeGrade(player.Shooting2)
		shooting3Grade := util.GetAttributeGrade(player.Shooting3)
		freeThrowGrade := util.GetAttributeGrade(player.FreeThrow)
		finishingGrade := util.GetAttributeGrade(player.Finishing)
		reboundingGrade := util.GetAttributeGrade(player.Rebounding)
		ballworkGrade := util.GetAttributeGrade(player.Ballwork)
		interiorDefenseGrade := util.GetAttributeGrade(player.InteriorDefense)
		perimeterDefenseGrade := util.GetAttributeGrade(player.PerimeterDefense)
		potentialGrade := util.GetPotentialGrade(player.Potential)
		overallGrade := util.GetPlayerOverallGrade(player.Overall)

		res := structs.CollegePlayerResponse{
			FirstName:             player.FirstName,
			LastName:              player.LastName,
			Position:              player.Position,
			Archetype:             player.Archetype,
			Age:                   player.Age,
			Year:                  player.Year,
			State:                 player.State,
			Country:               player.Country,
			Stars:                 player.Stars,
			Height:                player.Height,
			PotentialGrade:        potentialGrade,
			Shooting2Grade:        shooting2Grade,
			Shooting3Grade:        shooting3Grade,
			FreeThrowGrade:        freeThrowGrade,
			FinishingGrade:        finishingGrade,
			BallworkGrade:         ballworkGrade,
			ReboundingGrade:       reboundingGrade,
			InteriorDefenseGrade:  interiorDefenseGrade,
			PerimeterDefenseGrade: perimeterDefenseGrade,
			OverallGrade:          overallGrade,
			Stamina:               player.Stamina,
			PlaytimeExpectations:  player.PlaytimeExpectations,
			Minutes:               player.Minutes,
			Potential:             player.Potential,
			Personality:           player.Personality,
			RecruitingBias:        player.RecruitingBias,
			WorkEthic:             player.WorkEthic,
			AcademicBias:          player.AcademicBias,
			PlayerID:              player.PlayerID,
			TeamID:                player.TeamID,
			TeamAbbr:              player.TeamAbbr,
			IsRedshirting:         player.IsRedshirting,
			IsRedshirt:            player.IsRedshirt,
			PositionOne:           player.PositionOne,
			PositionTwo:           player.PositionTwo,
			PositionThree:         player.PositionThree,
			P1Minutes:             player.P1Minutes,
			P2Minutes:             player.P2Minutes,
			P3Minutes:             player.P3Minutes,
			InsideProportion:      player.InsideProportion,
			MidRangeProportion:    player.MidRangeProportion,
			ThreePointProportion:  player.ThreePointProportion,
		}

		csvRoster = append(csvRoster, res)
	}

	HeaderRow := []string{
		"Team", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Age", "Stars",
		"State", "Country", "Height", "Overall", "Finishing",
		"Shooting2", "Shooting3", "FreeThrow",
		"Ballwork", "Rebounding", "Interior Defense", "Perimeter Defense",
		"Playtime Expectations", "Stamina", "Potential",
		"Personality", "Recruiting Bias", "Work Ethic", "Academic Bias",
		"RedshirtingStatus",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, csvModel := range csvRoster {
		redshirtStatus := ""
		if csvModel.IsRedshirt {
			redshirtStatus = "Former Redshirt"
		} else if csvModel.IsRedshirting {
			redshirtStatus = "Currently Redshirting"
		}
		playerRow := []string{
			team.Team, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, strconv.Itoa(csvModel.Year), strconv.Itoa(csvModel.Age), strconv.Itoa(csvModel.Stars),
			csvModel.State, csvModel.Country, csvModel.Height, csvModel.OverallGrade, csvModel.FinishingGrade,
			csvModel.Shooting2Grade, csvModel.Shooting3Grade, csvModel.FreeThrowGrade,
			csvModel.BallworkGrade, csvModel.ReboundingGrade, csvModel.InteriorDefenseGrade, csvModel.PerimeterDefenseGrade,
			strconv.Itoa(csvModel.PlaytimeExpectations), strconv.Itoa(csvModel.Stamina), csvModel.PotentialGrade, csvModel.Personality,
			csvModel.RecruitingBias, csvModel.WorkEthic, csvModel.AcademicBias, redshirtStatus,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportNBARosterToCSV(TeamID string, w http.ResponseWriter) {
	// Get Team Data
	team := GetNBATeamByTeamID(TeamID)
	w.Header().Set("Content-Disposition", "attachment;filename="+team.Team+".csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetAllNBAPlayersByTeamID(TeamID)

	HeaderRow := []string{
		"Team", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Age", "Stars",
		"State", "Country", "Height", "Overall", "Finishing",
		"Shooting2", "Shooting3", "FreeThrow",
		"Ballwork", "Rebounding", "Interior Defense", "Perimeter Defense",
		"Playtime Expectations", "Stamina", "Potential",
		"Personality", "Free Agency Bias", "Work Ethic", "NBA Status",
		"Year 1", "Y1 Opt", "Year 2", "Y2 Opt", "Year 3", "Y3 Opt", "Year 4", "Y4 Opt", "Year 5", "Y5 Opt",
		"Contract Length", "Contract Type",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, csvModel := range players {
		nbaStatus := "Active"
		if csvModel.IsGLeague {
			nbaStatus = "G-League"
		} else if csvModel.IsTwoWay {
			nbaStatus = "Two-Way"
		} else if csvModel.IsOnTradeBlock {
			nbaStatus = "On Trade Block"
		} else if csvModel.IsInternational {
			nbaStatus = "International"
		}
		playerRow := []string{
			team.Team, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, strconv.Itoa(csvModel.Year), strconv.Itoa(csvModel.Age), strconv.Itoa(csvModel.Stars),
			csvModel.State, csvModel.Country, csvModel.Height, strconv.Itoa(csvModel.Overall), strconv.Itoa(csvModel.Finishing),
			strconv.Itoa(csvModel.Shooting2), strconv.Itoa(csvModel.Shooting3), strconv.Itoa(csvModel.FreeThrow),
			strconv.Itoa(csvModel.Ballwork), strconv.Itoa(csvModel.Rebounding), strconv.Itoa(csvModel.InteriorDefense), strconv.Itoa(csvModel.PerimeterDefense),
			strconv.Itoa(csvModel.PlaytimeExpectations), strconv.Itoa(csvModel.Stamina), csvModel.PotentialGrade, csvModel.Personality,
			csvModel.FreeAgency, csvModel.WorkEthic, nbaStatus,
			strconv.Itoa(int(csvModel.Contract.Year1Total)), strconv.FormatBool(csvModel.Contract.Year1Opt),
			strconv.Itoa(int(csvModel.Contract.Year2Total)), strconv.FormatBool(csvModel.Contract.Year2Opt),
			strconv.Itoa(int(csvModel.Contract.Year3Total)), strconv.FormatBool(csvModel.Contract.Year3Opt),
			strconv.Itoa(int(csvModel.Contract.Year4Total)), strconv.FormatBool(csvModel.Contract.Year4Opt),
			strconv.Itoa(int(csvModel.Contract.Year5Total)), strconv.FormatBool(csvModel.Contract.Year5Opt),
			strconv.Itoa(int(csvModel.Contract.YearsRemaining)), csvModel.Contract.ContractType,
		}

		err = writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportMatchResults(w http.ResponseWriter, seasonID, weekID, nbaWeekID, matchType string) {
	fileName := "wahoos_secret_results_list.csv"
	w.Header().Set("Content-Disposition", "attachment;"+fileName)
	w.Header().Set("Transfer-Encoding", "chunked")
	writer := csv.NewWriter(w)

	// Get All needed data
	matchChn := make(chan []structs.Match)
	nbaMatchChn := make(chan []structs.NBAMatch)
	collegeTeamChn := make(chan []structs.CollegeTeamResponse)
	nbaTeamChn := make(chan []structs.NBATeamResponse)
	collegeTeamMap := make(map[uint]structs.CollegeTeamResponse)
	nbaTeamMap := make(map[uint]structs.NBATeamResponse)

	go func() {
		matches := GetMatchesByWeekIdAndMatchType(weekID, seasonID, matchType)
		matchChn <- matches
	}()

	go func() {
		nbamatches := GetNBAMatchesByWeekIdAndMatchType(nbaWeekID, seasonID, matchType)
		nbaMatchChn <- nbamatches
	}()

	go func() {
		ct := GetAllActiveCollegeTeamsWithSeasonStats(seasonID, weekID, matchType, "WEEK")
		collegeTeamChn <- ct
	}()

	go func() {
		nt := GetAllActiveNBATeamsWithSeasonStats(seasonID, nbaWeekID, matchType, "WEEK")
		nbaTeamChn <- nt
	}()

	collegeMatches := <-matchChn
	close(matchChn)
	nbaMatches := <-nbaMatchChn
	close(nbaMatchChn)
	collegeTeamSeasonStats := <-collegeTeamChn
	close(collegeTeamChn)
	nbaTeamSeasonStats := <-nbaTeamChn
	close(nbaTeamChn)

	for _, t := range collegeTeamSeasonStats {
		collegeTeamMap[t.ID] = t
	}

	for _, t := range nbaTeamSeasonStats {
		nbaTeamMap[t.ID] = t
	}

	HeaderRow := []string{
		"League", "Week", "Match", "Home Team", "Home Coach", "Home Rank", "Home Score",
		"Home Possessions", "Away Team", "Away Coach", "Away Rank", "Away Score",
		"Away Possessions", "Neutral Site", "Conference", "Division",
		"Match Name", "Arena", "City", "State/Country",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, m := range collegeMatches {
		if !m.GameComplete {
			continue
		}
		homeTeam := collegeTeamMap[m.HomeTeamID]
		awayTeam := collegeTeamMap[m.AwayTeamID]
		neutralStr := "N"
		if m.IsNeutralSite {
			neutralStr = "Y"
		}
		confStr := "N"
		if m.IsConference {
			confStr = "Y"
		}
		divStr := "N"

		row := []string{
			"CBB", strconv.Itoa(int(m.Week)), m.MatchOfWeek, m.HomeTeam, m.AwayTeamCoach,
			strconv.Itoa(int(m.HomeTeamRank)), strconv.Itoa(int(m.HomeTeamScore)), strconv.Itoa(int(homeTeam.Stats.Possessions)),
			m.AwayTeam, m.AwayTeamCoach, strconv.Itoa(int(m.AwayTeamRank)), strconv.Itoa(int(m.AwayTeamScore)), strconv.Itoa(int(awayTeam.Stats.Possessions)),
			neutralStr, confStr, divStr, m.MatchName, m.Arena, m.City, m.State,
		}
		err = writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
	for _, m := range nbaMatches {
		if !m.GameComplete {
			continue
		}
		homeTeam := nbaTeamMap[m.HomeTeamID]
		awayTeam := nbaTeamMap[m.AwayTeamID]
		neutralStr := "N"
		if m.IsNeutralSite {
			neutralStr = "Y"
		}
		confStr := "N"
		if m.IsConference {
			confStr = "Y"
		}
		divStr := "N"
		if m.IsDivisional {
			divStr = "Y"
		}

		row := []string{
			homeTeam.League, strconv.Itoa(int(m.Week)), m.MatchOfWeek, m.HomeTeam, m.AwayTeamCoach,
			"N/A", strconv.Itoa(int(m.HomeTeamScore)), strconv.Itoa(int(homeTeam.Stats.Possessions)),
			m.AwayTeam, m.AwayTeamCoach, "N/A", strconv.Itoa(int(m.AwayTeamScore)), strconv.Itoa(int(awayTeam.Stats.Possessions)),
			neutralStr, confStr, divStr, m.MatchName, m.Arena, m.City, m.State,
		}
		err = writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func ExportStatsMain(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType, playerView string) {
	fileName := "babas_secret_" + league + "_" + viewType + "_" + playerView + "_stats_list.csv"
	w.Header().Set("Content-Disposition", "attachment;"+fileName)
	w.Header().Set("Transfer-Encoding", "chunked")

	if league == "cbb" {
		if playerView == "TEAM" && viewType == "SEASON" {
			ExportCollegeTeamSeasonStats(w, league, seasonID, weekID, matchType, viewType)
		} else if playerView == "TEAM" && viewType == "WEEK" {
			ExportCollegeTeamWeeklyStats(w, league, seasonID, weekID, matchType, viewType)
		} else if playerView == "PLAYER" && viewType == "SEASON" {
			ExportCollegePlayerSeasonStats(w, league, seasonID, weekID, matchType, viewType)
		} else if playerView == "PLAYER" && viewType == "WEEK" {
			ExportCollegePlayerWeeklyStats(w, league, seasonID, weekID, matchType, viewType)
		}
	} else {
		if playerView == "TEAM" && viewType == "SEASON" {
			ExportNBATeamSeasonStats(w, league, seasonID, weekID, matchType, viewType)
		} else if playerView == "TEAM" && viewType == "WEEK" {
			ExportNBATeamWeeklyStats(w, league, seasonID, weekID, matchType, viewType)
		} else if playerView == "PLAYER" && viewType == "SEASON" {
			ExportNBAPlayerSeasonStats(w, league, seasonID, weekID, matchType, viewType)
		} else if playerView == "PLAYER" && viewType == "WEEK" {
			ExportNBAPlayerWeeklyStats(w, league, seasonID, weekID, matchType, viewType)
		}
	}
}

// Export CBB Team Seasonal Data
func ExportCollegeTeamSeasonStats(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType string) {
	writer := csv.NewWriter(w)

	teamSeasonStats := GetAllActiveCollegeTeamsWithSeasonStats(seasonID, weekID, matchType, viewType)

	HeaderRow := []string{
		"Team", "Conference", "Games Played", "Points", "FieldGoalsMade",
		"FieldGoalAttempts", "FG Percent", "3PtsMade", "3PtAttempts", "3Pt Percent",
		"Free Throws Made", "FT Attempts", "FT Percent", "Offensive Rebounds",
		"Defensive Rebounds", "Assists", "Steals", "Blocks", "Turnovers",
		"Points Allowed", "FGM Allowed", "FGA Allowed", "FG%Allowed",
		"3PMAllowed", "3PAAllowed", "3P%Allowed", "FTMAllowed", "FTAAllowed", "FT%Allowed",
		"OReboundsAllowed", "DReboundsAllowed", "ReboundsAllowed", "AssistsAllowed", "StealsAllowed",
		"BlocksAllowed", "TurnoversAllowed",
		"PointDifferential", "FGMDiff", "FGADiff", "FG%Diff", "3PMDiff", "3PADiff", "3P%Diff",
		"FTMDiff", "FTADiff", "FT%Diff", "ORDiff", "DRDiff", "ReboundsDiff", "AssistsDiff",
		"StealsDiff", "BlocksDiff", "TurnoverDiff",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, team := range teamSeasonStats {

		stats := team.SeasonStats

		teamRow := getTeamSeasonRow(team.Team, team.Conference, stats)

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

// Export CBB Team Individual Data
func ExportCollegeTeamWeeklyStats(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType string) {
	writer := csv.NewWriter(w)

	teamSeasonStats := GetAllActiveCollegeTeamsWithSeasonStats(seasonID, weekID, matchType, viewType)

	HeaderRow := []string{
		"Team", "Conference", "Week", "Match",
		"Points", "1stHalf", "2ndHalf", "Overtime", "FieldGoalsMade",
		"FieldGoalAttempts", "FG Percent", "3PtsMade", "3PtAttempts", "3Pt Percent",
		"Free Throws Made", "FT Attempts", "FT Percent", "Offensive Rebounds",
		"Defensive Rebounds", "Assists", "Steals", "Blocks", "Turnovers", "Fouls",
		"Points Allowed", "FGM Allowed", "FGA Allowed", "FG%Allowed",
		"3PMAllowed", "3PAAllowed", "3P%Allowed", "FTMAllowed", "FTAAllowed", "FT%Allowed",
		"OReboundsAllowed", "DReboundsAllowed", "ReboundsAllowed", "AssistsAllowed", "StealsAllowed",
		"BlocksAllowed", "TurnoversAllowed",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, team := range teamSeasonStats {

		stats := team.Stats

		teamRow := getCollegeTeamWeeklyRow(team.Team, team.Conference, matchType, stats)

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

// Export CBB Player Seasonal Data
func ExportCollegePlayerSeasonStats(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType string) {
	writer := csv.NewWriter(w)

	players := GetAllCollegePlayersWithSeasonStats(seasonID, weekID, matchType, viewType)

	HeaderRow := []string{
		"First Name", "Last Name", "Position", "Age", "Year", "Team",
		"Games Played", "Minutes", "Possessions", "Points", "FieldGoalsMade",
		"FieldGoalAttempts", "FG Percent", "3PtsMade", "3PtAttempts", "3Pt Percent",
		"Free Throws Made", "FT Attempts", "FT Percent", "Offensive Rebounds",
		"Defensive Rebounds", "Assists", "Steals", "Blocks", "Turnovers",
		"MPG", "PossessionsPerGame", "PointsPerGame", "FGMPerGame", "FGA PerGame",
		"3PMPerGame", "3PAPerGame", "FTMPerGame", "FTAPerGame",
		"OReboundsPerGame", "DReboundsPerGame", "ReboundsPerGame", "AssistsPerGame", "StealsPerGame",
		"BlocksPerGame", "TurnoversPerGame",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {

		stats := player.SeasonStats

		teamRow := getCollegePlayerSeasonRow(player, stats)

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

// Export CBB player Individual Data
func ExportCollegePlayerWeeklyStats(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType string) {
	writer := csv.NewWriter(w)

	players := GetAllCollegePlayersWithSeasonStats(seasonID, weekID, matchType, viewType)

	HeaderRow := []string{
		"First Name", "Last Name", "Position", "Age", "Year", "Team",
		"Week", "Match", "Minutes", "Possessions", "Points", "FieldGoalsMade",
		"FieldGoalAttempts", "FG Percent", "3PtsMade", "3PtAttempts", "3Pt Percent",
		"Free Throws Made", "FT Attempts", "FT Percent", "Offensive Rebounds",
		"Defensive Rebounds", "Assists", "Steals", "Blocks", "Turnovers",
		"Fouls",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {

		stats := player.Stats

		teamRow := getCollegePlayerWeeklyRow(matchType, player, stats)

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

// Export NBA Team Seasonal Data
func ExportNBATeamSeasonStats(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType string) {
	writer := csv.NewWriter(w)

	teamSeasonStats := GetAllActiveNBATeamsWithSeasonStats(seasonID, weekID, matchType, viewType)

	HeaderRow := []string{
		"Team", "Conference", "Games Played", "Points", "FieldGoalsMade",
		"FieldGoalAttempts", "FG Percent", "3PtsMade", "3PtAttempts", "3Pt Percent",
		"Free Throws Made", "FT Attempts", "FT Percent", "Offensive Rebounds",
		"Defensive Rebounds", "Assists", "Steals", "Blocks", "Turnovers",
		"Points Allowed", "FGM Allowed", "FGA Allowed", "FG%Allowed",
		"3PMAllowed", "3PAAllowed", "3P%Allowed", "FTMAllowed", "FTAAllowed", "FT%Allowed",
		"OReboundsAllowed", "DReboundsAllowed", "ReboundsAllowed", "AssistsAllowed", "StealsAllowed",
		"BlocksAllowed", "TurnoversAllowed",
		"PointDifferential", "FGMDiff", "FGADiff", "FG%Diff", "3PMDiff", "3PADiff", "3P%Diff",
		"FTMDiff", "FTADiff", "FT%Diff", "ORDiff", "DRDiff", "ReboundsDiff", "AssistsDiff",
		"StealsDiff", "BlocksDiff", "TurnoverDiff",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, team := range teamSeasonStats {

		stats := team.SeasonStats

		teamRow := getTeamSeasonRow(team.Team, team.Conference, stats)

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

// Export NBA Team Individual Data
func ExportNBATeamWeeklyStats(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType string) {
	writer := csv.NewWriter(w)

	teamSeasonStats := GetAllActiveNBATeamsWithSeasonStats(seasonID, weekID, matchType, viewType)

	HeaderRow := []string{
		"Team", "Conference", "Week", "Match",
		"Points", "1stQuarter", "2ndQuarter", "3rdQuarter", "4thQuarter",
		"Overtime", "FieldGoalsMade",
		"FieldGoalAttempts", "FG Percent", "3PtsMade", "3PtAttempts", "3Pt Percent",
		"Free Throws Made", "FT Attempts", "FT Percent", "Offensive Rebounds",
		"Defensive Rebounds", "Assists", "Steals", "Blocks", "Turnovers", "Fouls",
		"Points Allowed", "FGM Allowed", "FGA Allowed", "FG%Allowed",
		"3PMAllowed", "3PAAllowed", "3P%Allowed", "FTMAllowed", "FTAAllowed", "FT%Allowed",
		"OReboundsAllowed", "DReboundsAllowed", "ReboundsAllowed", "AssistsAllowed", "StealsAllowed",
		"BlocksAllowed", "TurnoversAllowed",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, team := range teamSeasonStats {

		stats := team.Stats

		teamRow := getNBATeamWeeklyRow(team.Team, team.Conference, matchType, stats)

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

// Export NBA Player Seasonal Data
func ExportNBAPlayerSeasonStats(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType string) {
	writer := csv.NewWriter(w)

	players := GetAllNBAPlayersWithSeasonStats(seasonID, weekID, matchType, viewType)

	HeaderRow := []string{
		"First Name", "Last Name", "Position", "Age", "Year", "Team",
		"Games Played", "Minutes", "Possessions", "Points", "FieldGoalsMade",
		"FieldGoalAttempts", "FG Percent", "3PtsMade", "3PtAttempts", "3Pt Percent",
		"Free Throws Made", "FT Attempts", "FT Percent", "Offensive Rebounds",
		"Defensive Rebounds", "Assists", "Steals", "Blocks", "Turnovers",
		"MPG", "PossessionsPerGame", "PointsPerGame", "FGMPerGame", "FGA PerGame",
		"3PMPerGame", "3PAPerGame", "FTMPerGame", "FTAPerGame",
		"OReboundsPerGame", "DReboundsPerGame", "ReboundsPerGame", "AssistsPerGame", "StealsPerGame",
		"BlocksPerGame", "TurnoversPerGame",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {

		stats := player.SeasonStats

		teamRow := getNBAPlayerSeasonRow(player, stats)

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

// Export NBA player Individual Data
func ExportNBAPlayerWeeklyStats(w http.ResponseWriter, league, seasonID, weekID, matchType, viewType string) {
	writer := csv.NewWriter(w)

	players := GetAllNBAPlayersWithSeasonStats(seasonID, weekID, matchType, viewType)

	HeaderRow := []string{
		"First Name", "Last Name", "Position", "Age", "Year", "Team",
		"Week", "Match", "Minutes", "Possessions", "Points", "FieldGoalsMade",
		"FieldGoalAttempts", "FG Percent", "3PtsMade", "3PtAttempts", "3Pt Percent",
		"Free Throws Made", "FT Attempts", "FT Percent", "Offensive Rebounds",
		"Defensive Rebounds", "Assists", "Steals", "Blocks", "Turnovers",
		"Fouls",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {

		stats := player.Stats

		teamRow := getNBAPlayerWeeklyRow(matchType, player, stats)

		err = writer.Write(teamRow)
		if err != nil {
			log.Fatal("Cannot write croot row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

// get Team Row
func getTeamSeasonRow(team, conference string, stats structs.TeamSeasonStatsResponse) []string {
	return []string{
		team,
		conference,
		strconv.Itoa(int(stats.GamesPlayed)),
		strconv.Itoa(stats.Points),
		strconv.Itoa(stats.FGM),
		strconv.Itoa(stats.FGA),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMade),
		strconv.Itoa(stats.ThreePointAttempts),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(stats.FTM),
		strconv.Itoa(stats.FTA),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(stats.OffRebounds),
		strconv.Itoa(stats.DefRebounds),
		strconv.Itoa(stats.Assists),
		strconv.Itoa(stats.Steals),
		strconv.Itoa(stats.Blocks),
		strconv.Itoa(stats.TotalTurnovers),
		strconv.Itoa(stats.PointsAgainst),
		strconv.Itoa(stats.FGMAgainst),
		strconv.Itoa(stats.FGAAgainst),
		strconv.FormatFloat(stats.FGPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMadeAgainst),
		strconv.Itoa(stats.ThreePointAttemptsAgainst),
		strconv.FormatFloat(stats.ThreePointPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.FTMAgainst),
		strconv.Itoa(stats.FTAAgainst),
		strconv.FormatFloat(stats.FTPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.OffReboundsAllowed),
		strconv.Itoa(stats.DefReboundsAllowed),
		strconv.Itoa(stats.ReboundsAllowed),
		strconv.Itoa(stats.AssistsAllowed),
		strconv.Itoa(stats.StealsAllowed),
		strconv.Itoa(stats.BlocksAllowed),
		strconv.Itoa(stats.TurnoversAllowed),
		strconv.FormatFloat(stats.PointsDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.FGMDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.FGADiff, 'f', 2, 64),
		strconv.FormatFloat(stats.FGPercentDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.TPMDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.TPADiff, 'f', 2, 64),
		strconv.FormatFloat(stats.TPPercentDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.FTMDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.FTADiff, 'f', 2, 64),
		strconv.FormatFloat(stats.FTPercentDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.OReboundsDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.DReboundsDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.ReboundsDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.AssistsDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.StealsDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.BlocksDiff, 'f', 2, 64),
		strconv.FormatFloat(stats.TODiff, 'f', 2, 64),
	}
}

// Get College Team Weekly Row
func getCollegeTeamWeeklyRow(team, conference, matchType string, stats structs.TeamStats) []string {
	return []string{
		team,
		conference,
		strconv.Itoa(int(stats.Week)),
		matchType,
		strconv.Itoa(stats.Points),
		strconv.Itoa(stats.FirstHalfScore),
		strconv.Itoa(stats.SecondHalfScore),
		strconv.Itoa(stats.OvertimeScore),
		strconv.Itoa(stats.FGM),
		strconv.Itoa(stats.FGA),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMade),
		strconv.Itoa(stats.ThreePointAttempts),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(stats.FTM),
		strconv.Itoa(stats.FTA),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(stats.OffRebounds),
		strconv.Itoa(stats.DefRebounds),
		strconv.Itoa(stats.Assists),
		strconv.Itoa(stats.Steals),
		strconv.Itoa(stats.Blocks),
		strconv.Itoa(stats.TotalTurnovers),
		strconv.Itoa(stats.Fouls),
		strconv.Itoa(stats.PointsAgainst),
		strconv.Itoa(stats.FGMAgainst),
		strconv.Itoa(stats.FGAAgainst),
		strconv.FormatFloat(stats.FGPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMadeAgainst),
		strconv.Itoa(stats.ThreePointAttemptsAgainst),
		strconv.FormatFloat(stats.ThreePointPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.FTMAgainst),
		strconv.Itoa(stats.FTAAgainst),
		strconv.FormatFloat(stats.FTPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.OffReboundsAllowed),
		strconv.Itoa(stats.DefReboundsAllowed),
		strconv.Itoa(stats.ReboundsAllowed),
		strconv.Itoa(stats.AssistsAllowed),
		strconv.Itoa(stats.StealsAllowed),
		strconv.Itoa(stats.BlocksAllowed),
		strconv.Itoa(stats.TurnoversAllowed),
	}
}

// Get NBA Team Weekly Row
func getNBATeamWeeklyRow(team, conference, matchType string, stats structs.NBATeamStats) []string {
	return []string{
		team,
		conference,
		strconv.Itoa(int(stats.Week)),
		matchType,
		strconv.Itoa(stats.Points),
		strconv.Itoa(stats.FirstHalfScore),
		strconv.Itoa(stats.SecondQuarterScore),
		strconv.Itoa(stats.SecondHalfScore),
		strconv.Itoa(stats.FourthQuarterScore),
		strconv.Itoa(stats.OvertimeScore),
		strconv.Itoa(stats.FGM),
		strconv.Itoa(stats.FGA),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMade),
		strconv.Itoa(stats.ThreePointAttempts),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(stats.FTM),
		strconv.Itoa(stats.FTA),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(stats.OffRebounds),
		strconv.Itoa(stats.DefRebounds),
		strconv.Itoa(stats.Assists),
		strconv.Itoa(stats.Steals),
		strconv.Itoa(stats.Blocks),
		strconv.Itoa(stats.TotalTurnovers),
		strconv.Itoa(stats.Fouls),
		strconv.Itoa(stats.PointsAgainst),
		strconv.Itoa(stats.FGMAgainst),
		strconv.Itoa(stats.FGAAgainst),
		strconv.FormatFloat(stats.FGPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMadeAgainst),
		strconv.Itoa(stats.ThreePointAttemptsAgainst),
		strconv.FormatFloat(stats.ThreePointPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.FTMAgainst),
		strconv.Itoa(stats.FTAAgainst),
		strconv.FormatFloat(stats.FTPercentAgainst, 'f', 2, 64),
		strconv.Itoa(stats.OffReboundsAllowed),
		strconv.Itoa(stats.DefReboundsAllowed),
		strconv.Itoa(stats.ReboundsAllowed),
		strconv.Itoa(stats.AssistsAllowed),
		strconv.Itoa(stats.StealsAllowed),
		strconv.Itoa(stats.BlocksAllowed),
		strconv.Itoa(stats.TurnoversAllowed),
	}
}

// Get CollegePlayerSeasonRow
func getCollegePlayerSeasonRow(player structs.CollegePlayerResponse, stats structs.CollegePlayerSeasonStats) []string {
	return []string{
		player.FirstName,
		player.LastName,
		player.Position,
		strconv.Itoa(player.Age),
		strconv.Itoa(player.Year),
		player.TeamAbbr,
		strconv.Itoa(int(stats.GamesPlayed)),
		strconv.Itoa(stats.Minutes),
		strconv.Itoa(stats.Possessions),
		strconv.Itoa(stats.Points),
		strconv.Itoa(stats.FGM),
		strconv.Itoa(stats.FGA),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMade),
		strconv.Itoa(stats.ThreePointAttempts),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(stats.FTM),
		strconv.Itoa(stats.FTA),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(stats.OffRebounds),
		strconv.Itoa(stats.DefRebounds),
		strconv.Itoa(stats.Assists),
		strconv.Itoa(stats.Steals),
		strconv.Itoa(stats.Blocks),
		strconv.Itoa(stats.Turnovers),
		strconv.FormatFloat(stats.MinutesPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.PossessionsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.PPG, 'f', 2, 64),
		strconv.FormatFloat(stats.FGMPG, 'f', 2, 64),
		strconv.FormatFloat(stats.FGAPG, 'f', 2, 64),
		strconv.FormatFloat(stats.ThreePointsMadePerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.ThreePointAttemptsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.FTMPG, 'f', 2, 64),
		strconv.FormatFloat(stats.FTAPG, 'f', 2, 64),
		strconv.FormatFloat(stats.OffReboundsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.DefReboundsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.ReboundsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.AssistsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.BlocksPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.TurnoversPerGame, 'f', 2, 64),
	}
}

// Get CollegePlayerWeeklyRow
func getCollegePlayerWeeklyRow(matchType string, player structs.CollegePlayerResponse, stats structs.CollegePlayerStats) []string {
	return []string{
		player.FirstName,
		player.LastName,
		player.Position,
		strconv.Itoa(player.Age),
		strconv.Itoa(player.Year),
		player.TeamAbbr,
		strconv.Itoa(int(stats.Week)),
		matchType,
		strconv.Itoa(stats.Minutes),
		strconv.Itoa(stats.Possessions),
		strconv.Itoa(stats.Points),
		strconv.Itoa(stats.FGM),
		strconv.Itoa(stats.FGA),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMade),
		strconv.Itoa(stats.ThreePointAttempts),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(stats.FTM),
		strconv.Itoa(stats.FTA),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(stats.OffRebounds),
		strconv.Itoa(stats.DefRebounds),
		strconv.Itoa(stats.Assists),
		strconv.Itoa(stats.Steals),
		strconv.Itoa(stats.Blocks),
		strconv.Itoa(stats.Turnovers),
		strconv.Itoa(stats.Fouls),
	}
}

// Get NBAPlayerSeasonRow
func getNBAPlayerSeasonRow(player structs.NBAPlayerResponse, stats structs.NBAPlayerSeasonStats) []string {
	return []string{
		player.FirstName,
		player.LastName,
		player.Position,
		strconv.Itoa(player.Age),
		strconv.Itoa(player.Year),
		player.TeamAbbr,
		strconv.Itoa(int(stats.GamesPlayed)),
		strconv.Itoa(stats.Minutes),
		strconv.Itoa(stats.Possessions),
		strconv.Itoa(stats.Points),
		strconv.Itoa(stats.FGM),
		strconv.Itoa(stats.FGA),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMade),
		strconv.Itoa(stats.ThreePointAttempts),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(stats.FTM),
		strconv.Itoa(stats.FTA),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(stats.OffRebounds),
		strconv.Itoa(stats.DefRebounds),
		strconv.Itoa(stats.Assists),
		strconv.Itoa(stats.Steals),
		strconv.Itoa(stats.Blocks),
		strconv.Itoa(stats.Turnovers),
		strconv.FormatFloat(stats.MinutesPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.PossessionsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.PPG, 'f', 2, 64),
		strconv.FormatFloat(stats.FGMPG, 'f', 2, 64),
		strconv.FormatFloat(stats.FGAPG, 'f', 2, 64),
		strconv.FormatFloat(stats.ThreePointsMadePerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.ThreePointAttemptsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.FTMPG, 'f', 2, 64),
		strconv.FormatFloat(stats.FTAPG, 'f', 2, 64),
		strconv.FormatFloat(stats.OffReboundsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.DefReboundsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.ReboundsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.AssistsPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.BlocksPerGame, 'f', 2, 64),
		strconv.FormatFloat(stats.TurnoversPerGame, 'f', 2, 64),
	}
}

// Get CollegePlayerWeeklyRow
func getNBAPlayerWeeklyRow(matchType string, player structs.NBAPlayerResponse, stats structs.NBAPlayerStats) []string {
	return []string{
		player.FirstName,
		player.LastName,
		player.Position,
		strconv.Itoa(player.Age),
		strconv.Itoa(player.Year),
		player.TeamAbbr,
		strconv.Itoa(int(stats.Week)),
		matchType,
		strconv.Itoa(stats.Minutes),
		strconv.Itoa(stats.Possessions),
		strconv.Itoa(stats.Points),
		strconv.Itoa(stats.FGM),
		strconv.Itoa(stats.FGA),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(stats.ThreePointsMade),
		strconv.Itoa(stats.ThreePointAttempts),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(stats.FTM),
		strconv.Itoa(stats.FTA),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(stats.OffRebounds),
		strconv.Itoa(stats.DefRebounds),
		strconv.Itoa(stats.Assists),
		strconv.Itoa(stats.Steals),
		strconv.Itoa(stats.Blocks),
		strconv.Itoa(stats.Turnovers),
		strconv.Itoa(stats.Fouls),
	}
}
