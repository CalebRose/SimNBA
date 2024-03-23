package managers

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"github.com/jinzhu/gorm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CrootGenerator struct {
	firstNameMap      map[string][][]string
	lastNameMap       map[string][][]string
	collegePlayerList []structs.CollegePlayer
	coachList         []structs.CollegeCoach
	teamMap           map[uint]structs.Team
	positionList      []string
	CrootList         []structs.Recruit
	GlobalList        []structs.GlobalPlayer
	newID             uint
	count             int
	requiredPlayers   int
	star5             int
	star4             int
	star3             int
	star2             int
	star1             int
	ovr70             int
	ovr60             int
	ovr50             int
	ovr40             int
	ovr30             int
	highestOvr        int
	lowestOvr         int
	pickedEthnicity   string
	caser             cases.Caser
}

func (pg *CrootGenerator) GenerateRecruits() {
	for pg.count < pg.requiredPlayers {
		player, globalPlayer := pg.generatePlayer()
		pg.CrootList = append(pg.CrootList, player)
		pg.GlobalList = append(pg.GlobalList, globalPlayer)
		pg.updateStatistics(player) // A method to update player counts and statistics
		if player.RelativeType == 5 {
			twinPlayer, twinGlobal := pg.generateTwin(&player)
			pg.updateStatistics(twinPlayer)
			pg.CrootList = append(pg.CrootList, twinPlayer)
			pg.GlobalList = append(pg.GlobalList, twinGlobal)
			pg.count++
		}
		pg.count++
		pg.newID++
	}
}

func (pg *CrootGenerator) generatePlayer() (structs.Recruit, structs.GlobalPlayer) {
	cpLen := len(pg.collegePlayerList) - 1
	coachLen := len(pg.coachList)
	relativeType := 0
	relativeID := 0
	coachTeamID := 0
	coachTeamAbbr := ""
	notes := ""
	firstName := ""
	lastName := ""
	state := ""
	country := ""
	pg.pickedEthnicity = pickEthnicity()
	firstNameList := pg.firstNameMap[pg.pickedEthnicity]
	lastNameList := pg.lastNameMap[pg.pickedEthnicity]
	fName := getName(firstNameList)
	firstName = pg.caser.String(strings.ToLower(fName))
	// Roll for type of recruit generated
	// If num == 200, then create some flair
	roof := 100
	relativeRoll := util.GenerateIntFromRange(1, roof)
	relativeIdx := 0
	if relativeRoll == roof {
		relativeType = getRelativeType()
		if relativeType == 2 {
			// Brother of college player
			fmt.Println("BROTHER")
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			lastName = cp.LastName
			state = cp.State
			country = cp.Country
			notes = "Brother of " + cp.TeamAbbr + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 3 {
			fmt.Println("COUSIN")
			// Cousin
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			coinFlip := util.GenerateIntFromRange(1, 2)
			if coinFlip == 1 {
				lastName = cp.LastName
			} else {
				lName := getName(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			country = cp.Country
			notes = "Cousin of " + cp.TeamAbbr + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 4 {
			// Half Brother
			fmt.Println("HALF BROTHER GENERATED")
			relativeIdx = util.GenerateIntFromRange(0, cpLen)
			if relativeIdx < 0 || relativeIdx > len(pg.collegePlayerList) {
				relativeIdx = util.GenerateIntFromRange(0, cpLen)
			}
			cp := pg.collegePlayerList[relativeIdx]
			relativeID = int(cp.ID)
			coinFlip := util.GenerateIntFromRange(1, 3)
			if coinFlip < 3 {
				lastName = cp.LastName
			} else {
				lName := getName(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			country = cp.Country
			notes = "Half-Brother of " + cp.TeamAbbr + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		} else if relativeType == 5 {
			// Twin
			relativeType = 5
			relativeID = int(pg.newID)
		} else if relativeType == 6 {
			// Coach's Son
			fmt.Println("COACH'S SON")
			relativeIdx = util.GenerateIntFromRange(0, coachLen)
			if relativeIdx < 0 || relativeIdx > len(pg.coachList) {
				relativeIdx = util.GenerateIntFromRange(0, coachLen)
			}
			coach := pg.coachList[relativeIdx]
			relativeID = int(coach.ID)
			lastName = getCoachLastName(coach.Name)
			team := pg.teamMap[coach.TeamID]
			state = team.State
			country = "USA"
			notes = "Son of Coach " + coach.Name + " of " + coach.Team
			coachTeamID = int(coach.TeamID)
			coachTeamAbbr = coach.Team
		} else if relativeType == 7 {
			// Coach's Nephew
			fmt.Println("COACH'S NEPHEW")
			relativeIdx = util.GenerateIntFromRange(0, coachLen)
			if relativeIdx < 0 || relativeIdx > len(pg.coachList) {
				relativeIdx = util.GenerateIntFromRange(0, coachLen)
			}
			coach := pg.coachList[relativeIdx]
			relativeID = int(coach.ID)
			coachLastName := getCoachLastName(coach.Name)
			coinFlip := util.GenerateIntFromRange(1, 2)
			if coinFlip == 1 {
				lastName = coachLastName
			} else {
				lName := getName(lastNameList)
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			team := pg.teamMap[coach.TeamID]
			state = team.State
			country = "USA"
			notes = "Son of Coach " + coach.Name + " of " + coach.Team
			coachTeamID = int(coach.TeamID)
			coachTeamAbbr = coach.Team
		}
	} else {
		relativeType = 1
	}
	if relativeType == 1 || relativeType == 5 {
		lName := getName(lastNameList)
		lastName = pg.caser.String(strings.ToLower(lName))
		state = ""
		country = pickCountry(pg.pickedEthnicity)
		if country == "USA" {
			state = pickState()
		}
	}
	pickedPosition := util.PickFromStringList(pg.positionList)
	year := 1
	player := createRecruit(firstName, lastName, state, country, pg.pickedEthnicity, pickedPosition, year, pg.newID)
	player.AssignRelativeData(uint(relativeID), uint(relativeType), uint(coachTeamID), coachTeamAbbr, notes)
	globalPlayer := structs.GlobalPlayer{
		CollegePlayerID: pg.newID,
		RecruitID:       pg.newID,
		NBAPlayerID:     pg.newID,
	}

	globalPlayer.SetID(pg.newID)
	return player, globalPlayer
}

func (pg *CrootGenerator) generateTwin(player *structs.Recruit) (structs.Recruit, structs.GlobalPlayer) {
	fmt.Println("TWIN!!")
	// Generate Twin Record
	relativeID := int(pg.newID)
	pg.newID++
	twinRelativeID := relativeID
	relativeID = int(pg.newID)
	firstNameList := pg.firstNameMap[pg.pickedEthnicity]
	twinName := getName(firstNameList)
	twinN := pg.caser.String(strings.ToLower(twinName))
	twinPosition := ""
	if player.Position == "PF" {
		twinPosition = util.PickFromStringList([]string{"C", "SF"})
	} else if player.Position == "C" {
		twinPosition = "PF"
	} else if player.Position == "SF" {
		twinPosition = util.PickFromStringList([]string{"PF", "SG"})
	} else if player.Position == "SG" {
		twinPosition = util.PickFromStringList([]string{"SF", "PG"})
	} else {
		twinPosition = "SG"
	}
	twinNotes := "Twin Brother of " + strconv.Itoa(player.Stars) + " Star Recruit " + player.Position + " " + player.FirstName + " " + player.LastName
	twinPlayer := createRecruit(twinN, player.LastName, player.State, player.Country, pg.pickedEthnicity, twinPosition, player.Year, pg.newID)
	twinPlayer.AssignRelativeData(uint(twinRelativeID), 4, 0, "", twinNotes)
	notes := "Twin Brother of " + strconv.Itoa(twinPlayer.Stars) + " Star Recruit " + twinPlayer.Position + " " + twinPlayer.FirstName + " " + twinPlayer.LastName
	player.AssignRelativeData(uint(relativeID), 4, 0, "", notes)
	globalTwinPlayer := structs.GlobalPlayer{
		CollegePlayerID: pg.newID,
		RecruitID:       pg.newID,
		NBAPlayerID:     pg.newID,
	}
	globalTwinPlayer.SetID(pg.newID)
	player.AssignRelativeData(uint(relativeID), player.RelativeType, 0, "", notes)
	globalPlayer := structs.GlobalPlayer{
		CollegePlayerID: pg.newID,
		RecruitID:       pg.newID,
		NBAPlayerID:     pg.newID,
	}

	globalPlayer.SetID(pg.newID)
	return twinPlayer, globalPlayer
}

func (pg *CrootGenerator) updateStatistics(player structs.Recruit) {
	if player.Stars == 5 {
		pg.star5++
	} else if player.Stars == 4 {
		pg.star4++
	} else if player.Stars == 3 {
		pg.star3++
	} else if player.Stars == 2 {
		pg.star2++
	} else {
		pg.star1++
	}

	if player.Overall >= 70 {
		pg.ovr70++
	} else if player.Overall >= 60 {
		pg.ovr60++
	} else if player.Overall >= 50 {
		pg.ovr50++
	} else if player.Overall >= 40 {
		pg.ovr40++
	} else if player.Overall >= 30 {
		pg.ovr30++
	}

	if player.Overall > pg.highestOvr {
		pg.highestOvr = player.Overall
	}
	if player.Overall < pg.lowestOvr {
		pg.lowestOvr = player.Overall
	}
}

func (pg *CrootGenerator) OutputRecruitStats() {
	// Croot Stats
	fmt.Println("Total Recruit Count: ", pg.count)
	fmt.Println("Total Ovr 70  Count: ", pg.ovr70)
	fmt.Println("Total Ovr 60  Count: ", pg.ovr60)
	fmt.Println("Total Ovr 50  Count: ", pg.ovr50)
	fmt.Println("Total Ovr 40  Count: ", pg.ovr40)
	fmt.Println("Total Ovr 30  Count: ", pg.ovr30)
	fmt.Println("Total 5 Star  Count: ", pg.star5)
	fmt.Println("Total 4 Star  Count: ", pg.star4)
	fmt.Println("Total 3 Star  Count: ", pg.star3)
	fmt.Println("Total 2 Star  Count: ", pg.star2)
	fmt.Println("Total 1 Star  Count: ", pg.star1)
	fmt.Println("Highest Recruit Ovr: ", pg.highestOvr)
	fmt.Println("Lowest  Recruit Ovr: ", pg.lowestOvr)
}

func GenerateCoachesForAITeams() {
	db := dbprovider.GetInstance().GetDB()

	teams := GetOnlyAITeamRecruitingProfiles()
	firstNameMap, lastNameMap := getNameMaps()

	coachList := []structs.CollegeCoach{}
	allActiveCoaches := GetAllCollegeCoaches()
	collegeTeamMap := GetCollegeTeamMap()
	retiredPlayers := GetAllRetiredPlayers()
	retireeMap := make(map[uint]bool)
	coachMap := make(map[uint]bool)

	for _, coach := range allActiveCoaches {
		if coach.FormerPlayerID > 0 {
			coachMap[coach.FormerPlayerID] = true
		}
	}

	for _, team := range teams {
		// Skip over teams currently controlled by a user
		if !team.IsAI || team.ID != 368 {
			continue
		}

		pickedEthnicity := pickEthnicity()
		almaMater := pickAlmaMater(teams)
		coach := createCollegeCoach(team, almaMater.ID, almaMater.TeamAbbr, pickedEthnicity, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], retiredPlayers, &retireeMap, &coachMap)
		coachList = append(coachList, coach)
		collegeTeam := collegeTeamMap[team.TeamID]
		collegeTeam.AssignCoach(coach.Name)
		team.UpdateAIBehavior(true, true, coach.StarMax, coach.StarMin, coach.PointMin, coach.PointMax, coach.Scheme, coach.DefensiveScheme)
		team.AssignRecruiter(coach.Name)
		db.Save(&collegeTeam)
		db.Save(&team)
	}

	for _, coach := range coachList {
		db.Create(&coach)
	}
}

func GenerateTestPlayersForTP() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	newID := lastPlayerRecord.ID + 1
	firstNameMap, lastNameMap := getNameMaps()
	var positionList []string = []string{"PG", "SG", "PF", "SF", "C"}

	for i := 0; i < 15; i++ {
		pickedEthnicity := pickEthnicity()
		pickedPosition := util.PickFromStringList(positionList)
		year := util.GenerateIntFromRange(1, 3)
		emptyTeam := structs.Team{}
		player := createCollegePlayer(emptyTeam, pickedEthnicity, pickedPosition, year, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], newID)
		// playerList = append(playerList, player)
		err = db.Create(&player).Error
		if err != nil {
			log.Panicln("Could not save player record")
		}

		globalPlayer := structs.GlobalPlayer{
			Model:           gorm.Model{ID: newID},
			CollegePlayerID: newID,
			RecruitID:       newID,
			NBAPlayerID:     newID,
		}

		db.Create(&globalPlayer)
		newID++
	}
}

func GenerateNewTeams() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	// var playerList []structs.CollegePlayer

	newID := lastPlayerRecord.ID + 1

	teams := GetAllActiveCollegeTeams()
	firstNameMap, lastNameMap := getNameMaps()
	var positionList []string = []string{"PG", "SG", "PF", "SF", "C"}

	for _, team := range teams {
		if team.ID != 367 {
			continue
		}
		// Test Generation
		yearList := []int{}
		players := GetCollegePlayersByTeamId(strconv.Itoa(int(team.ID)))
		if len(players) > 0 {
			continue
		}
		seniors := 3
		juniors := 3
		sophomores := 3
		freshmen := 3
		sfCount := 0
		pfCount := 0
		pgCount := 0
		sgCount := 0
		cCount := 0
		requiredPlayers := 12
		count := 0

		requiredPlayers -= len(players)
		for _, player := range players {
			if player.Year == 4 && !player.IsRedshirt {
				seniors--
			} else if player.Year == 3 || (player.Year == 4 && player.IsRedshirt) {
				juniors--
			} else if player.Year == 2 || (player.Year == 3 && player.IsRedshirt) {
				sophomores--
			} else {
				freshmen--
			}
			if player.Position == "PF" {
				pfCount++
			} else if player.Position == "SF" {
				sfCount++
			} else if player.Position == "PG" {
				pgCount++
			} else if player.Position == "SG" {
				sgCount++
			} else {
				cCount++
			}
		}

		for i := 0; i < seniors; i++ {
			yearList = append(yearList, 4)
		}
		for i := 0; i < juniors; i++ {
			yearList = append(yearList, 3)
		}
		for i := 0; i < sophomores; i++ {
			yearList = append(yearList, 2)
		}
		for i := 0; i < freshmen; i++ {
			yearList = append(yearList, 1)
		}
		var positionQueue []string
		for i := 0; i < requiredPlayers; i++ {
			pickedPosition := util.PickFromStringList(positionList)
			if pickedPosition == "PF" && pfCount > 2 {
				quickList := []string{"PG", "SF", "SG", "C"}
				pickedPosition = util.PickFromStringList(quickList)
			} else if pickedPosition == "SF" && sfCount > 2 {
				quickList := []string{"PG", "PF", "SG", "C"}
				pickedPosition = util.PickFromStringList(quickList)
			} else if pickedPosition == "PG" && pgCount > 2 {
				quickList := []string{"PF", "SG", "SF", "C"}
				pickedPosition = util.PickFromStringList(quickList)
			} else if pickedPosition == "SG" && pgCount > 2 {
				quickList := []string{"PF", "PG", "SF", "C"}
				pickedPosition = util.PickFromStringList(quickList)
			} else if pickedPosition == "C" && cCount > 2 {
				quickList := []string{"SF", "PG", "SG", "PF"}
				pickedPosition = util.PickFromStringList(quickList)
			}

			if pickedPosition == "SF" {
				sfCount++
			} else if pickedPosition == "PF" {
				pfCount++
			} else if pickedPosition == "PG" {
				pgCount++
			} else if pickedPosition == "SG" {
				sgCount++
			} else {
				cCount++
			}

			positionQueue = append(positionQueue, pickedPosition)
		}

		rand.Shuffle(len(positionQueue), func(i, j int) {
			positionQueue[i], positionQueue[j] = positionQueue[j], positionQueue[i]
		})

		for count < requiredPlayers {
			pickedEthnicity := pickEthnicity()
			pickedPosition := positionQueue[count]
			year := yearList[count]
			player := createCollegePlayer(team, pickedEthnicity, pickedPosition, year, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], newID)
			// playerList = append(playerList, player)
			err := db.Create(&player).Error
			if err != nil {
				log.Panicln("Could not save player record")
			}

			globalPlayer := structs.GlobalPlayer{
				Model:           gorm.Model{ID: newID},
				CollegePlayerID: newID,
				RecruitID:       newID,
				NBAPlayerID:     newID,
			}

			db.Create(&globalPlayer)

			count++
			newID++
		}

	}
	// return playerList
}

func GenerateGlobalPlayerRecords() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}
	collegePlayerID := lastPlayerRecord.ID + 1

	for collegePlayerID <= 2476 {
		player := GetCollegePlayerByPlayerID(strconv.Itoa(int(collegePlayerID)))

		if player.ID > 0 {
			var globalPlayer structs.GlobalPlayer

			err := db.Where("id = ?", strconv.Itoa(int(player.ID))).Find(&globalPlayer).Error
			if err != nil {
				// Check
				fmt.Println("Record does not exist")
			}

			if globalPlayer.ID == 0 {
				globalPlayer = structs.GlobalPlayer{
					RecruitID:       collegePlayerID,
					CollegePlayerID: collegePlayerID,
					NBAPlayerID:     collegePlayerID,
				}
				globalPlayer.SetID(collegePlayerID)
			}

			err = db.Create(&globalPlayer).Error
			if err != nil {
				// Figure it out.
				log.Fatalln(err.Error())
			}
		}
		collegePlayerID++
	}
}

func GenerateCroots() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer
	ts := GetTimestamp()

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	// var playerList []structs.CollegePlayer
	fNameMap, lNameMap := getNameMaps()
	generator := CrootGenerator{
		firstNameMap:      fNameMap,
		lastNameMap:       lNameMap,
		collegePlayerList: GetAllCollegePlayers(),
		coachList:         GetAllActiveCollegeCoaches(),
		teamMap:           GetCollegeTeamMap(),
		positionList:      []string{"PG", "SG", "PF", "SF", "C"},
		newID:             lastPlayerRecord.ID + 1,
		requiredPlayers:   util.GenerateIntFromRange(1031, 1101),
		count:             0,
		star5:             0,
		star4:             0,
		star3:             0,
		star2:             0,
		star1:             0,
		ovr70:             0,
		ovr60:             0,
		ovr50:             0,
		ovr40:             0,
		ovr30:             0,
		highestOvr:        0,
		lowestOvr:         100000,
		CrootList:         []structs.Recruit{},
		GlobalList:        []structs.GlobalPlayer{},
		caser:             cases.Title(language.English),
		pickedEthnicity:   "",
	}

	// Test Generation
	// requiredPlayers := util.GenerateIntFromRange(203, 205)
	// 1061 is the number of open spots on teams in the league.
	// Currently 363 teams. 363 * 3 = 1089, the size of the average class.
	// The plan is to ensure that every recruit is signed
	generator.GenerateRecruits()
	// Croot Stats
	generator.OutputRecruitStats()

	// Import Batches
	for _, r := range generator.CrootList {
		repository.CreateRecruitRecord(r, db)
	}

	for _, g := range generator.GlobalList {
		repository.CreateGlobalPlayerRecord(g, db)
	}

	ts.ToggleGeneratedCroots()
	repository.SaveTimeStamp(ts, db)
	// return playerList

	AssignAllRecruitRanks()
}

func GenerateInternationalPlayers() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	// var playerList []structs.CollegePlayer

	newID := lastPlayerRecord.ID + 1

	firstNameMap, lastNameMap := getNameMaps()
	var positionList []string = []string{"PG", "SG", "PF", "SF", "C"}

	// Get all ISL teams
	allProfessionalTeams := GetAllActiveNBATeams()

	for _, team := range allProfessionalTeams {
		// If an NBA team, skip
		if team.LeagueID == 1 {
			continue
		}
		count := 0
		teamID := strconv.Itoa(int(team.ID))
		currentPlayers := GetAllNBAPlayersByTeamID(teamID)
		requiredPlayers := 13 - len(currentPlayers)

		teamNeedsMap := make(map[string]bool)
		positionCount := make(map[string]int)

		if _, ok := teamNeedsMap["PG"]; !ok {
			teamNeedsMap["PG"] = true
		}
		if _, ok := teamNeedsMap["SG"]; !ok {
			teamNeedsMap["SG"] = true
		}
		if _, ok := teamNeedsMap["SF"]; !ok {
			teamNeedsMap["SF"] = true
		}
		if _, ok := teamNeedsMap["PF"]; !ok {
			teamNeedsMap["PF"] = true
		}
		if _, ok := teamNeedsMap["C"]; !ok {
			teamNeedsMap["C"] = true
		}

		if _, ok := positionCount["PG"]; !ok {
			positionCount["PG"] = 0
		}
		if _, ok := positionCount["SG"]; !ok {
			positionCount["SG"] = 0
		}
		if _, ok := positionCount["SF"]; !ok {
			positionCount["SF"] = 0
		}
		if _, ok := positionCount["PF"]; !ok {
			positionCount["PF"] = 0
		}
		if _, ok := positionCount["C"]; !ok {
			positionCount["C"] = 0
		}

		for _, r := range currentPlayers {
			positionCount[r.Position] += 1
		}

		if positionCount["PG"] >= 3 {
			teamNeedsMap["PG"] = false
		} else if positionCount["SG"] >= 4 {
			teamNeedsMap["SG"] = false
		} else if positionCount["SF"] >= 4 {
			teamNeedsMap["SF"] = false
		} else if positionCount["PF"] >= 4 {
			teamNeedsMap["PF"] = false
		} else if positionCount["C"] >= 3 {
			teamNeedsMap["C"] = false
		}

		positionNeedList := []string{}

		for _, pos := range positionList {
			if (pos == "PG" || pos == "C") && positionCount[pos] >= 3 {
				continue
			}
			if (pos == "SG" || pos == "SF" || pos == "PF") && positionCount[pos] >= 4 {
				continue
			}
			positionNeedList = append(positionNeedList, pos)
		}

		// Generate two international players from the team's host country
		for count < requiredPlayers {
			pickedPosition := util.PickFromStringList(positionNeedList)
			pickedEthnicity := pickISLEthnicity(team.Country)
			year := 1
			player := createInternationalPlayer(team.ID, team.Team, team.Country, pickedEthnicity, pickedPosition, year, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], newID)

			year1Salary := 0.0
			year2Salary := 0.0
			year3Salary := 0.0
			year4Salary := 0.0
			year5Salary := 0.0
			yearsRemaining := 2
			if player.Age < 22 {
				yearsRemaining = 22 - player.Age
				if yearsRemaining > 5 {
					yearsRemaining = 5
				}
			}
			for i := 1; i <= yearsRemaining; i++ {
				if i == 1 {
					year1Salary = 0.5
				}
				if i == 2 {
					year2Salary = 0.5
				}
				if i == 3 {
					year3Salary = 0.5
				}
				if i == 4 {
					year4Salary = 0.5
				}
				if i == 5 {
					year5Salary = 0.5
				}
			}
			contract := structs.NBAContract{
				PlayerID:       player.PlayerID,
				TeamID:         player.TeamID,
				Team:           player.TeamAbbr,
				OriginalTeamID: player.TeamID,
				OriginalTeam:   player.TeamAbbr,
				YearsRemaining: uint(yearsRemaining),
				ContractType:   "International",
				TotalRemaining: year1Salary + year2Salary + year3Salary + year4Salary + year5Salary,
				Year1Total:     year1Salary,
				Year2Total:     year2Salary,
				Year3Total:     year3Salary,
				Year4Total:     year4Salary,
				Year5Total:     year5Salary,
				IsActive:       true,
			}

			err := db.Create(&player).Error
			if err != nil {
				log.Panicln("Could not save player record")
			}

			err = db.Create(&contract).Error
			if err != nil {
				log.Panicln("Could not save player record")
			}

			globalPlayer := structs.GlobalPlayer{
				CollegePlayerID: newID,
				RecruitID:       newID,
				NBAPlayerID:     newID,
			}

			globalPlayer.SetID(newID)

			db.Create(&globalPlayer)

			count++
			newID++
		}
	}
	// return playerList
}

func CleanUpRecruits() {
	db := dbprovider.GetInstance().GetDB()

	croots := GetAllUnsignedRecruits()

	for _, croot := range croots {
		if croot.PotentialGrade != "" && croot.ProPotentialGrade > 0 && croot.RecruitModifier > 0 {
			continue
		}
		potential := ""
		proPotential := 0
		recruitMod := 0
		if croot.PotentialGrade == "" {
			potential = util.GetWeightedPotentialGrade(croot.Potential)
		}

		if croot.ProPotentialGrade == 0 {
			proPotential = util.GenerateIntFromRange(1, 100)
		}

		if croot.RecruitModifier == 0 {
			recruitMod = GetRecruitModifier(croot.Stars)
		}

		croot.FixRecruit(potential, proPotential, recruitMod)

		err := db.Save(&croot).Error
		if err != nil {
			log.Panicln(err.Error())
		}
	}
}

func GenerateAttributeSpecs() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()
	croots := GetAllRecruitRecords()
	nbaPlayers := GetAllNBAPlayers()

	for _, cp := range collegePlayers {
		// Specialties
		specs := util.GetSpecialties(cp.Position)
		for _, spec := range specs {
			cp.ToggleSpecialties(spec)
		}
		if len(specs) > 0 {
			db.Save(&cp)
		}
	}

	for _, r := range croots {
		specs := util.GetSpecialties(r.Position)
		for _, spec := range specs {
			r.ToggleSpecialties(spec)
		}
		if len(specs) > 0 {
			db.Save(&r)
		}
	}

	for _, n := range nbaPlayers {
		specs := util.GetSpecialties(n.Position)
		for _, spec := range specs {
			n.ToggleSpecialties(spec)
		}
		if len(specs) > 0 {
			db.Save(&n)
		}
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

func GenerateDraftWarRooms() {
	db := dbprovider.GetInstance().GetDB()

	allProfessionalTeams := GetAllActiveNBATeams()

	for _, n := range allProfessionalTeams {
		if n.League != "SimNBA" {
			continue
		}
		room := GetNBAWarRoomByTeamID(strconv.Itoa(int(n.ID)))
		if room.ID > 0 {
			continue
		}
		warRoom := structs.NBAWarRoom{
			TeamID:         n.ID,
			Team:           n.Team + " " + n.Nickname,
			ScoutingPoints: 100,
			SpentPoints:    0,
		}
		db.Create(&warRoom)
	}
}

func GenerateNewAttributes() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()

	for _, c := range collegePlayers {
		if c.ID == 0 {
			continue
		}
		discipline := util.GenerateIntFromRange(1, 20)
		injuryRating := util.GenerateIntFromRange(1, 20)
		c.SetDisciplineAndIR(discipline, injuryRating)
		repository.SaveCollegePlayerRecord(c, db)
	}

	nbaPlayers := GetAllNBAPlayers()

	for _, n := range nbaPlayers {
		discipline := util.GenerateIntFromRange(1, 20)
		injuryRating := util.GenerateIntFromRange(1, 20)
		n.SetDisciplineAndIR(discipline, injuryRating)

		repository.SaveProfessionalPlayerRecord(n, db)
	}
}

// Private Methods
func createCollegePlayer(team structs.Team, ethnicity string, position string, year int, firstNameList [][]string, lastNameList [][]string, id uint) structs.CollegePlayer {
	fName := getName(firstNameList)
	lName := getName(lastNameList)
	caser := cases.Title(language.English)

	firstName := caser.String(strings.ToLower(fName))
	lastName := caser.String(strings.ToLower(lName))
	state := ""
	country := pickCountry(ethnicity)
	if country == "USA" {
		state = pickState()
	}
	height := getHeight(position)
	potential := util.GeneratePotential()
	proPotential := util.GeneratePotential()
	stamina := util.GenerateIntFromRange(25, 38)
	shooting2 := getAttribute(position, "Shooting2", true)
	shooting3 := getAttribute(position, "Shooting3", true)
	finishing := getAttribute(position, "Finishing", true)
	freeThrow := getAttribute(position, "FreeThrow", true)
	ballwork := getAttribute(position, "Ballwork", true)
	rebounding := getAttribute(position, "Rebounding", true)
	interiorDefense := getAttribute(position, "Interior Defense", true)
	perimeterDefense := getAttribute(position, "Perimeter Defense", true)

	overall := (int((shooting2 + shooting3 + freeThrow) / 3)) + finishing + ballwork + rebounding + int((interiorDefense+perimeterDefense)/2)
	stars := getStarRating(overall)

	expectations := util.GetPlaytimeExpectations(stars, year, overall)
	personality := util.GetPersonality()
	academicBias := util.GetAcademicBias()
	workEthic := util.GetWorkEthic()
	recruitingBias := util.GetRecruitingBias()
	freeAgency := util.GetFreeAgencyBias(0, 0)
	potentialGrade := util.GetWeightedPotentialGrade(potential)
	var basePlayer = structs.BasePlayer{
		FirstName:            firstName,
		LastName:             lastName,
		Position:             position,
		Age:                  19,
		Year:                 1,
		State:                state,
		Country:              country,
		Stars:                stars,
		Height:               height,
		Shooting2:            shooting2,
		Shooting3:            shooting3,
		FreeThrow:            freeThrow,
		Finishing:            finishing,
		Ballwork:             ballwork,
		Rebounding:           rebounding,
		InteriorDefense:      interiorDefense,
		PerimeterDefense:     perimeterDefense,
		Potential:            potential,
		ProPotentialGrade:    proPotential,
		Stamina:              stamina,
		PlaytimeExpectations: expectations,
		Minutes:              0,
		Overall:              overall,
		Personality:          personality,
		FreeAgency:           freeAgency,
		RecruitingBias:       recruitingBias,
		WorkEthic:            workEthic,
		AcademicBias:         academicBias,
		PotentialGrade:       potentialGrade,
	}

	var collegePlayer = structs.CollegePlayer{
		BasePlayer:    basePlayer,
		PlayerID:      id,
		TeamID:        team.ID,
		TeamAbbr:      team.Abbr,
		IsRedshirt:    false,
		IsRedshirting: false,
		HasGraduated:  false,
	}
	collegePlayer.SetID(id)

	// Specialties
	specs := util.GetSpecialties(collegePlayer.Position)
	for _, spec := range specs {
		collegePlayer.ToggleSpecialties(spec)
	}

	for i := 1; i < year && year > 1; i++ {
		collegePlayer = ProgressCollegePlayer(collegePlayer, 0, true)
	}

	return collegePlayer
}

func createRecruit(fName, lName, state, country, ethnicity, position string, year int, id uint) structs.Recruit {
	age := 18
	height := getHeight(position)
	potential := util.GeneratePotential()
	potentialGrade := util.GetWeightedPotentialGrade(potential)
	proPotential := util.GeneratePotential()
	stamina := util.GenerateIntFromRange(25, 38)
	discipline := util.GenerateIntFromRange(1, 20)
	injuryRating := util.GenerateIntFromRange(1, 20)
	personality := util.GetPersonality()
	academicBias := util.GetAcademicBias()
	workEthic := util.GetWorkEthic()
	recruitingBias := util.GetRecruitingBias()
	freeAgency := util.GetFreeAgencyBias(0, 0)

	var basePlayer = structs.BasePlayer{
		FirstName:         fName,
		LastName:          lName,
		Position:          position,
		Age:               age,
		Year:              year,
		State:             state,
		Country:           country,
		Height:            height,
		Potential:         potential,
		PotentialGrade:    potentialGrade,
		ProPotentialGrade: proPotential,
		Stamina:           stamina,
		Minutes:           0,
		Personality:       personality,
		FreeAgency:        freeAgency,
		RecruitingBias:    recruitingBias,
		WorkEthic:         workEthic,
		AcademicBias:      academicBias,
		InjuryRating:      injuryRating,
		Discipline:        discipline,
	}

	var croot = structs.Recruit{
		BasePlayer: basePlayer,
		PlayerID:   id,
		TeamID:     0,
		TeamAbbr:   "",
		IsSigned:   false,
		IsTransfer: false,
	}

	// Specialties
	specs := util.GetSpecialties(position)
	for _, spec := range specs {
		croot.ToggleSpecialties(spec)
	}

	shooting2 := util.GetAttributeNew(position, "Shooting2", croot.SpecShooting2)
	shooting3 := util.GetAttributeNew(position, "Shooting3", croot.SpecShooting3)
	finishing := util.GetAttributeNew(position, "Finishing", croot.SpecFinishing)
	freeThrow := util.GetAttributeNew(position, "FreeThrow", croot.SpecFreeThrow)
	ballwork := util.GetAttributeNew(position, "Ballwork", croot.SpecBallwork)
	rebounding := util.GetAttributeNew(position, "Rebounding", croot.SpecRebounding)
	interiorDefense := util.GetAttributeNew(position, "Interior Defense", croot.SpecInteriorDefense)
	perimeterDefense := util.GetAttributeNew(position, "Perimeter Defense", croot.SpecPerimeterDefense)

	overall := (int((shooting2 + shooting3 + freeThrow) / 3)) + finishing + ballwork + rebounding + int((interiorDefense+perimeterDefense)/2)
	stars := getStarRating(overall)
	recruitModifier := GetRecruitModifier(stars)
	expectations := util.GetPlaytimeExpectations(stars, year, overall)
	croot.SetID(id)
	croot.AssignRecruitModifier(recruitModifier)
	croot.SetAttributes(shooting2, shooting3, finishing, freeThrow, ballwork, rebounding, interiorDefense, perimeterDefense, overall, stars, expectations)

	croot.SetID(id)

	return croot
}

func createInternationalPlayer(teamID uint, team, country, ethnicity, position string, year int, firstNameList [][]string, lastNameList [][]string, id uint) structs.NBAPlayer {
	fName := getName(firstNameList)
	lName := getName(lastNameList)
	caser := cases.Title(language.English)

	firstName := caser.String(strings.ToLower(fName))
	lastName := caser.String(strings.ToLower(lName))
	age := util.GenerateISLAge()
	primeAge := util.GeneratePrimeAge()
	height := getHeight(position)
	potential := util.GeneratePotential()
	potentialGrade := util.GetWeightedPotentialGrade(potential)
	proPotential := util.GeneratePotential()
	stamina := util.GenerateStamina()

	personality := util.GetPersonality()
	academicBias := util.GetAcademicBias()
	workEthic := util.GetWorkEthic()
	recruitingBias := util.GetRecruitingBias()
	freeAgency := util.GetFreeAgencyBias(0, 0)

	var basePlayer = structs.BasePlayer{
		FirstName:         firstName,
		LastName:          lastName,
		Position:          position,
		Age:               age,
		Year:              year,
		State:             "",
		Country:           country,
		Height:            height,
		Potential:         potential,
		PotentialGrade:    potentialGrade,
		ProPotentialGrade: proPotential,
		Stamina:           stamina,
		Minutes:           0,
		Personality:       personality,
		FreeAgency:        freeAgency,
		RecruitingBias:    recruitingBias,
		WorkEthic:         workEthic,
		AcademicBias:      academicBias,
	}

	isNBAEligible := age > 21

	var player = structs.NBAPlayer{
		BasePlayer:      basePlayer,
		PlayerID:        id,
		TeamID:          teamID,
		TeamAbbr:        team,
		IsNBA:           isNBAEligible,
		IsInternational: true,
		PrimeAge:        uint(primeAge),
	}

	// Specialties
	specs := util.GetSpecialties(position)
	for _, spec := range specs {
		player.ToggleSpecialties(spec)
	}

	shooting2 := util.GetAttributeNew(position, "Shooting2", player.SpecShooting2)
	shooting3 := util.GetAttributeNew(position, "Shooting3", player.SpecShooting3)
	finishing := util.GetAttributeNew(position, "Finishing", player.SpecFinishing)
	freeThrow := util.GetAttributeNew(position, "FreeThrow", player.SpecFreeThrow)
	ballwork := util.GetAttributeNew(position, "Ballwork", player.SpecBallwork)
	rebounding := util.GetAttributeNew(position, "Rebounding", player.SpecRebounding)
	interiorDefense := util.GetAttributeNew(position, "Interior Defense", player.SpecInteriorDefense)
	perimeterDefense := util.GetAttributeNew(position, "Perimeter Defense", player.SpecPerimeterDefense)

	overall := (int((shooting2 + shooting3 + freeThrow) / 3)) + finishing + ballwork + rebounding + int((interiorDefense+perimeterDefense)/2)
	stars := getStarRating(overall)
	expectations := util.GetProfessionalPlaytimeExpectations(age, primeAge, overall)

	player.SetID(id)
	player.SetAttributes(shooting2, shooting3, finishing, freeThrow, ballwork, rebounding, interiorDefense, perimeterDefense, overall, stars, expectations)

	if age > 18 && age < 23 {
		diff := 22 - age

		for i := 0; i < diff; i++ {
			player = ProgressNBAPlayer(player, true)
		}
	}

	return player
}

func createCollegeCoach(team structs.TeamRecruitingProfile, almaMaterID uint, almaMater, ethnicity string, firstNameList, lastNameList [][]string, retiredPlayers []structs.RetiredPlayer, retireeMap, coachMap *map[uint]bool) structs.CollegeCoach {
	firstName := ""
	lastName := ""
	diceRoll := util.GenerateIntFromRange(1, 20)
	formerPlayerID := uint(0)
	almaID := almaMaterID
	alma := almaMater
	age := 32
	if diceRoll == 20 {
		// Get a former player as a coach
		idx := util.GenerateIntFromRange(0, len(retiredPlayers)-1)
		retiree := retiredPlayers[idx]
		for (*retireeMap)[retiree.ID] || (*coachMap)[retiree.ID] {
			idx = util.GenerateIntFromRange(0, len(retiredPlayers)-1)
			retiree = retiredPlayers[idx]
		}
		(*retireeMap)[retiree.ID] = true
		(*coachMap)[retiree.ID] = true
		formerPlayerID = retiree.ID
		almaID = retiree.CollegeID
		alma = retiree.College
		firstName = retiree.FirstName
		lastName = retiree.LastName
		age = retiree.Age + 1
	} else {
		fName := getName(firstNameList)
		lName := getName(lastNameList)
		caser := cases.Title(language.English)
		firstName = caser.String(strings.ToLower(fName))
		lastName = caser.String(strings.ToLower(lName))
		age = getCoachAge()
	}
	fullName := firstName + " " + lastName

	schoolQuality := team.AIQuality
	adminBehavior := team.AIBehavior
	goodHire := getGoodHire(schoolQuality, adminBehavior)
	starMin, starMax := getStarRange(schoolQuality, goodHire)
	pointMin, pointmax := getPointRange(schoolQuality, goodHire)
	odds1 := 0
	odds2 := 0
	odds3 := 0
	odds4 := 0
	odds5 := 0

	starList := make([]int, 5)
	for i := starMin; i <= starMax; i++ {
		starList = append(starList, i)
	}

	for _, star := range starList {
		if star == 1 {
			odds1 = 10
		} else if star == 2 {
			odds2 = 10
		} else if star == 3 {
			odds3 = 8
		} else if star == 4 {
			odds4 = 5
		} else if star == 5 {
			odds5 = 5
		}
	}

	schemeRoll := util.GenerateIntFromRange(1, 6)
	scheme := "Traditional"
	if schemeRoll == 6 {
		schemeList := []string{"Traditional", "Small Ball", "Microball", "Jumbo"}
		scheme = util.PickFromStringList(schemeList)
	}
	defensiveScheme := "Man-to-Man"
	defensiveSchemeList := []string{"Man-to-Man", "1-3-1 Zone", "3-2 Zone", "2-3 Zone", "Box-and-One Zone"}
	schemeRoll = util.GenerateIntFromRange(1, 6)
	if schemeRoll == 6 {
		defensiveScheme = util.PickFromStringList(defensiveSchemeList)
	}
	contractLength := util.GenerateIntFromRange(2, 5)
	startingPrestige := getStartingPrestige(goodHire)
	teamBuildingList := []string{"Recruiting", "Transfer", "Average"}
	teamBuildPref := util.PickFromStringList(teamBuildingList)
	careerPrefList := []string{"Average", "Prefers to Stay at Current Job", "Wants to coach Alma-Mater", "Wants a more competitive job", "Average"}
	careerPref := util.PickFromStringList(careerPrefList)
	promiseTendencyList := []string{"Average", "Under-Promise", "Over-Promise"}
	promiseTendency := util.PickFromStringList(promiseTendencyList)
	if goodHire {
		fmt.Println("Good hire for " + team.TeamAbbr + "!")
	}
	formerPlayer := formerPlayerID > 0

	if formerPlayer {
		fmt.Println("Former SimNBA Player " + fullName + " is committing to coach for " + team.TeamAbbr + "!")
	}

	coach := structs.CollegeCoach{
		Name:                   fullName,
		Age:                    age,
		TeamID:                 team.ID,
		Team:                   team.TeamAbbr,
		FormerPlayerID:         formerPlayerID,
		AlmaMaterID:            almaID,
		AlmaMater:              alma,
		Prestige:               startingPrestige,
		PointMin:               pointMin,
		PointMax:               pointmax,
		StarMin:                starMin,
		StarMax:                starMax,
		Odds1:                  odds1,
		Odds2:                  odds2,
		Odds3:                  odds3,
		Odds4:                  odds4,
		Odds5:                  odds5,
		Scheme:                 scheme,
		SchoolTenure:           0,
		CareerTenure:           0,
		ContractLength:         contractLength,
		YearsRemaining:         contractLength,
		IsRetired:              false,
		IsFormerPlayer:         formerPlayer,
		DefensiveScheme:        defensiveScheme,
		TeambuildingPreference: teamBuildPref,
		CareerPreference:       careerPref,
		PromiseTendency:        promiseTendency,
		PortalReputation:       100,
	}

	if startingPrestige > 1 {
		for i := 0; i < startingPrestige; i++ {
			selectStar := util.GenerateIntFromRange(starMin, starMax)
			coach.IncrementOdds(selectStar)
		}
	}

	return coach
}

func getNameList(ethnicity string, isFirstName bool) [][]string {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimNBA\\data"
	var fileName string
	if ethnicity == "Caucasian" {
		if isFirstName {
			fileName = "FNameW.csv"
		} else {
			fileName = "LNameW.csv"
		}
	} else if ethnicity == "African" {
		if isFirstName {
			fileName = "FNameB.csv"
		} else {
			fileName = "LNameB.csv"
		}
	} else if ethnicity == "Asian" {
		if isFirstName {
			fileName = "FNameA.csv"
		} else {
			fileName = "LNameA.csv"
		}
	} else if ethnicity == "NativeAmerican" {
		if isFirstName {
			fileName = "FNameN.csv"
		} else {
			fileName = "LNameN.csv"
		}
	} else {
		if isFirstName {
			fileName = "FNameH.csv"
		} else {
			fileName = "LNameH.csv"
		}
	}
	path = path + "\\" + fileName
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

func getNameMaps() (map[string][][]string, map[string][][]string) {
	CaucasianFirstNameList := getNameList("Caucasian", true)
	CaucasianLastNameList := getNameList("Caucasian", false)
	AfricanFirstNameList := getNameList("African", true)
	AfricanLastNameList := getNameList("African", false)
	HispanicFirstNameList := getNameList("Hispanic", true)
	HispanicLastNameList := getNameList("Hispanic", false)
	NativeFirstNameList := getNameList("NativeAmerican", true)
	NativeLastNameList := getNameList("NativeAmerican", false)
	AsianFirstNameList := getNameList("Asian", true)
	AsianLastNameList := getNameList("Asian", false)

	firstNameMap := make(map[string][][]string)
	firstNameMap["Caucasian"] = CaucasianFirstNameList
	firstNameMap["African"] = AfricanFirstNameList
	firstNameMap["Hispanic"] = HispanicFirstNameList
	firstNameMap["NativeAmerican"] = NativeFirstNameList
	firstNameMap["Asian"] = AsianFirstNameList

	lastNameMap := make(map[string][][]string)
	lastNameMap["Caucasian"] = CaucasianLastNameList
	lastNameMap["African"] = AfricanLastNameList
	lastNameMap["Hispanic"] = HispanicLastNameList
	lastNameMap["NativeAmerican"] = NativeLastNameList
	lastNameMap["Asian"] = AsianLastNameList

	return (firstNameMap), (lastNameMap)
}

func pickEthnicity() string {
	min := 0
	max := 10000
	num := util.GenerateIntFromRange(min, max)

	if num < 6000 {
		return "Caucasian"
	} else if num < 7800 {
		return "African"
	} else if num < 8900 {
		return "Hispanic"
	} else if num < 9975 {
		return "Asian"
	}
	return "NativeAmerican"
}

func pickISLEthnicity(country string) string {

	if country == "England" || country == "Scotland" || country == "Spain" ||
		country == "Italy" || country == "Latvia" || country == "Poland" ||
		country == "Estonia" || country == "Ukraine" || country == "France" ||
		country == "Germany" || country == "Belgium" || country == "Netherlands" ||
		country == "Turkey" || country == "Greece" || country == "Australia" ||
		country == "Israel" || country == "Lithuania" || country == "Serbia" {
		return "Caucasian"
	} else if country == "Morocco" || country == "Egypt" {
		return "African"
	} else if country == "Mexico" || country == "Argentina" || country == "Brazil" {
		return "Hispanic"
	} else if country == "China" || country == "Japan" || country == "South Korea" ||
		country == "Taiwan" || country == "Phillipines" || country == "New Zealand" {
		return "Asian"
	}
	return "NativeAmerican"
}

func pickCountry(ethnicity string) string {
	min := 0
	max := 10000
	num := util.GenerateIntFromRange(min, max)

	if num < 8001 {
		return "USA"
	} else if num < 8251 {
		if ethnicity == "African" {
			return "Dominican Republic"
		} else if ethnicity == "Hispanic" {
			return "Mexico"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "China"
		} else {
			return "Canada"
		}
	} else if num < 8301 {
		if ethnicity == "African" {
			return "The Bahamas"
		} else if ethnicity == "Hispanic" {
			return "Guatemala"
		} else if ethnicity == "NativeAmerican" {
			return "Russia"
		} else if ethnicity == "Asian" {
			return "China"
		} else {
			return "United Kingdom"
		}
	} else if num < 8351 {
		if ethnicity == "African" {
			return "Jamaica"
		} else if ethnicity == "Hispanic" {
			return "Costa Rica"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "China"
		} else {
			return "France"
		}
	} else if num < 8401 {
		if ethnicity == "African" {
			return "Democratic Republic of Congo"
		} else if ethnicity == "Hispanic" {
			return "Colombia"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "China"
		} else {
			return "Spain"
		}
	} else if num < 8451 {
		if ethnicity == "African" {
			return "South Africa"
		} else if ethnicity == "Hispanic" {
			return "Belize"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "China"
		} else {
			return "Ireland"
		}
	} else if num < 8501 {
		if ethnicity == "African" {
			return "Haiti"
		} else if ethnicity == "Hispanic" {
			return "El Salvador"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "China"
		} else {
			return "Spain"
		}
	} else if num < 8551 {
		if ethnicity == "African" {
			return "Ethiopia"
		} else if ethnicity == "Hispanic" {
			return "Cuba"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Japan"
		} else {
			return "Germany"
		}
	} else if num < 8601 {
		if ethnicity == "African" {
			return "Chad"
		} else if ethnicity == "Hispanic" {
			return "Honduras"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Japan"
		} else {
			return "Poland"
		}
	} else if num < 8651 {
		if ethnicity == "African" {
			return "Ghana"
		} else if ethnicity == "Hispanic" {
			return "Nicaragua"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Japan"
		} else {
			return "Sweden"
		}
	} else if num < 8701 {
		if ethnicity == "African" {
			return "Guinea"
		} else if ethnicity == "Hispanic" {
			return "Panama"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Vietnam"
		} else {
			return "Norway"
		}
	} else if num < 8751 {
		if ethnicity == "African" {
			return "Senegal"
		} else if ethnicity == "Hispanic" {
			return "Dominican Republic"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Vietnam"
		} else {
			return "Denmark"
		}
	} else if num < 8801 {
		if ethnicity == "African" {
			return "Morocco"
		} else if ethnicity == "Hispanic" {
			return "Mexico"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Indonesia"
		} else {
			return "Portugal"
		}
	} else if num < 8851 {
		if ethnicity == "African" {
			return "Algeria"
		} else if ethnicity == "Hispanic" {
			return "Mexico"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Indonesia"
		} else {
			return "Austria"
		}
	} else if num < 8901 {
		if ethnicity == "African" {
			return "Nigeria"
		} else if ethnicity == "Hispanic" {
			return "Venezuela"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Indonesia"
		} else {
			return "Hungary"
		}
	} else if num < 8951 {
		if ethnicity == "African" {
			return "Cameroon"
		} else if ethnicity == "Hispanic" {
			return "French Guiana"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Indonesia"
		} else {
			return "Croatia"
		}
	} else if num < 9001 {
		if ethnicity == "African" {
			return "Egypt"
		} else if ethnicity == "Hispanic" {
			return "Brazil"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Thailand"
		} else {
			return "Greece"
		}
	} else if num < 9051 {
		if ethnicity == "African" {
			return "Eritrea"
		} else if ethnicity == "Hispanic" {
			return "Brazil"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Thailand"
		} else {
			return "Israel"
		}
	} else if num < 9101 {
		if ethnicity == "African" {
			return "Kenya"
		} else if ethnicity == "Hispanic" {
			return "Guyana"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "South Korea"
		} else {
			return "Bulgaria"
		}
	} else if num < 9151 {
		if ethnicity == "African" {
			return "Liberia"
		} else if ethnicity == "Hispanic" {
			return "Ecuador"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Malaysia"
		} else {
			return "Romania"
		}
	} else if num < 9201 {
		if ethnicity == "African" {
			return "Tanzania"
		} else if ethnicity == "Hispanic" {
			return "Chile"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "India"
		} else {
			return "Montenegro"
		}
	} else if num < 9251 {
		if ethnicity == "African" {
			return "Zimbabwe"
		} else if ethnicity == "Hispanic" {
			return "Uruguay"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "India"
		} else {
			return "Turkey"
		}
	} else if num < 9301 {
		if ethnicity == "African" {
			return "Malawi"
		} else if ethnicity == "Hispanic" {
			return "Argentina"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "India"
		} else {
			return "Serbia"
		}
	} else if num < 9351 {
		if ethnicity == "African" {
			return "Senegal"
		} else if ethnicity == "Hispanic" {
			return "Argentina"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Israel"
		} else {
			return "Belgium"
		}
	} else if num < 9401 {
		if ethnicity == "African" {
			return "Senegal"
		} else if ethnicity == "Hispanic" {
			return "Argentina"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Bangladesh"
		} else {
			return "Ukraine"
		}
	} else if num < 9501 {
		if ethnicity == "African" {
			return "DCR"
		} else if ethnicity == "Hispanic" {
			return "Uruguay"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Philippines"
		} else {
			return "Ukraine"
		}
	} else if num < 9601 {
		if ethnicity == "African" {
			return "Nigeria"
		} else if ethnicity == "Hispanic" {
			return "Uruguay"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Philippines"
		} else {
			return "Russia"
		}
	} else if num < 9701 {
		if ethnicity == "African" {
			return "South Africa"
		} else if ethnicity == "Hispanic" {
			return "Chile"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Philippines"
		} else {
			return "Russia"
		}
	} else if num < 9801 {
		if ethnicity == "African" {
			return "South Africa"
		} else if ethnicity == "Hispanic" {
			return "Chile"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Singapore"
		} else {
			return "Lithuania"
		}
	} else if num < 9901 {
		if ethnicity == "African" {
			return "Uganda"
		} else if ethnicity == "Hispanic" {
			return "Peru"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Cambodia"
		} else {
			return "Estonia"
		}
	} else if num < 9951 {
		if ethnicity == "African" {
			return "Zambia"
		} else if ethnicity == "Hispanic" {
			return "Grenada"
		} else if ethnicity == "NativeAmerican" {
			return "USA"
		} else if ethnicity == "Asian" {
			return "Taiwan"
		} else {
			return "Finland"
		}
	} else if num < 9976 {
		if ethnicity == "African" {
			return "Tunisia"
		} else if ethnicity == "Hispanic" {
			return "Barbados"
		} else if ethnicity == "NativeAmerican" {
			return "Canada"
		} else if ethnicity == "Asian" {
			return "Myanmar"
		} else {
			return "Iceland"
		}
	} else {
		if ethnicity == "African" {
			return "Algeria"
		} else if ethnicity == "Hispanic" {
			return "Suriname"
		} else if ethnicity == "NativeAmerican" {
			return "Antarctica"
		} else if ethnicity == "Asian" {
			return "North Korea"
		} else {
			return "Luxembourg"
		}
	}
}

func pickState() string {
	states := []string{"Alabama", "Arkansas", "Arizona", "California", "Colorado", "Connecticut", "Delaware", "Florida", "Georgia", "Hawaii", "Idaho", "Illinois", "Indiana", "Iowa", "Kansas", "Kentucky", "Louisiana", "Maine", "Maryland", "Massachusetts", "Michigan", "Minnesota", "Mississippi", "Missouri", "Montana", "Nebraska", "Nevada", "New Hampshire", "New Jersey", "New Mexico", "New York", "North Carolina", "North Dakota", "Ohio", "Oklahoma", "Oregon", "Pennsylvania", "Rhode Island", "South Carolina", "South Dakota", "Tennessee", "Texas", "Utah", "Vermont", "Virginia", "Washington", "West Virginia", "Wisconsin", "Wyoming", "District of Columbia", "Guam", "Puerto Rico", "American Samoa"}

	return util.PickFromStringList(states)
}

func getHeight(position string) string {
	foot := 0
	inches := 0
	if position == "PG" || position == "SG" {
		footMin := 5
		footMax := 6
		foot = util.GenerateIntFromRange(footMin, footMax)

		if foot == 5 {
			inchesMin := 10
			inchesMax := 11
			inches = util.GenerateIntFromRange(inchesMin, inchesMax)
		} else {
			inchesMin := 0
			inchesMax := 5
			inches = util.GenerateIntFromRange(inchesMin, inchesMax)
		}
	} else if position == "PF" || position == "SF" {
		foot = 6
		inchesMin := 5
		inchesMax := 8
		inches = util.GenerateIntFromRange(inchesMin, inchesMax)
	} else {
		footMin := 6
		footMax := 7
		foot = util.GenerateIntFromRange(footMin, footMax)

		if foot == 6 {
			inchesMin := 9
			inchesMax := 11
			inches = util.GenerateIntFromRange(inchesMin, inchesMax)
		} else {
			inchesMin := 0
			inchesMax := 1
			inches = util.GenerateIntFromRange(inchesMin, inchesMax)
		}
	}
	height := strconv.Itoa(foot) + "-" + strconv.Itoa(inches)
	return height
}

func getAttribute(position string, attribute string, isGeneration bool) int {
	if position == "PG" || position == "SG" {
		if attribute == "Shooting2" {
			return util.GenerateIntFromRange(7, 17)
		} else if attribute == "Shooting3" {
			return util.GenerateIntFromRange(7, 17)
		} else if attribute == "Finishing" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "FreeThrow" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Ballwork" {
			return util.GenerateIntFromRange(7, 17)
		} else if attribute == "Rebounding" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Interior Defense" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Perimeter Defense" {
			return util.GenerateIntFromRange(1, 11)
		} else {
			return 1
		}
	} else if position == "PF" || position == "SF" {
		if attribute == "Shooting2" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Shooting3" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "FreeThrow" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Finishing" {
			return util.GenerateIntFromRange(6, 16)
		} else if attribute == "Ballwork" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Rebounding" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Interior Defense" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Perimeter Defense" {
			return util.GenerateIntFromRange(4, 14)
		} else {
			return 1
		}
	} else if position == "C" {
		if attribute == "Shooting2" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Shooting3" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "FreeThrow" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Finishing" {
			return util.GenerateIntFromRange(6, 16)
		} else if attribute == "Ballwork" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Rebounding" {
			return util.GenerateIntFromRange(6, 16)
		} else if attribute == "Interior Defense" {
			return util.GenerateIntFromRange(6, 16)
		} else if attribute == "Perimeter Defense" {
			return util.GenerateIntFromRange(6, 16)
		} else {
			return 1
		}
	}
	return 1
}

func getStarRating(overall int) int {
	if overall > 67 {
		return 5
	} else if overall > 61 {
		return 4
	} else if overall > 52 {
		return 3
	} else if overall > 45 {
		return 2
	} else {
		return 1
	}
}

func GetRecruitModifier(stars int) int {
	if stars == 5 {
		return util.GenerateIntFromRange(80, 117)
	} else if stars == 4 {
		return util.GenerateIntFromRange(100, 125)
	} else if stars == 3 {
		return util.GenerateIntFromRange(117, 150)
	} else if stars == 2 {
		return util.GenerateIntFromRange(125, 200)
	}
	return util.GenerateIntFromRange(150, 250)
}

func getName(list [][]string) string {
	endOfListWeight, err := strconv.Atoi(list[len(list)-1][1])
	if err != nil {
		log.Fatalln("Could not convert number from string")
	}
	name := ""
	num := util.GenerateIntFromRange(1, endOfListWeight)
	for i := 1; i < len(list); i++ {
		weight, err := strconv.Atoi(list[i][1])
		if err != nil {
			log.Fatalln("Could not convert number from string in name generator")
		}
		if num < weight {
			name = list[i][0]
			break
		}
	}
	return name
}

func pickAlmaMater(teams []structs.TeamRecruitingProfile) structs.TeamRecruitingProfile {
	start := 0
	end := len(teams) - 1
	idx := util.GenerateIntFromRange(start, end)
	return teams[idx]
}

func getCoachAge() int {
	num := util.GenerateIntFromRange(1, 100)

	if num < 10 {
		return util.GenerateIntFromRange(32, 36)
	} else if num < 25 {
		return util.GenerateIntFromRange(37, 39)
	} else if num < 55 {
		return util.GenerateIntFromRange(40, 49)
	} else if num < 80 {
		return util.GenerateIntFromRange(50, 59)
	} else if num < 95 {
		return util.GenerateIntFromRange(60, 65)
	} else {
		return util.GenerateIntFromRange(66, 70)
	}
}

func getGoodHire(schoolQuality, adminBehavior string) bool {
	diceRoll := util.GenerateIntFromRange(1, 20)
	mod := 0
	if schoolQuality == "P6" || schoolQuality == "Cinderella" {
		mod += 1
	} else if schoolQuality == "Blue Blood" {
		mod += 3
	}
	if adminBehavior == "Aggressive" {
		mod += 3
	} else if adminBehavior == "Conservative" {
		mod -= 3
	}

	sum := diceRoll + mod
	goodHire := sum > 12
	return goodHire
}

func getStarRange(schoolQuality string, goodHire bool) (int, int) {

	if schoolQuality == "Blue Blood" {
		if goodHire {
			return 3, 5
		} else {
			return 3, 4
		}
	} else if schoolQuality == "Cinderella" {
		if goodHire {
			return 2, 4
		} else {
			return 2, 3
		}
	} else if schoolQuality == "P6" {
		if goodHire {
			return 2, 4
		} else {
			return 2, 3
		}
	} else {
		if goodHire {
			return 1, 3
		} else {
			return 1, 2
		}
	}
}

func getPointRange(schoolQuality string, goodHire bool) (int, int) {
	min := 0
	max := 15
	if schoolQuality == "Blue Blood" {
		if goodHire {
			min = util.GenerateIntFromRange(7, 8)
			max = util.GenerateIntFromRange(12, 16)
		} else {
			min = util.GenerateIntFromRange(6, 7)
			max = util.GenerateIntFromRange(10, 13)
		}
	} else if schoolQuality == "Cinderella" {
		if goodHire {
			min = util.GenerateIntFromRange(5, 7)
			max = util.GenerateIntFromRange(10, 15)
		} else {
			min = util.GenerateIntFromRange(4, 6)
			max = util.GenerateIntFromRange(10, 12)
		}
	} else if schoolQuality == "P6" {
		if goodHire {
			min = util.GenerateIntFromRange(5, 8)
			max = util.GenerateIntFromRange(10, 14)
		} else {
			min = util.GenerateIntFromRange(4, 6)
			max = util.GenerateIntFromRange(8, 12)
		}
	} else {
		if goodHire {
			min = util.GenerateIntFromRange(3, 6)
			max = util.GenerateIntFromRange(8, 12)
		} else {
			min = 4
			max = util.GenerateIntFromRange(6, 8)
		}
	}
	return min, max
}

func getStartingPrestige(goodHire bool) int {
	if goodHire {
		return util.GenerateIntFromRange(3, 7)
	}
	return util.GenerateIntFromRange(1, 5)
}

func getRelativeType() int {
	roll := util.GenerateIntFromRange(1, 1000)
	// Brother of existing player
	if roll < 600 {
		return 2
	}
	// Cousin of existing player
	if roll < 800 {
		return 3
	}
	// Half brother of existing player
	if roll < 850 {
		return 4
	}
	// Twin
	if roll < 900 {
		return 5
	}
	// Best friend of another recruit
	if roll < 925 {
		return 8
	}
	// Best friend of a college player
	if roll < 950 {
		return 8
	}
	// Coach's Son
	if roll < 985 {
		return 6
	}
	// Coach's Nephew
	return 7
}

func getCoachLastName(fullName string) string {
	nameSplit := strings.Split(fullName, " ")
	return nameSplit[1]
}
