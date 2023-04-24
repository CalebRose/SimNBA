package managers

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func MigrateOldPlayerDataToNewTables() {
	db := dbprovider.GetInstance().GetDB()
	rand.Seed(time.Now().Unix())

	Players := GetAllCollegePlayersFromOldTable()

	for _, player := range Players {

		shooting := player.Shooting

		Shooting2 := util.GenerateIntFromRange(shooting-3, shooting+3)
		diff := Shooting2 - shooting
		Shooting3 := shooting - diff

		personality := util.GetPersonality()
		academicBias := util.GetAcademicBias()
		workEthic := util.GetWorkEthic()
		recruitingBias := util.GetRecruitingBias()
		freeAgency := util.GetFreeAgencyBias()

		abbr := ""
		teamId := 0

		var base = structs.BasePlayer{
			FirstName:            player.FirstName,
			LastName:             player.LastName,
			Position:             player.Position,
			Age:                  player.Age,
			Year:                 player.Year,
			State:                player.State,
			Country:              player.Country,
			Stars:                player.Stars,
			Height:               player.Height,
			Shooting2:            Shooting2,
			Shooting3:            Shooting3,
			Finishing:            player.Finishing,
			Ballwork:             player.Ballwork,
			Rebounding:           player.Rebounding,
			Defense:              player.Defense,
			Potential:            player.PotentialGrade,
			ProPotentialGrade:    player.ProPotentialGrade,
			Stamina:              player.Stamina,
			PlaytimeExpectations: player.PlaytimeExpectations,
			Minutes:              player.MinutesA,
			Overall:              player.Overall,
			Personality:          personality,
			FreeAgency:           freeAgency,
			RecruitingBias:       recruitingBias,
			WorkEthic:            workEthic,
			AcademicBias:         academicBias,
		}

		if teamId != player.TeamID {
			teamId = player.TeamID
			team := GetTeamByTeamID(strconv.Itoa(teamId))
			abbr = team.Abbr

			var collegePlayer = structs.CollegePlayer{
				BasePlayer:    base,
				PlayerID:      player.ID,
				TeamID:        uint(player.TeamID),
				TeamAbbr:      abbr,
				IsRedshirt:    player.IsRedshirt,
				IsRedshirting: player.IsRedshirting,
				HasGraduated:  false,
			}

			collegePlayer.SetID(player.ID)

			err := db.Save(&collegePlayer).Error
			if err != nil {
				log.Fatal("Could not save College Player " + player.FirstName + " " + player.LastName + " " + abbr)
			}
		} else {
			var recruit = structs.Recruit{
				BasePlayer: base,
				PlayerID:   player.ID,
				IsTransfer: true,
			}

			recruit.SetID(player.ID)

			err := db.Save(&recruit).Error
			if err != nil {
				log.Fatal("Could not save College Transfer " + player.FirstName + " " + player.LastName + " " + abbr)
			}
		}

		var globalPlayer = structs.GlobalPlayer{
			CollegePlayerID: player.ID,
			RecruitID:       player.ID,
			NBAPlayerID:     player.ID,
		}

		globalPlayer.SetID(player.ID)

		err := db.Save(&globalPlayer).Error
		if err != nil {
			log.Fatal("Could not save global record for College Player " + player.FirstName + " " + player.LastName + " " + abbr)
		}
	}
}

func MigrateRecruits() {
	db := dbprovider.GetInstance().GetDB()

	croots := GetAllRecruitRecords()

	for _, croot := range croots {
		// Convert to College Player Record
		cp := structs.CollegePlayer{}
		cp.MapFromRecruit(croot)

		// Save College Player Record
		err := db.Create(&cp).Error
		if err != nil {
			log.Panicln("Could not save new college player record")
		}

		// Delete Recruit Record
		db.Delete(&croot)
	}
}

func getPlayerData() [][]string {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimNBA\\data\\SimNBA_Players_2022.csv"
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+path, err)
	}

	return records
}
