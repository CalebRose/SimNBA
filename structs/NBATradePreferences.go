package structs

import "gorm.io/gorm"

type NBATradePreferences struct {
	gorm.Model
	NBATeamID                uint
	PointGuards              bool
	PointGuardSpecialties    string
	ShootingGuards           bool
	ShootingGuardSpecialties string
	PowerForwards            bool
	PowerForwardSpecialties  string
	SmallForwards            bool
	SmallForwardSpecialties  string
	Centers                  bool
	CenterSpecialties        string
	DraftPicks               bool
	DraftPickType            string
}

type NBATradePreferencesDTO struct {
	NBATeamID                uint
	PointGuards              bool
	PointGuardSpecialties    string
	ShootingGuards           bool
	ShootingGuardSpecialties string
	PowerForwards            bool
	PowerForwardSpecialties  string
	SmallForwards            bool
	SmallForwardSpecialties  string
	Centers                  bool
	CenterSpecialties        string
	DraftPicks               bool
	DraftPickType            string
}

func (tp *NBATradePreferences) UpdatePreferences(pref NBATradePreferencesDTO) {
	tp.PointGuards = pref.PointGuards
	if tp.PointGuards {
		tp.PointGuardSpecialties = pref.PointGuardSpecialties
	}
	tp.ShootingGuards = pref.ShootingGuards
	if tp.ShootingGuards {
		tp.ShootingGuardSpecialties = pref.ShootingGuardSpecialties
	}
	tp.PowerForwards = pref.PowerForwards
	if tp.PowerForwards {
		tp.PowerForwardSpecialties = pref.PowerForwardSpecialties
	}
	tp.SmallForwards = pref.SmallForwards
	if tp.SmallForwards {
		tp.SmallForwardSpecialties = pref.SmallForwardSpecialties
	}
	tp.Centers = pref.Centers
	if tp.Centers {
		tp.CenterSpecialties = pref.CenterSpecialties
	}
	tp.DraftPicks = pref.DraftPicks
	if tp.DraftPicks {
		tp.DraftPickType = pref.DraftPickType
	}
}
