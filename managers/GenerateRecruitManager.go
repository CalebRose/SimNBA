package managers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

func getLatestRecord(db *gorm.DB) uint {
	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	return lastPlayerRecord.ID + 1
}

type CrootGenerator struct {
	firstNameMap      map[string][][]string
	lastNameMap       map[string][][]string
	nameMap           map[string]map[string][]string
	faceDataBlob      map[string][]string
	collegePlayerList []structs.CollegePlayer
	teamMap           map[uint]structs.Team
	positionList      []string
	CrootList         []structs.Recruit
	GlobalList        []structs.GlobalPlayer
	FacesList         []structs.FaceData
	attributeBlob     map[string]map[string]map[string]map[string]interface{}
	usCrootLocations  map[string][]structs.CrootLocation
	newID             uint
	count             int
	requiredPlayers   int
	star5             int
	star4             int
	star3             int
	star2             int
	star1             int
	ovr35             int
	ovr30             int
	ovr25             int
	ovr20             int
	ovr15             int
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
	relativeType := 0
	relativeID := 0
	notes := ""
	star := util.GetStarRating(false, false)
	lastName := ""
	state := ""
	pg.pickedEthnicity = pickEthnicity()
	country := pickCountry(pg.pickedEthnicity)
	locale := pickLocale(country)
	switch country {
	case "USA":
		state = util.PickState()
	}
	nameList := pg.nameMap[locale]
	fName := util.PickFromStringList(nameList["first_names"])
	firstName := pg.caser.String(strings.ToLower(fName))
	// Roll for type of recruit generated
	// If num == 200, then create some flair
	roof := 100
	relativeRoll := util.GenerateIntFromRange(1, roof)
	relativeIdx := 0
	if relativeRoll == roof {
		relativeType = getRelativeType()
		switch relativeType {
		case 2:
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
			notes = "Brother of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		case 3:
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
				lName := util.PickFromStringList(nameList["last_names"])
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			country = cp.Country
			notes = "Cousin of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		case 4:
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
				lName := util.PickFromStringList(nameList["last_names"])
				lastName = pg.caser.String(strings.ToLower(lName))
			}
			state = cp.State
			country = cp.Country
			notes = "Half-Brother of " + cp.Team + " " + cp.Position + " " + cp.FirstName + " " + cp.LastName
		case 5:
			// Twin
			relativeType = 5
			relativeID = int(pg.newID)
		}
	} else {
		relativeType = 1
	}
	if relativeType == 1 || relativeType == 5 {
		lName := util.PickFromStringList(nameList["last_names"])
		lastName = pg.caser.String(strings.ToLower(lName))
	}
	pickedPosition := util.PickPositionFromList()
	player := createRecruit(pickedPosition, firstName, lastName, state, country, "", "", 1, int(star), pg.newID, pg.attributeBlob, pg.usCrootLocations[state])
	player.AssignRelativeData(uint(relativeID), uint(relativeType), 0, "", notes)
	globalPlayer := structs.GlobalPlayer{
		CollegePlayerID: pg.newID,
		RecruitID:       pg.newID,
		NBAPlayerID:     pg.newID,
	}
	skinColor := getSkinColor(player.Country)

	face := getFace(pg.newID, 238, skinColor, pg.faceDataBlob)

	pg.FacesList = append(pg.FacesList, face)
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
	country := player.Country
	locale := pickLocale(country)
	firstNameList := pg.nameMap[locale]["first_names"]
	twinName := util.PickFromStringList(firstNameList)
	twinN := pg.caser.String(strings.ToLower(twinName))
	twinPosition := ""
	switch player.Position {
	case "F":
		twinPosition = util.PickFromStringList([]string{"C", "F"})
	case "C":
		twinPosition = "F"
	case "G":
		twinPosition = util.PickFromStringList([]string{"G", "F"})
	default:
		twinPosition = "G"
	}
	twinNotes := "Twin Brother of " + strconv.Itoa(int(player.Stars)) + " Star Recruit " + player.Position + " " + player.FirstName + " " + player.LastName
	twinPlayer := createRecruit(twinPosition, twinN, player.LastName, player.State, player.Country, "", "", 1, int(player.Stars), pg.newID, pg.attributeBlob, pg.usCrootLocations[player.State])
	twinPlayer.AssignRelativeData(uint(twinRelativeID), 4, 0, "", twinNotes)
	notes := "Twin Brother of " + strconv.Itoa(int(twinPlayer.Stars)) + " Star Recruit " + twinPlayer.Position + " " + twinPlayer.FirstName + " " + twinPlayer.LastName
	player.AssignRelativeData(uint(relativeID), 4, 0, "", notes)
	globalTwinPlayer := structs.GlobalPlayer{
		CollegePlayerID: pg.newID,
		RecruitID:       pg.newID,
		NBAPlayerID:     pg.newID,
	}
	globalTwinPlayer.SetID(pg.newID)
	player.AssignRelativeData(uint(relativeID), uint(player.RelativeType), 0, "", notes)
	globalPlayer := structs.GlobalPlayer{
		CollegePlayerID: pg.newID,
		RecruitID:       pg.newID,
		NBAPlayerID:     pg.newID,
	}
	skinColor := getSkinColor(player.Country)

	face := getFace(pg.newID, 238, skinColor, pg.faceDataBlob)

	pg.FacesList = append(pg.FacesList, face)
	globalPlayer.SetID(pg.newID)
	return twinPlayer, globalPlayer
}

func (pg *CrootGenerator) updateStatistics(player structs.Recruit) {
	switch player.Stars {
	case 5:
		pg.star5++
	case 4:
		pg.star4++
	case 3:
		pg.star3++
	case 2:
		pg.star2++
	default:
		pg.star1++
	}

	if player.Overall >= 35 {
		pg.ovr35++
	} else if player.Overall >= 30 {
		pg.ovr30++
	} else if player.Overall >= 25 {
		pg.ovr25++
	} else if player.Overall >= 20 {
		pg.ovr20++
	} else if player.Overall >= 15 {
		pg.ovr15++
	}

	if int(player.Overall) > pg.highestOvr {
		pg.highestOvr = int(player.Overall)
	}
	if int(player.Overall) < pg.lowestOvr {
		pg.lowestOvr = int(player.Overall)
	}
}

func (pg *CrootGenerator) OutputRecruitStats() {
	// Croot Stats
	fmt.Println("Total Recruit Count: ", pg.count)
	fmt.Println("Total Ovr 35  Count: ", pg.ovr35)
	fmt.Println("Total Ovr 30  Count: ", pg.ovr30)
	fmt.Println("Total Ovr 25  Count: ", pg.ovr25)
	fmt.Println("Total Ovr 20  Count: ", pg.ovr20)
	fmt.Println("Total Ovr 15  Count: ", pg.ovr15)
	fmt.Println("Total 5 Star  Count: ", pg.star5)
	fmt.Println("Total 4 Star  Count: ", pg.star4)
	fmt.Println("Total 3 Star  Count: ", pg.star3)
	fmt.Println("Total 2 Star  Count: ", pg.star2)
	fmt.Println("Total 1 Star  Count: ", pg.star1)
	fmt.Println("Highest Recruit Ovr: ", pg.highestOvr)
	fmt.Println("Lowest  Recruit Ovr: ", pg.lowestOvr)
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
	firstNameMap, lastNameMap := getNameMaps()
	teams := GetAllActiveCollegeTeams()

	generator := CrootGenerator{
		attributeBlob: getAttributeBlob(),
	}

	var positionList []string = []string{"G", "G", "F", "F", "C"}

	for _, team := range teams {
		if team.ID < 370 {
			continue
		}
		// Test Generation
		yearList := []int{}
		players := GetCollegePlayersByTeamId(strconv.Itoa(int(team.ID)))
		seniors := 3
		juniors := 3
		sophomores := 3
		freshmen := 3
		fCount := 0
		gCount := 0
		cCount := 0
		requiredPlayers := 13
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
			switch player.Position {
			case "F":
				fCount++
			case "G":
				gCount++
			default:
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
			if pickedPosition == "F" && fCount > 5 {
				quickList := []string{"G", "C"}
				pickedPosition = util.PickFromStringList(quickList)
			} else if pickedPosition == "G" && gCount > 5 {
				quickList := []string{"F", "C"}
				pickedPosition = util.PickFromStringList(quickList)
			} else if pickedPosition == "C" && cCount > 2 {
				quickList := []string{"F", "G"}
				pickedPosition = util.PickFromStringList(quickList)
			}

			switch pickedPosition {
			case "F":
				fCount++
			case "G":
				gCount++
			default:
				cCount++
			}

			positionQueue = append(positionQueue, pickedPosition)
		}

		rand.Shuffle(len(positionQueue), func(i, j int) {
			positionQueue[i], positionQueue[j] = positionQueue[j], positionQueue[i]
		})

		rand.Shuffle(len(yearList), func(i, j int) {
			yearList[i], yearList[j] = yearList[j], yearList[i]
		})

		for count < requiredPlayers {
			year := yearList[count]
			pickedEthnicity := pickEthnicity()
			pickedPosition := positionQueue[count]
			player := createCollegePlayer(team, pickedEthnicity, pickedPosition, year, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], newID, false, generator.attributeBlob)
			globalPlayer := structs.GlobalPlayer{
				Model:           gorm.Model{ID: newID},
				CollegePlayerID: newID,
				RecruitID:       newID,
				NBAPlayerID:     newID,
			}
			// playerList = append(playerList, player)
			err := db.Create(&player).Error
			if err != nil {
				log.Panicln("Could not save player record")
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
	generator := CrootGenerator{
		nameMap:           getInternationalNameMap(),
		collegePlayerList: GetAllCollegePlayers(),
		teamMap:           GetCollegeTeamMap(),
		positionList:      []string{"G", "G", "F", "F", "C"},
		newID:             lastPlayerRecord.ID + 1,
		requiredPlayers:   util.GenerateIntFromRange(1500, 1600),
		usCrootLocations:  getCrootLocations("HS"),
		attributeBlob:     getAttributeBlob(),
		faceDataBlob:      getFaceDataBlob(),
		count:             0,
		star5:             0,
		star4:             0,
		star3:             0,
		star2:             0,
		star1:             0,
		highestOvr:        0,
		lowestOvr:         100000,
		CrootList:         []structs.Recruit{},
		GlobalList:        []structs.GlobalPlayer{},
		FacesList:         []structs.FaceData{},
		caser:             cases.Title(language.English),
		pickedEthnicity:   "",
	}

	// Test Generation
	// requiredPlayers := util.GenerateIntFromRange(203, 205)
	// 1061 is the number of open spots on teams in the league.
	// Currently 378 teams. 363 * 3 = 1089, the size of the average class.
	// The plan is to ensure that every recruit is signed
	generator.GenerateRecruits()
	// Croot Stats
	generator.OutputRecruitStats()

	// Import Batches
	repository.CreateRecruitRecordsBatch(db, generator.CrootList, 500)
	repository.CreateGlobalRecordsBatch(db, generator.GlobalList, 500)
	repository.CreateFaceRecordsBatch(db, generator.FacesList, 500)
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

	nameMap := getInternationalNameMap()
	requiredLimit := 1300
	poolCount := GetYouthDevelopmentPlayerCount()

	blob := getAttributeBlob()

	for poolCount < requiredLimit {
		pickedPosition := util.PickPositionFromList()
		country := pickISLCountry()
		pickedEthnicity := pickLocale(country)
		year := 1
		countryNames := nameMap[pickedEthnicity]
		player := createInternationalPlayer(0, "", country, pickedEthnicity, pickedPosition, year, countryNames["first_names"], countryNames["last_names"], newID, blob)
		repository.CreateProfessionalPlayerRecord(player, db)

		globalPlayer := structs.GlobalPlayer{
			CollegePlayerID: newID,
			RecruitID:       newID,
			NBAPlayerID:     newID,
		}

		globalPlayer.SetID(newID)

		repository.CreateGlobalPlayerRecord(globalPlayer, db)

		poolCount++
		newID++
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
		discipline := util.GenerateNormalizedIntFromRange(1, 20)
		injuryRating := util.GenerateNormalizedIntFromRange(1, 20)
		c.SetDisciplineAndIR(discipline, injuryRating)
		repository.SaveCollegePlayerRecord(c, db)
	}

	nbaPlayers := GetAllNBAPlayers()

	for _, n := range nbaPlayers {
		discipline := util.GenerateNormalizedIntFromRange(1, 20)
		injuryRating := util.GenerateNormalizedIntFromRange(1, 20)
		n.SetDisciplineAndIR(discipline, injuryRating)
		repository.SaveProfessionalPlayerRecord(n, db)
	}
}

// Private Methods
func createCollegePlayer(team structs.Team, ethnicity string, position string, year int, firstNameList [][]string, lastNameList [][]string, id uint, isWalkon bool, attributeBlob map[string]map[string]map[string]map[string]interface{}) structs.CollegePlayer {
	fName := getName(firstNameList)
	lName := getName(lastNameList)
	caser := cases.Title(language.English)
	archetype := util.GetArchetype(position)

	firstName := caser.String(strings.ToLower(fName))
	lastName := caser.String(strings.ToLower(lName))
	state := ""
	country := pickCountry(ethnicity)
	if country == "USA" {
		state = util.PickState()
	}
	height := 0
	weight := 0
	potential := util.GeneratePotential()
	potentialGrade := util.GetWeightedPotentialGrade(potential)
	proPotential := util.GeneratePotential()
	stamina := util.GenerateIntFromRange(25, 38)
	discipline := util.GenerateNormalizedIntFromRange(1, 20)
	injuryRating := util.GenerateNormalizedIntFromRange(1, 20)
	personality := util.GetPersonality()
	academicBias := util.GetAcademicBias()
	workEthic := util.GetWorkEthic()
	recruitingBias := util.GetRecruitingBias()
	freeAgency := util.GetFreeAgencyBias(0, 0)

	stars := util.GetStarRating(false, false)
	// Preferences
	program := util.GenerateNormalizedIntFromRange(1, 9)
	profDevelopment := util.GenerateNormalizedIntFromRange(1, 9)
	traditions := util.GenerateNormalizedIntFromRange(1, 9)
	facilities := util.GenerateNormalizedIntFromRange(1, 9)
	atmosphere := util.GenerateNormalizedIntFromRange(1, 9)
	academics := util.GenerateNormalizedIntFromRange(1, 9)
	conferencePrestige := util.GenerateNormalizedIntFromRange(1, 9)
	coachPref := util.GenerateNormalizedIntFromRange(1, 9)
	seasonMomentumPref := util.GenerateNormalizedIntFromRange(1, 9)
	campusLife := util.GenerateNormalizedIntFromRange(1, 9)
	// Attributes
	insideShooting := getAttributeValue(position, archetype, stars, "InsideShooting", attributeBlob)
	shooting2 := getAttributeValue(position, archetype, stars, "MidRangeShooting", attributeBlob)
	shooting3 := getAttributeValue(position, archetype, stars, "ThreePointShooting", attributeBlob)
	freeThrow := getAttributeValue(position, archetype, stars, "FreeThrow", attributeBlob)
	ballwork := getAttributeValue(position, archetype, stars, "Ballwork", attributeBlob)
	rebounding := getAttributeValue(position, archetype, stars, "Rebounding", attributeBlob)
	agility := getAttributeValue(position, archetype, stars, "Agility", attributeBlob)
	stealing := getAttributeValue(position, archetype, stars, "Stealing", attributeBlob)
	blocking := getAttributeValue(position, archetype, stars, "Blocking", attributeBlob)
	interiorDefense := getAttributeValue(position, archetype, stars, "InteriorDefense", attributeBlob)
	perimeterDefense := getAttributeValue(position, archetype, stars, "PerimeterDefense", attributeBlob)
	// Specialties
	specInsideShooting := util.GenerateSpecialty(position, archetype, "InsideShooting")
	specMidRangeShooting := util.GenerateSpecialty(position, archetype, "MidRangeShooting")
	specThreePointShooting := util.GenerateSpecialty(position, archetype, "ThreePointShooting")
	specFreeThrow := util.GenerateSpecialty(position, archetype, "FreeThrow")
	specBallwork := util.GenerateSpecialty(position, archetype, "Ballwork")
	specAgility := util.GenerateSpecialty(position, archetype, "Agility")
	specStealing := util.GenerateSpecialty(position, archetype, "Stealing")
	specBlocking := util.GenerateSpecialty(position, archetype, "Blocking")
	specRebounding := util.GenerateSpecialty(position, archetype, "Rebounding")
	specInteriorDefense := util.GenerateSpecialty(position, archetype, "InteriorDefense")
	specPerimeterDefense := util.GenerateSpecialty(position, archetype, "PerimeterDefense")
	height = getAttributeValue(position, archetype, stars, "Height", attributeBlob)
	weight = getAttributeValue(position, archetype, stars, "Weight", attributeBlob)
	if isWalkon {
		stars = 0
	}
	expectations := util.GetPlaytimeExpectations(stars, year, 10)
	basePlayer := structs.BasePlayer{
		FirstName:              firstName,
		LastName:               lastName,
		Position:               position,
		Archetype:              archetype,
		Age:                    19,
		Year:                   uint8(year),
		State:                  state,
		Country:                country,
		Height:                 uint8(height),
		Weight:                 uint16(weight),
		Agility:                uint8(agility),
		InsideShooting:         uint8(insideShooting),
		MidRangeShooting:       uint8(shooting2),
		ThreePointShooting:     uint8(shooting3),
		FreeThrow:              uint8(freeThrow),
		Ballwork:               uint8(ballwork),
		Stealing:               uint8(stealing),
		Blocking:               uint8(blocking),
		Rebounding:             uint8(rebounding),
		InteriorDefense:        uint8(interiorDefense),
		PerimeterDefense:       uint8(perimeterDefense),
		SpecInsideShooting:     specInsideShooting,
		SpecMidRangeShooting:   specMidRangeShooting,
		SpecThreePointShooting: specThreePointShooting,
		SpecFreeThrow:          specFreeThrow,
		SpecBallwork:           specBallwork,
		SpecAgility:            specAgility,
		SpecStealing:           specStealing,
		SpecBlocking:           specBlocking,
		SpecRebounding:         specRebounding,
		SpecInteriorDefense:    specInteriorDefense,
		SpecPerimeterDefense:   specPerimeterDefense,
		Potential:              potential,
		PotentialGrade:         potentialGrade,
		ProPotentialGrade:      proPotential,
		PlaytimeExpectations:   uint8(expectations),
		Stamina:                uint8(stamina),
		Personality:            personality,
		FreeAgency:             freeAgency,
		RecruitingBias:         recruitingBias,
		WorkEthic:              workEthic,
		AcademicBias:           academicBias,
		InjuryRating:           uint8(injuryRating),
		Discipline:             uint8(discipline),
		PlayerID:               id,
		TeamID:                 team.ID,
		Team:                   team.Abbr,
		PlayerPreferences: structs.PlayerPreferences{
			ProgramPref:        uint8(program),
			ProfDevPref:        uint8(profDevelopment),
			TraditionsPref:     uint8(traditions),
			FacilitiesPref:     uint8(facilities),
			AtmospherePref:     uint8(atmosphere),
			AcademicsPref:      uint8(academics),
			ConferencePref:     uint8(conferencePrestige),
			CoachPref:          uint8(coachPref),
			SeasonMomentumPref: uint8(seasonMomentumPref),
			CampusLifePref:     uint8(campusLife),
		},
	}

	collegePlayer := structs.CollegePlayer{
		BasePlayer:    basePlayer,
		IsRedshirt:    false,
		IsRedshirting: false,
		HasGraduated:  false,
	}
	collegePlayer.SetID(id)
	return collegePlayer
}

func createRecruit(position, fName, lName, state, country, cit, hs string, year uint8, stars int, id uint, blob map[string]map[string]map[string]map[string]interface{}, hsBlob []structs.CrootLocation) structs.Recruit {
	age := 18
	city, highSchool := cit, hs
	if state != "" && len(hsBlob) > 0 && country == "USA" {
		city, highSchool = getCityAndHighSchool(hsBlob)
	}
	archetype := util.GetArchetype(position)

	height := getAttributeValue(position, archetype, stars, "Height", blob)
	weight := getAttributeValue(position, archetype, stars, "Weight", blob)
	potential := util.GeneratePotential()
	potentialGrade := util.GetWeightedPotentialGrade(potential)
	proPotential := util.GeneratePotential()
	stamina := util.GenerateIntFromRange(25, 38)
	discipline := util.GenerateNormalizedIntFromRange(1, 20)
	injuryRating := util.GenerateNormalizedIntFromRange(1, 20)
	personality := util.GetPersonality()
	academicBias := util.GetAcademicBias()
	workEthic := util.GetWorkEthic()
	recruitingBias := util.GetRecruitingBias()
	freeAgency := util.GetFreeAgencyBias(0, 0)

	// Preferences
	program := util.GenerateNormalizedIntFromRange(1, 9)
	profDevelopment := util.GenerateNormalizedIntFromRange(1, 9)
	traditions := util.GenerateNormalizedIntFromRange(1, 9)
	facilities := util.GenerateNormalizedIntFromRange(1, 9)
	atmosphere := util.GenerateNormalizedIntFromRange(1, 9)
	academics := util.GenerateNormalizedIntFromRange(1, 9)
	conferencePrestige := util.GenerateNormalizedIntFromRange(1, 9)
	coachPref := util.GenerateNormalizedIntFromRange(1, 9)
	seasonMomentumPref := util.GenerateNormalizedIntFromRange(1, 9)
	campusLife := util.GenerateNormalizedIntFromRange(1, 9)

	// Attributes
	insideShooting := getAttributeValue(position, archetype, stars, "InsideShooting", blob)
	midRangeShooting := getAttributeValue(position, archetype, stars, "MidRangeShooting", blob)
	threePointShooting := getAttributeValue(position, archetype, stars, "ThreePointShooting", blob)
	freeThrow := getAttributeValue(position, archetype, stars, "FreeThrow", blob)
	ballwork := getAttributeValue(position, archetype, stars, "Ballwork", blob)
	rebounding := getAttributeValue(position, archetype, stars, "Rebounding", blob)
	interiorDefense := getAttributeValue(position, archetype, stars, "InteriorDefense", blob)
	perimeterDefense := getAttributeValue(position, archetype, stars, "PerimeterDefense", blob)
	agility := getAttributeValue(position, archetype, stars, "Agility", blob)
	stealing := getAttributeValue(position, archetype, stars, "Stealing", blob)
	blocking := getAttributeValue(position, archetype, stars, "Blocking", blob)
	// Specialties
	specInsideShooting := util.GenerateSpecialty(position, archetype, "InsideShooting")
	specMidRangeShooting := util.GenerateSpecialty(position, archetype, "MidRangeShooting")
	specThreePointShooting := util.GenerateSpecialty(position, archetype, "ThreePointShooting")
	specFreeThrow := util.GenerateSpecialty(position, archetype, "FreeThrow")
	specBallwork := util.GenerateSpecialty(position, archetype, "Ballwork")
	specAgility := util.GenerateSpecialty(position, archetype, "Agility")
	specStealing := util.GenerateSpecialty(position, archetype, "Stealing")
	specBlocking := util.GenerateSpecialty(position, archetype, "Blocking")
	specRebounding := util.GenerateSpecialty(position, archetype, "Rebounding")
	specInteriorDefense := util.GenerateSpecialty(position, archetype, "InteriorDefense")
	specPerimeterDefense := util.GenerateSpecialty(position, archetype, "PerimeterDefense")

	starRating := stars
	if starRating > 5 {
		starRating = 5
	}

	var basePlayer = structs.BasePlayer{
		FirstName:              fName,
		LastName:               lName,
		Position:               position,
		Archetype:              archetype,
		Age:                    uint8(age),
		Year:                   uint8(year),
		Stars:                  uint8(starRating),
		City:                   city,
		HighSchool:             highSchool,
		State:                  state,
		Country:                country,
		Height:                 uint8(height),
		Weight:                 uint16(weight),
		Agility:                uint8(agility),
		InsideShooting:         uint8(insideShooting),
		MidRangeShooting:       uint8(midRangeShooting),
		ThreePointShooting:     uint8(threePointShooting),
		FreeThrow:              uint8(freeThrow),
		Ballwork:               uint8(ballwork),
		Stealing:               uint8(stealing),
		Blocking:               uint8(blocking),
		Rebounding:             uint8(rebounding),
		InteriorDefense:        uint8(interiorDefense),
		PerimeterDefense:       uint8(perimeterDefense),
		SpecInsideShooting:     specInsideShooting,
		SpecMidRangeShooting:   specMidRangeShooting,
		SpecThreePointShooting: specThreePointShooting,
		SpecFreeThrow:          specFreeThrow,
		SpecBallwork:           specBallwork,
		SpecAgility:            specAgility,
		SpecStealing:           specStealing,
		SpecBlocking:           specBlocking,
		SpecRebounding:         specRebounding,
		SpecInteriorDefense:    specInteriorDefense,
		SpecPerimeterDefense:   specPerimeterDefense,
		Potential:              potential,
		PotentialGrade:         potentialGrade,
		ProPotentialGrade:      proPotential,
		Stamina:                uint8(stamina),
		Personality:            personality,
		FreeAgency:             freeAgency,
		RecruitingBias:         recruitingBias,
		WorkEthic:              workEthic,
		AcademicBias:           academicBias,
		InjuryRating:           uint8(injuryRating),
		Discipline:             uint8(discipline),
		PlayerPreferences: structs.PlayerPreferences{
			ProgramPref:        uint8(program),
			ProfDevPref:        uint8(profDevelopment),
			TraditionsPref:     uint8(traditions),
			FacilitiesPref:     uint8(facilities),
			AtmospherePref:     uint8(atmosphere),
			AcademicsPref:      uint8(academics),
			ConferencePref:     uint8(conferencePrestige),
			CoachPref:          uint8(coachPref),
			SeasonMomentumPref: uint8(seasonMomentumPref),
			CampusLifePref:     uint8(campusLife),
		},
	}

	var croot = structs.Recruit{
		BasePlayer: basePlayer,
		PlayerID:   id,
		IsSigned:   false,
		IsTransfer: false,
	}

	croot.GetOverall()

	recruitModifier := GetRecruitModifier(uint8(starRating))
	expectations := util.GetPlaytimeExpectations(starRating, int(year), int(croot.Overall))
	croot.SetExpectations(uint8(expectations))
	croot.SetID(id)
	croot.AssignRecruitModifier(recruitModifier)
	croot.SetID(id)

	return croot
}

func createInternationalPlayer(teamID uint, team, country, ethnicity, position string, year int, firstNameList, lastNameList []string, id uint, blob map[string]map[string]map[string]map[string]interface{}) structs.NBAPlayer {
	if len(firstNameList) == 0 {
		fmt.Println(country)
	}
	fName := util.PickFromStringList(firstNameList)
	lName := util.PickFromStringList(lastNameList)
	caser := cases.Title(language.English)
	firstName := caser.String(strings.ToLower(fName))
	lastName := caser.String(strings.ToLower(lName))
	age := util.GenerateISLAge()
	primeAge := util.GeneratePrimeAge()
	height := 0
	weight := 0
	potential := util.GeneratePotential()
	potentialGrade := util.GetWeightedPotentialGrade(potential)
	proPotential := util.GeneratePotential()
	stamina := util.GenerateStamina()

	personality := util.GetPersonality()
	academicBias := util.GetAcademicBias()
	workEthic := util.GetWorkEthic()
	recruitingBias := util.GetRecruitingBias()
	freeAgency := util.GetFreeAgencyBias(0, 0)
	archetype := util.GetArchetype(position)

	stars := util.GetStarRating(false, true)

	// Specialties
	insideShooting := getAttributeValue(position, archetype, stars, "InsideShooting", blob)
	midRangeShooting := getAttributeValue(position, archetype, stars, "MidRangeShooting", blob)
	threePointShooting := getAttributeValue(position, archetype, stars, "ThreePointShooting", blob)
	freeThrow := getAttributeValue(position, archetype, stars, "FreeThrow", blob)
	ballwork := getAttributeValue(position, archetype, stars, "Ballwork", blob)
	rebounding := getAttributeValue(position, archetype, stars, "Rebounding", blob)
	interiorDefense := getAttributeValue(position, archetype, stars, "InteriorDefense", blob)
	perimeterDefense := getAttributeValue(position, archetype, stars, "PerimeterDefense", blob)
	agility := getAttributeValue(position, archetype, stars, "Agility", blob)
	stealing := getAttributeValue(position, archetype, stars, "Stealing", blob)
	blocking := getAttributeValue(position, archetype, stars, "Blocking", blob)
	// Specialties
	specInsideShooting := util.GenerateSpecialty(position, archetype, "InsideShooting")
	specMidRangeShooting := util.GenerateSpecialty(position, archetype, "MidRangeShooting")
	specThreePointShooting := util.GenerateSpecialty(position, archetype, "ThreePointShooting")
	specFreeThrow := util.GenerateSpecialty(position, archetype, "FreeThrow")
	specBallwork := util.GenerateSpecialty(position, archetype, "Ballwork")
	specAgility := util.GenerateSpecialty(position, archetype, "Agility")
	specStealing := util.GenerateSpecialty(position, archetype, "Stealing")
	specBlocking := util.GenerateSpecialty(position, archetype, "Blocking")
	specRebounding := util.GenerateSpecialty(position, archetype, "Rebounding")
	specInteriorDefense := util.GenerateSpecialty(position, archetype, "InteriorDefense")
	specPerimeterDefense := util.GenerateSpecialty(position, archetype, "PerimeterDefense")

	var basePlayer = structs.BasePlayer{
		FirstName:              firstName,
		LastName:               lastName,
		Position:               position,
		Age:                    uint8(age),
		Year:                   uint8(year),
		State:                  "",
		Country:                country,
		Height:                 uint8(height),
		Weight:                 uint16(weight),
		Agility:                uint8(agility),
		InsideShooting:         uint8(insideShooting),
		MidRangeShooting:       uint8(midRangeShooting),
		ThreePointShooting:     uint8(threePointShooting),
		FreeThrow:              uint8(freeThrow),
		Ballwork:               uint8(ballwork),
		Stealing:               uint8(stealing),
		Blocking:               uint8(blocking),
		Rebounding:             uint8(rebounding),
		InteriorDefense:        uint8(interiorDefense),
		PerimeterDefense:       uint8(perimeterDefense),
		SpecInsideShooting:     specInsideShooting,
		SpecMidRangeShooting:   specMidRangeShooting,
		SpecThreePointShooting: specThreePointShooting,
		SpecFreeThrow:          specFreeThrow,
		SpecBallwork:           specBallwork,
		SpecAgility:            specAgility,
		SpecStealing:           specStealing,
		SpecBlocking:           specBlocking,
		SpecRebounding:         specRebounding,
		SpecInteriorDefense:    specInteriorDefense,
		SpecPerimeterDefense:   specPerimeterDefense,
		Potential:              potential,
		PotentialGrade:         potentialGrade,
		ProPotentialGrade:      proPotential,
		Stamina:                uint8(stamina),
		Personality:            personality,
		FreeAgency:             freeAgency,
		RecruitingBias:         recruitingBias,
		WorkEthic:              workEthic,
		AcademicBias:           academicBias, PlayerID: id,
		TeamID: teamID,
		Team:   team,
	}

	isNBAEligible := age > 22

	var player = structs.NBAPlayer{
		BasePlayer:      basePlayer,
		IsNBA:           isNBAEligible,
		IsInternational: true,
		IsIntGenerated:  true,
	}
	player.GetOverall()
	discipline := util.GenerateNormalizedIntFromRange(1, 20)
	injuryRating := util.GenerateNormalizedIntFromRange(1, 20)
	expectations := util.GetProfessionalPlaytimeExpectations(uint8(age), uint8(primeAge), uint8(player.Overall))
	player.SetExpectations(uint8(expectations))

	player.SetID(id)
	player.SetDisciplineAndIR(discipline, injuryRating)
	if age > 18 && age < 23 {
		diff := age - 18
		if diff > 3 {
			diff = 3
		}

		for i := 0; i < diff; i++ {
			player = ProgressNBAPlayer(player, true)
		}
	}
	player.SetAge(age)

	return player
}

func getNameList(ethnicity string, isFirstName bool) [][]string {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimNBA\\data"
	var fileName string
	switch ethnicity {
	case "Caucasian":
		if isFirstName {
			fileName = "FNameW.csv"
		} else {
			fileName = "LNameW.csv"
		}
	case "African":
		if isFirstName {
			fileName = "FNameB.csv"
		} else {
			fileName = "LNameB.csv"
		}
	case "Asian":
		if isFirstName {
			fileName = "FNameA.csv"
		} else {
			fileName = "LNameA.csv"
		}
	case "NativeAmerican":
		if isFirstName {
			fileName = "FNameN.csv"
		} else {
			fileName = "LNameN.csv"
		}
	default:
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

func getInternationalNameMap() map[string]map[string][]string {
	path := filepath.Join(os.Getenv("ROOT"), "data", "unique_male_names_by_country.json")
	content := util.ReadJson(path)
	var payload map[string]map[string][]string

	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatalln("Error during unmarshal: ", err)
	}

	return payload
}

func pickEthnicity() string {
	min := 0
	max := 10000
	num := util.GenerateIntFromRange(min, max)

	if num < 5000 {
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

func pickLocale(country string) string {
	countryMap := map[string][]string{
		"USA":                {"en_US", "es_ES"},
		"Antarctica":         {"en_US", "es_ES", "en_GB", "fr_CA", "en_CA", "en_IE", "it_IT", "ru_RU", "zh_CN"},
		"England":            {"en_GB", "en_US"},
		"United Kingdom":     {"en_GB", "en_US"},
		"Scotland":           {"en_GB", "en_IE"},
		"Austria":            {"de_AT"},
		"Canada":             {"fr_CA", "en_CA"},
		"Ireland":            {"en_IE"},
		"Wales":              {"en_GB", "en_IE"},
		"Spain":              {"es_ES"},
		"Malta":              {"es_ES"},
		"Italy":              {"it_IT"},
		"Portugal":           {"pt_PT"},
		"France":             {"fr_FR", "fr_CA"},
		"Switzerland":        {"fr_FR", "de_DE"},
		"Andorra":            {"fr_FR", "es_ES"},
		"Germany":            {"de_AT", "de_CH", "de_DE"},
		"Belgium":            {"nl_BE", "nl_NL", "fr_FR"},
		"Luxembourg":         {"fr_FR", "de_DE"},
		"Netherlands":        {"nl_BE", "nl_NL", "de_DE"},
		"Lithuania":          {"lt_LT"},
		"Latvia":             {"lv_LV", "lt_LT"},
		"Poland":             {"pl_PL"},
		"Finland":            {"sv_SE", "fi_FI"},
		"Denmark":            {"dk_DK", "no_NO"},
		"Sweden":             {"sv_SE", "no_NO"},
		"Iceland":            {"sv_SE", "no_NO"},
		"Norway":             {"no_NO"},
		"Bulgaria":           {"bg_BG", "ro_RO"},
		"Serbia":             {"bs_BA", "sl_SI", "ro_RO", "bg_BG"},
		"Croatia":            {"hu_HU", "sl_SI", "hr_HR"},
		"Hungary":            {"sl_SI", "hu_HU"},
		"Bosnia":             {"bs_BA", "ro_RO", "sl_SI"},
		"Czech Republic":     {"cs_CZ", "bg_BG"},
		"Slovakia":           {"cs_CZ", "bg_BG"},
		"Estonia":            {"et_EE", "lt_LT"},
		"Kosovo":             {"sl_SI", "ro_RO"},
		"Montenegro":         {"sl_SI", "ro_RO"},
		"Romania":            {"sl_SI", "ru_RU", "ro_RO", "bg_BG"},
		"Moldova":            {"uk_UA", "ru_RU", "ro_RO"},
		"Slovenia":           {"sl_SI", "ro_RO", "bg_BG"},
		"Cyprus":             {"el_GR", "tr_TR"},
		"Turkey":             {"tr_TR"},
		"Greece":             {"el_GR", "tr_TR"},
		"Albania":            {"el_GR"},
		"North Macedonia":    {"el_GR"},
		"Mexico":             {"es_MX"},
		"Argentina":          {"es_MX"},
		"Brazil":             {"es_MX", "pt_BR"},
		"China":              {"zh_CN"},
		"HK":                 {"zh_CN"},
		"Japan":              {"ja_JP"},
		"South Korea":        {"ko_KR"},
		"North Korea":        {"ko_KR"},
		"Taiwan":             {"zh_TW"},
		"Philippines":        {"en_PH", "es_ES"},
		"Phillipines":        {"en_PH", "es_ES"},
		"Indonesia":          {"ms_MY", "zh_CN"},
		"Malaysia":           {"ms_MY", "vi_VN", "th_TH", "zh_CN"},
		"Singapore":          {"zh_CN", "th_TH"},
		"Laos":               {"zh_CN", "vi_VN"},
		"Myanmar":            {"zh_CN", "th_TH"},
		"Cambodia":           {"zh_CN", "vi_VN"},
		"Thailand":           {"en_TH"},
		"Vietnam":            {"vi_VN"},
		"Papua New Guinea":   {"en_PH", "en_NZ"},
		"Solomon Islands":    {"en_PH", "en_NZ"},
		"New Caledonia":      {"en_PH", "en_NZ"},
		"Fiji":               {"en_PH", "en_NZ"},
		"French Polynesia":   {"en_PH", "en_NZ"},
		"Vanuatu":            {"en_PH", "en_NZ"},
		"Australia":          {"en_AU"},
		"New Zealand":        {"en_NZ"},
		"Chile":              {"es_MX"},
		"Colombia":           {"es_MX"},
		"Guatemala":          {"es_MX"},
		"Dominican Republic": {"es_MX"},
		"Grenada":            {"es_MX"},
		"Barbados":           {"es_MX", "zu_ZA"},
		"El Salvador":        {"es_MX"},
		"Belize":             {"es_MX"},
		"Honduras":           {"es_MX"},
		"Trinidad":           {"es_MX"},
		"French Guiana":      {"es_MX", "fr_FR"},
		"Jamaica":            {"es_MX", "zu_ZA"},
		"Haiti":              {"es_MX", "zu_ZA"},
		"The Bahamas":        {"es_MX", "zu_ZA"},
		"Costa Rica":         {"es_MX"},
		"Nicaragua":          {"es_MX"},
		"Panama":             {"es_MX"},
		"Cuba":               {"es_MX"},
		"Puerto Rico":        {"es_MX"},
		"Venezuela":          {"es_MX"},
		"Guyana":             {"es_MX"},
		"Ecuador":            {"es_MX"},
		"Suriname":           {"nl_NL", "es_MX", "pt_BR"},
		"Peru":               {"es_MX"},
		"Paraguay":           {"es_MX"},
		"Sierra Leone":       {"es_MX"},
		"Uruguay":            {"pt_PT", "es_MX", "pt_BR"},
		"Azerbaijan":         {"uk_UA", "hy_AM", "az_AZ"},
		"Georgia":            {"uk_UA", "hy_AM", "az_AZ"},
		"Armenia":            {"hy_AM", "az_AZ"},
		"Ukraine":            {"uk_UA"},
		"Russia":             {"ru_RU"},
		"Tajikistan":         {"ar_SA", "ru_RU"},
		"Kyrgyzstan":         {"zh_CN", "ru_RU"},
		"Kazakhstan":         {"tr_TR", "ru_RU"},
		"Uzbekistan":         {"tr_TR", "ru_RU"},
		"Turkmenistan":       {"ar_SA", "ru_RU"},
		"Mongolia":           {"ru_RU", "zh_CN"},
		"Nepal":              {"zh_CN"},
		"Bangladesh":         {"en_IN"},
		"India":              {"en_IN"},
		"Pakistan":           {"id_ID", "en_IN"},
		"Ethiopia":           {"sa_SA", "zu_ZA"},
		"Chad":               {"sa_SA"},
		"Senegal":            {"sa_SA"},
		"Algeria":            {"sa_SA", "ar_EG"},
		"Togo":               {"sa_SA"},
		"Cameroon":           {"sa_SA"},
		"Eritrea":            {"sa_SA"},
		"Liberia":            {"sa_SA"},
		"Libya":              {"sa_SA", "ar_EG"},
		"Tanzania":           {"sa_SA"},
		"Guinea":             {"sa_SA"},
		"The Gambia":         {"sa_SA"},
		"Mali":               {"sa_SA"},
		"Niger":              {"sa_SA"},
		"Nigeria":            {"sa_SA"},
		"Benin":              {"sa_SA"},
		"Gabon":              {"sa_SA"},
		"Angola":             {"sa_SA"},
		"Malawi":             {"sa_SA"},
		"Namibia":            {"sa_SA"},
		"Botswana":           {"sa_SA"},
		"South Africa":       {"sa_SA"},
		"Zimbabwe":           {"sa_SA"},
		"Mozambique":         {"sa_SA"},
		"Madagascar":         {"sa_SA"},
		"Kenya":              {"sa_SA"},
		"Somalia":            {"sa_SA"},
		"Djibouti":           {"sa_SA"},
		"Sudan":              {"sa_SA"},
		"Rwanda":             {"sa_SA"},
		"Uganda":             {"sa_SA"},
		"DRC":                {"sa_SA"},
		"South Sudan":        {"sa_SA"},
		"Burundi":            {"sa_SA"},
		"Ivory Coast":        {"sa_SA"},
		"Ghana":              {"sa_SA"},
		"Tunisia":            {"sa_SA", "ar_EG"},
		"Zambia":             {"sa_SA"},
		"Morocco":            {"ar_EG", "sa_SA"},
		"Egypt":              {"ar_EG"},
		"Palestine":          {"ar_EG", "ar_SA"},
		"Israel":             {"ar_JO"},
		"Jordan":             {"ar_JO"},
		"Lebanon":            {"ar_EG", "ar_SA", "ar_JO"},
		"Iraq":               {"ar_EG", "ar_SA"},
		"Iran":               {"ar_EG", "ar_SA"},
		"Saudi Arabia":       {"ar_EG", "ar_SA"},
		"Kuwait":             {"ar_EG", "ar_SA"},
		"Syria":              {"ar_EG", "ar_SA"},
		"Bahrain":            {"ar_EG", "ar_SA"},
		"Qatar":              {"ar_EG", "ar_SA"},
		"UAE":                {"ar_EG", "ar_SA"},
		"Yemen":              {"ar_EG", "ar_SA"},
	}
	selectedCountryCodes := countryMap[country]
	if len(selectedCountryCodes) == 0 {
		fmt.Println(country)
	}
	code := util.PickFromStringList(countryMap[country])
	return code
}

func pickISLCountry() string {
	randomNum := util.GenerateIntFromRange(1, 100)

	if randomNum < 96 {
		countries := []structs.ISLCountry{
			{Name: "Spain", Weight: 4}, {Name: "France", Weight: 4}, {Name: "Germany", Weight: 3}, {Name: "Turkey", Weight: 2}, {Name: "England", Weight: 2}, {Name: "Czech Republic", Weight: 2}, {Name: "Scotland", Weight: 1},
			{Name: "Andorra", Weight: 1}, {Name: "Belgium", Weight: 1}, {Name: "Netherlands", Weight: 1}, {Name: "Portugal", Weight: 1}, {Name: "China", Weight: 8}, {Name: "Japan", Weight: 5}, {Name: "South Korea", Weight: 4},
			{Name: "Russia", Weight: 4}, {Name: "HK", Weight: 1}, {Name: "Kazakhstan", Weight: 1}, {Name: "Mali", Weight: 2}, {Name: "Mozambique", Weight: 2}, {Name: "Nigeria", Weight: 2}, {Name: "Algeria", Weight: 1},
			{Name: "Angola", Weight: 1}, {Name: "Cameroon", Weight: 1}, {Name: "DRC", Weight: 1}, {Name: "Guinea", Weight: 1}, {Name: "Ivory Coast", Weight: 1}, {Name: "Madagascar", Weight: 1}, {Name: "Morocco", Weight: 1},
			{Name: "Rwanda", Weight: 1}, {Name: "Senegal", Weight: 1}, {Name: "South Africa", Weight: 1}, {Name: "South Sudan", Weight: 1}, {Name: "Tunisia", Weight: 1}, {Name: "Uganda", Weight: 1}, {Name: "Argentina", Weight: 5},
			{Name: "Brazil", Weight: 5}, {Name: "Mexico", Weight: 3}, {Name: "Chile", Weight: 2}, {Name: "Colombia", Weight: 1}, {Name: "Nicaragua", Weight: 1}, {Name: "Panama", Weight: 1}, {Name: "Puerto Rico", Weight: 1},
			{Name: "Uruguay", Weight: 1}, {Name: "Italy", Weight: 3}, {Name: "Serbia", Weight: 2}, {Name: "Israel", Weight: 2}, {Name: "Greece", Weight: 2}, {Name: "Cyprus", Weight: 2}, {Name: "Bulgaria", Weight: 1}, {Name: "Croatia", Weight: 1},
			{Name: "Hungary", Weight: 1}, {Name: "Kosovo", Weight: 1}, {Name: "Montenegro", Weight: 1}, {Name: "Romania", Weight: 1}, {Name: "Slovenia", Weight: 1}, {Name: "Ukraine", Weight: 3}, {Name: "Lithuania", Weight: 2},
			{Name: "Denmark", Weight: 2}, {Name: "Finland", Weight: 2}, {Name: "Iceland", Weight: 2}, {Name: "Norway", Weight: 2}, {Name: "Sweden", Weight: 2}, {Name: "Latvia", Weight: 1}, {Name: "Poland", Weight: 1}, {Name: "Australia", Weight: 5},
			{Name: "New Zealand", Weight: 3}, {Name: "Vietnam", Weight: 3}, {Name: "Philippines", Weight: 2}, {Name: "Taiwan", Weight: 2}, {Name: "Indonesia", Weight: 2}, {Name: "Malaysia", Weight: 1}, {Name: "Singapore", Weight: 1},
			{Name: "Thailand", Weight: 1}, {Name: "Egypt", Weight: 3}, {Name: "Bahrain", Weight: 2}, {Name: "Iran", Weight: 2}, {Name: "Kuwait", Weight: 2}, {Name: "Lebanon", Weight: 2}, {Name: "Saudi Arabia", Weight: 2},
			{Name: "Syria", Weight: 2}, {Name: "Azerbaijan", Weight: 1}, {Name: "Palestine", Weight: 1}, {Name: "Iraq", Weight: 1}, {Name: "Qatar", Weight: 1}, {Name: "UAE", Weight: 1}, {Name: "Slovakia", Weight: 1},
		}
		// Calculate the total weight
		totalWeight := 0
		for _, country := range countries {
			totalWeight += country.Weight
		}

		// Generate a random number between 0 and totalWeight-1
		randomWeight := util.GenerateIntFromRange(0, totalWeight)
		for _, country := range countries {
			if randomWeight < country.Weight {
				return country.Name
			}
			randomWeight -= country.Weight
		}
	}
	return util.PickFromStringList([]string{"Dominican Republic", "Canada", "The Bahamas", "Guatemala",
		"Ireland", "Wales", "Jamaica", "Costa Rica", "Belgium", "Colombia", "Belize", "Haiti", "El Salvador", "Ethiopia",
		"Cuba", "Laos", "Papua New Guinea", "Chad", "Honduras", "Nicaragua", "Panama", "Finland", "Senegal", "Algeria", "Togo",
		"Austria", "Hungary", "Venezuela", "Cameroon", "French Guiana", "Trinidad", "Croatia", "Eritrea", "Guyana",
		"Liberia", "Libya", "Tanzania", "Peru", "Paraguay", "Sierra Leone", "Guinea", "The Gambia", "Mali",
		"Niger", "Benin", "Gabon", "Angola", "Malawi", "Namibia", "Botswana", "Nepal", "India", "Bangladesh", "Myanmar", "Laos",
		"Cambodia", "Tajikistan", "Kyrgyzstan", "Pakistan", "Yemen", "Uzbekistan", "Turkmenistan", "Mongolia", "Solomon Islands",
		"New Caledonia", "Fiji", "French Polynesia", "Vanuatu", "Switzerland", "Malta", "Albania", "North Macedonia", "Moldova",
		"Georgia", "Armenia"})
}

func pickCountry(ethnicity string) string {
	min := 0
	max := 10000
	num := util.GenerateIntFromRange(min, max)

	if num < 8201 {
		return "USA"
	} else if num < 8251 {
		switch ethnicity {
		case "African":
			return "Dominican Republic"
		case "Hispanic":
			return "Mexico"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "China"
		default:
			return "Canada"
		}
	} else if num < 8301 {
		switch ethnicity {
		case "African":
			return "The Bahamas"
		case "Hispanic":
			return "Guatemala"
		case "NativeAmerican":
			return "Russia"
		case "Asian":
			return "China"
		default:
			return "United Kingdom"
		}
	} else if num < 8351 {
		switch ethnicity {
		case "African":
			return "Jamaica"
		case "Hispanic":
			return "Costa Rica"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "China"
		default:
			return "France"
		}
	} else if num < 8401 {
		switch ethnicity {
		case "African":
			return "DRC"
		case "Hispanic":
			return "Colombia"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "China"
		default:
			return "Spain"
		}
	} else if num < 8451 {
		switch ethnicity {
		case "African":
			return "South Africa"
		case "Hispanic":
			return "Belize"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "China"
		default:
			return "Ireland"
		}
	} else if num < 8501 {
		switch ethnicity {
		case "African":
			return "Haiti"
		case "Hispanic":
			return "El Salvador"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "China"
		default:
			return "Spain"
		}
	} else if num < 8551 {
		switch ethnicity {
		case "African":
			return "Ethiopia"
		case "Hispanic":
			return "Cuba"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Japan"
		default:
			return "Germany"
		}
	} else if num < 8601 {
		switch ethnicity {
		case "African":
			return "Chad"
		case "Hispanic":
			return "Honduras"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Japan"
		default:
			return "Poland"
		}
	} else if num < 8651 {
		switch ethnicity {
		case "African":
			return "Ghana"
		case "Hispanic":
			return "Nicaragua"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Japan"
		default:
			return "Sweden"
		}
	} else if num < 8701 {
		switch ethnicity {
		case "African":
			return "Guinea"
		case "Hispanic":
			return "Panama"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Vietnam"
		default:
			return "Norway"
		}
	} else if num < 8751 {
		switch ethnicity {
		case "African":
			return "Senegal"
		case "Hispanic":
			return "Dominican Republic"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Vietnam"
		default:
			return "Denmark"
		}
	} else if num < 8801 {
		switch ethnicity {
		case "African":
			return "Morocco"
		case "Hispanic":
			return "Mexico"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Indonesia"
		default:
			return "Portugal"
		}
	} else if num < 8851 {
		switch ethnicity {
		case "African":
			return "Algeria"
		case "Hispanic":
			return "Mexico"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Indonesia"
		default:
			return "Austria"
		}
	} else if num < 8901 {
		switch ethnicity {
		case "African":
			return "Nigeria"
		case "Hispanic":
			return "Venezuela"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Indonesia"
		default:
			return "Hungary"
		}
	} else if num < 8951 {
		switch ethnicity {
		case "African":
			return "Cameroon"
		case "Hispanic":
			return "French Guiana"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Indonesia"
		default:
			return "Croatia"
		}
	} else if num < 9001 {
		switch ethnicity {
		case "African":
			return "Egypt"
		case "Hispanic":
			return "Brazil"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Thailand"
		default:
			return "Greece"
		}
	} else if num < 9051 {
		switch ethnicity {
		case "African":
			return "Eritrea"
		case "Hispanic":
			return "Brazil"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Thailand"
		default:
			return "Israel"
		}
	} else if num < 9101 {
		switch ethnicity {
		case "African":
			return "Kenya"
		case "Hispanic":
			return "Guyana"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "South Korea"
		default:
			return "Bulgaria"
		}
	} else if num < 9151 {
		switch ethnicity {
		case "African":
			return "Liberia"
		case "Hispanic":
			return "Ecuador"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Malaysia"
		default:
			return "Romania"
		}
	} else if num < 9201 {
		switch ethnicity {
		case "African":
			return "Tanzania"
		case "Hispanic":
			return "Chile"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "India"
		default:
			return "Montenegro"
		}
	} else if num < 9251 {
		switch ethnicity {
		case "African":
			return "Zimbabwe"
		case "Hispanic":
			return "Uruguay"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "India"
		default:
			return "Turkey"
		}
	} else if num < 9301 {
		switch ethnicity {
		case "African":
			return "Malawi"
		case "Hispanic":
			return "Argentina"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "India"
		default:
			return "Serbia"
		}
	} else if num < 9351 {
		switch ethnicity {
		case "African":
			return "Senegal"
		case "Hispanic":
			return "Argentina"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Israel"
		default:
			return "Belgium"
		}
	} else if num < 9401 {
		switch ethnicity {
		case "African":
			return "Senegal"
		case "Hispanic":
			return "Argentina"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Bangladesh"
		default:
			return "Ukraine"
		}
	} else if num < 9501 {
		switch ethnicity {
		case "African":
			return "DRC"
		case "Hispanic":
			return "Uruguay"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Philippines"
		default:
			return "Ukraine"
		}
	} else if num < 9601 {
		switch ethnicity {
		case "African":
			return "Nigeria"
		case "Hispanic":
			return "Uruguay"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Philippines"
		default:
			return "Russia"
		}
	} else if num < 9701 {
		switch ethnicity {
		case "African":
			return "South Africa"
		case "Hispanic":
			return "Chile"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Philippines"
		default:
			return "Russia"
		}
	} else if num < 9801 {
		switch ethnicity {
		case "African":
			return "South Africa"
		case "Hispanic":
			return "Chile"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Singapore"
		default:
			return "Lithuania"
		}
	} else if num < 9901 {
		switch ethnicity {
		case "African":
			return "Uganda"
		case "Hispanic":
			return "Peru"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Cambodia"
		default:
			return "Estonia"
		}
	} else if num < 9951 {
		switch ethnicity {
		case "African":
			return "Zambia"
		case "Hispanic":
			return "Grenada"
		case "NativeAmerican":
			return "USA"
		case "Asian":
			return "Taiwan"
		default:
			return "Finland"
		}
	} else if num < 9976 {
		switch ethnicity {
		case "African":
			return "Tunisia"
		case "Hispanic":
			return "Barbados"
		case "NativeAmerican":
			return "Canada"
		case "Asian":
			return "Myanmar"
		default:
			return "Iceland"
		}
	} else {
		switch ethnicity {
		case "African":
			return "Algeria"
		case "Hispanic":
			return "Suriname"
		case "NativeAmerican":
			return "Antarctica"
		case "Asian":
			return "North Korea"
		default:
			return "Luxembourg"
		}
	}
}

func GetRecruitModifier(stars uint8) int {
	switch stars {
	case 5:
		return util.GenerateIntFromRange(80, 117)
	case 4:
		return util.GenerateIntFromRange(100, 125)
	case 3:
		return util.GenerateIntFromRange(117, 150)
	case 2:
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
	if roll < 900 {
		return 4
	}
	// Twin
	if roll < 950 {
		return 5
	}
	return 1
}

func GenerateInternationalPlayersByTeam() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	// var playerList []structs.CollegePlayer

	newID := lastPlayerRecord.ID + 1

	nameMap := getInternationalNameMap()

	blob := getAttributeBlob()

	islTeams := GetInternationalTeams()
	facesBlob := getFaceDataBlob()
	faces := []structs.FaceData{}
	positions := []string{"PG", "SG", "SF", "PF", "C"}
	globalPlayers := []structs.GlobalPlayer{}
	newIntPlayers := []structs.NBAPlayer{}
	contracts := []structs.NBAContract{}
	for _, t := range islTeams {
		teamID := strconv.Itoa(int(t.ID))
		roster := GetAllNBAPlayersByTeamID(teamID)
		count := len(roster)
		if count > 14 {
			continue
		}

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

		for _, r := range roster {
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

		positionList := []string{}

		for _, p := range positions {
			if !teamNeedsMap[p] {
				continue
			}
			maxCount := 4
			posCount := positionCount[p]
			if p == "PG" || p == "C" {
				maxCount = 3
			}

			diff := maxCount - posCount
			for i := 1; i <= diff; i++ {
				positionList = append(positionList, p)
				positionCount[p] += 1
				if p == "PG" || p == "C" {
					teamNeedsMap[p] = positionCount[p] < 3
				} else {
					teamNeedsMap[p] = positionCount[p] < 4
				}
			}
		}

		country := t.Country
		pickedEthnicity := pickLocale(country)
		countryNames := nameMap[pickedEthnicity]
		year := 1

		rand.Shuffle(len(positionList), func(i, j int) {
			positionList[i], positionList[j] = positionList[j], positionList[i]
		})

		for _, pos := range positionList {
			if count > 13 {
				break
			}
			player := createInternationalPlayer(0, "", country, pickedEthnicity, pos, year, countryNames["first_names"], countryNames["last_names"], newID, blob)
			if player.Overall > 79 {
				fmt.Printf("PING! %d", player.Overall)
			}
			player.SignWithTeam(t.ID, t.Team, false, 3)
			globalPlayer := structs.GlobalPlayer{
				CollegePlayerID: newID,
				RecruitID:       newID,
				NBAPlayerID:     newID,
			}

			globalPlayer.SetID(newID)
			globalPlayers = append(globalPlayers, globalPlayer)
			newIntPlayers = append(newIntPlayers, player)
			skinColor := getSkinColor(country)
			face := getFace(newID, 238, skinColor, facesBlob)
			faces = append(faces, face)
			c := structs.NBAContract{
				PlayerID:       newID,
				TeamID:         t.ID,
				Team:           t.Team,
				OriginalTeamID: t.ID,
				OriginalTeam:   t.Team,
				YearsRemaining: 3,
				ContractType:   "ISL",
				TotalRemaining: 3,
				Year1Total:     1,
				Year2Total:     1,
				Year3Total:     1,
				IsActive:       true,
			}
			contracts = append(contracts, c)
			newID++
			count++
		}
	}
	repository.CreateFaceRecordsBatch(db, faces, 100)
	repository.CreateGlobalRecordsBatch(db, globalPlayers, 100)
	repository.CreateNBAPlayerRecordsBatch(db, newIntPlayers, 100)
	repository.CreateProContractRecordsBatch(db, contracts, 100)
}

func GenerateCollegeWalkons() {
	db := dbprovider.GetInstance().GetDB()
	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	// var playerList []structs.CollegePlayer

	newID := lastPlayerRecord.ID + 1
	firstNameMap, lastNameMap := getNameMaps()
	// Get All User Teams
	teams := GetAllActiveCollegeTeams()
	collegePlayers := GetAllCollegePlayers()
	collegePlayerMapByTeamID := MakeCollegePlayerMapByTeamID(collegePlayers, true)
	positions := []string{"G", "F", "C"}
	facesBlob := getFaceDataBlob()
	faces := []structs.FaceData{}
	globalPlayers := []structs.GlobalPlayer{}
	collegePlayersToUpload := []structs.CollegePlayer{}

	generator := CrootGenerator{
		attributeBlob: getAttributeBlob(),
	}

	for _, team := range teams {
		if !team.IsUserCoached {
			continue
		}

		roster := collegePlayerMapByTeamID[team.ID]

		if len(roster) > 9 {
			continue
		}

		count := 0
		playersNeeded := 10 - len(roster)

		teamNeedsMap := make(map[string]bool)
		positionCount := make(map[string]int)

		if _, ok := teamNeedsMap["G"]; !ok {
			teamNeedsMap["G"] = true
		}
		if _, ok := teamNeedsMap["F"]; !ok {
			teamNeedsMap["F"] = true
		}
		if _, ok := teamNeedsMap["C"]; !ok {
			teamNeedsMap["C"] = true
		}

		if _, ok := positionCount["G"]; !ok {
			positionCount["PG"] = 0
		}
		if _, ok := positionCount["F"]; !ok {
			positionCount["SF"] = 0
		}
		if _, ok := positionCount["C"]; !ok {
			positionCount["C"] = 0
		}

		for _, r := range roster {
			positionCount[r.Position] += 1
		}

		if positionCount["G"] >= 5 {
			teamNeedsMap["G"] = false
		} else if positionCount["F"] >= 5 {
			teamNeedsMap["F"] = false
		} else if positionCount["C"] >= 3 {
			teamNeedsMap["C"] = false
		}

		positionList := []string{}

		for _, p := range positions {
			if !teamNeedsMap[p] {
				continue
			}
			maxCount := 4
			posCount := positionCount[p]
			if p == "C" {
				maxCount = 3
			}

			diff := maxCount - posCount
			for i := 1; i <= diff; i++ {
				positionList = append(positionList, p)
				positionCount[p] += 1
				if p == "C" {
					teamNeedsMap[p] = positionCount[p] < 3
				} else {
					teamNeedsMap[p] = positionCount[p] < 4
				}
			}
		}

		year := 1

		rand.Shuffle(len(positionList), func(i, j int) {
			positionList[i], positionList[j] = positionList[j], positionList[i]
		})

		for _, position := range positionList {
			if count >= playersNeeded {
				break
			}
			pickedEthnicity := pickEthnicity()
			player := createCollegePlayer(team, pickedEthnicity, position, year, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], newID, true, generator.attributeBlob)
			globalPlayer := structs.GlobalPlayer{
				Model:           gorm.Model{ID: newID},
				CollegePlayerID: newID,
				RecruitID:       newID,
				NBAPlayerID:     newID,
			}

			if player.ID == 0 {
				player.SetID(newID)
			}

			collegePlayersToUpload = append(collegePlayersToUpload, player)
			globalPlayers = append(globalPlayers, globalPlayer)
			skinColor := getSkinColor(player.Country)
			face := getFace(newID, 238, skinColor, facesBlob)
			faces = append(faces, face)

			newID++
			count++
		}
	}
	repository.CreateCollegePlayersRecordBatch(db, collegePlayersToUpload, 500)
	repository.CreateGlobalRecordsBatch(db, globalPlayers, 500)
	repository.CreateFaceRecordsBatch(db, faces, 500)
}

func CreateCustomCroots() {
	db := dbprovider.GetInstance().GetDB()
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimNBA\\data\\2025\\2025_Custom_Croot_Class.csv"
	crootCSV := util.ReadCSV(path)
	latestID := getLatestRecord(db)

	generator := CrootGenerator{
		attributeBlob:    getAttributeBlob(),
		usCrootLocations: getCrootLocations("HS"),
	}

	crootList := []structs.Recruit{}

	for idx, row := range crootCSV {
		if idx < 1 {
			continue
		}
		if row[0] == "" {
			break
		}
		firstName := row[0]
		lastName := row[1]
		position := row[2]
		state := row[4]
		country := row[5]
		attr1 := row[6]
		attr2 := row[7]
		crootFor := row[8]
		star := util.GetStarRating(true, false)
		croot := createRecruit(position, firstName, lastName, state, country, "", "", 1, star, latestID, generator.attributeBlob, generator.usCrootLocations[state])
		croot.SetID(latestID)
		croot.SetCustomCroot(crootFor)
		croot.SetCustomAttribute(attr1)
		croot.SetCustomAttribute(attr2)
		croot.GetOverall()
		croot.Stars = uint8(star)
		latestID++
		crootList = append(crootList, croot)
	}

	for _, croot := range crootList {
		gp := structs.GlobalPlayer{
			CollegePlayerID: croot.ID,
			NBAPlayerID:     croot.ID,
			RecruitID:       croot.ID,
		}

		gp.SetID(croot.ID)

		repository.CreateRecruitRecord(croot, db)
		repository.CreateGlobalPlayerRecord(gp, db)
	}
	AssignAllRecruitRanks()
}

func getCrootLocations(locale string) map[string][]structs.CrootLocation {
	path := filepath.Join(os.Getenv("ROOT"), "data", locale+".json")

	content := util.ReadJson(path)

	var payload map[string][]structs.CrootLocation
	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err)
	}

	return payload
}

func getAttributeBlob() map[string]map[string]map[string]map[string]interface{} {
	path := filepath.Join(os.Getenv("ROOT"), "data", "attributeBlob.json")

	content := util.ReadJson(path)

	var payload map[string]map[string]map[string]map[string]interface{}
	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err)
	}

	return payload
}

func getCityAndHighSchool(schools []structs.CrootLocation) (string, string) {
	if len(schools) == 0 {
		fmt.Println("NO SCHOOLS?!")
		return "", ""
	}
	randInt := util.GenerateIntFromRange(0, len(schools)-1)

	return schools[randInt].City, schools[randInt].HighSchool
}

func getAttributeValue(pos string, arch string, star int, attr string, blob map[string]map[string]map[string]map[string]interface{}) int {
	starStr := strconv.Itoa(star)
	switch pos {
	case "C":
		switch arch {
		case "Rim Protector":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Post Scorer":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Stretch Center":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "All-Around":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	case "F":
		switch arch {
		case "Power Forward":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Small Forward":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Point Forward":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Swingman":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Two-Way":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "All-Around":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	case "G":
		switch arch {
		case "Point Guard":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Shooting Guard":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Combo Guard":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "Slasher":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "3-and-D":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		case "All-Around":
			switch attr {
			case "InsideShooting", "MidRangeShooting", "ThreePointShooting", "FreeThrow",
				"Agility", "Ballwork", "Rebounding", "Stealing", "Blocking",
				"InteriorDefense", "PerimeterDefense":
				return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
			}
		}
		return getValueFromInterfaceRange(starStr, blob[pos][arch][attr])
	}
	return util.GenerateIntFromRange(5, 15)
}

func getValueFromInterfaceRange(star string, starMap map[string]interface{}) int {
	// Check if the key exists in the map
	u, exists := starMap[star]
	if !exists {
		fmt.Printf("Key '%s' not found in starMap.\n", star)
		return 0 // Return a default value
	}

	// Check if the value can be asserted as a slice of interfaces
	minMax, ok := u.([]interface{})
	if !ok {
		fmt.Printf("Value for key '%s' is not a slice of interfaces.\n", star)
		return 0 // Return a default value
	}

	// Ensure the slice has at least two elements
	if len(minMax) < 2 {
		fmt.Printf("Value for key '%s' does not have enough elements (expected at least 2).\n", star)
		return 0 // Return a default value
	}

	// Check if the first element is a float64
	min, ok := minMax[0].(float64)
	if !ok {
		fmt.Printf("First element of '%s' is not a float64.\n", star)
		return 0 // Return a default value
	}

	// Check if the second element is a float64
	max, ok := minMax[1].(float64)
	if !ok {
		fmt.Printf("Second element of '%s' is not a float64.\n", star)
		return 0 // Return a default value
	}

	// Generate a random value in the range [min, max]
	return util.GenerateIntFromRange(int(min), int(max))
}
