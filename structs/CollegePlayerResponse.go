package structs

type CollegePlayerResponse struct {
	FirstName             string
	LastName              string
	Archetype             string
	Position              string
	Age                   uint8
	Year                  uint8
	State                 string
	Country               string
	Stars                 uint8
	Height                uint8
	Weight                uint8
	PotentialGrade        string
	Shooting2Grade        string
	Shooting3Grade        string
	FreeThrowGrade        string
	FinishingGrade        string
	BallworkGrade         string
	AgilityGrade          string
	StealingGrade         string
	BlockingGrade         string
	ReboundingGrade       string
	InteriorDefenseGrade  string
	PerimeterDefenseGrade string
	Stamina               uint8
	PlaytimeExpectations  uint8
	InsideProportion      float64
	MidRangeProportion    float64
	ThreePointProportion  float64
	PositionOne           string
	PositionTwo           string
	PositionThree         string
	P1Minutes             uint8
	P2Minutes             uint8
	P3Minutes             uint8
	Minutes               uint8
	Potential             uint8
	OverallGrade          string
	Personality           string
	RecruitingBias        string
	RecruitingBiasValue   string
	WorkEthic             string
	AcademicBias          string
	PlayerID              uint
	TeamID                uint
	Team                  string
	IsRedshirting         bool
	IsRedshirt            bool
	PreviousTeamID        uint
	PreviousTeam          string
	TransferStatus        uint8  // 1 == Intends, 2 == Is Transferring
	TransferLikeliness    string // Low, Medium, High
	LegacyID              uint   // Either a legacy school or a legacy coach
	DisciplineGrade       string
	InjuryGrade           string
	IsInjured             bool
	InjuryName            string
	InjuryType            string
	WeeksOfRecovery       uint8
	InjuryReserve         bool
	SeasonStats           CollegePlayerSeasonStats
	Stats                 CollegePlayerStats
}

type NBAPlayerResponse struct {
	FirstName             string
	LastName              string
	Position              string
	Age                   uint8
	Year                  uint8
	State                 string
	Country               string
	Stars                 uint8
	Height                uint8
	Weight                uint8
	PotentialGrade        string
	Shooting2Grade        string
	Shooting3Grade        string
	FreeThrowGrade        string
	FinishingGrade        string
	BallworkGrade         string
	AgilityGrade          string
	StealingGrade         string
	BlockingGrade         string
	ReboundingGrade       string
	InteriorDefenseGrade  string
	PerimeterDefenseGrade string
	Stamina               uint8
	PlaytimeExpectations  uint8
	Minutes               uint8
	Potential             uint8
	OverallGrade          string
	Personality           string
	RecruitingBias        string
	WorkEthic             string
	AcademicBias          string
	PlayerID              uint
	TeamID                uint
	Team                  string
	IsRedshirting         bool
	IsRedshirt            bool
	SeasonStats           NBAPlayerSeasonStats
	Stats                 NBAPlayerStats
}

type TransferPlayerResponse struct {
	ID                   uint
	FirstName            string
	LastName             string
	Archetype            string
	Position             string
	Age                  uint8
	Year                 uint8
	State                string
	Country              string
	Stars                uint8
	Height               uint8
	Weight               uint8
	PotentialGrade       string
	Overall              string
	Shooting2            string
	Shooting3            string
	FreeThrow            string
	Finishing            string
	Ballwork             string
	Rebounding           string
	InteriorDefense      string
	PerimeterDefense     string
	Stamina              uint8
	PlaytimeExpectations uint8
	Minutes              uint8
	OverallGrade         string
	Personality          string
	RecruitingBias       string
	RecruitingBiasValue  string
	WorkEthic            string
	AcademicBias         string
	PlayerID             uint
	TeamID               uint
	Team                 string
	IsRedshirting        bool
	IsRedshirt           bool
	PreviousTeamID       uint
	PreviousTeam         string
	TransferStatus       uint8  // 1 == Intends, 2 == Is Transferring
	TransferLikeliness   string // Low, Medium, High
	LegacyID             uint   // Either a legacy school or a legacy coach
	SeasonStats          CollegePlayerSeasonStats
	Stats                CollegePlayerStats
	LeadingTeams         []LeadingTeams
}

func (c *TransferPlayerResponse) Map(r CollegePlayer, ovr string) {
	c.ID = r.ID
	c.PlayerID = r.PlayerID
	c.TeamID = r.TeamID
	c.FirstName = r.FirstName
	c.LastName = r.LastName
	c.Position = r.Position
	c.Height = r.Height
	c.Stars = r.Stars
	c.Shooting2 = attributeMapper(r.MidRangeShooting)
	c.Shooting3 = attributeMapper(r.ThreePointShooting)
	c.Finishing = attributeMapper(r.InsideShooting)
	c.FreeThrow = attributeMapper(r.FreeThrow)
	c.Ballwork = attributeMapper(r.Ballwork)
	c.Rebounding = attributeMapper(r.Rebounding)
	c.InteriorDefense = attributeMapper(r.InteriorDefense)
	c.PerimeterDefense = attributeMapper(r.PerimeterDefense)
	c.Stamina = r.Stamina
	c.OverallGrade = ovr
	c.PotentialGrade = r.PotentialGrade
	c.Personality = r.Personality
	c.RecruitingBias = r.RecruitingBias
	c.AcademicBias = r.AcademicBias
	c.WorkEthic = r.WorkEthic
	c.State = r.State
	c.Country = r.Country
	c.Team = r.Team
	c.PreviousTeam = r.PreviousTeam
	c.PreviousTeamID = r.PreviousTeamID
	c.Year = r.Year
	c.IsRedshirt = r.IsRedshirt
	c.IsRedshirting = r.IsRedshirting
}
