package managers

import (
	"log"
	"strconv"
	"sync"

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
}

type BootstrapDataTwo struct {
	CollegeNews      []structs.NewsLog
	TeamProfileMap   map[string]*structs.TeamRecruitingProfile
	CollegeStandings []structs.CollegeStandings
	ProStandings     []structs.NBAStandings
	CapsheetMap      map[uint]structs.NBACapsheet
	ProRosterMap     map[uint][]structs.NBAPlayer
	TopNBAPoints     []structs.NBAPlayer
	TopNBAAssists    []structs.NBAPlayer
	TopNBARebounds   []structs.NBAPlayer
	ProInjuryReport  []structs.NBAPlayer
}

type BootstrapDataThree struct {
	Recruits        []structs.Croot
	FreeAgency      structs.FreeAgencyResponse
	ProNews         []structs.NewsLog
	AllCollegeGames []structs.Match
	AllProGames     []structs.NBAMatch
}

func GetBootstrapData(collegeID, proID string) BootstrapData {
	log.Println("GetBootstrapData called with collegeID:", collegeID, "and proID:", proID)

	var wg sync.WaitGroup
	var mu sync.Mutex

	var (
		collegeTeam         structs.Team
		allCollegeTeams     []structs.Team
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
		allProTeams         []structs.NBATeam
		proNotifications    []structs.Notification
		nbaGameplan         structs.NBAGameplan
	)

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))

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

	if len(collegeID) > 0 {
		wg.Add(3)
		go func() {
			defer wg.Done()
			collegeTeam = GetTeamByTeamID(collegeID)
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

			mu.Lock()
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
	if len(proID) > 0 {
		wg.Add(3)
		go func() {
			defer wg.Done()
			nbaTeam = GetNBATeamByTeamID(proID)
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
	return BootstrapData{
		CollegeTeam:          collegeTeam,
		AllCollegeTeams:      allCollegeTeams,
		CollegeRosterMap:     collegeRosterMap,
		PortalPlayers:        portalPlayers,
		CollegeInjuryReport:  collegeInjuryReport,
		CollegeNotifications: collegeNoti,
		CollegeGameplan:      collegeGameplan,
		TopCBBPoints:         cbbPoints,
		TopCBBAssists:        cbbAssists,
		TopCBBRebounds:       cbbRebounds,
		AllProTeams:          allProTeams,
		NBATeam:              nbaTeam,
		ProNotifications:     proNotifications,
		NBAGameplan:          nbaGameplan,
	}
}

func GetSecondBootstrapData(collegeID, proID string) BootstrapDataTwo {
	log.Println("GetSecondBootstrapData called with collegeID:", collegeID, "and proID:", proID)

	var wg sync.WaitGroup
	var mu sync.Mutex

	var (
		collegeNews      []structs.NewsLog
		teamProfileMap   map[string]*structs.TeamRecruitingProfile
		collegeStandings []structs.CollegeStandings
		proStandings     []structs.NBAStandings
		capsheetMap      map[uint]structs.NBACapsheet
		proRosterMap     map[uint][]structs.NBAPlayer
		topNBAPoints     []structs.NBAPlayer
		topNBAAssists    []structs.NBAPlayer
		topNBARebounds   []structs.NBAPlayer
		proInjuryReport  []structs.NBAPlayer
	)

	ts := GetTimestamp()
	log.Println("Timestamp:", ts)

	seasonID := strconv.Itoa(int(ts.SeasonID))

	if len(collegeID) > 0 {
		wg.Add(3)
		go func() {
			defer wg.Done()
			log.Println("Fetching College News Logs...")
			collegeNews = GetAllCBBNewsLogs()
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
	if len(proID) > 0 {
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
			topNBAPoints = getNFLOrderedListByStatType("POINTS", uint(nbaTeamID), nbaStats, nbaPlayerMap)
			topNBAAssists = getNFLOrderedListByStatType("ASSISTS", uint(nbaTeamID), nbaStats, nbaPlayerMap)
			topNBARebounds = getNFLOrderedListByStatType("REBOUNDS", uint(nbaTeamID), nbaStats, nbaPlayerMap)
			mu.Unlock()
			log.Println("Fetched NFL Players, roster count:", len(proRosterMap), "injured count:", len(proInjuryReport))
		}()

		wg.Wait()
		log.Println("Completed all Pro data queries.")
	}

	return BootstrapDataTwo{
		CollegeNews:      collegeNews,
		TeamProfileMap:   teamProfileMap,
		CollegeStandings: collegeStandings,
		ProStandings:     proStandings,
		CapsheetMap:      capsheetMap,
		ProRosterMap:     proRosterMap,
		TopNBAPoints:     topNBAPoints,
		TopNBAAssists:    topNBAAssists,
		TopNBARebounds:   topNBARebounds,
		ProInjuryReport:  proInjuryReport,
	}
}

func GetThirdBootstrapData(collegeID, proID string) BootstrapDataThree {
	var wg sync.WaitGroup
	var (
		recruits        []structs.Croot
		freeAgency      structs.FreeAgencyResponse
		proNews         []structs.NewsLog
		allCollegeGames []structs.Match
		allProGames     []structs.NBAMatch
	)
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	freeAgencyCh := make(chan structs.FreeAgencyResponse, 1)

	if len(collegeID) > 0 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			recruits = GetAllCollegeRecruits()
		}()

		go func() {
			defer wg.Done()
			allCollegeGames = GetCBBMatchesBySeasonID(seasonID)
		}()
		wg.Wait()
	}
	if len(proID) > 0 {
		wg.Add(3)

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
			GetAllAvailableNBAPlayersViaChan(proID, freeAgencyCh)
		}()

		freeAgency = <-freeAgencyCh
		wg.Wait()
	}
	return BootstrapDataThree{
		Recruits:        recruits,
		FreeAgency:      freeAgency,
		ProNews:         proNews,
		AllCollegeGames: allCollegeGames,
		AllProGames:     allProGames,
	}
}
