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

func MigrateNBAPlayersToTables() {
	db := dbprovider.GetInstance().GetDB()
	rand.Seed(time.Now().Unix())

	playersCSV := getPlayerData()

	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}
	LatestNBAPlayerID := lastPlayerRecord.ID + 1

	nbaTeams := GetAllActiveNBATeams()
	collegeTeams := GetAllActiveCollegeTeams()

	teamMap := make(map[string]uint)
	for _, team := range nbaTeams {
		teamMap[team.Abbr] = team.ID
	}

	for _, team := range collegeTeams {
		teamMap[team.Abbr] = team.ID
	}

	teamMap["OMAH"] = 143
	teamMap["Japan"] = 0
	teamMap["Australia"] = 0
	teamMap["Canada"] = 0
	teamMap["Colombia"] = 0
	teamMap["Croatia"] = 0
	teamMap["Czech Republic"] = 0
	teamMap["Germany"] = 0
	teamMap["Ghana"] = 0
	teamMap["Great Britain"] = 0
	teamMap["Jamaica"] = 0
	teamMap["Japan"] = 0
	teamMap["Kenya"] = 0
	teamMap["Lithuania"] = 0
	teamMap["Netherlands"] = 0
	teamMap["New Zealand"] = 0
	teamMap["Senegal"] = 0
	teamMap["Serbia"] = 0
	teamMap["Slovenia"] = 0
	teamMap["South Africa"] = 0
	teamMap["Spain"] = 0
	teamMap["Ukraine"] = 0

	for idx, row := range playersCSV {
		if idx < 2 {
			continue
		}
		teamAbbr := row[0]
		teamID := teamMap[teamAbbr]
		IsGLeague := util.ConvertStringToBool(row[1])
		IsTwoWay := util.ConvertStringToBool(row[2])
		IsWaived := util.ConvertStringToBool(row[3])
		Position := row[4]
		FirstName := row[5]
		LastName := row[6]
		Height := row[7]
		Age := util.ConvertStringToInt(row[8])
		YearsInNBA := util.ConvertStringToInt(row[9])
		College := row[10]
		CollegeID := teamMap[College]
		if College == "OMAH" {
			College = "UNOM"
		}
		shooting := util.ConvertStringToInt(row[12])
		Shooting2 := util.GenerateIntFromRange(shooting-3, shooting+3)
		diff := Shooting2 - shooting
		Shooting3 := shooting - diff
		finishing := util.ConvertStringToInt(row[13])
		freeThrow := util.GenerateIntFromRange(Shooting2-2, Shooting2+2)
		ballwork := util.ConvertStringToInt(row[15])
		rebounding := util.ConvertStringToInt(row[16])
		defense := util.ConvertStringToInt(row[17])
		interiorDefense := util.GenerateIntFromRange(defense-3, defense+3)
		diff = interiorDefense - defense
		perimeterDefense := defense - diff
		stamina := util.ConvertStringToInt(row[20])
		potentialGrade := row[21]
		NBAProgression := util.GetNBAProgressionRatingFromGrade(potentialGrade)
		overall := util.ConvertStringToInt(row[22])
		FirstAllNBA := util.ConvertStringToBool(row[24])
		DPOY := util.ConvertStringToBool(row[25])
		MVP := util.ConvertStringToBool(row[26])
		MaxRequested := util.ConvertStringToBool(row[27])
		SuperMax := util.ConvertStringToBool(row[28])
		ContractType := row[29]
		Year1Salary := util.ConvertStringToFloat(row[30])
		Year1Opt := util.ConvertStringToBool(row[31])
		Year2Salary := util.ConvertStringToFloat(row[32])
		Year2Opt := util.ConvertStringToBool(row[33])
		Year3Salary := util.ConvertStringToFloat(row[34])
		Year3Opt := util.ConvertStringToBool(row[35])
		Year4Salary := util.ConvertStringToFloat(row[36])
		Year4Opt := util.ConvertStringToBool(row[37])
		Year5Salary := util.ConvertStringToFloat(row[38])
		Year5Opt := util.ConvertStringToBool(row[39])
		RemainingContract := util.ConvertStringToFloat(row[40])
		YearsRemaining := util.ConvertStringToInt(row[41])
		InsidePercentage := util.ConvertStringToInt(row[42])
		MidPercentage := util.ConvertStringToInt(row[43])
		ThreePointPercentage := util.ConvertStringToInt(row[44])
		Position1 := row[46]
		Position1Minutes := util.ConvertStringToInt(row[47])
		Position2 := row[48]
		Position2Minutes := util.ConvertStringToInt(row[49])
		Position3 := row[50]
		Position3Minutes := util.ConvertStringToInt(row[51])
		IsRetiring := util.ConvertStringToBool(row[52])
		PrimeAge := util.ConvertStringToInt(row[53])

		// Check to see if player exists as a Historical Player
		// If so, use the player IDs from the draftee record for the player
		// Else, we will need a new one
		playerID := 0
		hpr := GetNBADrafteeByNameAndCollege(FirstName, LastName, College)
		if hpr.ID == 0 {
			playerID = int(LatestNBAPlayerID)
			LatestNBAPlayerID++
		} else {
			playerID = int(hpr.ID)
		}

		globalPlayer := structs.GlobalPlayer{
			NBAPlayerID:     uint(playerID),
			CollegePlayerID: uint(playerID),
		}
		globalPlayer.SetID(uint(playerID))

		basePlayer := structs.BasePlayer{
			FirstName:            FirstName,
			LastName:             LastName,
			Position:             Position,
			Age:                  Age,
			Year:                 YearsInNBA,
			Height:               Height,
			Shooting2:            Shooting2,
			Shooting3:            Shooting3,
			FreeThrow:            freeThrow,
			Finishing:            finishing,
			Ballwork:             ballwork,
			Rebounding:           rebounding,
			Defense:              defense,
			InteriorDefense:      interiorDefense,
			PerimeterDefense:     perimeterDefense,
			Potential:            NBAProgression,
			PotentialGrade:       potentialGrade,
			ProPotentialGrade:    NBAProgression,
			Stamina:              stamina,
			PlaytimeExpectations: 0,
			Minutes:              0,
			Overall:              overall,
		}

		nbaPlayer := structs.NBAPlayer{
			PlayerID:             uint(playerID),
			TeamID:               teamID,
			TeamAbbr:             teamAbbr,
			CollegeID:            CollegeID,
			College:              College,
			DraftPickID:          0,
			DraftPick:            0,
			DraftedTeamID:        teamID,
			DraftedTeamAbbr:      teamAbbr,
			PrimeAge:             uint(PrimeAge),
			IsNBA:                true,
			IsSuperMaxQualified:  SuperMax,
			IsFreeAgent:          false,
			IsGLeague:            IsGLeague,
			IsTwoWay:             IsTwoWay,
			IsWaived:             IsWaived,
			IsOnTradeBlock:       false,
			IsFirstTeamANBA:      FirstAllNBA,
			IsDPOY:               DPOY,
			IsMVP:                MVP,
			IsInternational:      false,
			IsRetiring:           IsRetiring,
			PositionOne:          Position1,
			PositionTwo:          Position2,
			PositionThree:        Position3,
			Position1Minutes:     uint(Position1Minutes),
			Position2Minutes:     uint(Position2Minutes),
			Position3Minutes:     uint(Position3Minutes),
			BasePlayer:           basePlayer,
			MaxRequested:         MaxRequested,
			InsidePercentage:     uint(InsidePercentage),
			MidPercentage:        uint(MidPercentage),
			ThreePointPercentage: uint(ThreePointPercentage),
		}

		nbaPlayer.SetID(uint(playerID))

		nbaContract := structs.NBAContract{
			PlayerID:       uint(playerID),
			TeamID:         teamID,
			Team:           teamAbbr,
			OriginalTeamID: teamID,
			OriginalTeam:   teamAbbr,
			YearsRemaining: uint(YearsRemaining),
			ContractType:   ContractType,
			TotalRemaining: RemainingContract,
			Year1Total:     Year1Salary,
			Year2Total:     Year2Salary,
			Year3Total:     Year3Salary,
			Year4Total:     Year4Salary,
			Year5Total:     Year5Salary,
			Year1Opt:       Year1Opt,
			Year2Opt:       Year2Opt,
			Year3Opt:       Year3Opt,
			Year4Opt:       Year4Opt,
			Year5Opt:       Year5Opt,
			IsDeadCap:      false,
			IsActive:       true,
			IsComplete:     false,
		}

		db.Create(&nbaPlayer)
		db.Create(&nbaContract)
		db.Create(&globalPlayer)
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
