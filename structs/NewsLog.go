package structs

import "gorm.io/gorm"

type NewsLog struct {
	gorm.Model
	WeekID      uint
	Week        uint
	SeasonID    uint
	Season      uint
	MessageType string
	Message     string
	League      string
}
