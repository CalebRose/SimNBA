package managers

import (
	"encoding/csv"
	"log"
	"net/http"
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
	w.Header().Set("Content-Disposition", "attachment;filename=sagebows_secret_player_list.csv")
	w.Header().Set("Transfer-Encoding", "chunked")

	writer := csv.NewWriter(w)

	players := GetAllCollegePlayers()

	HeaderRow := []string{
		"College", "First Name", "Last Name", "Position", "Year", "Is_Redshirt", "Age",
		"Stars", "State", "Country", "Height",
		"Overall", "Shooting 2s", "Shooting 3s", "Free Throwing", "Finishing",
		"Ballwork", "Rebounding", "InteriorDefense", "PerimeterDefense", "Potential Grade",
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

		playerRow := []string{
			player.TeamAbbr, player.FirstName, player.LastName, player.Position, strconv.Itoa(player.Year), strconv.FormatBool(player.IsRedshirt), strconv.Itoa(player.Age),
			strconv.Itoa(player.Stars), player.State, player.Country, player.Height,
			overallGrade, shooting2Grade, shooting3Grade, freeThrowGrade, finishingGrade,
			ballworkGrade, reboundingGrade, interiorDefenseGrade, perimeterDefenseGrade, potentialGrade,
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

		db.Order("minutes desc").Where("team_id = ?", teamID).Find(&players)

		offenseRank := 0.0
		defenseRank := 0.0
		overall := 0.0
		starPower := 0.0
		count := 0
		for _, player := range players {
			if player.Minutes == 0 {
				continue
			}
			o := ((float64(player.Shooting2) + float64(player.Shooting3) + float64(player.Finishing) + float64(player.FreeThrow)) / 4) * (float64(player.Minutes) / float64(player.Stamina))
			d := ((float64(player.Ballwork) + float64(player.Rebounding) + float64(player.InteriorDefense) + float64(player.PerimeterDefense)) / 4) * (float64(player.Minutes) / float64(player.Stamina))
			offenseRank += o
			defenseRank += d
			starPower += float64(player.Stars)
			count++
		}

		if count == 0 {
			db.Order("overall desc").Where("team_id = ?", teamID).Find(&players)

			for idx, player := range players {
				if idx > 8 {
					break
				}
				o := ((float64(player.Shooting2) + float64(player.Shooting3) + float64(player.Finishing) + float64(player.FreeThrow)) / 4) * 0.87
				d := ((float64(player.Ballwork) + float64(player.Rebounding) + float64(player.InteriorDefense) + float64(player.PerimeterDefense)) / 4) * 0.87
				offenseRank += o
				defenseRank += d
				starPower += float64(player.Stars)
				count++
			}
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
