package models

type MatchType string

const (
	Practice      MatchType = "Practice"
	Qualification MatchType = "Qualification"
	QuarterFinal  MatchType = "QuarterFinal"
	SemiFinal     MatchType = "SemiFinal"
	Final         MatchType = "Final"
)

type Match struct {
	ID uint

	MatchType   MatchType
	MatchNumber uint

	Won     bool
	Carried bool

	TeamId uint
}
