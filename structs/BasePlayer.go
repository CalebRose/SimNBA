package structs

import "math/rand"

var archetypeWeights = map[string]map[string]map[string]float64{
	"C": {
		"Rim Protector": {
			"InsideShooting":     1.1,
			"MidRangeShooting":   0.85,
			"ThreePointShooting": 0.7,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.0,
			"Rebounding":         1.2,
			"Stealing":           0.7,
			"Blocking":           1.3,
			"InteriorDefense":    1.3,
			"PerimeterDefense":   0.8,
		},
		"Post Scorer": {
			"InsideShooting":     1.2,
			"MidRangeShooting":   1.1,
			"ThreePointShooting": 0.75,
			"FreeThrow":          1.1,
			"Agility":            1.0,
			"Ballwork":           1.0,
			"Rebounding":         1.15,
			"Stealing":           0.75,
			"Blocking":           1.15,
			"InteriorDefense":    1.15,
			"PerimeterDefense":   0.7,
		},
		"Stretch Center": {
			"InsideShooting":     0.9,
			"MidRangeShooting":   1.1,
			"ThreePointShooting": 1.2,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.0,
			"Rebounding":         0.9,
			"Stealing":           0.7,
			"Blocking":           0.9,
			"InteriorDefense":    1.2,
			"PerimeterDefense":   1.0,
		},
		"All-Around": {
			"InsideShooting":     1.0,
			"MidRangeShooting":   1.0,
			"ThreePointShooting": 1.0,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.0,
			"Rebounding":         1.0,
			"Stealing":           1.0,
			"Blocking":           1.0,
			"InteriorDefense":    1.0,
			"PerimeterDefense":   1.0,
		},
	},
	"F": {
		"Power Forward": {
			"InsideShooting":     1.2,
			"MidRangeShooting":   1.2,
			"ThreePointShooting": 1.0,
			"FreeThrow":          0.75,
			"Agility":            1.0,
			"Ballwork":           0.75,
			"Rebounding":         1.1,
			"Stealing":           1.0,
			"Blocking":           1.2,
			"InteriorDefense":    1.1,
			"PerimeterDefense":   0.75,
		},
		"Small Forward": {
			"InsideShooting":     1.0,
			"MidRangeShooting":   1.0,
			"ThreePointShooting": 1.0,
			"FreeThrow":          1.0,
			"Agility":            1.1,
			"Ballwork":           1.0,
			"Rebounding":         1.0,
			"Stealing":           1.15,
			"Blocking":           0.8,
			"InteriorDefense":    1.0,
			"PerimeterDefense":   1.0,
		},
		"Point Forward": {
			"InsideShooting":     1.1,
			"MidRangeShooting":   1.0,
			"ThreePointShooting": 0.85,
			"FreeThrow":          1.1,
			"Agility":            1.0,
			"Ballwork":           1.3,
			"Rebounding":         1.0,
			"Stealing":           1.0,
			"Blocking":           0.8,
			"InteriorDefense":    0.8,
			"PerimeterDefense":   1.0,
		},
		"Swingman": {
			"InsideShooting":     0.85,
			"MidRangeShooting":   1.0,
			"ThreePointShooting": 1.0,
			"FreeThrow":          1.0,
			"Agility":            1.2,
			"Ballwork":           1.0,
			"Rebounding":         0.85,
			"Stealing":           1.15,
			"Blocking":           0.85,
			"InteriorDefense":    0.9,
			"PerimeterDefense":   1.0,
		},
		"Two-Way": {
			"InsideShooting":     0.9,
			"MidRangeShooting":   0.9,
			"ThreePointShooting": 0.9,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           0.9,
			"Rebounding":         1.0,
			"Stealing":           1.15,
			"Blocking":           1.1,
			"InteriorDefense":    1.05,
			"PerimeterDefense":   1.05,
		},
		"All-Around": {
			"InsideShooting":     1.0,
			"MidRangeShooting":   1.0,
			"ThreePointShooting": 1.0,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.0,
			"Rebounding":         1.0,
			"Stealing":           1.0,
			"Blocking":           1.0,
			"InteriorDefense":    1.0,
			"PerimeterDefense":   1.0,
		},
	},
	"G": {
		"Point Guard": {
			"InsideShooting":     0.7,
			"MidRangeShooting":   1.1,
			"ThreePointShooting": 1.0,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.3,
			"Rebounding":         0.7,
			"Stealing":           1.0,
			"Blocking":           0.7,
			"InteriorDefense":    1.0,
			"PerimeterDefense":   1.15,
		},
		"Shooting Guard": {
			"InsideShooting":     1.2,
			"MidRangeShooting":   1.2,
			"ThreePointShooting": 1.2,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.0,
			"Rebounding":         0.75,
			"Stealing":           1.0,
			"Blocking":           0.7,
			"InteriorDefense":    0.8,
			"PerimeterDefense":   0.85,
		},
		"Combo Guard": {
			"InsideShooting":     1.0,
			"MidRangeShooting":   1.1,
			"ThreePointShooting": 1.1,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.15,
			"Rebounding":         0.8,
			"Stealing":           1.0,
			"Blocking":           0.75,
			"InteriorDefense":    0.9,
			"PerimeterDefense":   1.0,
		},
		"Slasher": {
			"InsideShooting":     1.2,
			"MidRangeShooting":   0.75,
			"ThreePointShooting": 0.7,
			"FreeThrow":          1.0,
			"Agility":            1.2,
			"Ballwork":           0.8,
			"Rebounding":         0.75,
			"Stealing":           1.2,
			"Blocking":           1.1,
			"InteriorDefense":    1.0,
			"PerimeterDefense":   1.0,
		},
		"3-and-D": {
			"InsideShooting":     0.7,
			"MidRangeShooting":   0.9,
			"ThreePointShooting": 1.2,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.0,
			"Rebounding":         0.8,
			"Stealing":           1.0,
			"Blocking":           0.7,
			"InteriorDefense":    1.15,
			"PerimeterDefense":   1.2,
		},
		"All-Around": {
			"InsideShooting":     1.0,
			"MidRangeShooting":   1.0,
			"ThreePointShooting": 1.0,
			"FreeThrow":          1.0,
			"Agility":            1.0,
			"Ballwork":           1.0,
			"Rebounding":         1.0,
			"Stealing":           1.0,
			"Blocking":           1.0,
			"InteriorDefense":    1.0,
			"PerimeterDefense":   1.0,
		},
	},

	// Add more archetypes as needed
}

type BasePlayer struct {
	TeamID                 uint
	Team                   string
	PlayerID               uint
	FirstName              string
	LastName               string
	Position               string
	Archetype              string
	Age                    uint8
	PrimeAge               uint8
	Year                   uint8
	City                   string
	HighSchool             string
	State                  string
	Country                string
	Stars                  uint8
	Height                 uint8
	Weight                 uint16
	InsideShooting         uint8
	SpecInsideShooting     bool
	MidRangeShooting       uint8
	SpecMidRangeShooting   bool
	ThreePointShooting     uint8
	SpecThreePointShooting bool
	FreeThrow              uint8
	SpecFreeThrow          bool
	Agility                uint8
	SpecAgility            bool
	Ballwork               uint8
	SpecBallwork           bool
	Rebounding             uint8
	SpecRebounding         bool
	Stealing               uint8
	SpecStealing           bool
	Blocking               uint8
	SpecBlocking           bool
	InteriorDefense        uint8
	SpecInteriorDefense    bool
	PerimeterDefense       uint8
	SpecPerimeterDefense   bool
	Potential              uint8
	PotentialGrade         string
	ProPotentialGrade      uint8
	Stamina                uint8
	Discipline             uint8
	InjuryRating           uint8
	IsInjured              bool
	InjuryName             string
	InjuryType             string
	WeeksOfRecovery        uint8
	InjuryReserve          bool
	PlaytimeExpectations   uint8
	Overall                uint8
	SpecCount              uint8
	Personality            string
	FreeAgency             string
	RecruitingBias         string
	RecruitingBiasValue    string
	WorkEthic              string
	AcademicBias           string
	PreviousTeamID         uint
	PreviousTeam           string
	RelativeID             uint8
	RelativeType           uint8
	Notes                  string
	PlayerPreferences
}

func (b *BasePlayer) GetOverall() {
	weights := archetypeWeights[b.Position][b.Archetype]
	totalWeight := 0.0
	weightedSum := 0.0

	for attr, weight := range weights {
		var value uint8
		switch attr {
		case "InsideShooting":
			value = b.InsideShooting
		case "MidRangeShooting":
			value = b.MidRangeShooting
		case "ThreePointShooting":
			value = b.ThreePointShooting
		case "FreeThrow":
			value = b.FreeThrow
		case "Agility":
			value = b.Agility
		case "Ballwork":
			value = b.Ballwork
		case "Rebounding":
			value = b.Rebounding
		case "Stealing":
			value = b.Stealing
		case "Blocking":
			value = b.Blocking
		case "InteriorDefense":
			value = b.InteriorDefense
		case "PerimeterDefense":
			value = b.PerimeterDefense
		}
		weightedSum += float64(value) * weight
		if value > 0 {
			totalWeight += weight
		}
	}

	// Normalize to 1–50 range
	if totalWeight > 0 {
		b.Overall = uint8((weightedSum / totalWeight)) // * 50.0
	} else {
		b.Overall = 0
	}
}

func (cp *BasePlayer) GetPotentialGrade() {
	adjust := rand.Intn(20) - 10
	if adjust == 0 {
		test := rand.Intn(2000) - 1000

		if test > 0 {
			adjust += 1
		} else if test < 0 {
			adjust -= 1
		} else {
			adjust = 0
		}
	}
	potential := cp.Potential + uint8(adjust)
	if potential > 80 {
		cp.PotentialGrade = "A+"
	} else if potential > 70 {
		cp.PotentialGrade = "A"
	} else if potential > 65 {
		cp.PotentialGrade = "A-"
	} else if potential > 60 {
		cp.PotentialGrade = "B+"
	} else if potential > 55 {
		cp.PotentialGrade = "B"
	} else if potential > 50 {
		cp.PotentialGrade = "B-"
	} else if potential > 40 {
		cp.PotentialGrade = "C+"
	} else if potential > 30 {
		cp.PotentialGrade = "C"
	} else if potential > 25 {
		cp.PotentialGrade = "C-"
	} else if potential > 20 {
		cp.PotentialGrade = "D+"
	} else if potential > 15 {
		cp.PotentialGrade = "D"
	} else if potential > 10 {
		cp.PotentialGrade = "D-"
	} else {
		cp.PotentialGrade = "F"
	}
}

func (cp *BasePlayer) GetNBAPotentialGrade() {
	adjust := rand.Intn(20) - 10
	if adjust == 0 {
		test := rand.Intn(2000) - 1000

		if test > 0 {
			adjust += 1
		} else if test < 0 {
			adjust -= 1
		} else {
			adjust = 0
		}
	}
	potential := cp.ProPotentialGrade + uint8(adjust)
	if potential > 80 {
		cp.PotentialGrade = "A+"
	} else if potential > 70 {
		cp.PotentialGrade = "A"
	} else if potential > 65 {
		cp.PotentialGrade = "A-"
	} else if potential > 60 {
		cp.PotentialGrade = "B+"
	} else if potential > 55 {
		cp.PotentialGrade = "B"
	} else if potential > 50 {
		cp.PotentialGrade = "B-"
	} else if potential > 40 {
		cp.PotentialGrade = "C+"
	} else if potential > 30 {
		cp.PotentialGrade = "C"
	} else if potential > 25 {
		cp.PotentialGrade = "C-"
	} else if potential > 20 {
		cp.PotentialGrade = "D+"
	} else if potential > 15 {
		cp.PotentialGrade = "D"
	} else if potential > 10 {
		cp.PotentialGrade = "D-"
	} else {
		cp.PotentialGrade = "F"
	}
}

func (cp *BasePlayer) GetSpecCount() {
	count := 0
	if cp.SpecInsideShooting {
		count++
	}
	if cp.SpecMidRangeShooting {
		count++
	}
	if cp.SpecThreePointShooting {
		count++
	}
	if cp.SpecFreeThrow {
		count++
	}
	if cp.SpecAgility {
		count++
	}
	if cp.SpecBallwork {
		count++
	}
	if cp.SpecRebounding {
		count++
	}
	if cp.SpecStealing {
		count++
	}
	if cp.SpecBlocking {
		count++
	}
	if cp.SpecInteriorDefense {
		count++
	}
	if cp.SpecPerimeterDefense {
		count++
	}
	cp.SpecCount = uint8(count)
}

func (p *BasePlayer) SetDisciplineAndIR(val, val2 int) {
	p.Discipline = uint8(val)
	p.InjuryRating = uint8(val2)
}

func (p *BasePlayer) SetExpectations(val uint8) {
	p.PlaytimeExpectations = val
}

func (c *BasePlayer) SetWorkEthic(ethic string) {
	c.WorkEthic = ethic
}

func (c *BasePlayer) SetPersonality(personality string) {
	c.Personality = personality
}

func (c *BasePlayer) SetFreeAgencyBias(faBias string) {
	c.FreeAgency = faBias
}

func (c *BasePlayer) SetRecruitingBias(recBias string) {
	c.RecruitingBias = recBias
}

func (c *BasePlayer) SetAge(age int) {
	c.Age = uint8(age)
}

func (c *BasePlayer) SetMinutesExpectations(min uint8) {
	c.PlaytimeExpectations = min
}

func (np *BasePlayer) SetAttributes(s2, s3, fn, ft, bl, rb, id, pd, ovr, stars, exp int) {
	np.MidRangeShooting = uint8(s2)
	np.ThreePointShooting = uint8(s3)
	np.InsideShooting = uint8(fn)
	np.FreeThrow = uint8(ft)
	np.Ballwork = uint8(bl)
	np.Rebounding = uint8(rb)
	np.InteriorDefense = uint8(id)
	np.PerimeterDefense = uint8(pd)
	np.Overall = uint8(ovr)
	np.Stars = uint8(stars)
	np.PlaytimeExpectations = uint8(exp)
	np.IsInjured = false
	np.WeeksOfRecovery = 0
	np.InjuryName = ""
	np.InjuryType = ""
}

func (bp *BasePlayer) SetInjury(ijName, ijType string, wor uint8) {
	bp.IsInjured = true
	bp.InjuryName = ijName
	bp.InjuryType = ijType
	bp.WeeksOfRecovery = wor
}
func (bp *BasePlayer) RunRecovery() {
	bp.WeeksOfRecovery -= 1
	if bp.WeeksOfRecovery == 0 || bp.WeeksOfRecovery > 100 {
		bp.IsInjured = false
		bp.InjuryName = ""
		bp.InjuryType = ""
		bp.WeeksOfRecovery = 0
	}
}
