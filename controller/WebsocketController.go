package controller

import (
	"github.com/CalebRose/SimNBA/managers"
	"github.com/CalebRose/SimNBA/structs"
)

func GetUpdatedTimestamp() structs.Timestamp {
	ts := managers.GetTimestamp()
	return ts
}
