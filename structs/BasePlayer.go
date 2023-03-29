package structs

type BasePlayer struct {
	FirstName            string
	LastName             string
	Position             string
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
	Overall              int
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
