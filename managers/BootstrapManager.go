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
	AllCollegeTeams        []structs.Team
	CollegeTeam            structs.Team
	CollegeRosterMap       map[uint][]structs.CollegePlayer
	PortalPlayers          []structs.TransferPlayerResponse
	HistoricCollegePlayers []structs.HistoricCollegePlayer
	RetiredPlayers         []structs.RetiredPlayer
	CollegeInjuryReport    []structs.CollegePlayer
	CollegeNotifications   []structs.Notification
	CollegeGameplan        structs.Gameplan
	TopCBBPoints           []structs.CollegePlayer
	TopCBBAssists          []structs.CollegePlayer
	TopCBBRebounds         []structs.CollegePlayer
	NBATeam                structs.NBATeam
	AllProTeams            []structs.NBATeam
	ProNotifications       []structs.Notification
	NBAGameplan            structs.NBAGameplan
	FaceData               map[uint]structs.FaceDataResponse
	CollegeNews            []structs.NewsLog
	CollegePolls           []structs.CollegePollOfficial
	TeamProfileMap         map[string]*structs.TeamRecruitingProfile
	CollegeStandings       []structs.CollegeStandings
	ProStandings           []structs.NBAStandings
	CapsheetMap            map[uint]structs.NBACapsheet
	ProRosterMap           map[uint][]structs.NBAPlayer
	TopNBAPoints           []structs.NBAPlayer
	TopNBAAssists          []structs.NBAPlayer
	TopNBARebounds         []structs.NBAPlayer
	ProInjuryReport        []structs.NBAPlayer
	GLeaguePlayers         []structs.NBAPlayer
	InternationalPlayers   []structs.NBAPlayer
	Recruits               []structs.Croot
	RecruitProfiles        []structs.PlayerRecruitProfile
	FreeAgentOffers        []structs.NBAContractOffer
	WaiverOffers           []structs.NBAWaiverOffer
	ProNews                []structs.NewsLog
	AllCollegeGames        []structs.Match
	AllProGames            []structs.NBAMatch
	ContractMap            map[uint]structs.NBAContract
	ExtensionMap           map[uint]structs.NBAExtensionOffer
	CollegePromises        []structs.CollegePromise
	TradeProposals         structs.NBATeamProposals
	TradePreferencesMap    map[uint]structs.NBATradePreferences
	DraftPicks             []structs.DraftPick
	FreeAgents             []structs.NBAPlayer
	WaiverPlayers          []structs.NBAPlayer
	PollSubmission         structs.CollegePollSubmission
	NBADraftees            []structs.NBADraftee
	WarRoomMap             map[uint]structs.NBAWarRoom
	ScoutingProfileMap     map[uint]structs.ScoutingProfile
	TransferPortalProfiles []structs.TransferPortalProfile
	CollegeGameplanMap     map[uint]structs.Gameplan
	ProGameplanMap         map[uint]structs.NBAGameplan
}

type BootstrapDataNews struct {
	CollegeNews []structs.NewsLog
	ProNews     []structs.NewsLog
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

func GetBootstrapDataLanding(collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var (
		collegeTeam            structs.Team
		collegePlayers         []structs.CollegePlayer
		collegeRosterMap       map[uint][]structs.CollegePlayer
		historicCollegePlayers []structs.HistoricCollegePlayer
		portalPrep             []structs.CollegePlayer
		portalPlayers          []structs.TransferPlayerResponse
		collegeInjuryReport    []structs.CollegePlayer
		collegeNoti            []structs.Notification
		collegeGameplan        structs.Gameplan
		cbbPoints              []structs.CollegePlayer
		cbbAssists             []structs.CollegePlayer
		cbbRebounds            []structs.CollegePlayer
		collegeStandings       []structs.CollegeStandings
		allCollegeGames        []structs.Match
		collegePolls           []structs.CollegePollOfficial
		nbaTeam                structs.NBATeam
		nbaGameplan            structs.NBAGameplan
		proNotifications       []structs.Notification
		topNBAPoints           []structs.NBAPlayer
		topNBAAssists          []structs.NBAPlayer
		topNBARebounds         []structs.NBAPlayer
		proRosterMap           map[uint][]structs.NBAPlayer
		gLeaguePlayers         []structs.NBAPlayer
		internationalPlayers   []structs.NBAPlayer
		proInjuryReport        []structs.NBAPlayer
		capsheetMap            map[uint]structs.NBACapsheet
		retiredPlayers         []structs.RetiredPlayer
		allProGames            []structs.NBAMatch
		proStandings           []structs.NBAStandings
	)

	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(6)
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

		go func() {
			defer wg.Done()
			collegeGameplan = GetGameplansByTeam(collegeID)
		}()

		go func() {
			defer wg.Done()
			collegePolls = GetAllCollegePolls()
		}()

		go func() {
			defer wg.Done()
			collegeStandings = GetAllConferenceStandingsBySeasonID(seasonID)
		}()

		wg.Wait()

		portalProfileMap := MakeFullTransferPortalProfileMap()
		portalPlayers = MakeTransferPortalPlayerResponseList(portalPrep, portalProfileMap)

		wg.Add(3)
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
			allCollegeGames = GetCBBMatchesBySeasonID(seasonID)
		}()

		go func() {
			defer wg.Done()
			historicCollegePlayers = GetAllHistoricCollegePlayers()
		}()

		wg.Wait()
	}

	if len(proID) > 0 && proID != "0" {
		nbaTeamID := util.ConvertStringToInt(proID)

		wg.Add(5)
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

		go func() {
			defer wg.Done()
			proStandings = GetNBAStandingsBySeasonID(seasonID)
		}()

		go func() {
			defer wg.Done()
			mu.Lock()
			capsheetMap = GetCapsheetMap()
			mu.Unlock()
		}()

		wg.Wait()

		wg.Add(3)

		go func() {
			defer wg.Done()
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
		}()

		go func() {
			defer wg.Done()
			allProGames = GetNBAMatchesBySeasonID(seasonID)
		}()
		go func() {
			defer wg.Done()
			retiredPlayers = GetAllRetiredPlayers()
		}()

		wg.Wait()
	}

	return BootstrapData{
		CollegeTeam:            collegeTeam,
		CollegeRosterMap:       collegeRosterMap,
		PortalPlayers:          portalPlayers,
		CollegeInjuryReport:    collegeInjuryReport,
		CollegeNotifications:   collegeNoti,
		CollegeGameplan:        collegeGameplan,
		TopCBBPoints:           cbbPoints,
		TopCBBAssists:          cbbAssists,
		TopCBBRebounds:         cbbRebounds,
		CollegePolls:           collegePolls,
		CollegeStandings:       collegeStandings,
		AllCollegeGames:        allCollegeGames,
		NBATeam:                nbaTeam,
		ProNotifications:       proNotifications,
		NBAGameplan:            nbaGameplan,
		ProStandings:           proStandings,
		CapsheetMap:            capsheetMap,
		ProRosterMap:           proRosterMap,
		TopNBAPoints:           topNBAPoints,
		TopNBAAssists:          topNBAAssists,
		TopNBARebounds:         topNBARebounds,
		ProInjuryReport:        proInjuryReport,
		GLeaguePlayers:         gLeaguePlayers,
		InternationalPlayers:   internationalPlayers,
		AllProGames:            allProGames,
		HistoricCollegePlayers: historicCollegePlayers,
		RetiredPlayers:         retiredPlayers,
	}
}

func GetBootstrapDataTeamRoster(collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup

	var (
		contractMap         map[uint]structs.NBAContract
		extensionMap        map[uint]structs.NBAExtensionOffer
		collegePromises     []structs.CollegePromise
		tradeProposals      structs.NBATeamProposals
		tradePreferencesMap map[uint]structs.NBATradePreferences
		draftPicks          []structs.DraftPick
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			promises := GetAllCollegePromisesByTeamID(collegeID)
			collegePromises = promises
		}()
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(5)
		go func() {
			defer wg.Done()
			contractMap = GetContractMap()
		}()

		go func() {
			defer wg.Done()
			extensionMap = GetExtensionMap()
		}()

		go func() {
			defer wg.Done()
			tradeProposals = GetTradeProposalsByNBAID(proID)
		}()

		go func() {
			defer wg.Done()
			tradePreferencesMap = GetTradePreferencesMap()
		}()
		go func() {
			defer wg.Done()
			draftPicks = GetAllRelevantDraftPicks()
		}()
	}

	wg.Wait()
	return BootstrapData{
		ContractMap:         contractMap,
		ExtensionMap:        extensionMap,
		CollegePromises:     collegePromises,
		TradeProposals:      tradeProposals,
		TradePreferencesMap: tradePreferencesMap,
		DraftPicks:          draftPicks,
	}
}

func GetBootstrapDataRecruiting(collegeID string) BootstrapData {
	var wg sync.WaitGroup

	var (
		teamProfileMap  map[string]*structs.TeamRecruitingProfile
		recruits        []structs.Croot
		recruitProfiles []structs.PlayerRecruitProfile
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			recruits = GetAllCollegeRecruits()
		}()

		go func() {
			defer wg.Done()
			recruitProfiles = repository.FindRecruitPlayerProfileRecords(collegeID, "", false, false, true)
		}()

		go func() {
			defer wg.Done()
			teamProfileMap = GetTeamProfileMap()
		}()
	}

	wg.Wait()
	return BootstrapData{
		Recruits:        recruits,
		RecruitProfiles: recruitProfiles,
		TeamProfileMap:  teamProfileMap,
	}
}

func GetBootstrapDataFreeAgency(proID string) BootstrapData {
	var wg sync.WaitGroup

	var (
		freeAgents      []structs.NBAPlayer
		waiverPlayers   []structs.NBAPlayer
		freeAgentOffers []structs.NBAContractOffer
		waiverOffers    []structs.NBAWaiverOffer
	)

	if len(proID) > 0 && proID != "0" {
		wg.Add(4)

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
			freeAgents = GetAllFreeAgents()
		}()

		go func() {
			defer wg.Done()
			waiverPlayers = GetAllWaiverWirePlayers()
		}()

	}

	wg.Wait()
	return BootstrapData{
		FreeAgentOffers: freeAgentOffers,
		WaiverOffers:    waiverOffers,
		FreeAgents:      freeAgents,
		WaiverPlayers:   waiverPlayers,
	}
}

func GetBootstrapDataScheduling(username, collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup
	var (
		officialPolls  []structs.CollegePollOfficial
		pollSubmission structs.CollegePollSubmission
	)
	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(2)
		go func() {
			defer wg.Done()
			officialPolls = GetAllCollegePolls()
		}()
		go func() {
			defer wg.Done()
			pollSubmission = GetPollSubmissionByUsernameWeekAndSeason(username)
		}()
	}

	wg.Wait()
	return BootstrapData{
		PollSubmission: pollSubmission,
		CollegePolls:   officialPolls,
	}
}

func GetBootstrapDataDraft(proID string) BootstrapData {
	var wg sync.WaitGroup

	var (
		nbaDraftees        []structs.NBADraftee
		warRoomMap         map[uint]structs.NBAWarRoom      // BY TEAM
		scoutingProfileMap map[uint]structs.ScoutingProfile // By TEAM
	)

	if len(proID) > 0 && proID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			nbaDraftees = GetAllNBADraftees()
		}()

		go func() {
			defer wg.Done()
			nbaWarRooms := GetAllNBAWarRooms()
			warRoomMap = MakeNBAWarRoomMap(nbaWarRooms)
		}()

		go func() {
			defer wg.Done()
			scoutingProfiles := GetAllScoutingProfiles()
			scoutingProfileMap = MakeScoutingProfileMapByTeam(scoutingProfiles)

		}()

		log.Println("Initiated all Pro data queries.")
	}
	wg.Wait()
	return BootstrapData{
		NBADraftees:        nbaDraftees,
		WarRoomMap:         warRoomMap,
		ScoutingProfileMap: scoutingProfileMap,
	}
}

func GetBootstrapDataPortal(collegeID string) BootstrapData {
	var wg sync.WaitGroup
	var (
		teamProfileMap         map[string]*structs.TeamRecruitingProfile // Get Just in Case because this page also uses this data
		transferPortalProfiles []structs.TransferPortalProfile
		collegePromises        []structs.CollegePromise
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(3)
		go func() {
			defer wg.Done()
			transferPortalProfiles = GetActiveTransferPortalProfiles()
		}()

		go func() {
			defer wg.Done()
			teamProfileMap = GetTeamProfileMap()
		}()

		go func() {
			defer wg.Done()
			promises := GetAllCollegePromisesByTeamID(collegeID)
			collegePromises = promises
		}()

	}

	wg.Wait()
	return BootstrapData{
		TransferPortalProfiles: transferPortalProfiles,
		TeamProfileMap:         teamProfileMap,
		CollegePromises:        collegePromises,
	}
}

func GetBootstrapDataGameplan(collegeID, proID string) BootstrapData {
	var wg sync.WaitGroup

	var (
		collegeGameplanMap map[uint]structs.Gameplan
		proGameplanMap     map[uint]structs.NBAGameplan
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			collegeGameplans := GetAllCollegeGameplans()
			collegeGameplanMap = MakeCollegeGameplanMap(collegeGameplans)
		}()
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gameplans := GetAllNBAGameplans()
			proGameplanMap = MakeNBAGameplanMap(gameplans)
		}()
	}

	wg.Wait()

	return BootstrapData{
		CollegeGameplanMap: collegeGameplanMap,
		ProGameplanMap:     proGameplanMap,
	}
}

func GetBootstrapDataStats(collegeID, proID string) BootstrapData {
	return BootstrapData{}
}

func GetNewsBootstrap(collegeID, proID string) BootstrapDataNews {
	var wg sync.WaitGroup

	var (
		collegeNews []structs.NewsLog
		proNews     []structs.NewsLog
	)

	if len(collegeID) > 0 && collegeID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Fetching College News Logs...")
			collegeNews = GetAllCBBNewsLogs()
			log.Println("Fetched College News Logs, count:", len(collegeNews))
		}()
		log.Println("Initiated all College data queries.")
	}

	if len(proID) > 0 && proID != "0" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			proNews = GetAllNBANewsLogs()
		}()

	}

	wg.Wait()

	return BootstrapDataNews{
		CollegeNews: collegeNews,
		ProNews:     proNews,
	}
}
