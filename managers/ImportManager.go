package managers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/secrets"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"github.com/jinzhu/gorm"
)

func ImportMatchResultsToDB(Results structs.ImportMatchResultsDTO) {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println("Starting import...")

	tsChn := make(chan structs.Timestamp)

	go func() {
		ts := GetTimestamp()
		tsChn <- ts
	}()

	timestamp := <-tsChn
	close(tsChn)

	matchType := ""

	if !timestamp.GamesARan {
		matchType = "A"
	} else if !timestamp.GamesBRan {
		matchType = "B"
	} else if !timestamp.GamesCRan {
		matchType = "C"
	} else if !timestamp.GamesDRan {
		matchType = "D"
	}

	var teamStats []structs.TeamStats
	var nbaTeamStats []structs.NBATeamStats

	// Import College Game Results
	for _, dto := range Results.CBBResults {
		record := make(chan structs.Match)
		go func() {
			asyncMatch := GetMatchByMatchId(dto.GameID)
			record <- asyncMatch
		}()

		gameRecord := <-record
		close(record)

		var playerStats []structs.CollegePlayerStats

		homeTeamChn := make(chan structs.Team)
		go func() {
			homeTeam := GetTeamByTeamID(strconv.Itoa(dto.TeamOne.ID))
			homeTeamChn <- homeTeam
		}()

		ht := <-homeTeamChn
		close(homeTeamChn)

		matchID, err := strconv.Atoi(dto.GameID)
		if err != nil {
			log.Fatalln("Could not convert string to int")
		}

		homeTeam := mapToCollegeTeamStatsObject(ht.ID, uint(matchID), timestamp.CollegeWeekID, uint(timestamp.NBAWeek), timestamp.SeasonID, matchType, dto.TeamOne, dto.TeamTwo)

		teamStats = append(teamStats, homeTeam)

		awayTeamChn := make(chan structs.Team)
		go func() {
			awayTeam := GetTeamByTeamID(strconv.Itoa(dto.TeamTwo.ID))
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		awayTeam := mapToCollegeTeamStatsObject(at.ID, uint(matchID), timestamp.CollegeWeekID, uint(timestamp.NBAWeek), timestamp.SeasonID, matchType, dto.TeamTwo, dto.TeamOne)

		teamStats = append(teamStats, awayTeam)

		for _, player := range dto.RosterOne {
			id := player.ID
			collegePlayerStats := mapToCBBPlayerStatsObject(player, id, matchID, timestamp.SeasonID, timestamp.CollegeWeekID, matchType)
			playerStats = append(playerStats, collegePlayerStats)
		}

		for _, player := range dto.RosterTwo {
			id := player.ID
			collegePlayerStats := mapToCBBPlayerStatsObject(player, id, matchID, timestamp.SeasonID, timestamp.CollegeWeekID, matchType)
			playerStats = append(playerStats, collegePlayerStats)
		}

		gameRecord.UpdateScore(dto.TeamOne.Stats.Points, dto.TeamTwo.Stats.Points)

		err = db.Save(&gameRecord).Error
		if err != nil {
			log.Panicln("Could not save Game " + strconv.Itoa(int(gameRecord.ID)) + "Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
		}

		err = db.CreateInBatches(&playerStats, len(playerStats)).Error
		if err != nil {
			log.Panicln("Could not save player stats from week " + strconv.Itoa(timestamp.CollegeWeek))
		}

		fmt.Println("Finished Game " + strconv.Itoa(int(gameRecord.ID)) + " Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
	}

	// Import NBA Game Results
	for _, dto := range Results.NBAResults {
		record := make(chan structs.NBAMatch)
		go func() {
			asyncMatch := GetNBAMatchByMatchId(dto.GameID)
			record <- asyncMatch
		}()

		gameRecord := <-record
		close(record)

		var playerStats []structs.NBAPlayerStats

		homeTeamChn := make(chan structs.NBATeam)
		go func() {
			homeTeam := GetNBATeamByTeamID(strconv.Itoa(dto.TeamOne.ID))
			homeTeamChn <- homeTeam
		}()

		ht := <-homeTeamChn
		close(homeTeamChn)

		matchID := util.ConvertStringToInt(dto.GameID)

		homeTeam := mapToNBATeamStatsObject(ht.ID, uint(matchID), timestamp.NBAWeekID, uint(timestamp.NBAWeek), timestamp.SeasonID, matchType, dto.TeamOne, dto.TeamTwo)

		nbaTeamStats = append(nbaTeamStats, homeTeam)

		awayTeamChn := make(chan structs.NBATeam)
		go func() {
			awayTeam := GetNBATeamByTeamID(strconv.Itoa(dto.TeamTwo.ID))
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		awayTeam := mapToNBATeamStatsObject(at.ID, uint(matchID), timestamp.NBAWeekID, uint(timestamp.NBAWeek), timestamp.SeasonID, matchType, dto.TeamTwo, dto.TeamOne)

		nbaTeamStats = append(nbaTeamStats, awayTeam)

		for _, player := range dto.RosterOne {
			id := player.ID
			nbaPlayerStats := mapToNBAPlayerStatsObject(player, id, matchID, timestamp.SeasonID, timestamp.NBAWeekID, matchType)
			playerStats = append(playerStats, nbaPlayerStats)
		}

		for _, player := range dto.RosterTwo {
			id := player.ID
			nbaPlayerStats := mapToNBAPlayerStatsObject(player, id, matchID, timestamp.SeasonID, timestamp.NBAWeekID, matchType)
			playerStats = append(playerStats, nbaPlayerStats)
		}

		gameRecord.UpdateScore(dto.TeamOne.Stats.Points, dto.TeamTwo.Stats.Points)

		err := db.Save(&gameRecord).Error
		if err != nil {
			log.Panicln("Could not save Game " + strconv.Itoa(int(gameRecord.ID)) + "Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
		}

		err = db.CreateInBatches(&playerStats, len(playerStats)).Error
		if err != nil {
			log.Panicln("Could not save player stats from week " + strconv.Itoa(timestamp.CollegeWeek))
		}

		fmt.Println("Finished Game " + strconv.Itoa(int(gameRecord.ID)) + " Between " + gameRecord.HomeTeam + " and " + gameRecord.AwayTeam)
	}

	// Import all college team stats
	for _, stats := range teamStats {
		err := db.Create(&stats).Error
		if err != nil {
			log.Panicln("Could not save team stats!")
		}
	}
	// Import All nba team stats
	for _, stats := range nbaTeamStats {
		err := db.Create(&stats).Error
		if err != nil {
			log.Panicln("Could not save team stats!")
		}
	}
	fmt.Println("Finished Import for all games")
}

func ImportNBATeamsAndArenas() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["nbateams"]
	arenapath := secrets.GetPath()["arenas"]
	nbaTeamsCSV := util.ReadCSV(path)

	for idx, row := range nbaTeamsCSV {
		if idx < 2 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		team := row[1]
		nickname := row[2]
		abbr := row[3]
		city := row[4]
		state := row[5]
		country := row[6]
		conferenceID := util.ConvertStringToInt(row[7])
		conference := row[8]
		divisionID := util.ConvertStringToInt(row[9])
		division := row[10]
		arenaID := util.ConvertStringToInt(row[11])
		arena := row[12]

		nbaTeam := structs.NBATeam{
			Team:         team,
			Nickname:     nickname,
			Abbr:         abbr,
			City:         city,
			State:        state,
			Country:      country,
			ConferenceID: uint(conferenceID),
			Conference:   conference,
			DivisionID:   uint(divisionID),
			Division:     division,
			ArenaID:      uint(arenaID),
			Arena:        arena,
			IsActive:     true,
		}

		nbaTeam.AssignID(uint(id))

		db.Create(&nbaTeam)
	}

	arenasCSV := util.ReadCSV(arenapath)

	for idx, row := range arenasCSV {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		name := row[1]
		city := row[2]
		state := row[3]
		country := row[4]
		capacity := util.ConvertStringToInt(row[5])
		hometeam := row[6]

		arena := structs.Arena{
			ArenaName: name,
			City:      city,
			State:     state,
			Country:   country,
			Capacity:  uint(capacity),
			HomeTeam:  hometeam,
		}

		arena.AssignID(uint(id))

		db.Create(&arena)
	}
}

func SyncContractValues() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	SalaryCap := ts.Y1Capspace

	players := GetAllNBAPlayers()

	for _, player := range players {
		val := 0.0

		// Check if Max or Supermax qualified
		if player.Overall > 109 && !player.MaxRequested {
			player.ToggleMaxRequested()
		}

		if (player.IsDPOY || player.IsMVP || player.IsFirstTeamANBA) && !player.IsSuperMaxQualified {
			// SUPER MAX
			player.ToggleSuperMax()
		}

		if player.IsSuperMaxQualified {
			if player.Year > 9 {
				val = SalaryCap * 0.35
			} else if player.Year > 6 {
				val = SalaryCap * 0.3
			} else {
				val = SalaryCap * 0.25
			}
		} else if player.MaxRequested {
			if player.Year > 9 {
				val = SalaryCap * 0.3
			} else if player.Year > 6 {
				val = SalaryCap * 0.25
			} else {
				val = SalaryCap * 0.2
			}
		} else {
			if player.Year > 9 {
				val = 2.5
			} else if player.Year > 8 {
				val = 2.4
			} else if player.Year > 7 {
				val = 2.3
			} else if player.Year > 6 {
				val = 2.2
			} else if player.Year > 5 {
				val = 2.0
			} else if player.Year > 4 {
				val = 1.9
			} else if player.Year > 3 {
				val = 1.8
			} else if player.Year > 2 {
				val = 1.7
			} else if player.Year > 1 {
				val = 1.6
			} else if player.Year > 0 {
				val = 1.5
			} else {
				val = 0.9
			}
		}
		player.AssignMinimumContractValue(val)

		db.Save(&player)
	}
}

func ImportFAPreferences() {
	fmt.Println(time.Now().UnixNano())
	db := dbprovider.GetInstance().GetDB()
	nbaPlayers := GetAllNBAPlayers()

	for _, p := range nbaPlayers {
		NegotiationRound := 0
		if p.Overall > 95 {
			NegotiationRound = util.GenerateIntFromRange(2, 4)
		} else {
			NegotiationRound = util.GenerateIntFromRange(3, 6)
		}

		SigningRound := NegotiationRound + util.GenerateIntFromRange(2, 5)
		if SigningRound > 10 {
			SigningRound = 10
		}

		p.AssignFAPreferences(uint(NegotiationRound), uint(SigningRound))

		db.Save(&p)
	}
}

func ImportNBAPersonalities() {
	fmt.Println(time.Now().UnixNano())
	db := dbprovider.GetInstance().GetDB()
	nbaPlayers := GetAllNBAPlayers()

	for _, p := range nbaPlayers {
		freeAgency := util.GetFreeAgencyBias(p.Age, p.Overall)
		workEthic := util.GetWorkEthic()
		personality := util.GetPersonality()

		p.SetFreeAgencyBias(freeAgency)
		p.SetWorkEthic(workEthic)
		p.SetPersonality(personality)

		db.Save(&p)
	}
}

func ImportCBBGames() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["cbbmatches"]
	collegeMatches := util.ReadCSV(path)

	collegeTeams := GetAllActiveCollegeTeams()
	collegeMap := make(map[string]structs.Team)

	for _, t := range collegeTeams {
		collegeMap[t.Abbr] = t
	}

	for idx, row := range collegeMatches {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		season := util.ConvertStringToInt(row[1])
		seasonID := season - 2020
		week := util.ConvertStringToInt(row[2])
		weekID := week + 20
		timeSlot := row[3]
		matchType := row[4]
		isConf := false
		if matchType == "Conf" {
			isConf = true
		}
		homeTeamAbbr := row[5]
		awayTeamAbbr := row[6]
		homeTeam := collegeMap[homeTeamAbbr]
		awayTeam := collegeMap[awayTeamAbbr]
		gameTitle := row[25]
		nextGameID := util.ConvertStringToInt(row[18])
		hoA := row[19]
		neutralSite := util.ConvertStringToBool(row[12])
		invitational := util.ConvertStringToBool(row[13])
		conferenceTournament := util.ConvertStringToBool(row[14])
		nit := util.ConvertStringToBool(row[15])
		tournament := util.ConvertStringToBool(row[16])
		nationalChamp := util.ConvertStringToBool(row[17])
		arena := row[22]
		city := row[23]
		state := row[24]
		homeCoach := homeTeam.Coach
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.Coach
		if awayCoach == "" {
			awayCoach = "AI"
		}

		match := structs.Match{
			Model:                  gorm.Model{ID: uint(id)},
			SeasonID:               uint(seasonID),
			WeekID:                 uint(weekID),
			Week:                   uint(week),
			MatchOfWeek:            timeSlot,
			IsConference:           isConf,
			HomeTeam:               homeTeamAbbr,
			HomeTeamID:             homeTeam.ID,
			AwayTeamID:             awayTeam.ID,
			HomeTeamCoach:          homeCoach,
			AwayTeam:               awayTeamAbbr,
			AwayTeamCoach:          awayCoach,
			MatchName:              gameTitle,
			NextGameID:             uint(nextGameID),
			NextGameHOA:            hoA,
			IsNeutralSite:          neutralSite,
			IsInvitational:         invitational,
			IsConferenceTournament: conferenceTournament,
			IsNITGame:              nit,
			IsPlayoffGame:          tournament,
			IsNationalChampionship: nationalChamp,
			Arena:                  arena,
			City:                   city,
			State:                  state,
		}

		db.Create(&match)
	}
}

func ImportNBAGames() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["nbamatches"]
	professionalMatches := util.ReadCSV(path)

	professionalTeams := GetAllActiveNBATeams()
	teamMap := make(map[string]structs.NBATeam)

	for _, t := range professionalTeams {
		teamMap[t.Team+" "+t.Nickname] = t
	}

	for idx, row := range professionalMatches {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		season := util.ConvertStringToInt(row[1])
		seasonID := season - 2020
		week := util.ConvertStringToInt(row[2])
		weekID := week
		timeSlot := row[3]
		homeTeamStr := row[6]
		awayTeamStr := row[7]
		homeTeam := teamMap[homeTeamStr]
		awayTeam := teamMap[awayTeamStr]
		gameTitle := row[26]
		nextGameID := util.ConvertStringToInt(row[18])
		hoA := row[19]
		conference := util.ConvertStringToBool(row[12])
		divisional := util.ConvertStringToBool(row[13])
		international := util.ConvertStringToBool(row[14])
		playoff := util.ConvertStringToBool(row[16])
		finals := util.ConvertStringToBool(row[17])
		arena := row[22]
		city := row[23]
		state := row[24]
		country := row[25]
		homeCoach := homeTeam.NBACoachName
		if homeCoach == "" {
			homeCoach = "AI"
		}
		awayCoach := awayTeam.NBACoachName
		if awayCoach == "" {
			awayCoach = "AI"
		}

		match := structs.NBAMatch{
			Model:           gorm.Model{ID: uint(id)},
			SeasonID:        uint(seasonID),
			WeekID:          uint(weekID),
			Week:            uint(week),
			MatchOfWeek:     timeSlot,
			IsConference:    conference,
			IsDivisional:    divisional,
			HomeTeam:        homeTeamStr,
			HomeTeamID:      homeTeam.ID,
			AwayTeamID:      awayTeam.ID,
			HomeTeamCoach:   homeCoach,
			AwayTeam:        awayTeamStr,
			AwayTeamCoach:   awayCoach,
			MatchName:       gameTitle,
			NextGameID:      uint(nextGameID),
			NextGameHOA:     hoA,
			IsNeutralSite:   conference,
			IsPlayoffGame:   playoff,
			IsTheFinals:     finals,
			IsInternational: international,
			Arena:           arena,
			City:            city,
			State:           state,
			Country:         country,
		}

		db.Create(&match)
	}
}

func ImportNBAStandings() {
	db := dbprovider.GetInstance().GetDB()
	path := secrets.GetPath()["nbastandings"]
	nbaStandingsCSV := util.ReadCSV(path)

	for idx, row := range nbaStandingsCSV {
		if idx < 1 {
			continue
		}

		id := util.ConvertStringToInt(row[0])
		teamID := util.ConvertStringToInt(row[1])
		team := row[2]
		seasonID := util.ConvertStringToInt(row[3])
		season := util.ConvertStringToInt(row[4])
		leagueID := util.ConvertStringToInt(row[5])
		league := row[6]
		conferenceID := util.ConvertStringToInt(row[7])
		conference := row[8]
		divisionID := util.ConvertStringToInt(row[9])
		division := row[10]
		postSeasonStatus := row[11]
		isConferenceChampion := util.ConvertStringToBool(row[12])
		totalWins := util.ConvertStringToInt(row[13])
		totalLosses := util.ConvertStringToInt(row[14])

		standings := structs.NBAStandings{
			TeamID:               uint(teamID),
			TeamName:             team,
			SeasonID:             uint(seasonID),
			Season:               season,
			LeagueID:             uint(leagueID),
			League:               league,
			ConferenceID:         uint(conferenceID),
			ConferenceName:       conference,
			DivisionID:           uint(divisionID),
			DivisionName:         division,
			PostSeasonStatus:     postSeasonStatus,
			IsConferenceChampion: isConferenceChampion,
			BaseStandings: structs.BaseStandings{
				TotalWins:   totalWins,
				TotalLosses: totalLosses,
			},
		}

		standings.AssignID(uint(id))

		db.Create(&standings)
	}
}

func ImportNewPositions() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()
	recruits := GetAllRecruitRecords()

	for _, c := range collegePlayers {
		if c.Position == "C" {
			continue
		}
		shooting := (c.Shooting2 + c.Shooting3) / 2
		if c.Position == "G" {
			if shooting > c.Ballwork || c.Archetype == "Floor General" {
				c.SetNewPosition("PG")
			} else {
				c.SetNewPosition("SG")
			}
		} else {
			if c.Rebounding > shooting || c.Archetype == "Point Forward" {
				c.SetNewPosition("PF")
			} else {
				c.SetNewPosition("SF")
			}
		}
		db.Save(&c)
	}

	for _, r := range recruits {
		if r.Position == "C" {
			continue
		}
		shooting := (r.Shooting2 + r.Shooting3) / 2
		if r.Position == "G" {
			if shooting > r.Ballwork || r.Archetype == "Floor General" {
				r.SetNewPosition("PG")
			} else {
				r.SetNewPosition("SG")
			}
		} else {
			if r.Rebounding > shooting || r.Archetype == "Point Forward" {
				r.SetNewPosition("PF")
			} else {
				r.SetNewPosition("SF")
			}
		}
		db.Save(&r)
	}
}

func ImportNewTeams() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	path := secrets.GetPath()["teams"]
	teams := util.ReadCSV(path)

	for idx, row := range teams {
		if idx == 0 {
			continue
		}
		teamID := util.ConvertStringToInt(row[0])
		team := row[1]
		nickname := row[2]
		abbr := row[3]
		city := row[4]
		state := row[5]
		country := row[6]
		conferenceID := util.ConvertStringToInt(row[7])
		conference := row[8]
		season := row[10]
		isActive := util.ConvertStringToBool(row[13])

		t := structs.Team{
			Team:         team,
			Nickname:     nickname,
			Abbr:         abbr,
			City:         city,
			State:        state,
			Country:      country,
			ConferenceID: uint(conferenceID),
			Conference:   conference,
			IsActive:     isActive,
			IsNBA:        false,
			FirstSeason:  season,
			Model:        gorm.Model{ID: uint(teamID)},
		}

		db.Create(&t)

		standings := structs.CollegeStandings{
			TeamID:         uint(teamID),
			TeamName:       team,
			TeamAbbr:       abbr,
			SeasonID:       ts.SeasonID,
			Season:         ts.Season,
			ConferenceID:   uint(conferenceID),
			ConferenceName: conference,
		}

		db.Create(&standings)
	}
}

func ImportDraftPicks() {
	db := dbprovider.GetInstance().GetDB()

	path := secrets.GetPath()["draftpicks"]
	picks := util.ReadCSV(path)

	for idx, row := range picks {
		if idx == 0 {
			continue
		}
		id := util.ConvertStringToInt(row[0])
		round := util.ConvertStringToInt(row[1])
		draftNumber := util.ConvertStringToInt(row[2])
		season_id := util.ConvertStringToInt(row[3])
		season := util.ConvertStringToInt(row[4])
		drafteeID := 0
		teamID := util.ConvertStringToInt(row[5])
		team := row[6]
		originalTeamID := util.ConvertStringToInt(row[7])
		originalTeam := row[8]
		notes := ""

		draftPick := structs.DraftPick{
			DraftRound:     uint(round),
			DraftNumber:    uint(draftNumber),
			DrafteeID:      uint(drafteeID),
			TeamID:         uint(teamID),
			Team:           team,
			OriginalTeamID: uint(originalTeamID),
			OriginalTeam:   originalTeam,
			Season:         uint(season),
			Notes:          notes,
			SeasonID:       uint(season_id),
			Model:          gorm.Model{ID: uint(id)},
		}

		db.Create(&draftPick)
	}
}

func filterLotteryPicks(list []structs.DraftLottery, id uint) []structs.DraftLottery {
	newList := []structs.DraftLottery{}
	for _, l := range list {
		if l.ID != id {
			newList = append(newList, l)
		}
	}
	return newList
}

func mapToCollegeTeamStatsObject(teamID, matchID, weekID, week, seasonID uint, matchType string, TeamOne, TeamTwo structs.TeamResultsDTO) structs.TeamStats {
	return structs.TeamStats{
		TeamID:                    teamID,
		MatchID:                   matchID,
		WeekID:                    weekID,
		SeasonID:                  seasonID,
		MatchType:                 matchType,
		Points:                    TeamOne.Stats.Points,
		Possessions:               TeamOne.Stats.Possessions,
		FGM:                       TeamOne.Stats.FGM,
		FGA:                       TeamOne.Stats.FGA,
		FGPercent:                 TeamOne.Stats.FGPercent,
		ThreePointsMade:           TeamOne.Stats.ThreePointsMade,
		ThreePointAttempts:        TeamOne.Stats.ThreePointAttempts,
		ThreePointPercent:         TeamOne.Stats.ThreePointPercent,
		FTM:                       TeamOne.Stats.FTM,
		FTA:                       TeamOne.Stats.FTA,
		FTPercent:                 TeamOne.Stats.FTPercent,
		Rebounds:                  TeamOne.Stats.Rebounds,
		OffRebounds:               TeamOne.Stats.OffRebounds,
		DefRebounds:               TeamOne.Stats.DefRebounds,
		Assists:                   TeamOne.Stats.Assists,
		Steals:                    TeamOne.Stats.Steals,
		Blocks:                    TeamOne.Stats.Blocks,
		TotalTurnovers:            TeamOne.Stats.TotalTurnovers,
		LargestLead:               TeamOne.Stats.LargestLead,
		FirstHalfScore:            TeamOne.Stats.FirstHalfScore,
		SecondHalfScore:           TeamOne.Stats.SecondHalfScore,
		OvertimeScore:             TeamOne.Stats.OvertimeScore,
		Fouls:                     TeamOne.Stats.Fouls,
		PointsAgainst:             TeamTwo.Stats.Points,
		FGMAgainst:                TeamTwo.Stats.FGM,
		FGAAgainst:                TeamTwo.Stats.FGA,
		FGPercentAgainst:          TeamTwo.Stats.FGPercent,
		ThreePointsMadeAgainst:    TeamTwo.Stats.ThreePointsMade,
		ThreePointAttemptsAgainst: TeamTwo.Stats.ThreePointAttempts,
		ThreePointPercentAgainst:  TeamTwo.Stats.ThreePointPercent,
		FTMAgainst:                TeamTwo.Stats.FTM,
		FTAAgainst:                TeamTwo.Stats.FTA,
		FTPercentAgainst:          TeamTwo.Stats.FTPercent,
		ReboundsAllowed:           TeamTwo.Stats.Rebounds,
		OffReboundsAllowed:        TeamTwo.Stats.OffRebounds,
		DefReboundsAllowed:        TeamTwo.Stats.DefRebounds,
		AssistsAllowed:            TeamTwo.Stats.Assists,
		StealsAllowed:             TeamTwo.Stats.Steals,
		BlocksAllowed:             TeamTwo.Stats.Blocks,
		TurnoversAllowed:          TeamTwo.Stats.TotalTurnovers,
	}
}

func mapToNBATeamStatsObject(teamID, matchID, weekID, week, seasonID uint, matchType string, TeamOne, TeamTwo structs.TeamResultsDTO) structs.NBATeamStats {
	return structs.NBATeamStats{
		TeamID:                    teamID,
		MatchID:                   matchID,
		WeekID:                    weekID,
		SeasonID:                  seasonID,
		MatchType:                 matchType,
		Points:                    TeamOne.Stats.Points,
		Possessions:               TeamOne.Stats.Possessions,
		FGM:                       TeamOne.Stats.FGM,
		FGA:                       TeamOne.Stats.FGA,
		FGPercent:                 TeamOne.Stats.FGPercent,
		ThreePointsMade:           TeamOne.Stats.ThreePointsMade,
		ThreePointAttempts:        TeamOne.Stats.ThreePointAttempts,
		ThreePointPercent:         TeamOne.Stats.ThreePointPercent,
		FTM:                       TeamOne.Stats.FTM,
		FTA:                       TeamOne.Stats.FTA,
		FTPercent:                 TeamOne.Stats.FTPercent,
		Rebounds:                  TeamOne.Stats.Rebounds,
		OffRebounds:               TeamOne.Stats.OffRebounds,
		DefRebounds:               TeamOne.Stats.DefRebounds,
		Assists:                   TeamOne.Stats.Assists,
		Steals:                    TeamOne.Stats.Steals,
		Blocks:                    TeamOne.Stats.Blocks,
		TotalTurnovers:            TeamOne.Stats.TotalTurnovers,
		LargestLead:               TeamOne.Stats.LargestLead,
		FirstHalfScore:            TeamOne.Stats.FirstHalfScore,
		SecondHalfScore:           TeamOne.Stats.SecondHalfScore,
		OvertimeScore:             TeamOne.Stats.OvertimeScore,
		Fouls:                     TeamOne.Stats.Fouls,
		PointsAgainst:             TeamTwo.Stats.Points,
		FGMAgainst:                TeamTwo.Stats.FGM,
		FGAAgainst:                TeamTwo.Stats.FGA,
		FGPercentAgainst:          TeamTwo.Stats.FGPercent,
		ThreePointsMadeAgainst:    TeamTwo.Stats.ThreePointsMade,
		ThreePointAttemptsAgainst: TeamTwo.Stats.ThreePointAttempts,
		ThreePointPercentAgainst:  TeamTwo.Stats.ThreePointPercent,
		FTMAgainst:                TeamTwo.Stats.FTM,
		FTAAgainst:                TeamTwo.Stats.FTA,
		FTPercentAgainst:          TeamTwo.Stats.FTPercent,
		ReboundsAllowed:           TeamTwo.Stats.Rebounds,
		OffReboundsAllowed:        TeamTwo.Stats.OffRebounds,
		DefReboundsAllowed:        TeamTwo.Stats.DefRebounds,
		AssistsAllowed:            TeamTwo.Stats.Assists,
		StealsAllowed:             TeamTwo.Stats.Steals,
		BlocksAllowed:             TeamTwo.Stats.Blocks,
		TurnoversAllowed:          TeamTwo.Stats.TotalTurnovers,
	}
}

func mapToCBBPlayerStatsObject(player structs.PlayerDTO, id, matchID int, seasonID, weekID uint, matchType string) structs.CollegePlayerStats {
	return structs.CollegePlayerStats{
		CollegePlayerID:    uint(id),
		MatchID:            uint(matchID),
		SeasonID:           seasonID,
		MatchType:          matchType,
		Minutes:            player.Stats.Minutes,
		Possessions:        player.Stats.Possessions,
		FGM:                player.Stats.FGM,
		FGA:                player.Stats.FGA,
		FGPercent:          player.Stats.FGPercent,
		ThreePointsMade:    player.Stats.ThreePointsMade,
		ThreePointAttempts: player.Stats.ThreePointAttempts,
		ThreePointPercent:  player.Stats.ThreePointPercent,
		FTM:                player.Stats.FTM,
		FTA:                player.Stats.FTA,
		FTPercent:          player.Stats.FTPercent,
		Points:             player.Stats.Points,
		TotalRebounds:      player.Stats.TotalRebounds,
		OffRebounds:        player.Stats.OffRebounds,
		DefRebounds:        player.Stats.DefRebounds,
		Assists:            player.Stats.Assists,
		Steals:             player.Stats.Steals,
		Blocks:             player.Stats.Blocks,
		Turnovers:          player.Stats.Turnovers,
		Fouls:              player.Stats.Fouls,
	}
}

func mapToNBAPlayerStatsObject(player structs.PlayerDTO, id, matchID int, seasonID, weekID uint, matchType string) structs.NBAPlayerStats {
	return structs.NBAPlayerStats{
		NBAPlayerID:        uint(id),
		MatchID:            uint(matchID),
		SeasonID:           seasonID,
		WeekID:             weekID,
		MatchType:          matchType,
		Minutes:            player.Stats.Minutes,
		Possessions:        player.Stats.Possessions,
		FGM:                player.Stats.FGM,
		FGA:                player.Stats.FGA,
		FGPercent:          player.Stats.FGPercent,
		ThreePointsMade:    player.Stats.ThreePointsMade,
		ThreePointAttempts: player.Stats.ThreePointAttempts,
		ThreePointPercent:  player.Stats.ThreePointPercent,
		FTM:                player.Stats.FTM,
		FTA:                player.Stats.FTA,
		FTPercent:          player.Stats.FTPercent,
		Points:             player.Stats.Points,
		TotalRebounds:      player.Stats.TotalRebounds,
		OffRebounds:        player.Stats.OffRebounds,
		DefRebounds:        player.Stats.DefRebounds,
		Assists:            player.Stats.Assists,
		Steals:             player.Stats.Steals,
		Blocks:             player.Stats.Blocks,
		Turnovers:          player.Stats.Turnovers,
		Fouls:              player.Stats.Fouls,
	}
}
