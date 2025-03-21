package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"gorm.io/gorm"
)

func GetTimestamp() structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()

	var timeStamp structs.Timestamp

	err := db.Find(&timeStamp).Error
	if err != nil {
		log.Fatal(err)
	}

	return timeStamp
}

func LockRecruiting() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()

	ts.ToggleLockRecruiting()

	err := db.Save(&ts).Error
	if err != nil {
		log.Fatal(err)
	}
}

func MoveUpInOffseasonFreeAgency() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	if ts.IsNBAOffseason || ts.NBASeasonOver {
		ts.MoveUpFreeAgencyRound()
	}
	db.Save(&ts)
}

func SyncToNextWeek() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	if ts.CollegeWeek < 21 || !ts.IsOffSeason {
		ResetCollegeStandingsRanks()
	}

	ts.SyncToNextWeek()

	if ts.CollegeWeek < 21 || !ts.CollegeSeasonOver {
		SyncCollegePollSubmissionForCurrentWeek(uint(ts.CollegeWeek), ts.CollegeWeekID, ts.SeasonID)
		ts.TogglePollRan()
	}
	if ts.NBAWeek > 21 && !ts.IsNBAOffseason {
		// Update bools so that teams can't trade in middle of next season's post season
		db.Model(&structs.NBATeam{}).Where("id < ?", "33").Update("can_trade", false)
		GenerateNBAPlayoffGames(db, ts)
		IndicateWhetherTeamCanTradeInPostSeason()
	}
	if ts.NBASeasonOver && ts.CollegeSeasonOver && ts.FreeAgencyRound > 2 {
		ts.MoveUpASeason()
	}
	repository.SaveTimeStamp(ts, db)
}

func GenerateNBAPlayoffGames(db *gorm.DB, ts structs.Timestamp) {

	nbaWeekID := strconv.Itoa(int(ts.NBAWeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	nbaGames := GetNBAMatchesByWeekIdAndMatchType(nbaWeekID, seasonID, "A")
	teamMap := GetProfessionalTeamMap()
	finalsTally := 0
	if len(nbaGames) == 0 {
		// Get active NBA Series
		activeNBASeries := GetAllActiveNBASeries()
		for _, s := range activeNBASeries {
			// Skip series that are not ready
			if s.HomeTeamID == 0 || s.AwayTeamID == 0 {
				continue
			}
			if s.IsTheFinals && s.SeriesComplete {
				finalsTally += 1
				// Officially End the season
				if finalsTally == 2 {
					ts.EndTheProfessionalSeason()
					repository.SaveTimeStamp(ts, db)
					break
				}

			}
			gameCount := strconv.Itoa(int(s.GameCount))
			homeTeam := ""
			homeTeamID := 0
			homeTeamCoach := ""
			homeTeamRank := 0
			awayTeam := ""
			awayTeamID := 0
			awayTeamCoach := ""
			awayTeamRank := 0
			arena := ""
			city := ""
			state := ""
			country := ""
			seriesName := s.SeriesName
			matchName := seriesName + " Game: " + gameCount
			if gameCount == "1" || gameCount == "2" || gameCount == "5" || gameCount == "7" {
				ht := teamMap[s.HomeTeamID]
				homeTeam = s.HomeTeam
				homeTeamID = int(s.HomeTeamID)
				homeTeamCoach = s.HomeTeamCoach
				homeTeamRank = int(s.HomeTeamRank)
				awayTeam = s.AwayTeam
				awayTeamID = int(s.AwayTeamID)
				awayTeamCoach = s.AwayTeamCoach
				awayTeamRank = int(s.AwayTeamRank)
				arena = ht.Arena
				city = ht.City
				state = ht.State
				country = ht.Country
			} else {
				ht := teamMap[s.AwayTeamID]
				homeTeam = s.AwayTeam
				homeTeamID = int(s.AwayTeamID)
				homeTeamCoach = s.AwayTeamCoach
				homeTeamRank = int(s.AwayTeamRank)
				awayTeam = s.HomeTeam
				awayTeamID = int(s.HomeTeamID)
				awayTeamCoach = s.HomeTeamCoach
				awayTeamRank = int(s.HomeTeamRank)
				arena = ht.Arena
				city = ht.City
				state = ht.State
				country = ht.Country
			}
			nbaMatch := structs.NBAMatch{
				MatchName:       matchName,
				WeekID:          ts.NBAWeekID,
				Week:            uint(ts.NBAWeek),
				SeasonID:        ts.SeasonID,
				SeriesID:        s.ID,
				HomeTeamID:      uint(homeTeamID),
				HomeTeam:        homeTeam,
				HomeTeamCoach:   homeTeamCoach,
				HomeTeamRank:    uint(homeTeamRank),
				AwayTeamID:      uint(awayTeamID),
				AwayTeam:        awayTeam,
				AwayTeamCoach:   awayTeamCoach,
				AwayTeamRank:    uint(awayTeamRank),
				Arena:           arena,
				City:            city,
				State:           state,
				Country:         country,
				IsPlayoffGame:   s.IsPlayoffGame,
				IsTheFinals:     s.IsTheFinals,
				IsInternational: s.IsInternational,
				MatchOfWeek:     "A",
			}
			db.Create(&nbaMatch)
		}
	}
}

func IndicateWhetherTeamCanTradeInPostSeason() {
	db := dbprovider.GetInstance().GetDB()
	nbaTeamMap := GetProfessionalTeamMap()
	nbaIsInPlayoffsMap := make(map[uint]bool)
	activeNBASeries := GetAllActiveNBASeries()

	for i := 1; i < 33; i++ {
		idx := uint(i)
		nbaTeam := nbaTeamMap[idx]
		if nbaTeam.CanTrade {
			continue
		}
		nbaIsInPlayoffsMap[idx] = false
	}

	for _, s := range activeNBASeries {
		if s.SeriesComplete {
			continue
		}
		nbaIsInPlayoffsMap[s.HomeTeamID] = true
		nbaIsInPlayoffsMap[s.AwayTeamID] = true
	}

	for i := 1; i < 33; i++ {
		idx := uint(i)
		nbaTeam := nbaTeamMap[idx]
		if nbaTeam.CanTrade {
			continue
		}
		playoffsBool := nbaIsInPlayoffsMap[idx]
		if !playoffsBool {
			nbaTeam.ActivateTradeAbility()
			db.Save(&nbaTeam)
		}
	}
}

func ShowGames() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	matchType := ""
	if !ts.GamesARan {
		matchType = "A"
	} else if !ts.GamesBRan {
		matchType = "B"
	} else if !ts.GamesCRan {
		matchType = "C"
	} else if !ts.GamesDRan {
		matchType = "D"
	}
	if matchType == "" {
		log.Fatalln("Cannot sync results!")
	}
	UpdateStandings(ts, matchType)
	UpdateSeasonStats(ts, matchType)
	ts.ToggleGames(matchType)
	repository.SaveTimeStamp(ts, db)
}

func RegressGames(match string) {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	RegressStandings(ts, match)
	RegressSeasonStats(ts, match)
	ts.ToggleGamesARan()
	repository.SaveTimeStamp(ts, db)
}
