package structs

import (
	"gorm.io/gorm"
)

type CollegePollSubmission struct {
	gorm.Model
	Username     string
	SeasonID     uint
	WeekID       uint
	Week         uint
	RankOne      string
	RankOneID    uint
	RankTwo      string
	RankTwoID    uint
	RankThree    string
	RankThreeID  uint
	RankFour     string
	RankFourID   uint
	RankFive     string
	RankFiveID   uint
	RankSix      string
	RankSixID    uint
	RankSeven    string
	RankSevenID  uint
	RankEight    string
	RankEightID  uint
	RankNine     string
	RankNineID   uint
	RankTen      string
	RankTenID    uint
	RankEleven   string
	RankElevenID uint
	Rank12       string
	Rank12ID     uint
	Rank13       string
	Rank13ID     uint
	Rank14       string
	Rank14ID     uint
	Rank15       string
	Rank15ID     uint
	Rank16       string
	Rank16ID     uint
	Rank17       string
	Rank17ID     uint
	Rank18       string
	Rank18ID     uint
	Rank19       string
	Rank19ID     uint
	Rank20       string
	Rank20ID     uint
	Rank21       string
	Rank21ID     uint
	Rank22       string
	Rank22ID     uint
	Rank23       string
	Rank23ID     uint
	Rank24       string
	Rank24ID     uint
	Rank25       string
	Rank25ID     uint
}

type CollegePollOfficial struct {
	gorm.Model
	SeasonID           uint
	WeekID             uint
	Week               uint
	RankOne            string
	RankOneID          uint
	RankOneVotes       uint
	RankOneNo1Votes    uint
	RankTwo            string
	RankTwoID          uint
	RankTwoVotes       uint
	RankTwoNo1Votes    uint
	RankThree          string
	RankThreeID        uint
	RankThreeVotes     uint
	RankThreeNo1Votes  uint
	RankFour           string
	RankFourID         uint
	RankFourVotes      uint
	RankFourNo1Votes   uint
	RankFive           string
	RankFiveID         uint
	RankFiveVotes      uint
	RankFiveNo1Votes   uint
	RankSix            string
	RankSixID          uint
	RankSixVotes       uint
	RankSixNo1Votes    uint
	RankSeven          string
	RankSevenID        uint
	RankSevenVotes     uint
	RankSevenNo1Votes  uint
	RankEight          string
	RankEightID        uint
	RankEightVotes     uint
	RankEightNo1Votes  uint
	RankNine           string
	RankNineID         uint
	RankNineVotes      uint
	RankNineNo1Votes   uint
	RankTen            string
	RankTenID          uint
	RankTenVotes       uint
	RankTenNo1Votes    uint
	RankEleven         string
	RankElevenID       uint
	RankElevenVotes    uint
	RankElevenNo1Votes uint
	Rank12             string
	Rank12ID           uint
	Rank12Votes        uint
	Rank12No1Votes     uint
	Rank13             string
	Rank13ID           uint
	Rank13Votes        uint
	Rank13No1Votes     uint
	Rank14             string
	Rank14ID           uint
	Rank14Votes        uint
	Rank14No1Votes     uint
	Rank15             string
	Rank15ID           uint
	Rank15Votes        uint
	Rank15No1Votes     uint
	Rank16             string
	Rank16ID           uint
	Rank16Votes        uint
	Rank16No1Votes     uint
	Rank17             string
	Rank17ID           uint
	Rank17Votes        uint
	Rank17No1Votes     uint
	Rank18             string
	Rank18ID           uint
	Rank18Votes        uint
	Rank18No1Votes     uint
	Rank19             string
	Rank19ID           uint
	Rank19Votes        uint
	Rank19No1Votes     uint
	Rank20             string
	Rank20ID           uint
	Rank20Votes        uint
	Rank20No1Votes     uint
	Rank21             string
	Rank21ID           uint
	Rank21Votes        uint
	Rank21No1Votes     uint
	Rank22             string
	Rank22ID           uint
	Rank22Votes        uint
	Rank22No1Votes     uint
	Rank23             string
	Rank23ID           uint
	Rank23Votes        uint
	Rank23No1Votes     uint
	Rank24             string
	Rank24ID           uint
	Rank24Votes        uint
	Rank24No1Votes     uint
	Rank25             string
	Rank25ID           uint
	Rank25Votes        uint
	Rank25No1Votes     uint
}

func (s *CollegePollSubmission) AssignID(id uint) {
	s.ID = id
}

func (c *CollegePollOfficial) AssignRank(idx int, vote TeamVote) {
	if idx == 0 {
		c.RankOne = vote.Team
		c.RankOneVotes = vote.TotalVotes
		c.RankOneID = vote.TeamID
		c.RankOneNo1Votes = vote.Number1Votes
	} else if idx == 1 {
		c.RankTwo = vote.Team
		c.RankTwoVotes = vote.TotalVotes
		c.RankTwoID = vote.TeamID
		c.RankTwoNo1Votes = vote.Number1Votes
	} else if idx == 2 {
		c.RankThree = vote.Team
		c.RankThreeVotes = vote.TotalVotes
		c.RankThreeID = vote.TeamID
		c.RankThreeNo1Votes = vote.Number1Votes
	} else if idx == 3 {
		c.RankFour = vote.Team
		c.RankFourVotes = vote.TotalVotes
		c.RankFourID = vote.TeamID
		c.RankFourNo1Votes = vote.Number1Votes
	} else if idx == 4 {
		c.RankFive = vote.Team
		c.RankFiveVotes = vote.TotalVotes
		c.RankFiveID = vote.TeamID
		c.RankFiveNo1Votes = vote.Number1Votes
	} else if idx == 5 {
		c.RankSix = vote.Team
		c.RankSixVotes = vote.TotalVotes
		c.RankSixID = vote.TeamID
		c.RankSixNo1Votes = vote.Number1Votes
	} else if idx == 6 {
		c.RankSeven = vote.Team
		c.RankSevenVotes = vote.TotalVotes
		c.RankSevenID = vote.TeamID
		c.RankSevenNo1Votes = vote.Number1Votes
	} else if idx == 7 {
		c.RankEight = vote.Team
		c.RankEightVotes = vote.TotalVotes
		c.RankEightID = vote.TeamID
		c.RankEightNo1Votes = vote.Number1Votes
	} else if idx == 8 {
		c.RankNine = vote.Team
		c.RankNineVotes = vote.TotalVotes
		c.RankNineID = vote.TeamID
		c.RankNineNo1Votes = vote.Number1Votes
	} else if idx == 9 {
		c.RankTen = vote.Team
		c.RankTenVotes = vote.TotalVotes
		c.RankTenID = vote.TeamID
		c.RankTenNo1Votes = vote.Number1Votes
	} else if idx == 10 {
		c.RankEleven = vote.Team
		c.RankElevenVotes = vote.TotalVotes
		c.RankElevenID = vote.TeamID
		c.RankElevenNo1Votes = vote.Number1Votes
	} else if idx == 11 {
		c.Rank12 = vote.Team
		c.Rank12Votes = vote.TotalVotes
		c.Rank12ID = vote.TeamID
		c.Rank12No1Votes = vote.Number1Votes
	} else if idx == 12 {
		c.Rank13 = vote.Team
		c.Rank13Votes = vote.TotalVotes
		c.Rank13ID = vote.TeamID
		c.Rank13No1Votes = vote.Number1Votes
	} else if idx == 13 {
		c.Rank14 = vote.Team
		c.Rank14Votes = vote.TotalVotes
		c.Rank14ID = vote.TeamID
		c.Rank14No1Votes = vote.Number1Votes
	} else if idx == 14 {
		c.Rank15 = vote.Team
		c.Rank15Votes = vote.TotalVotes
		c.Rank15ID = vote.TeamID
		c.Rank15No1Votes = vote.Number1Votes
	} else if idx == 15 {
		c.Rank16 = vote.Team
		c.Rank16Votes = vote.TotalVotes
		c.Rank16ID = vote.TeamID
		c.Rank16No1Votes = vote.Number1Votes
	} else if idx == 16 {
		c.Rank17 = vote.Team
		c.Rank17Votes = vote.TotalVotes
		c.Rank17ID = vote.TeamID
		c.Rank17No1Votes = vote.Number1Votes
	} else if idx == 17 {
		c.Rank18 = vote.Team
		c.Rank18Votes = vote.TotalVotes
		c.Rank18ID = vote.TeamID
		c.Rank18No1Votes = vote.Number1Votes
	} else if idx == 18 {
		c.Rank19 = vote.Team
		c.Rank19Votes = vote.TotalVotes
		c.Rank19ID = vote.TeamID
		c.Rank19No1Votes = vote.Number1Votes
	} else if idx == 19 {
		c.Rank20 = vote.Team
		c.Rank20Votes = vote.TotalVotes
		c.Rank20ID = vote.TeamID
		c.Rank20No1Votes = vote.Number1Votes
	} else if idx == 20 {
		c.Rank21 = vote.Team
		c.Rank21Votes = vote.TotalVotes
		c.Rank21ID = vote.TeamID
		c.Rank21No1Votes = vote.Number1Votes
	} else if idx == 21 {
		c.Rank22 = vote.Team
		c.Rank22Votes = vote.TotalVotes
		c.Rank22ID = vote.TeamID
		c.Rank22No1Votes = vote.Number1Votes
	} else if idx == 22 {
		c.Rank23 = vote.Team
		c.Rank23Votes = vote.TotalVotes
		c.Rank23ID = vote.TeamID
		c.Rank23No1Votes = vote.Number1Votes
	} else if idx == 23 {
		c.Rank24 = vote.Team
		c.Rank24Votes = vote.TotalVotes
		c.Rank24ID = vote.TeamID
		c.Rank24No1Votes = vote.Number1Votes
	} else if idx == 24 {
		c.Rank25 = vote.Team
		c.Rank25Votes = vote.TotalVotes
		c.Rank25ID = vote.TeamID
		c.Rank25No1Votes = vote.Number1Votes
	}
}

type TeamVote struct {
	Team         string
	TeamID       uint
	TotalVotes   uint
	Number1Votes uint
}

func (t *TeamVote) AddVotes(num uint) {
	t.TotalVotes += (26 - num)
	if num == 1 {
		t.Number1Votes++
	}
}
