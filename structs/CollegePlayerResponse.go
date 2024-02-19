package structs

type CollegePlayerResponse struct {
	FirstName             string
	LastName              string
	Archetype             string
	Position              string
	Age                   int
	Year                  int
	State                 string
	Country               string
	Stars                 int
	Height                string
	PotentialGrade        string
	Shooting2Grade        string
	Shooting3Grade        string
	FreeThrowGrade        string
	FinishingGrade        string
	BallworkGrade         string
	ReboundingGrade       string
	InteriorDefenseGrade  string
	PerimeterDefenseGrade string
	Stamina               int
	PlaytimeExpectations  int
	InsideProportion      float64
	MidRangeProportion    float64
	ThreePointProportion  float64
	PositionOne           string
	PositionTwo           string
	PositionThree         string
	P1Minutes             int
	P2Minutes             int
	P3Minutes             int
	Minutes               int
	Potential             int
	OverallGrade          string
	Personality           string
	RecruitingBias        string
	WorkEthic             string
	AcademicBias          string
	PlayerID              uint
	TeamID                uint
	TeamAbbr              string
	IsRedshirting         bool
	IsRedshirt            bool
	TransferStatus        int    // 1 == Intends, 2 == Is Transferring
	TransferLikeliness    string // Low, Medium, High
	LegacyID              uint   // Either a legacy school or a legacy coach
	SeasonStats           CollegePlayerSeasonStats
	Stats                 CollegePlayerStats
}

type NBAPlayerResponse struct {
	FirstName             string
	LastName              string
	Position              string
	Age                   int
	Year                  int
	State                 string
	Country               string
	Stars                 int
	Height                string
	PotentialGrade        string
	Shooting2Grade        string
	Shooting3Grade        string
	FreeThrowGrade        string
	FinishingGrade        string
	BallworkGrade         string
	ReboundingGrade       string
	InteriorDefenseGrade  string
	PerimeterDefenseGrade string
	Stamina               int
	PlaytimeExpectations  int
	Minutes               int
	Potential             int
	OverallGrade          string
	Personality           string
	RecruitingBias        string
	WorkEthic             string
	AcademicBias          string
	PlayerID              uint
	TeamID                uint
	TeamAbbr              string
	IsRedshirting         bool
	IsRedshirt            bool
	SeasonStats           NBAPlayerSeasonStats
	Stats                 NBAPlayerStats
}
