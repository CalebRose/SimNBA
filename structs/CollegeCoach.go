package structs

import (
	"gorm.io/gorm"
)

type CollegeCoach struct {
	gorm.Model
	Name                   string // Name of Coach
	Age                    int    // Age of Coach. Anywhere between 34 and 76
	TeamID                 uint
	Team                   string
	AlmaMaterID            uint
	AlmaMater              string // The School They Attended
	FormerPlayerID         uint
	Prestige               int    // Level system. Every 10 wins, every playoff win & conference tourney championship nets 1 point which can then be applied towards one of the five odds categories, if they qualify
	PointMin               int    // Minimum number of points the coach will put towards a player
	PointMax               int    // Maximum number of points the coach will place on a player
	StarMin                int    // Minimum star rating they will target on a croot (floor)
	StarMax                int    // Maximum star rating they will target on a croot (ceiling)
	Odds1                  int    // Modifier towards adding 1 star croots to board
	Odds2                  int    // Modifier towards adding 2 star croots to board
	Odds3                  int    // Modifier towards adding 3 star croots to board
	Odds4                  int    // Modifier towards adding 4 star croots to board
	Odds5                  int    // Modifier towards adding 5 star croots to board
	Scheme                 string // Desired scheme the coach wants to run -- will recruit based on the desired scheme
	DefensiveScheme        string // Desired defensive scheme the coach wants to run -- will recruit based on the desired scheme
	TeambuildingPreference string // Coaches that prefer to recruit vs Coaches that will utilize the transfer portal
	CareerPreference       string // "Prefers to stay at their current job", "Wants to coach at Alma Mater", "Wants a more competitive job", "Average"
	PromiseTendency        string // Coach will either under-promise, over-promise, or be average on promises within transfer portal
	PortalReputation       int    // A value between 1-100 signifying the coach's reputation and behavior in the transfer portal.
	SchoolTenure           int    // Number of years coach is participating on the team
	CareerTenure           int    // Number of years the coach is actively coaching in the college sim
	ContractLength         int    // Number of years of the current contract
	YearsRemaining         int    // Years left on current contract
	IsRetired              bool
	IsFormerPlayer         bool
}

func (c *CollegeCoach) IncrementOdds(star int) {
	switch star {
	case 1:
		c.Odds1 += 1
	case 2:
		c.Odds2 += 1
	case 3:
		c.Odds3 += 1
	case 4:
		c.Odds4 += 1
	case 5:
		c.Odds5 += 1
	}
}
