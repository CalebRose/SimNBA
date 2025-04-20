package managers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

func GetAllFaces() map[uint]structs.FaceDataResponse {
	faces := repository.FindFaceDataRecords(repository.FaceDataQuery{})
	faceBlob := getFaceDataBlob()

	// Precompute common blob lookups.
	accessoriesBlob := faceBlob["accessories"]
	bodyBlob := faceBlob["body"]
	earBlob := faceBlob["ear"]
	eyeBlob := faceBlob["eye"]
	eyeLineBlob := faceBlob["eyeLine"]
	eyebrowBlob := faceBlob["eyebrow"]
	facialHairBlob := faceBlob["facialHair"]
	glassesBlob := faceBlob["glasses"]
	hairBgBlob := faceBlob["hairBg"]
	headBlob := faceBlob["head"]
	jerseyBlob := faceBlob["jersey"]
	miscLineBlob := faceBlob["miscLine"]
	mouthBlob := faceBlob["mouth"]
	noseBlob := faceBlob["nose"]
	smileLineBlob := faceBlob["smileLine"]

	numFaces := len(faces)
	// We'll gather results in a channel and merge later.
	type faceResult struct {
		playerID uint
		data     structs.FaceDataResponse
	}
	resultCh := make(chan faceResult, numFaces)

	// Determine worker count and chunk size.
	numWorkers := runtime.NumCPU()
	chunkSize := (numFaces + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup

	// Process faces in parallel.
	for i := 0; i < numFaces; i += chunkSize {
		end := min(i+chunkSize, numFaces)
		wg.Add(1)
		// Capture the slice chunk.
		go func(facesChunk []structs.FaceData) {
			defer wg.Done()
			// Each goroutine gets its own buffer.
			buf := make([]byte, 0, 32)
			for _, face := range facesChunk {
				// Precompute dynamic blob lookups.
				// (Assuming face.SkinTone is a string field.)
				skinBlob := faceBlob[face.SkinTone+"Skin"]
				hairColorBlob := faceBlob[face.SkinTone+"HairColor"]
				hairBlob := faceBlob[face.SkinTone+"Hair"]
				// Build facialHairShave string using no-allocation methods.
				buf = buf[:0] // reset buffer
				buf = append(buf, "rgba(0,0,0,0."...)
				buf = strconv.AppendInt(buf, int64(face.FacialHairShave), 10)
				buf = append(buf, ')')
				facialHairShave := string(buf)

				resultCh <- faceResult{
					playerID: face.PlayerID,
					data: structs.FaceDataResponse{
						PlayerID:        face.PlayerID,
						Accessories:     accessoriesBlob[face.Accessories],
						Body:            bodyBlob[face.Body],
						BodySize:        face.BodySize,
						Ear:             earBlob[face.Ear],
						Eye:             eyeBlob[face.Eye],
						EyeLine:         eyeLineBlob[face.EyeLine],
						Eyebrow:         eyebrowBlob[face.Eyebrow],
						FacialHair:      facialHairBlob[face.FacialHair],
						Glasses:         glassesBlob[face.Glasses],
						Hair:            hairBlob[face.Hair],
						HairBG:          hairBgBlob[face.HairBG],
						HairFlip:        face.HairFlip,
						Head:            headBlob[face.Head],
						Jersey:          jerseyBlob[face.Jersey],
						MiscLine:        miscLineBlob[face.MiscLine],
						Mouth:           mouthBlob[face.Mouth],
						MouthFlip:       face.MouthFlip,
						Nose:            noseBlob[face.Nose],
						NoseFlip:        face.NoseFlip,
						SmileLine:       smileLineBlob[face.SmileLine],
						EarSize:         face.EarSize,
						EyeAngle:        face.EyeAngle,
						EyeBrowAngle:    face.EyeBrowAngle,
						FaceSize:        face.FaceSize,
						FacialHairShave: facialHairShave,
						NoseSize:        face.NoseSize,
						SmileLineSize:   face.SmileLineSize,
						SkinColor:       skinBlob[face.SkinColor],
						HairColor:       hairColorBlob[face.HairColor],
					},
				}
			}
		}(faces[i:end])
	}

	wg.Wait()
	close(resultCh)

	// Merge all results into the final map.
	finalMap := make(map[uint]structs.FaceDataResponse, numFaces)
	for res := range resultCh {
		finalMap[res.playerID] = res.data
	}
	return finalMap
}

func MigrateFaceDataToRecruits() {
	db := dbprovider.GetInstance().GetDB()
	// Get Recruits
	recruits := GetAllRecruitRecords()
	// Get Full Name Lists
	faceDataBlob := getFaceDataBlob()
	faceDataList := []structs.FaceData{}
	// Initialize List
	for _, r := range recruits {
		skinColor := getSkinColor(r.Country)
		// Store data

		face := getFace(r.ID, 238, skinColor, faceDataBlob)

		faceDataList = append(faceDataList, face)
	}

	repository.CreateFaceRecordsBatch(db, faceDataList, 500)
}

func MigrateFaceDataToHistoricCollegePlayers() {
	db := dbprovider.GetInstance().GetDB()
	// Get Recruits
	players := GetAllHistoricCollegePlayers()
	// Get Full Name Lists
	faceDataBlob := getFaceDataBlob()
	faceDataList := []structs.FaceData{}
	// Initialize List
	for _, p := range players {
		skinColor := getSkinColor(p.Country)
		// Store data

		face := getFace(p.ID, 238, skinColor, faceDataBlob)

		faceDataList = append(faceDataList, face)
	}

	repository.CreateFaceRecordsBatch(db, faceDataList, 500)
}

func MigrateFaceDataToCollegePlayers() {
	db := dbprovider.GetInstance().GetDB()
	// Get Recruits
	players := GetAllCollegePlayers()
	// Get Full Name Lists
	faceDataBlob := getFaceDataBlob()
	faceDataList := []structs.FaceData{}
	// Initialize List
	for _, p := range players {
		skinColor := getSkinColor(p.Country)
		// Store data

		face := getFace(p.ID, 238, skinColor, faceDataBlob)

		faceDataList = append(faceDataList, face)
	}

	repository.CreateFaceRecordsBatch(db, faceDataList, 500)
}

func MigrateFaceDataToProPlayers() {
	db := dbprovider.GetInstance().GetDB()
	// Get Recruits
	players := GetAllNBAPlayers()
	// Get Full Name Lists
	faceDataBlob := getFaceDataBlob()
	faceDataList := []structs.FaceData{}
	// Initialize List
	for _, p := range players {
		skinColor := getSkinColor(p.Country)
		// Store data

		face := getFace(p.ID, 238, skinColor, faceDataBlob)

		faceDataList = append(faceDataList, face)
	}

	repository.CreateFaceRecordsBatch(db, faceDataList, 500)
}

func getFace(id uint, weight int, ethnicity string, faceDataBlob map[string][]string) structs.FaceData {
	hairColorIdx := uint8(0)
	hairColorLen := len(faceDataBlob[ethnicity+"HairColor"]) - 1
	if hairColorLen > 0 {
		hairColorIdx = uint8(util.GenerateIntFromRange(0, len(faceDataBlob[ethnicity+"HairColor"])-1))
	}
	skinColorIdx := uint8(0)
	skinColorLen := len(faceDataBlob[ethnicity+"Skin"]) - 1
	if skinColorLen > 0 {
		skinColorIdx = uint8(util.GenerateIntFromRange(0, len(faceDataBlob[ethnicity+"Skin"])-1))
	}
	face := structs.FaceData{
		PlayerID:        id,
		Accessories:     uint8(util.GenerateIntFromRange(0, len(faceDataBlob["accessories"])-1)),
		Body:            getBodySize(weight),
		BodySize:        getBodyFat(weight),
		Ear:             uint8(util.GenerateIntFromRange(0, len(faceDataBlob["ear"])-1)),
		EarSize:         float32(util.GenerateFloatFromRange(0.5, 1.5)),
		Eye:             uint8(util.GenerateIntFromRange(0, len(faceDataBlob["eye"])-1)),
		EyeLine:         uint8(util.GenerateIntFromRange(0, len(faceDataBlob["eyeLine"])-1)),
		EyeAngle:        int8(util.GenerateIntFromRange(-10, 15)),
		Eyebrow:         uint8(util.GenerateIntFromRange(0, len(faceDataBlob["eyebrow"])-1)),
		EyeBrowAngle:    int8(util.GenerateIntFromRange(-15, 20)),
		FaceSize:        getFaceSize(weight),
		FacialHair:      getFacialHair(len(faceDataBlob["facialHair"]) - 1),
		FacialHairShave: getShaveStyle(),
		Glasses:         0,
		Hair:            uint8(util.GenerateIntFromRange(0, len(faceDataBlob[ethnicity+"Hair"])-1)),
		HairBG:          getHairBackground(),
		HairColor:       uint8(hairColorIdx),
		HairFlip:        util.GenerateIntFromRange(1, 2) == 1,
		Head:            uint8(util.GenerateIntFromRange(0, len(faceDataBlob["head"])-1)),
		Jersey:          uint8(util.GenerateIntFromRange(0, len(faceDataBlob["jersey"])-1)),
		MiscLine:        uint8(util.GenerateIntFromRange(0, len(faceDataBlob["miscLine"])-1)),
		Mouth:           uint8(util.GenerateIntFromRange(0, len(faceDataBlob["mouth"])-1)),
		MouthFlip:       util.GenerateIntFromRange(1, 2) == 1,
		Nose:            uint8(util.GenerateIntFromRange(0, len(faceDataBlob["nose"])-1)),
		NoseFlip:        util.GenerateIntFromRange(1, 2) == 1,
		NoseSize:        float32(util.GenerateFloatFromRange(0.5, 1.25)),
		SkinTone:        ethnicity,
		SkinColor:       skinColorIdx,
		SmileLine:       uint8(util.GenerateIntFromRange(0, len(faceDataBlob["smileLine"])-1)),
		SmileLineSize:   float32(util.GenerateFloatFromRange(0.25, 2.25)),
	}

	return face
}

func getSkinColor(country string) string {
	if len(country) == 0 {
		return util.PickFromStringList([]string{"asian", "black", "brown", "white"})
	}
	if country == "Sweden" || country == "Russia" || country == "Poland" || country == "Portugal" ||
		country == "Greenland" || country == "Finland" || country == "Czech Republic" ||
		country == "Switzerland" || country == "Austria" || country == "Albania" ||
		country == "Andorra" || country == "Armenia" || country == "Bulgaria" ||
		country == "Slovakia" || country == "Germany" || country == "Latvia" || country == "Norway" ||
		country == "Denmark" || country == "Netherlands" || country == "Belarus" || country == "Ukraine" ||
		country == "Croatia" || country == "Cyprus" || country == "Estonia" ||
		country == "Georgia" || country == "Iceland" || country == "Ireland" ||
		country == "Italy" || country == "Kosovo" || country == "Lithuania" ||
		country == "Luxembourg" || country == "Malta" || country == "Montenegro" ||
		country == "Romania" || country == "Scotland" || country == "Serbia" ||
		country == "Slovenia" || country == "Spain" {
		return "white"
	}

	if country == "Algeria" || country == "Angola" || country == "Cameroon" ||
		country == "Chad" || country == "DCR" || country == "Democratic Republic of Congo" ||
		country == "DRC" || country == "Ethiopia" || country == "Ghana" ||
		country == "Guinea" || country == "Ivory Coast" || country == "Kenya" ||
		country == "Liberia" || country == "Madagascar" || country == "Malawi" ||
		country == "Mali" || country == "Mozambique" || country == "Namibia" ||
		country == "Niger" || country == "Nigeria" || country == "Rwanda" ||
		country == "Senegal" || country == "Sierra Leone" || country == "South Sudan" ||
		country == "Tanzania" || country == "The Gambia" || country == "Togo" ||
		country == "Uganda" || country == "Zambia" || country == "Zimbabwe" {
		return "black"
	}

	if country == "Brazil" || country == "Ecuador" || country == "India" ||
		country == "Mexico" || country == "Peru" || country == "Argentina" ||
		country == "Belize" || country == "Chile" || country == "Costa Rica" ||
		country == "Cuba" || country == "Dominican Republic" || country == "Egypt" ||
		country == "El Salvador" || country == "Gabon" || country == "Guatemala" ||
		country == "Iraq" || country == "Iran" || country == "Kuwait" ||
		country == "Pakistan" || country == "Qatar" || country == "Saudi Arabia" ||
		country == "Syria" || country == "UAE" || country == "Uruguay" ||
		country == "Venezuela" || country == "Yemen" {
		return "brown"
	}

	if country == "China" || country == "Japan" || country == "Kazakhstan" ||
		country == "Taiwan" || country == "Cambodia" || country == "HK" ||
		country == "Indonesia" || country == "Laos" || country == "Malaysia" ||
		country == "Mongolia" || country == "Myanmar" || country == "Nepal" ||
		country == "Philippines" || country == "Phillipines" || country == "Singapore" ||
		country == "South Korea" || country == "Thailand" || country == "Uzbekistan" ||
		country == "Vietnam" || country == "North Korea" {
		return "asian"
	}

	if country == "Colombia" || country == "Greece" || country == "New Zealand" ||
		country == "Paraguay" {
		return util.PickFromStringList([]string{"white", "brown", "black"})
	}

	if country == "French Guiana" || country == "Guyana" || country == "Haiti" ||
		country == "Honduras" || country == "Jamaica" || country == "Libya" ||
		country == "Morocco" || country == "New Caledonia" || country == "Nicaragua" ||
		country == "Panama" || country == "Papua New Guinea" || country == "The Bahamas" ||
		country == "Tunisia" || country == "Vanuatu" {
		return util.PickFromStringList([]string{"brown", "black"})
	}

	if country == "Azerbaijan" || country == "Israel" || country == "Lebanon" ||
		country == "Palestine" || country == "Puerto Rico" || country == "Turkey" {
		return util.PickFromStringList([]string{"white", "brown"})
	}

	if country == "France" || country == "South Africa" {
		return util.PickFromStringList([]string{"black", "white"})
	}

	if country == "Australia" || country == "England" || country == "UK" {
		return util.PickFromStringList([]string{"asian", "black", "brown", "white"})
	}
	if country == "USA" || country == "Canada" {
		return util.PickFromStringList([]string{"asian", "black", "black", "black", "black", "black", "brown", "white"})

	}
	return util.PickFromStringList([]string{"asian", "black", "black", "black", "black", "black", "brown", "white"})
}

func getFaceDataBlob() map[string][]string {
	path := filepath.Join(os.Getenv("ROOT"), "data", "faceData.json")

	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}

	var payload map[string][]string
	err = json.Unmarshal(f, &payload)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err)
	}

	return payload
}

func getHairBackground() uint8 {
	dr := util.GenerateNormalizedIntFromRange(1, 100)

	if dr < 94 {
		return 1
	}
	if dr < 98 {
		return 0
	}
	return 2
}

func getBodyFat(weight int) float32 {
	if weight < 240 {
		return float32(util.GenerateFloatFromRange(0.8, 0.96))
	}
	return float32(util.GenerateFloatFromRange(0.97, 1.2))
}

func getFaceSize(weight int) float32 {
	if weight < 240 {
		return float32(util.GenerateFloatFromRange(0, 0.6))
	}
	return float32(util.GenerateFloatFromRange(0.60001, 1))
}

func getShaveStyle() uint8 {
	dr := util.GenerateIntFromRange(1, 100)
	if dr < 60 {
		return 0
	}
	return uint8(util.GenerateIntFromRange(0, 5))
}

func getBodySize(weight int) uint8 {
	if weight < 241 {
		return uint8(util.GenerateIntFromRange(0, 2))
	}
	return uint8(util.GenerateIntFromRange(0, 4))
}

func getFacialHair(len int) uint8 {
	dr := util.GenerateIntFromRange(1, 100)
	if dr < 70 {
		return 0
	}
	return uint8(util.GenerateIntFromRange(0, len))
}
