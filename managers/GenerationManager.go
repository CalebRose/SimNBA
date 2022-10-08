package managers

import (
	"encoding/csv"
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
	potential := util.GenerateIntFromRange(25, 100)
	proPotential := util.GenerateIntFromRange(15, 100)
	stamina := util.GenerateIntFromRange(25, 38)
	shooting2 := getAttribute(position, "Shooting2")
	shooting3 := getAttribute(position, "Shooting3")
	finishing := getAttribute(position, "Finishing")
	ballwork := getAttribute(position, "Ballwork")
	rebounding := getAttribute(position, "Rebounding")
	defense := getAttribute(position, "Defense")

	overall := (int((shooting2 + shooting3) / 2)) + finishing + ballwork + rebounding + defense
	stars := getStarRating(overall)
	if stars == 5 {
		potential -= 25
		if potential < 0 {
			potential = 25
		}
	}
	expectations := getPlaytimeExpectations(stars, year)
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
		Finishing:            finishing,
		Ballwork:             ballwork,
		Rebounding:           rebounding,
		Defense:              defense,
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
		collegePlayer = ProgressPlayer(collegePlayer)
	}

	return collegePlayer
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

	if num < 7000 {
		return "USA"
	} else if num < 7100 {
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
	} else if num < 7200 {
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
	} else if num < 7300 {
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
	} else if num < 7400 {
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
	} else if num < 7500 {
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
	} else if num < 7600 {
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
	} else if num < 7700 {
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
	} else if num < 7800 {
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
	} else if num < 7900 {
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
	} else if num < 8000 {
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
	} else if num < 8100 {
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
	} else if num < 8200 {
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
	} else if num < 8300 {
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
	} else if num < 8400 {
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
	} else if num < 8500 {
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
	} else if num < 8600 {
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
	} else if num < 8700 {
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
	} else if num < 8800 {
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
	} else if num < 8900 {
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
	} else if num < 9000 {
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
	} else if num < 9100 {
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
	} else if num < 9200 {
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
	} else if num < 9300 {
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
	} else if num < 9400 {
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
	} else if num < 9500 {
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
	} else if num < 9600 {
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
	} else if num < 9700 {
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
	} else if num < 9800 {
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
	} else if num < 9900 {
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
	} else if num < 9950 {
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
	} else if num < 9975 {
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

func getAttribute(position string, attribute string) int {
	if position == "G" {
		if attribute == "Shooting2" {
			return util.GenerateIntFromRange(7, 17)
		} else if attribute == "Shooting3" {
			return util.GenerateIntFromRange(7, 17)
		} else if attribute == "Finishing" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Ballwork" {
			return util.GenerateIntFromRange(7, 17)
		} else if attribute == "Rebounding" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Defense" {
			return util.GenerateIntFromRange(1, 11)
		} else {
			return 1
		}
	} else if position == "F" {
		if attribute == "Shooting2" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Shooting3" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Finishing" {
			return util.GenerateIntFromRange(6, 16)
		} else if attribute == "Ballwork" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Rebounding" {
			return util.GenerateIntFromRange(4, 14)
		} else if attribute == "Defense" {
			return util.GenerateIntFromRange(4, 14)
		} else {
			return 1
		}
	} else if position == "C" {
		if attribute == "Shooting2" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Shooting3" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Finishing" {
			return util.GenerateIntFromRange(6, 16)
		} else if attribute == "Ballwork" {
			return util.GenerateIntFromRange(1, 11)
		} else if attribute == "Rebounding" {
			return util.GenerateIntFromRange(6, 16)
		} else if attribute == "Defense" {
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

func getPlaytimeExpectations(stars int, year int) int {
	if stars == 5 {
		if year == 4 {
			return util.GenerateIntFromRange(10, 29)
		} else if year == 3 {
			return util.GenerateIntFromRange(10, 25)
		} else if year == 2 {
			return util.GenerateIntFromRange(10, 20)
		}
		return util.GenerateIntFromRange(10, 15)
	} else if stars == 4 {
		if year == 4 {
			return util.GenerateIntFromRange(10, 25)
		} else if year == 3 {
			return util.GenerateIntFromRange(9, 20)
		} else if year == 2 {
			return util.GenerateIntFromRange(5, 17)
		}
		return util.GenerateIntFromRange(5, 15)
	} else if stars == 3 {
		if year == 4 {
			return util.GenerateIntFromRange(7, 21)
		} else if year == 3 {
			return util.GenerateIntFromRange(3, 17)
		} else if year == 2 {
			return util.GenerateIntFromRange(2, 13)
		}
		return util.GenerateIntFromRange(0, 10)
	} else if stars == 2 {
		if year == 4 {
			return util.GenerateIntFromRange(0, 13)
		} else if year == 3 {
			return util.GenerateIntFromRange(0, 13)
		} else if year == 2 {
			return util.GenerateIntFromRange(0, 9)
		}
		return util.GenerateIntFromRange(0, 6)
	} else {
		return util.GenerateIntFromRange(0, 5)
	}
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
