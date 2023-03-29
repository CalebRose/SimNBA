package structs

type FreeAgencyResponse struct {
	FreeAgents    []NBAPlayer
	WaiverPlayers []NBAPlayer
	TeamOffers    []NBAContractOffer
}
