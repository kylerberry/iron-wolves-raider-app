package main

// CharacterActivity is the top level wrapper for the activity response
type CharacterActivity struct {
	Response Response
}

// CharactersActivities is a map of Profile.Data.CharacterId -> CharacterActivity
type CharactersActivities map[string]CharacterActivity

// Response is the requested data from bungo
type Response struct {
	Activities []CharacterActivityDetail
}

// CharacterActivityDetail is an event played by a character
type CharacterActivityDetail struct {
	ActivityDetails CharacterActivityRef
	Values          map[string]CharacterActivityStat
}

// CharacterActivityRef is useful information about that activity
type CharacterActivityRef struct {
	ReferenceID int
	Mode        int
}

// CharacterActivityStat is a stat performed during the activity
type CharacterActivityStat struct {
	StatID string
	Basic  StatValue
}

// StatValue is the activity stat values
type StatValue struct {
	Value        float64
	DisplayValue string
}
