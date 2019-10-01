package main

import (
	"fmt"

	"github.com/kylerberry/iron-wolves-raider-app/server/member"
)

var reportedStats = map[string]bool{"assists": true, "kills": true, "deaths": true, "timePlayedSeconds": true, "completed": true}

type playerReport struct {
	Member   member.DestinyUserInfo
	Activity CharacterActivityRef
	Stats    map[string]CharacterActivityStat
}

// NewPlayerReport creates a report of activityStats
func NewPlayerReport(player member.DestinyUserInfo, activity CharacterActivityRef, activityRuns []CharacterActivityDetail) playerReport {
	if len(activityRuns) == 0 {
		return playerReport{}
	}
	stats := map[string]CharacterActivityStat{}
	for _, a := range activityRuns {
		for stat, value := range a.Values {
			if _, reported := reportedStats[stat]; !reported {
				continue
			}
			if v, ok := stats[stat]; ok {
				existing := &v
				existing.Basic.Value += value.Basic.Value
				existing.Basic.DisplayValue = fmt.Sprintf("%f", v.Basic.Value)
			} else {
				stats[stat] = value
			}
		}
	}

	// derived stats
	stats["killsDeathRatio"] = derivedStat("killsDeathRatio", stats, activityRuns)
	stats["completedRatio"] = derivedStat("completedRatio", stats, activityRuns)
	stats["fastestTimeInSeconds"] = derivedStat("fastestTimeInSeconds", stats, activityRuns)
	return playerReport{player, activity, stats}
}

func derivedStat(stat string, stats map[string]CharacterActivityStat, activityRuns []CharacterActivityDetail) CharacterActivityStat {
	var v float64
	switch stat {
	case "completedRatio":
		v = float64(stats["completed"].Basic.Value / float64(len(activityRuns)))
	case "killsDeathRatio":
		v = stats["kills"].Basic.Value / stats["deaths"].Basic.Value
	case "fastestTimeInSeconds":
		v = fastestTime(activityRuns)
	}

	displayV := fmt.Sprintf("%f", v)
	return CharacterActivityStat{stat, StatValue{v, displayV}}
}

func fastestTime(activities []CharacterActivityDetail) float64 {
	var fastestTime float64
	for _, a := range activities {
		timeInSeconds := a.Values["timePlayedSeconds"]
		if fastestTime == 0.0 || timeInSeconds.Basic.Value < fastestTime {
			fastestTime = timeInSeconds.Basic.Value
		}
	}
	return fastestTime
}
