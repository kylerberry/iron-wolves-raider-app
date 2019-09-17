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

// GetActivity finds an activity for a character
// func (ca CharacterActivity) GetActivity(id int) (CharacterActivityDetail, error) {
// 	for _, a := range ca.Response.Activities {
// 		if a.ActivityDetails.ReferenceID == id {
// 			return a, nil
// 		}
// 	}
// 	return CharacterActivityDetail{}, errors.New("Unable to find activity matching reference: " + string(id))
// }

// GetStat will return the stat information for an activity
// func (ca CharacterActivityDetail) GetStat(stat string) (CharacterActivityStat, error) {
// 	s, ok := ca.Values[stat]
// 	if ok {
// 		return s, nil
// 	}
// 	return CharacterActivityStat{}, errors.New("Unable to find activity stat: " + stat)
// }
