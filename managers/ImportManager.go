package managers

import (
	"fmt"
	"log"
	"sort"
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

	var teamStats []structs.TeamStats

	for _, dto := range Results.Results {
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

		homeTeam := structs.TeamStats{
			TeamID:                    uint(ht.ID),
			MatchID:                   uint(matchID),
			WeekID:                    gameRecord.WeekID,
			SeasonID:                  uint(gameRecord.SeasonID),
			Points:                    dto.TeamOne.Stats.Points,
			Possessions:               dto.TeamOne.Stats.Possessions,
			FGM:                       dto.TeamOne.Stats.FGM,
			FGA:                       dto.TeamOne.Stats.FGA,
			FGPercent:                 dto.TeamOne.Stats.FGPercent,
			ThreePointsMade:           dto.TeamOne.Stats.ThreePointsMade,
			ThreePointAttempts:        dto.TeamOne.Stats.ThreePointAttempts,
			ThreePointPercent:         dto.TeamOne.Stats.ThreePointPercent,
			FTM:                       dto.TeamOne.Stats.FTM,
			FTA:                       dto.TeamOne.Stats.FTA,
			FTPercent:                 dto.TeamOne.Stats.FTPercent,
			Rebounds:                  dto.TeamOne.Stats.Rebounds,
			OffRebounds:               dto.TeamOne.Stats.OffRebounds,
			DefRebounds:               dto.TeamOne.Stats.DefRebounds,
			Assists:                   dto.TeamOne.Stats.Assists,
			Steals:                    dto.TeamOne.Stats.Steals,
			Blocks:                    dto.TeamOne.Stats.Blocks,
			TotalTurnovers:            dto.TeamOne.Stats.TotalTurnovers,
			LargestLead:               dto.TeamOne.Stats.LargestLead,
			FirstHalfScore:            dto.TeamOne.Stats.FirstHalfScore,
			SecondHalfScore:           dto.TeamOne.Stats.SecondHalfScore,
			OvertimeScore:             dto.TeamOne.Stats.OvertimeScore,
			Fouls:                     dto.TeamOne.Stats.Fouls,
			PointsAgainst:             dto.TeamTwo.Stats.Points,
			FGMAgainst:                dto.TeamTwo.Stats.FGM,
			FGAAgainst:                dto.TeamTwo.Stats.FGA,
			FGPercentAgainst:          dto.TeamTwo.Stats.FGPercent,
			ThreePointsMadeAgainst:    dto.TeamTwo.Stats.ThreePointsMade,
			ThreePointAttemptsAgainst: dto.TeamTwo.Stats.ThreePointAttempts,
			ThreePointPercentAgainst:  dto.TeamTwo.Stats.ThreePointPercent,
			FTMAgainst:                dto.TeamTwo.Stats.FTM,
			FTAAgainst:                dto.TeamTwo.Stats.FTA,
			FTPercentAgainst:          dto.TeamTwo.Stats.FTPercent,
			ReboundsAllowed:           dto.TeamTwo.Stats.Rebounds,
			OffReboundsAllowed:        dto.TeamTwo.Stats.OffRebounds,
			DefReboundsAllowed:        dto.TeamTwo.Stats.DefRebounds,
			AssistsAllowed:            dto.TeamTwo.Stats.Assists,
			StealsAllowed:             dto.TeamTwo.Stats.Steals,
			BlocksAllowed:             dto.TeamTwo.Stats.Blocks,
			TurnoversAllowed:          dto.TeamTwo.Stats.TotalTurnovers,
		}

		teamStats = append(teamStats, homeTeam)

		awayTeamChn := make(chan structs.Team)
		go func() {
			awayTeam := GetTeamByTeamID(strconv.Itoa(dto.TeamTwo.ID))
			awayTeamChn <- awayTeam
		}()

		at := <-awayTeamChn
		close(awayTeamChn)

		awayTeam := structs.TeamStats{
			TeamID:                    at.ID,
			MatchID:                   uint(matchID),
			WeekID:                    gameRecord.WeekID,
			SeasonID:                  gameRecord.SeasonID,
			Points:                    dto.TeamTwo.Stats.Points,
			Possessions:               dto.TeamTwo.Stats.Possessions,
			FGM:                       dto.TeamTwo.Stats.FGM,
			FGA:                       dto.TeamTwo.Stats.FGA,
			FGPercent:                 dto.TeamTwo.Stats.FGPercent,
			ThreePointsMade:           dto.TeamTwo.Stats.ThreePointsMade,
			ThreePointAttempts:        dto.TeamTwo.Stats.ThreePointAttempts,
			ThreePointPercent:         dto.TeamTwo.Stats.ThreePointPercent,
			FTM:                       dto.TeamTwo.Stats.FTM,
			FTA:                       dto.TeamTwo.Stats.FTA,
			FTPercent:                 dto.TeamTwo.Stats.FTPercent,
			Rebounds:                  dto.TeamTwo.Stats.Rebounds,
			OffRebounds:               dto.TeamTwo.Stats.OffRebounds,
			DefRebounds:               dto.TeamTwo.Stats.DefRebounds,
			Assists:                   dto.TeamTwo.Stats.Assists,
			Steals:                    dto.TeamTwo.Stats.Steals,
			Blocks:                    dto.TeamTwo.Stats.Blocks,
			TotalTurnovers:            dto.TeamTwo.Stats.TotalTurnovers,
			LargestLead:               dto.TeamTwo.Stats.LargestLead,
			FirstHalfScore:            dto.TeamTwo.Stats.FirstHalfScore,
			SecondHalfScore:           dto.TeamTwo.Stats.SecondHalfScore,
			OvertimeScore:             dto.TeamTwo.Stats.OvertimeScore,
			Fouls:                     dto.TeamTwo.Stats.Fouls,
			PointsAgainst:             dto.TeamOne.Stats.Points,
			FGMAgainst:                dto.TeamOne.Stats.FGM,
			FGAAgainst:                dto.TeamOne.Stats.FGA,
			FGPercentAgainst:          dto.TeamOne.Stats.FGPercent,
			ThreePointsMadeAgainst:    dto.TeamOne.Stats.ThreePointsMade,
			ThreePointAttemptsAgainst: dto.TeamOne.Stats.ThreePointAttempts,
			ThreePointPercentAgainst:  dto.TeamOne.Stats.ThreePointPercent,
			FTMAgainst:                dto.TeamOne.Stats.FTM,
			FTAAgainst:                dto.TeamOne.Stats.FTA,
			FTPercentAgainst:          dto.TeamOne.Stats.FTPercent,
			ReboundsAllowed:           dto.TeamOne.Stats.Rebounds,
			OffReboundsAllowed:        dto.TeamOne.Stats.OffRebounds,
			DefReboundsAllowed:        dto.TeamOne.Stats.DefRebounds,
			AssistsAllowed:            dto.TeamOne.Stats.Assists,
			StealsAllowed:             dto.TeamOne.Stats.Steals,
			BlocksAllowed:             dto.TeamOne.Stats.Blocks,
			TurnoversAllowed:          dto.TeamOne.Stats.TotalTurnovers,
		}

		teamStats = append(teamStats, awayTeam)

		for _, player := range dto.RosterOne {
			id := player.ID
			collegePlayerStats := structs.CollegePlayerStats{
				CollegePlayerID:    uint(id),
				MatchID:            uint(matchID),
				SeasonID:           timestamp.SeasonID,
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
			playerStats = append(playerStats, collegePlayerStats)
		}

		for _, player := range dto.RosterTwo {
			id := player.ID
			collegePlayerStats := structs.CollegePlayerStats{
				CollegePlayerID:    uint(id),
				MatchID:            uint(matchID),
				SeasonID:           timestamp.SeasonID,
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

	for _, stats := range teamStats {
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

func ConductDraftLottery() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	path := secrets.GetPath()["draftlottery"]
	lotteryCSV := util.ReadCSV(path)
	ts := GetTimestamp()
	lotteryBalls := []structs.DraftLottery{}
	draftPicks := []structs.DraftPick{}

	for idx, row := range lotteryCSV {
		if idx < 1 {
			continue
		}
		pickNumber := idx + 1
		teamID := util.ConvertStringToInt(row[0])
		team := row[1]

		if idx < 15 {
			chances := util.GetLotteryChances(idx)
			lottery := structs.DraftLottery{
				ID:      uint(teamID),
				Team:    team,
				Chances: chances,
			}
			lotteryBalls = append(lotteryBalls, lottery)

			draftPick := structs.DraftPick{
				SeasonID:       ts.SeasonID,
				Season:         uint(ts.Season),
				DraftRound:     2,
				DraftNumber:    uint(32 + pickNumber),
				TeamID:         uint(teamID),
				Team:           team,
				OriginalTeamID: uint(teamID),
				OriginalTeam:   team,
				DraftValue:     0,
			}
			draftPicks = append(draftPicks, draftPick)
		} else {
			firstRoundPick := structs.DraftPick{
				SeasonID:       ts.SeasonID,
				Season:         uint(ts.Season),
				DraftRound:     2,
				DraftNumber:    uint(32 + pickNumber),
				TeamID:         uint(teamID),
				Team:           team,
				OriginalTeamID: uint(teamID),
				OriginalTeam:   team,
				DraftValue:     0,
			}
			secondRoundPick := structs.DraftPick{
				SeasonID:       ts.SeasonID,
				Season:         uint(ts.Season),
				DraftRound:     2,
				DraftNumber:    uint(32 + pickNumber),
				TeamID:         uint(teamID),
				Team:           team,
				OriginalTeamID: uint(teamID),
				OriginalTeam:   team,
				DraftValue:     0,
			}
			draftPicks = append(draftPicks, firstRoundPick)
			draftPicks = append(draftPicks, secondRoundPick)
		}
	}
	lotteryPicks := 16
	draftOrder := []structs.DraftLottery{}
	for i := 0; i < lotteryPicks; i++ {
		sum := 0
		for _, l := range lotteryBalls {
			sum += int(l.Chances)
		}

		chance := util.GenerateIntFromRange(1, sum)
		sum2 := 0
		for _, l := range lotteryBalls {
			sum2 += int(l.Chances)
			if chance < sum2 {
				draftOrder = append(draftOrder, l)

				lotteryBalls = filterLotteryPicks(lotteryBalls, l.ID)
				break
			}
		}
	}

	for idx, do := range draftOrder {
		pick := idx + 1
		draftPick := structs.DraftPick{
			SeasonID:       ts.SeasonID,
			Season:         uint(ts.Season),
			DraftRound:     1,
			DraftNumber:    uint(pick),
			TeamID:         do.ID,
			Team:           do.Team,
			OriginalTeamID: do.ID,
			OriginalTeam:   do.Team,
			DraftValue:     0,
		}
		draftPicks = append(draftPicks, draftPick)
	}

	sort.Sort(structs.ByDraftNumber(draftPicks))

	for _, pick := range draftPicks {
		db.Create(&pick)
	}
}

func GenerateGameplans() {
	db := dbprovider.GetInstance().GetDB()

	allProfessionalTeams := GetAllActiveNBATeams()

	for _, n := range allProfessionalTeams {
		gp := GetNBAGameplanByTeam(strconv.Itoa(int(n.ID)))
		if gp.ID > 0 {
			continue
		}
		gameplan := structs.NBAGameplan{
			TeamID:             n.ID,
			Game:               "A",
			Pace:               "Balanced",
			FocusPlayer:        "",
			OffensiveFormation: "Balanced",
			DefensiveFormation: "Man-to-Man",
			OffensiveStyle:     "Traditional",
		}
		db.Create(&gameplan)
	}
}

func GeneratePlaytimeExpectations() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()

	for _, c := range collegePlayers {
		newExpectations := util.GetPlaytimeExpectations(c.Stars, c.Year, c.Overall)

		c.SetExpectations(newExpectations)

		db.Save(&c)
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
