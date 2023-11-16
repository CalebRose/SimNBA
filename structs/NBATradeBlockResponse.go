package structs

type NBATradeBlockResponse struct {
	Team                   NBATeam
	TradablePlayers        []NBAPlayer
	DraftPicks             []DraftPick
	SentTradeProposals     []NBATradeProposalDTO
	ReceivedTradeProposals []NBATradeProposalDTO
	TradePreferences       NBATradePreferences
}
