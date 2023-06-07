package managers

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		"College", "First Name", "Last Name", "Position",
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
			player.TeamAbbr, player.FirstName, player.LastName, player.Position,
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
