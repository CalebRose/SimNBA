package managers

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func MigrateOldPlayerDataToNewTables() {

	db := dbprovider.GetInstance().GetDB()

	Players := GetAllCollegePlayers()

	for _, player := range Players {

		shooting := player.Shooting

		Shooting2 := util.GenerateIntFromRange(shooting-3, shooting+3)
		diff := Shooting2 - shooting
		Shooting3 := shooting + diff

		personality := util.GetPersonality()
		academicBias := util.GetAcademicBias()
		workEthic := util.GetWorkEthic()
		recruitingBias := util.GetRecruitingBias()
		freeAgency := util.GetFreeAgencyBias()

		abbr := ""
		teamId := 0

		if teamId != player.TeamID {
			teamId = player.TeamID
			team := GetTeamByTeamID(strconv.Itoa(teamId))
			abbr = team.Abbr
		}

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

		var collegePlayer = structs.CollegePlayer{
			BasePlayer:    base,
			PlayerID:      int(player.ID),
			TeamID:        player.TeamID,
			TeamAbbr:      abbr,
			IsRedshirt:    player.IsRedshirt,
			IsRedshirting: player.IsRedshirting,
			HasGraduated:  false,
		}

		var globalPlayer = structs.GlobalPlayer{
			CollegePlayerID: int(player.ID),
		}

		err := db.Save(&collegePlayer).Error
		if err != nil {
			log.Fatal("Could not save College Player " + player.FirstName + " " + player.LastName + " " + abbr)
		}

		err = db.Save(&globalPlayer).Error
		if err != nil {
			log.Fatal("Could not save global record for College Player " + player.FirstName + " " + player.LastName + " " + abbr)
		}
	}
}
