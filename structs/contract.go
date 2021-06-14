package structs

import "github.com/jinzhu/gorm"

// Contract - The contract of which the player is obligated to
type Contract struct {
	gorm.Model
	PlayerID       int
	TeamID         int
	CurrentYear    int
	ContractLength int
	ContractValue  int
	ContractType   string
	IsActive       bool
	// Do we want to keep track of the year?
}
