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
		"Stars", "College", "High School", "City", "State", "Country", "Height",
		"Weight", "Program Pref.", "Prof Dev Pref.", "Traditions Pref.", "Facilities Pref.", "Atmosphere Pref.",
		"Academics Pref.", "Campus Life Pref.", "Conference Pref.", "Coach Pref.", "Season Momentum Pref.", "Leading Teams",
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
			croot.FirstName, croot.LastName, croot.Position, strconv.Itoa(int(croot.Stars)),
			croot.College, croot.HighSchool, croot.City, croot.State, croot.Country, strconv.Itoa(int(croot.Height)),
			strconv.Itoa(int(croot.Weight)),
			strconv.Itoa(int(croot.ProgramPref)), strconv.Itoa(int(croot.ProfDevPref)), strconv.Itoa(int(croot.TraditionsPref)), strconv.Itoa(int(croot.FacilitiesPref)), strconv.Itoa(int(croot.AtmospherePref)),
			strconv.Itoa(int(croot.AcademicsPref)), strconv.Itoa(int(croot.CampusLifePref)), strconv.Itoa(int(croot.ConferencePref)), strconv.Itoa(int(croot.CoachPref)), strconv.Itoa(int(croot.SeasonMomentumPref)),
			strings.Join(leadingAbbr, ", "),
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
		"ID", "College", "First Name", "Last Name", "Position", "Year", "Is_Redshirt", "Age",
		"Stars", "High School", "City", "State", "Country", "Height", "Weight",
		"Overall", "Inside Shooting", "MidRange Shooting", "Three Point Shooting", "Free Throwing", "Agility",
		"Ballwork", "Stealing", "Blocking", "Rebounding", "InteriorDefense", "PerimeterDefense", "Stamina", "Potential Grade",
		"Personality", "RecruitingBias", "Work Ethic", "Previous Team",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {

		shooting2Grade := util.GetAttributeGrade(player.MidRangeShooting, int(player.Year))
		shooting3Grade := util.GetAttributeGrade(player.ThreePointShooting, int(player.Year))
		freeThrowGrade := util.GetAttributeGrade(player.FreeThrow, int(player.Year))
		finishingGrade := util.GetAttributeGrade(player.InsideShooting, int(player.Year))
		reboundingGrade := util.GetAttributeGrade(player.Rebounding, int(player.Year))
		ballworkGrade := util.GetAttributeGrade(player.Ballwork, int(player.Year))
		interiorDefenseGrade := util.GetAttributeGrade(player.InteriorDefense, int(player.Year))
		perimeterDefenseGrade := util.GetAttributeGrade(player.PerimeterDefense, int(player.Year))
		agilityGrade := util.GetAttributeGrade(player.Agility, int(player.Year))
		stealingGrade := util.GetAttributeGrade(player.Stealing, int(player.Year))
		blockingGrade := util.GetAttributeGrade(player.Blocking, int(player.Year))
		potentialGrade := util.GetPotentialGrade(player.Potential)
		overallGrade := util.GetAttributeGrade(player.Overall, int(player.Year))
		sta := strconv.Itoa(int(player.Stamina))

		playerRow := []string{
			strconv.Itoa(int(player.ID)), player.Team, player.FirstName, player.LastName, player.Position, strconv.Itoa(int(player.Year)), strconv.FormatBool(player.IsRedshirt), strconv.Itoa(int(player.Age)),
			strconv.Itoa(int(player.Stars)), player.HighSchool, player.City, player.State, player.Country, strconv.Itoa(int(player.Height)), strconv.Itoa(int(player.Weight)),
			overallGrade, finishingGrade, shooting2Grade, shooting3Grade, freeThrowGrade, agilityGrade,
			ballworkGrade, stealingGrade, blockingGrade, reboundingGrade, interiorDefenseGrade, perimeterDefenseGrade, sta, potentialGrade,
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

func ExportTransferPortalToCSV(w http.ResponseWriter) {
	// Get Team Data
	w.Header().Set("Content-Disposition", "attachment;filename=Official_CBB_Portal_List.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetTransferPortalPlayers()

	HeaderRow := []string{
		"College", "First Name", "Last Name", "Position", "Year", "Is_Redshirt", "Age",
		"Stars", "State", "Country", "Height",
		"Overall", "Shooting 2s", "Shooting 3s", "Free Throwing", "Finishing", "Agility",
		"Ballwork", "Rebounding", "Stealing", "Blocking", "InteriorDefense", "PerimeterDefense", "Stamina", "Potential Grade",
		"Personality", "RecruitingBias", "Work Ethic", "Previous Team",
	}

	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, player := range players {
		shooting2Grade := util.GetAttributeGrade(player.MidRangeShooting, int(player.Year))
		shooting3Grade := util.GetAttributeGrade(player.ThreePointShooting, int(player.Year))
		freeThrowGrade := util.GetAttributeGrade(player.FreeThrow, int(player.Year))
		finishingGrade := util.GetAttributeGrade(player.InsideShooting, int(player.Year))
		reboundingGrade := util.GetAttributeGrade(player.Rebounding, int(player.Year))
		ballworkGrade := util.GetAttributeGrade(player.Ballwork, int(player.Year))
		interiorDefenseGrade := util.GetAttributeGrade(player.InteriorDefense, int(player.Year))
		perimeterDefenseGrade := util.GetAttributeGrade(player.PerimeterDefense, int(player.Year))
		agilityGrade := util.GetAttributeGrade(player.Agility, int(player.Year))
		stealingGrade := util.GetAttributeGrade(player.Stealing, int(player.Year))
		blockingGrade := util.GetAttributeGrade(player.Blocking, int(player.Year))
		potentialGrade := util.GetPotentialGrade(player.Potential)
		overallGrade := util.GetAttributeGrade(player.Overall, int(player.Year))
		sta := strconv.Itoa(int(player.Stamina))

		playerRow := []string{
			player.Team, player.FirstName, player.LastName, player.Position, strconv.Itoa(int(player.Year)), strconv.FormatBool(player.IsRedshirt), strconv.Itoa(int(player.Age)),
			strconv.Itoa(int(player.Stars)), player.State, player.Country, strconv.Itoa(int(player.Height)),
			overallGrade, shooting2Grade, shooting3Grade, freeThrowGrade, finishingGrade, agilityGrade,
			ballworkGrade, reboundingGrade, stealingGrade, blockingGrade, interiorDefenseGrade, perimeterDefenseGrade, sta, potentialGrade,
			player.Personality, player.RecruitingBias, player.WorkEthic, player.PreviousTeam,
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
			o := ((float64(player.MidRangeShooting) + float64(player.ThreePointShooting) + float64(player.InsideShooting) + float64(player.FreeThrow)) / 4)
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

		teamRow := []string{t.Abbr, t.Conference, util.ConvertFloatToString(offenseRank),
			util.ConvertFloatToString(defenseRank),
			util.ConvertFloatToString(overall),
			util.ConvertFloatToString(starPower)}

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
		shooting2Grade := util.GetAttributeGrade(player.MidRangeShooting, int(player.Year))
		shooting3Grade := util.GetAttributeGrade(player.ThreePointShooting, int(player.Year))
		freeThrowGrade := util.GetAttributeGrade(player.FreeThrow, int(player.Year))
		finishingGrade := util.GetAttributeGrade(player.InsideShooting, int(player.Year))
		reboundingGrade := util.GetAttributeGrade(player.Rebounding, int(player.Year))
		ballworkGrade := util.GetAttributeGrade(player.Ballwork, int(player.Year))
		interiorDefenseGrade := util.GetAttributeGrade(player.InteriorDefense, int(player.Year))
		perimeterDefenseGrade := util.GetAttributeGrade(player.PerimeterDefense, int(player.Year))
		agilityGrade := util.GetAttributeGrade(player.Agility, int(player.Year))
		stealingGrade := util.GetAttributeGrade(player.Stealing, int(player.Year))
		blockingGrade := util.GetAttributeGrade(player.Blocking, int(player.Year))
		potentialGrade := util.GetPotentialGrade(player.Potential)
		overallGrade := util.GetAttributeGrade(player.Overall, int(player.Year))

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
			AgilityGrade:          agilityGrade,
			StealingGrade:         stealingGrade,
			BlockingGrade:         blockingGrade,
			OverallGrade:          overallGrade,
			Stamina:               player.Stamina,
			PlaytimeExpectations:  player.PlaytimeExpectations,
			Potential:             player.Potential,
			Personality:           player.Personality,
			RecruitingBias:        player.RecruitingBias,
			WorkEthic:             player.WorkEthic,
			AcademicBias:          player.AcademicBias,
			PlayerID:              player.PlayerID,
			TeamID:                player.TeamID,
			Team:                  player.Team,
			IsRedshirting:         player.IsRedshirting,
			IsRedshirt:            player.IsRedshirt,
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
			csvModel.Archetype, strconv.Itoa(int(csvModel.Year)), strconv.Itoa(int(csvModel.Age)), strconv.Itoa(int(csvModel.Stars)),
			csvModel.State, csvModel.Country, strconv.Itoa(int(csvModel.Height)), csvModel.OverallGrade, csvModel.FinishingGrade,
			csvModel.Shooting2Grade, csvModel.Shooting3Grade, csvModel.FreeThrowGrade,
			csvModel.BallworkGrade, csvModel.ReboundingGrade, csvModel.InteriorDefenseGrade, csvModel.PerimeterDefenseGrade,
			strconv.Itoa(int(csvModel.PlaytimeExpectations)), strconv.Itoa(int(csvModel.Stamina)), csvModel.PotentialGrade, csvModel.Personality,
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

	WriteNBAPlayersToCSV(writer, team.Team, players)
}

func ExportNBAFreeAgentsToCSV(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition", "attachment;filename=NBA_Free_Agents.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	writer := csv.NewWriter(w)

	// Get Players
	players := GetAllFreeAgents()
	WriteNBAPlayersToCSV(writer, "FA", players)
}

// WriteNBAPlayersToCSV writes a header row and a player row for each NBA player
// to the provided csv.Writer. teamName is used as the "Team" column value for
// every row, allowing callers to pass any subset of NBA players.
func WriteNBAPlayersToCSV(writer *csv.Writer, teamName string, players []structs.NBAPlayer) {
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
			teamName, csvModel.FirstName, csvModel.LastName, csvModel.Position,
			csvModel.Archetype, strconv.Itoa(int(csvModel.Year)), strconv.Itoa(int(csvModel.Age)), strconv.Itoa(int(csvModel.Stars)),
			csvModel.State, csvModel.Country, strconv.Itoa(int(csvModel.Height)), strconv.Itoa(int(csvModel.Overall)), strconv.Itoa(int(csvModel.InsideShooting)),
			strconv.Itoa(int(csvModel.MidRangeShooting)), strconv.Itoa(int(csvModel.ThreePointShooting)), strconv.Itoa(int(csvModel.FreeThrow)),
			strconv.Itoa(int(csvModel.Ballwork)), strconv.Itoa(int(csvModel.Rebounding)), strconv.Itoa(int(csvModel.InteriorDefense)), strconv.Itoa(int(csvModel.PerimeterDefense)),
			strconv.Itoa(int(csvModel.PlaytimeExpectations)), strconv.Itoa(int(csvModel.Stamina)), csvModel.PotentialGrade, csvModel.Personality,
			csvModel.FreeAgency, csvModel.WorkEthic, nbaStatus,
			util.ConvertFloatToString(csvModel.Contract.Year1Total), strconv.FormatBool(csvModel.Contract.Year1Opt),
			util.ConvertFloatToString(csvModel.Contract.Year2Total), strconv.FormatBool(csvModel.Contract.Year2Opt),
			util.ConvertFloatToString(csvModel.Contract.Year3Total), strconv.FormatBool(csvModel.Contract.Year3Opt),
			util.ConvertFloatToString(csvModel.Contract.Year4Total), strconv.FormatBool(csvModel.Contract.Year4Opt),
			util.ConvertFloatToString(csvModel.Contract.Year5Total), strconv.FormatBool(csvModel.Contract.Year5Opt),
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
	ts := GetTimestamp()

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
		if m.Week == uint(ts.CollegeWeek) && ((matchType == "A" && !ts.GamesARan) || (matchType == "B" && !ts.GamesBRan) || (matchType == "C" && !ts.GamesCRan) || (matchType == "D" && !ts.GamesDRan)) {
			m.HideScore()
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

		if m.Week == uint(ts.NBAWeek) && ((matchType == "A" && !ts.GamesARan) || (matchType == "B" && !ts.GamesBRan) || (matchType == "C" && !ts.GamesCRan) || (matchType == "D" && !ts.GamesDRan)) {
			m.HideScore()
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
		strconv.Itoa(int(player.Age)),
		strconv.Itoa(int(player.Year)),
		player.Team,
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
		strconv.Itoa(int(player.Age)),
		strconv.Itoa(int(player.Year)),
		player.Team,
		strconv.Itoa(int(stats.Week)),
		matchType,
		strconv.Itoa(int(stats.Minutes)),
		strconv.Itoa(int(stats.Possessions)),
		strconv.Itoa(int(stats.Points)),
		strconv.Itoa(int(stats.FGM)),
		strconv.Itoa(int(stats.FGA)),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.ThreePointsMade)),
		strconv.Itoa(int(stats.ThreePointAttempts)),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.FTM)),
		strconv.Itoa(int(stats.FTA)),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.OffRebounds)),
		strconv.Itoa(int(stats.DefRebounds)),
		strconv.Itoa(int(stats.Assists)),
		strconv.Itoa(int(stats.Steals)),
		strconv.Itoa(int(stats.Blocks)),
		strconv.Itoa(int(stats.Turnovers)),
		strconv.Itoa(int(stats.Fouls)),
	}
}

// Get NBAPlayerSeasonRow
func getNBAPlayerSeasonRow(player structs.NBAPlayerResponse, stats structs.NBAPlayerSeasonStats) []string {
	return []string{
		player.FirstName,
		player.LastName,
		player.Position,
		strconv.Itoa(int(player.Age)),
		strconv.Itoa(int(player.Year)),
		player.Team,
		strconv.Itoa(int(stats.GamesPlayed)),
		strconv.Itoa(int(stats.Minutes)),
		strconv.Itoa(int(stats.Possessions)),
		strconv.Itoa(int(stats.Points)),
		strconv.Itoa(int(stats.FGM)),
		strconv.Itoa(int(stats.FGA)),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.ThreePointsMade)),
		strconv.Itoa(int(stats.ThreePointAttempts)),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.FTM)),
		strconv.Itoa(int(stats.FTA)),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.OffRebounds)),
		strconv.Itoa(int(stats.DefRebounds)),
		strconv.Itoa(int(stats.Assists)),
		strconv.Itoa(int(stats.Steals)),
		strconv.Itoa(int(stats.Blocks)),
		strconv.Itoa(int(stats.Turnovers)),
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
		strconv.Itoa(int(player.Age)),
		strconv.Itoa(int(player.Year)),
		player.Team,
		strconv.Itoa(int(stats.Week)),
		matchType,
		strconv.Itoa(int(stats.Minutes)),
		strconv.Itoa(int(stats.Possessions)),
		strconv.Itoa(int(stats.Points)),
		strconv.Itoa(int(stats.FGM)),
		strconv.Itoa(int(stats.FGA)),
		strconv.FormatFloat(stats.FGPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.ThreePointsMade)),
		strconv.Itoa(int(stats.ThreePointAttempts)),
		strconv.FormatFloat(stats.ThreePointPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.FTM)),
		strconv.Itoa(int(stats.FTA)),
		strconv.FormatFloat(stats.FTPercent, 'f', 2, 64),
		strconv.Itoa(int(stats.OffRebounds)),
		strconv.Itoa(int(stats.DefRebounds)),
		strconv.Itoa(int(stats.Assists)),
		strconv.Itoa(int(stats.Steals)),
		strconv.Itoa(int(stats.Blocks)),
		strconv.Itoa(int(stats.Turnovers)),
		strconv.Itoa(int(stats.Fouls)),
	}
}
