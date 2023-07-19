package structs

import "math/rand"

type BasePlayer struct {
	FirstName            string
	LastName             string
	Position             string
	Archetype            string
	Age                  int
	Year                 int
	State                string
	Country              string
	Stars                int
	Height               string
	Shooting2            int
	SpecShooting2        bool
	Shooting3            int
	SpecShooting3        bool
	Finishing            int
	SpecFinishing        bool
	FreeThrow            int
	SpecFreeThrow        bool
	Ballwork             int
	SpecBallwork         bool
	Rebounding           int
	SpecRebounding       bool
	Defense              int
	InteriorDefense      int
	SpecInteriorDefense  bool
	PerimeterDefense     int
	SpecPerimeterDefense bool
	Potential            int
	PotentialGrade       string
	ProPotentialGrade    int
	Stamina              int
	PlaytimeExpectations int
	Minutes              int
	InsideProportion     float64
	MidRangeProportion   float64
	ThreePointProportion float64
	Overall              int
	PositionOne          string
	PositionTwo          string
	PositionThree        string
	P1Minutes            int
	P2Minutes            int
	P3Minutes            int
	SpecCount            int
	Personality          string
	FreeAgency           string
	RecruitingBias       string
	WorkEthic            string
	AcademicBias         string
}

func (b *BasePlayer) ToggleSpecialties(str string) {
	if str == "SpecShooting2" {
		b.SpecShooting2 = true
	} else if str == "SpecShooting3" {
		b.SpecShooting3 = true
	} else if str == "SpecFreeThrow" {
		b.SpecFreeThrow = true
	} else if str == "SpecFinishing" {
		b.SpecFinishing = true
	} else if str == "SpecBallwork" {
		b.SpecBallwork = true
	} else if str == "SpecRebounding" {
		b.SpecRebounding = true
	} else if str == "SpecInteriorDefense" {
		b.SpecInteriorDefense = true
	} else if str == "SpecPerimeterDefense" {
		b.SpecPerimeterDefense = true
	}
	b.SpecCount++
}

func (b *BasePlayer) AssignArchetype() {
	pos := b.Position
	if b.SpecCount > 5 {
		b.Archetype = "All-Around"
		return
	}
	if pos == "G" {
		if b.SpecBallwork && !b.SpecShooting2 && !b.SpecShooting3 || (b.Ballwork > b.Shooting2 && b.Ballwork > b.Shooting3) {
			b.Archetype = "Floor General"
		} else if (b.SpecShooting2 && b.SpecShooting3) && (!b.SpecBallwork || !b.SpecRebounding) {
			b.Archetype = "Sharp Shooter"
		} else if b.SpecBallwork && (b.SpecShooting2 || b.SpecFreeThrow) && (!b.SpecShooting3 || !b.SpecFinishing) {
			b.Archetype = "Mid-Range Magician"
		} else if b.SpecRebounding && (b.SpecInteriorDefense || b.SpecPerimeterDefense) && (!b.SpecShooting2 || !b.SpecShooting3) {
			b.Archetype = "Defensive Dawg"
		} else if b.SpecShooting3 && (b.SpecInteriorDefense || b.SpecPerimeterDefense) {
			b.Archetype = "3-and-D"
		} else if b.SpecFinishing && b.SpecBallwork {
			b.Archetype = "Dunk Specialist"
		} else if b.SpecShooting2 && b.SpecShooting3 && b.SpecFinishing && (!b.SpecBallwork || !b.SpecInteriorDefense || !b.SpecPerimeterDefense) {
			b.Archetype = "Microwave"
		}
	} else if pos == "F" {
		if b.SpecShooting3 && (b.SpecInteriorDefense || b.SpecPerimeterDefense) {
			b.Archetype = "Two-Way Wing"
		} else if (!b.SpecShooting2 || !b.SpecShooting3) && b.SpecFinishing && (b.SpecInteriorDefense || b.SpecPerimeterDefense) {
			b.Archetype = "Slasher"
		} else if b.SpecShooting2 && b.SpecFinishing && (b.SpecInteriorDefense || b.SpecPerimeterDefense) {
			b.Archetype = "Traditional Forward"
		} else if b.SpecShooting2 && b.SpecShooting3 && b.SpecFinishing && (!b.SpecInteriorDefense || !b.SpecPerimeterDefense) {
			b.Archetype = "Offensive Weapon"
		} else if b.SpecRebounding && (b.SpecInteriorDefense || b.SpecPerimeterDefense) && (!b.SpecShooting2 || !b.SpecShooting3 || !b.SpecFinishing) ||
			((b.Rebounding > b.Shooting2 && b.Rebounding > b.Shooting3) || (b.InteriorDefense > b.Shooting2 && b.InteriorDefense > b.Shooting3) || (b.PerimeterDefense > b.Shooting2 && b.PerimeterDefense > b.Shooting3)) {
			b.Archetype = "Pure Defender"
		} else if b.SpecBallwork && (b.SpecInteriorDefense || b.SpecPerimeterDefense) && (!b.SpecShooting2 || !b.SpecShooting3 || !b.SpecFinishing) {
			b.Archetype = "Point Forward"
		}
	} else if pos == "C" {
		if b.SpecRebounding && !b.SpecFinishing {
			b.Archetype = "Rim Protector"
		} else if b.SpecShooting3 || (b.Shooting3 > b.Shooting2 && (b.InteriorDefense > b.Shooting2 || b.PerimeterDefense > b.Shooting2)) {
			b.Archetype = "Stretch Bigs"
		} else if b.SpecFinishing && b.SpecRebounding {
			b.Archetype = "Lob Threat"
		}
	}
}

func (p *BasePlayer) SetMidShotProportion(mid float64) {
	p.MidRangeProportion = mid
}

func (p *BasePlayer) SetInsideProportion(ins float64) {
	p.InsideProportion = ins
}

func (p *BasePlayer) SetThreePointProportion(tp float64) {
	p.ThreePointProportion = tp
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
	potential := cp.Potential + adjust
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
	potential := cp.ProPotentialGrade + adjust
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

func (p *BasePlayer) SetExpectations(val int) {
	p.PlaytimeExpectations = val
}

func (p *BasePlayer) SetMinutes() {
	p.Minutes = p.P1Minutes + p.P2Minutes + p.P3Minutes
}

func (c *BasePlayer) UpdatePlayer(p BasePlayer) {
	c.P1Minutes = p.P1Minutes
	c.P2Minutes = p.P2Minutes
	c.P3Minutes = p.P3Minutes
	c.PositionOne = p.PositionOne
	c.PositionTwo = p.PositionTwo
	c.PositionThree = p.PositionThree
	c.InsideProportion = p.InsideProportion
	c.MidRangeProportion = p.MidRangeProportion
	c.ThreePointProportion = p.ThreePointProportion
	c.SetMinutes()
}

func (c *BasePlayer) SetP1Minutes(p int, pos string) {
	c.PositionOne = pos
	c.P1Minutes = p
	c.SetMinutes()
}

func (c *BasePlayer) SetP2Minutes(p int, pos string) {
	c.PositionTwo = pos
	c.P2Minutes = p
	c.SetMinutes()
}

func (c *BasePlayer) SetP3Minutes(p int, pos string) {
	c.PositionThree = pos
	c.P3Minutes = p
	c.SetMinutes()
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
