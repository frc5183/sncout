package models

type Team struct {
	ID uint

	Name   string
	Number uint

	Robot Robot

	Matches []Match
}
