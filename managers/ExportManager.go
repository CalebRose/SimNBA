package managers

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		"Ballwork", "Rebounding", "Defense", "Potential Grade",
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
			croot.Ballwork, croot.Rebounding, croot.Defense, croot.PotentialGrade,
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
