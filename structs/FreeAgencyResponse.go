package structs

type FreeAgencyResponse struct {
	FreeAgents     []NBAPlayer
	WaiverPlayers  []NBAPlayer
	GLeaguePlayers []NBAPlayer
	ISLPlayers     []NBAPlayer
	TeamOffers     []NBAContractOffer
	RosterCount    uint
}
