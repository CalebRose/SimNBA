package structs

import (
	"strconv"

	"gorm.io/gorm"
)

type CollegePoll struct {
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
	SeasonID        uint
	WeekID          uint
	Week            uint
	RankOne         string
	RankOneID       uint
	RankOneVotes    uint
	RankTwo         string
	RankTwoID       uint
	RankTwoVotes    uint
	RankThree       string
	RankThreeID     uint
	RankThreeVotes  uint
	RankFour        string
	RankFourID      uint
	RankFourVotes   uint
	RankFive        string
	RankFiveID      uint
	RankFiveVotes   uint
	RankSix         string
	RankSixID       uint
	RankSixVotes    uint
	RankSeven       string
	RankSevenID     uint
	RankSevenVotes  uint
	RankEight       string
	RankEightID     uint
	RankEightVotes  uint
	RankNine        string
	RankNineID      uint
	RankNineVotes   uint
	RankTen         string
	RankTenID       uint
	RankTenVotes    uint
	RankEleven      string
	RankElevenID    uint
	RankElevenVotes uint
	Rank12          string
	Rank12ID        uint
	Rank12Votes     uint
	Rank13          string
	Rank13ID        uint
	Rank13Votes     uint
	Rank14          string
	Rank14ID        uint
	Rank14Votes     uint
	Rank15          string
	Rank15ID        uint
	Rank15Votes     uint
	Rank16          string
	Rank16ID        uint
	Rank16Votes     uint
	Rank17          string
	Rank17ID        uint
	Rank17Votes     uint
	Rank18          string
	Rank18ID        uint
	Rank18Votes     uint
	Rank19          string
	Rank19ID        uint
	Rank19Votes     uint
	Rank20          string
	Rank20ID        uint
	Rank20Votes     uint
	Rank21          string
	Rank21ID        uint
	Rank21Votes     uint
	Rank22          string
	Rank22ID        uint
	Rank22Votes     uint
	Rank23          string
	Rank23ID        uint
	Rank23Votes     uint
	Rank24          string
	Rank24ID        uint
	Rank24Votes     uint
	Rank25          string
	Rank25ID        uint
	Rank25Votes     uint
}

func (c *CollegePollOfficial) AssignRank(idx int, vote TeamVote) {
	if idx == 0 {
		c.RankOne = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankOneVotes = vote.TotalVotes
		c.RankOneID = vote.TeamID

	} else if idx == 1 {
		c.RankTwo = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankTwoVotes = vote.TotalVotes
		c.RankTwoID = vote.TeamID
	} else if idx == 2 {
		c.RankThree = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankThreeVotes = vote.TotalVotes
		c.RankThreeID = vote.TeamID
	} else if idx == 3 {
		c.RankFour = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankFourVotes = vote.TotalVotes
		c.RankFourID = vote.TeamID
	} else if idx == 4 {
		c.RankFive = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankFiveVotes = vote.TotalVotes
		c.RankFiveID = vote.TeamID
	} else if idx == 5 {
		c.RankSix = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankSixVotes = vote.TotalVotes
		c.RankSixID = vote.TeamID
	} else if idx == 6 {
		c.RankSeven = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankSevenVotes = vote.TotalVotes
		c.RankSevenID = vote.TeamID
	} else if idx == 7 {
		c.RankEight = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankEightVotes = vote.TotalVotes
		c.RankEightID = vote.TeamID
	} else if idx == 8 {
		c.RankNine = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankNineVotes = vote.TotalVotes
		c.RankNineID = vote.TeamID
	} else if idx == 9 {
		c.RankTen = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankTenVotes = vote.TotalVotes
		c.RankTenID = vote.TeamID
	} else if idx == 10 {
		c.RankEleven = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.RankElevenVotes = vote.TotalVotes
		c.RankElevenID = vote.TeamID
	} else if idx == 11 {
		c.Rank12 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank12Votes = vote.TotalVotes
		c.Rank12ID = vote.TeamID
	} else if idx == 12 {
		c.Rank13 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank13Votes = vote.TotalVotes
		c.Rank13ID = vote.TeamID
	} else if idx == 13 {
		c.Rank14 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank14Votes = vote.TotalVotes
		c.Rank14ID = vote.TeamID
	} else if idx == 14 {
		c.Rank15 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank15Votes = vote.TotalVotes
		c.Rank15ID = vote.TeamID
	} else if idx == 15 {
		c.Rank16 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank16Votes = vote.TotalVotes
		c.Rank16ID = vote.TeamID
	} else if idx == 16 {
		c.Rank17 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank17Votes = vote.TotalVotes
		c.Rank17ID = vote.TeamID
	} else if idx == 17 {
		c.Rank18 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank18Votes = vote.TotalVotes
		c.Rank18ID = vote.TeamID
	} else if idx == 18 {
		c.Rank19 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank19Votes = vote.TotalVotes
		c.Rank19ID = vote.TeamID
	} else if idx == 19 {
		c.Rank20 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank20Votes = vote.TotalVotes
		c.Rank20ID = vote.TeamID
	} else if idx == 20 {
		c.Rank21 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank21Votes = vote.TotalVotes
		c.Rank21ID = vote.TeamID
	} else if idx == 21 {
		c.Rank22 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank22Votes = vote.TotalVotes
		c.Rank22ID = vote.TeamID
	} else if idx == 22 {
		c.Rank23 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank23Votes = vote.TotalVotes
		c.Rank23ID = vote.TeamID
	} else if idx == 23 {
		c.Rank24 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank24Votes = vote.TotalVotes
		c.Rank24ID = vote.TeamID
	} else if idx == 24 {
		c.Rank25 = vote.Team + "(" + strconv.Itoa(int(vote.Number1Votes)) + ")"
		c.Rank25Votes = vote.TotalVotes
		c.Rank25ID = vote.TeamID
	}
}

type TeamVote struct {
	Team         string
	TeamID       uint
	TotalVotes   uint
	Number1Votes uint
}
