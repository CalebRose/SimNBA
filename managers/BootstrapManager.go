package managers

import (
	"log"
	"strconv"
	"sync"

	"github.com/CalebRose/SimNBA/dbprovider"
	"github.com/CalebRose/SimNBA/repository"
	"github.com/CalebRose/SimNBA/structs"
	"github.com/CalebRose/SimNBA/util"
)

type BootstrapData struct {
	AllCollegeTeams      []structs.Team
	CollegeTeam          structs.Team
	CollegeRosterMap     map[uint][]structs.CollegePlayer
	PortalPlayers        []structs.TransferPlayerResponse
	CollegeInjuryReport  []structs.CollegePlayer
	CollegeNotifications []structs.Notification
	CollegeGameplan      structs.Gameplan
	TopCBBPoints         []structs.CollegePlayer
	TopCBBAssists        []structs.CollegePlayer
	TopCBBRebounds       []structs.CollegePlayer
	NBATeam              structs.NBATeam
	AllProTeams          []structs.NBATeam
	ProNotifications     []structs.Notification
	NBAGameplan          structs.NBAGameplan
	FaceData             map[uint]structs.FaceDataResponse
}

type BootstrapDataTwo struct {
	CollegeNews          []structs.NewsLog
	CollegePolls         []structs.CollegePollOfficial
	TeamProfileMap       map[string]*structs.TeamRecruitingProfile
	CollegeStandings     []structs.CollegeStandings
	ProStandings         []structs.NBAStandings
	CapsheetMap          map[uint]structs.NBACapsheet
	ProRosterMap         map[uint][]structs.NBAPlayer
	TopNBAPoints         []structs.NBAPlayer
	TopNBAAssists        []structs.NBAPlayer
	TopNBARebounds       []structs.NBAPlayer
	ProInjuryReport      []structs.NBAPlayer
	GLeaguePlayers       []structs.NBAPlayer
	InternationalPlayers []structs.NBAPlayer
}

type BootstrapDataThree struct {
	Recruits        []structs.Croot
	RecruitProfiles []structs.PlayerRecruitProfile
	FreeAgentOffers []structs.NBAContractOffer
	WaiverOffers    []structs.NBAWaiverOffer
	ProNews         []structs.NewsLog
	AllCollegeGames []structs.Match
	AllProGames     []structs.NBAMatch
	ContractMap     map[uint]structs.NBAContract
	ExtensionMap    map[uint]structs.NBAExtensionOffer
}

func GetBootstrapTeams() BootstrapData {
	var wg sync.WaitGroup
	var (
		allCollegeTeams []structs.Team
		allProTeams     []structs.NBATeam
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		allCollegeTeams = GetAllActiveCollegeTeams()
	}()
	go func() {
		defer wg.Done()
		allProTeams = GetAllActiveNBATeams()
	}()

	wg.Wait()

	return BootstrapData{
		AllCollegeTeams: allCollegeTeams,
		AllProTeams:     allProTeams,
	}
}

func GetBootstrapData(collegeID, proID string) BootstrapData {
	log.Println("GetBootstrapData called with collegeID:", collegeID, "and proID:", proID)

	var wg sync.WaitGroup
	var mu sync.Mutex

	var (
		collegeTeam         structs.Team
		collegePlayers      []structs.CollegePlayer
		collegeRosterMap    map[uint][]structs.CollegePlayer
		portalPrep          []structs.CollegePlayer
		portalPlayers       []structs.TransferPlayerResponse
		collegeInjuryReport []structs.CollegePlayer
		collegeNoti         []structs.Notification
		collegeGameplan     structs.Gameplan
		cbbPoints           []structs.CollegePlayer
		cbbAssists          []structs.CollegePlayer
		cbbRebounds         []structs.CollegePlayer
		nbaTeam             structs.NBATeam
		proNotifications    []structs.Notification
		nbaGameplan         structs.NBAGameplan
		faceDataMap         map[uint]structs.FaceDataResponse
	)

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			mu.Lock()
			collegeTeam = GetTeamByTeamID(collegeID)
			collegeTeam.UpdateLatestInstance()
			repository.SaveCollegeTeamRecord(collegeTeam, dbprovider.GetInstance().GetDB())
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			collegePlayers = GetAllCollegePlayers()
			mu.Lock()
			collegeRosterMap = MakeCollegePlayerMapByTeamID(collegePlayers, true)
			collegeInjuryReport = MakeCollegeInjuryList(collegePlayers)
			portalPrep = MakeCollegePortalList(collegePlayers)
			mu.Unlock()

		}()
		go func() {
			defer wg.Done()
			collegeNoti = GetNotificationByTeamIDAndLeague("CBB", collegeID)
		}()

		wg.Wait()

		portalProfileMap := MakeFullTransferPortalProfileMap(portalPrep)
		portalPlayers = MakeTransferPortalPlayerResponseList(portalPrep, portalProfileMap)

		wg.Add(2)
		go func() {
			defer wg.Done()
			cbbStats := GetCollegePlayerSeasonStatsBySeason(seasonID)
			historicPlayers := GetAllHistoricCollegePlayers()
			cpFromHistoric := makeHistoricPlayerList(historicPlayers)
			mu.Lock()
			collegePlayers = append(collegePlayers, cpFromHistoric...)
			collegePlayerMap := MakeCollegePlayerMap(collegePlayers)
			cbbPoints = GetCBBOrderedListByStatType("POINTS", collegeTeam.ID, cbbStats, collegePlayerMap)
			cbbAssists = GetCBBOrderedListByStatType("ASSISTS", collegeTeam.ID, cbbStats, collegePlayerMap)
			cbbRebounds = GetCBBOrderedListByStatType("REBOUNDS", collegeTeam.ID, cbbStats, collegePlayerMap)
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			collegeGameplan = GetGameplansByTeam(collegeID)
		}()
	}
	if len(proID) > 0 && proID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			mu.Lock()
			nbaTeam = GetNBATeamByTeamID(proID)
			nbaTeam.UpdateLatestInstance()
			repository.SaveNBATeamRecord(nbaTeam, dbprovider.GetInstance().GetDB())
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			proNotifications = GetNotificationByTeamIDAndLeague("NFL", proID)
		}()
		go func() {
			defer wg.Done()
			nbaGameplan = GetNBAGameplanByTeam(proID)
		}()
		wg.Wait()
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		faceDataMap = GetAllFaces()
	}()

	wg.Wait()
	return BootstrapData{
		CollegeTeam:          collegeTeam,
		CollegeRosterMap:     collegeRosterMap,
		PortalPlayers:        portalPlayers,
		CollegeInjuryReport:  collegeInjuryReport,
		CollegeNotifications: collegeNoti,
		CollegeGameplan:      collegeGameplan,
		TopCBBPoints:         cbbPoints,
		TopCBBAssists:        cbbAssists,
		TopCBBRebounds:       cbbRebounds,
		NBATeam:              nbaTeam,
		ProNotifications:     proNotifications,
		NBAGameplan:          nbaGameplan,
		FaceData:             faceDataMap,
	}
}

func GetSecondBootstrapData(collegeID, proID string) BootstrapDataTwo {
	log.Println("GetSecondBootstrapData called with collegeID:", collegeID, "and proID:", proID)

	var wg sync.WaitGroup
	var mu sync.Mutex

	var (
		collegeNews          []structs.NewsLog
		collegePolls         []structs.CollegePollOfficial
		teamProfileMap       map[string]*structs.TeamRecruitingProfile
		collegeStandings     []structs.CollegeStandings
		proStandings         []structs.NBAStandings
		capsheetMap          map[uint]structs.NBACapsheet
		proRosterMap         map[uint][]structs.NBAPlayer
		gLeaguePlayers       []structs.NBAPlayer
		internationalPlayers []structs.NBAPlayer
		topNBAPoints         []structs.NBAPlayer
		topNBAAssists        []structs.NBAPlayer
		topNBARebounds       []structs.NBAPlayer
		proInjuryReport      []structs.NBAPlayer
	)

	ts := GetTimestamp()
	log.Println("Timestamp:", ts)

	seasonID := strconv.Itoa(int(ts.SeasonID))

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(4)
		go func() {
			defer wg.Done()
			log.Println("Fetching College News Logs...")
			collegeNews = GetAllCBBNewsLogs()
			log.Println("Fetched College News Logs, count:", len(collegeNews))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching College News Logs...")
			collegePolls = GetAllCollegePolls()
			log.Println("Fetched College News Logs, count:", len(collegeNews))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching Team Profile Map...")
			teamProfileMap = GetTeamProfileMap()
			log.Println("Fetched Team Profile Map, count:", len(teamProfileMap))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching College Standings for seasonID:", ts.SeasonID)
			collegeStandings = GetAllConferenceStandingsBySeasonID(seasonID)
			log.Println("Fetched College Standings, count:", len(collegeStandings))
		}()
		wg.Wait()
		log.Println("Completed all College data queries.")

	}
	if len(proID) > 0 && proID != "0" {
		nbaTeamID := util.ConvertStringToInt(proID)
		wg.Add(3)
		go func() {
			defer wg.Done()
			log.Println("Fetching NFL Standings for seasonID:", ts.SeasonID)
			proStandings = GetNBAStandingsBySeasonID(seasonID)
			log.Println("Fetched NFL Standings, count:", len(proStandings))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching Capsheet Map...")
			mu.Lock()
			capsheetMap = GetCapsheetMap()
			mu.Unlock()
			log.Println("Fetched Capsheet Map, count:", len(capsheetMap))
		}()
		go func() {
			defer wg.Done()
			log.Println("Fetching NFL Players for roster mapping...")
			proPlayers := GetAllNBAPlayers()
			nbaStats := GetNBAPlayerSeasonStatsBySeason(seasonID)

			mu.Lock()
			nbaPlayerMap := MakeNBAPlayerMap(proPlayers)
			proRosterMap = MakeNBAPlayerMapByTeamID(proPlayers, true)
			proInjuryReport = MakeProInjuryList(proPlayers)
			gLeaguePlayers = MakeGLeagueList(proPlayers)
			internationalPlayers = MakeInternationalList(proPlayers)
			topNBAPoints = getNBAOrderedListByStatType("POINTS", uint(nbaTeamID), nbaStats, nbaPlayerMap)
			topNBAAssists = getNBAOrderedListByStatType("ASSISTS", uint(nbaTeamID), nbaStats, nbaPlayerMap)
			topNBARebounds = getNBAOrderedListByStatType("REBOUNDS", uint(nbaTeamID), nbaStats, nbaPlayerMap)
			mu.Unlock()
			log.Println("Fetched NFL Players, roster count:", len(proRosterMap), "injured count:", len(proInjuryReport))
		}()

		wg.Wait()
		log.Println("Completed all Pro data queries.")
	}

	return BootstrapDataTwo{
		CollegeNews:          collegeNews,
		TeamProfileMap:       teamProfileMap,
		CollegeStandings:     collegeStandings,
		CollegePolls:         collegePolls,
		ProStandings:         proStandings,
		CapsheetMap:          capsheetMap,
		ProRosterMap:         proRosterMap,
		TopNBAPoints:         topNBAPoints,
		TopNBAAssists:        topNBAAssists,
		TopNBARebounds:       topNBARebounds,
		ProInjuryReport:      proInjuryReport,
		GLeaguePlayers:       gLeaguePlayers,
		InternationalPlayers: internationalPlayers,
	}
}

func GetThirdBootstrapData(collegeID, proID string) BootstrapDataThree {
	var wg sync.WaitGroup
	var (
		recruits        []structs.Croot
		recruitProfiles []structs.PlayerRecruitProfile
		freeAgentOffers []structs.NBAContractOffer
		waiverOffers    []structs.NBAWaiverOffer
		proNews         []structs.NewsLog
		allCollegeGames []structs.Match
		allProGames     []structs.NBAMatch
		contractMap     map[uint]structs.NBAContract
		extensionMap    map[uint]structs.NBAExtensionOffer
	)
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			recruits = GetAllCollegeRecruits()
		}()

		go func() {
			defer wg.Done()
			recruitProfiles = GetRecruitingProfilesByTeamId(collegeID)
		}()

		go func() {
			defer wg.Done()
			allCollegeGames = GetCBBMatchesBySeasonID(seasonID)
		}()
		wg.Wait()
	}
	if len(proID) > 0 && proID != "0" {
		wg.Add(6)

		go func() {
			defer wg.Done()
			allProGames = GetNBAMatchesBySeasonID(seasonID)
		}()
		go func() {
			defer wg.Done()
			proNews = GetAllNBANewsLogs()
		}()

		go func() {
			defer wg.Done()
			freeAgentOffers = repository.FindAllFreeAgentOffers(repository.FreeAgencyQuery{IsActive: true})
		}()

		go func() {
			defer wg.Done()
			waiverOffers = repository.FindAllWaiverOffers(repository.FreeAgencyQuery{IsActive: true})
		}()

		go func() {
			defer wg.Done()
			contractMap = GetContractMap()
		}()

		go func() {
			defer wg.Done()
			extensionMap = GetExtensionMap()
		}()

		wg.Wait()
	}
	return BootstrapDataThree{
		Recruits:        recruits,
		RecruitProfiles: recruitProfiles,
		ProNews:         proNews,
		FreeAgentOffers: freeAgentOffers,
		WaiverOffers:    waiverOffers,
		AllCollegeGames: allCollegeGames,
		AllProGames:     allProGames,
		ContractMap:     contractMap,
		ExtensionMap:    extensionMap,
	}
}
