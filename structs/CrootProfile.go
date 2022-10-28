package structs

type CrootProfile struct {
	ID                 uint
	SeasonID           uint
	RecruitID          uint
	ProfileID          uint
	TotalPoints        float64
	AdjustedPoints     float64
	CurrentWeeksPoints int
	SpendingCount      int
	Scholarship        bool
	ScholarshipRevoked bool
	TeamAbbreviation   string
	RemovedFromBoard   bool
	IsSigned           bool
	IsLocked           bool
	CaughtCheating     bool
	SigningStatus      string
	Recruit            Croot
}

func (cp *CrootProfile) Map(rp PlayerRecruitProfile, c Croot) {
	cp.ID = rp.ID
	cp.SeasonID = rp.SeasonID
	cp.RecruitID = rp.RecruitID
	cp.ProfileID = rp.ProfileID
	cp.TotalPoints = rp.TotalPoints
	cp.CurrentWeeksPoints = rp.CurrentWeeksPoints
	cp.SpendingCount = rp.SpendingCount
	cp.Scholarship = rp.Scholarship
	cp.ScholarshipRevoked = rp.ScholarshipRevoked
	cp.TeamAbbreviation = rp.TeamAbbreviation
	cp.RemovedFromBoard = rp.RemovedFromBoard
	cp.IsSigned = rp.IsSigned
	cp.IsLocked = rp.IsLocked
	cp.Recruit = c
}

// Sorting Funcs
type ByCrootProfileTotal []CrootProfile

func (rp ByCrootProfileTotal) Len() int      { return len(rp) }
func (rp ByCrootProfileTotal) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp ByCrootProfileTotal) Less(i, j int) bool {
	return rp[i].TotalPoints > rp[j].TotalPoints
}
