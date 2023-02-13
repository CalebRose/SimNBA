package managers

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func GenerateNewTeams() {
	db := dbprovider.GetInstance().GetDB()
	rand.Seed(time.Now().Unix())

	var lastPlayerRecord structs.Player

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	// var playerList []structs.CollegePlayer

	newID := lastPlayerRecord.ID + 1

	teams := GetAllActiveCollegeTeams()
	firstNameMap, lastNameMap := getNameMaps()
	var positionList []string = []string{"G", "F", "C"}

	for _, team := range teams {
		// Test Generation
		yearList := []int{}
		players := GetCollegePlayersByTeamId(strconv.Itoa(int(team.ID)))
		seniors := 3
		juniors := 3
		sophomores := 3
		freshmen := 4
		fCount := 0
		gCount := 0
		cCount := 0
		requiredPlayers := 13
		count := 0
		if len(players) > 0 {
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
				if player.Position == "F" {
					fCount++
				} else if player.Position == "G" {
					gCount++
				} else {
					cCount++
				}
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
			if pickedPosition == "F" && fCount > 3 {
				quickList := []string{"G", "C"}
				pickedPosition = util.PickFromStringList(quickList)
			} else if pickedPosition == "G" && gCount > 3 {
				quickList := []string{"F", "C"}
				pickedPosition = util.PickFromStringList(quickList)
			} else if pickedPosition == "C" && cCount > 3 {
				quickList := []string{"F", "G"}
				pickedPosition = util.PickFromStringList(quickList)
			}

			if pickedPosition == "F" {
				fCount++
			} else if pickedPosition == "G" {
				gCount++
			} else {
				cCount++
			}

			positionQueue = append(positionQueue, pickedPosition)
		}

		for count < requiredPlayers {
			pickedEthnicity := pickEthnicity()
			pickedPosition := positionQueue[count]
			year := yearList[count]
			player := createCollegePlayer(team, pickedEthnicity, pickedPosition, year, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], newID)
			// playerList = append(playerList, player)
			err := db.Save(&player).Error
			if err != nil {
				log.Panicln("Could not save player record")
			}
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

	rand.Seed(time.Now().Unix())

	var lastPlayerRecord structs.GlobalPlayer

	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		log.Fatalln("Could not grab last player record from players table...")
	}

	// var playerList []structs.CollegePlayer

	newID := lastPlayerRecord.ID + 1

	firstNameMap, lastNameMap := getNameMaps()
	var positionList []string = []string{"G", "F", "C"}

	// Test Generation
	requiredPlayers := util.GenerateIntFromRange(621, 860)
	// 531 is the number of Seniors && Redshirt Seniors currently in the league
	// Currently 172 teams. 172 * 5 = 860, the max number of recruits that can be signed.
	// 172 * 3 = 516, the minimum required to sign.
	// 531 + 89 recruits left over = 620.
	count := 0

	for count < requiredPlayers {
		pickedPosition := util.PickFromStringList(positionList)
		pickedEthnicity := pickEthnicity()
		year := 1
		player := createRecruit(pickedEthnicity, pickedPosition, year, firstNameMap[pickedEthnicity], lastNameMap[pickedEthnicity], newID)
		// playerList = append(playerList, player)
		err := db.Create(&player).Error
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

// Private Methods
func createCollegePlayer(team structs.Team, ethnicity string, position string, year int, firstNameList [][]string, lastNameList [][]string, id uint) structs.CollegePlayer {
	fName := getName(firstNameList)
	lName := getName(lastNameList)

	firstName := strings.Title(strings.ToLower(fName))
	lastName := strings.Title(strings.ToLower(lName))
	age := 19
	if year == 4 {
		age = 22
	} else if year == 3 {
		age = 21
	} else if year == 2 {
		age = 20
	}
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

	potential -= util.GenerateIntFromRange(20, 30)

	if potential < 0 {
		potential = util.GenerateIntFromRange(5, 25)
	}

	expectations := util.GetPlaytimeExpectations(stars, year)
	personality := util.GetPersonality()
	academicBias := util.GetAcademicBias()
	workEthic := util.GetWorkEthic()
	recruitingBias := util.GetRecruitingBias()
	freeAgency := util.GetFreeAgencyBias()

	var basePlayer = structs.BasePlayer{
		FirstName:            firstName,
		LastName:             lastName,
		Position:             position,
		Age:                  age,
		Year:                 year,
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

	for i := 0; i < year && year > 1; i++ {
		collegePlayer = ProgressCollegePlayer(collegePlayer)
	}

	return collegePlayer
}

func createRecruit(ethnicity string, position string, year int, firstNameList [][]string, lastNameList [][]string, id uint) structs.Recruit {
	fName := getName(firstNameList)
	lName := getName(lastNameList)

	firstName := strings.Title(strings.ToLower(fName))
	lastName := strings.Title(strings.ToLower(lName))
	age := 18
	state := ""
	country := pickCountry(ethnicity)
	if country == "USA" {
		state = pickState()
	}
	height := getHeight(position)
	potential := util.GeneratePotential()
	potentialGrade := util.GetWeightedPotentialGrade(potential)
	proPotential := util.GeneratePotential()
	stamina := util.GenerateIntFromRange(25, 38)
	shooting2 := getAttribute(position, "Shooting2", false)
	shooting3 := getAttribute(position, "Shooting3", false)
	finishing := getAttribute(position, "Finishing", false)
	freeThrow := getAttribute(position, "FreeThrow", false)
	ballwork := getAttribute(position, "Ballwork", false)
	rebounding := getAttribute(position, "Rebounding", false)
	interiorDefense := getAttribute(position, "Interior Defense", false)
	perimeterDefense := getAttribute(position, "Perimeter Defense", false)

	overall := (int((shooting2 + shooting3 + freeThrow) / 3)) + finishing + ballwork + rebounding + int((interiorDefense+perimeterDefense)/2)
	stars := getStarRating(overall)
	recruitModifier := GetRecruitModifier(stars)
	expectations := util.GetPlaytimeExpectations(stars, year)
	personality := util.GetPersonality()
	academicBias := util.GetAcademicBias()
	workEthic := util.GetWorkEthic()
	recruitingBias := util.GetRecruitingBias()
	freeAgency := util.GetFreeAgencyBias()

	var basePlayer = structs.BasePlayer{
		FirstName:            firstName,
		LastName:             lastName,
		Position:             position,
		Age:                  age,
		Year:                 year,
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
		PotentialGrade:       potentialGrade,
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
	}

	var croot = structs.Recruit{
		BasePlayer:      basePlayer,
		PlayerID:        id,
		TeamID:          0,
		TeamAbbr:        "",
		RecruitModifier: recruitModifier,
		IsSigned:        false,
		IsTransfer:      false,
	}
	croot.SetID(id)

	return croot
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

func pickCountry(ethnicity string) string {
	min := 0
	max := 10000
	num := rand.Intn(max-min+1) + min

	if num < 7001 {
		return "USA"
	} else if num < 7101 {
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
	} else if num < 7201 {
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
	} else if num < 7301 {
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
	} else if num < 7401 {
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
	} else if num < 7501 {
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
	} else if num < 7601 {
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
	} else if num < 7701 {
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
	} else if num < 7801 {
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
	} else if num < 7901 {
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
	} else if num < 8001 {
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
	} else if num < 8101 {
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
	} else if num < 8201 {
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
	} else if num < 8301 {
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
	} else if num < 8401 {
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
	} else if num < 8501 {
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
	} else if num < 8601 {
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
	} else if num < 8701 {
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
	} else if num < 8801 {
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
	} else if num < 8901 {
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
	} else if num < 9001 {
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
	} else if num < 9101 {
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
	} else if num < 9201 {
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
	} else if num < 9301 {
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
	if position == "G" {
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
	} else if position == "F" {
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
	if isGeneration {
		return util.GenerateIntFromRange(1, 11)
	}
	if position == "G" {
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
	} else if position == "F" {
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
	if overall > 60 {
		return 5
	} else if overall > 52 {
		return 4
	} else if overall > 44 {
		return 3
	} else if overall > 36 {
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
