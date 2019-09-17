package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var raidID = map[string]int{
	"crownofsorrows":   3333172150,
	"scourgeofthepast": 548750096,
	// "spireofstars":,
	// "eaterofworlds":,
	// "leviathan":,
	"lastwish": 2122313384,
}

var apiKey string
var apiRoot string
var clanID = "882490"
var bungo BungoRequester

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Wolves!")
}

func clanStats(w http.ResponseWriter, r *http.Request) {
	m := getClanMembers()
	// just a test value for now. In the real thing we'll loop through the members list
	me, nf := m.FindByName("Kypothesis")
	if nf != nil {
		log.Fatal(nf)
	}
	// mode 4 == raid activity
	modes := []string{"4"}
	statGroups := []string{"General"}
	// day start
	// day end
	// *can only fetch 31 days of activity
	memberType := strconv.Itoa(me.DestinyUserInfo.MembershipType)
	body := bungo.Get(apiRoot + "/Destiny2/" + memberType + "/Account/" + me.DestinyUserInfo.MembershipID + "/Character/0/Stats?modes=" + strings.Join(modes, ",") + "&groups=" + strings.Join(statGroups, ","))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// clanMembers gets the clan members and wraps them in a struct
func getClanMembers() ClanMembers {
	// all members fit on 1 page atm
	page := "1"
	body := bungo.Get(apiRoot + "/GroupV2/" + clanID + "/Members?currentPage=" + page)
	m := ClanMembers{}
	jsonErr := json.Unmarshal(body, &m)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return m
}

func getActivities(m Member) CharactersActivities {
	p := getProfile(m)
	memberType := strconv.Itoa(m.DestinyUserInfo.MembershipType)
	ca := CharactersActivities{}
	for _, cid := range p.Characters() {
		body := bungo.Get(apiRoot + "/Destiny2/" + memberType + "/Account/" + m.DestinyUserInfo.MembershipID + "/Character/" + cid + "/Stats/Activities?mode=4")
		c := CharacterActivity{}
		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		if len(c.Response.Activities) > 0 {
			ca[cid] = c
		}
	}
	return ca
}

func getProfile(m Member) MemberProfile {
	memberType := strconv.Itoa(m.DestinyUserInfo.MembershipType)
	body := bungo.Get(apiRoot + "/Destiny2/" + memberType + "/Profile/" + m.DestinyUserInfo.MembershipID + "?components=100")
	mp := MemberProfile{}
	jsonErr := json.Unmarshal(body, &mp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return mp
}

// func clanLeaderboard(w http.ResponseWriter, r *http.Request) {
// }

// gets the clan leaderboard for an exact raid activity
func clanRaidLeaderboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clan := getClanMembers()
	rid := raidID[vars["raidSlug"]]
	membersLength := strconv.Itoa(len(clan.All()))
	playerReport := map[string]interface{}{}
	// all members
	for i, member := range clan.All() {
		index := strconv.Itoa(i + 1)
		log.Println("fetching activities for " + member.DestinyUserInfo.DisplayName + " (" + index + "/" + membersLength + ")\n")
		// @todo: this needs to look back more than 31 days
		raidRuns := []CharacterActivityDetail{}
		charActivities := getActivities(member)
		// all characters for each member
		for _, activities := range charActivities {
			// all activities
			for _, activity := range activities.Response.Activities {
				// hold onto the activity if it matches our desired raid id
				if rid == activity.ActivityDetails.ReferenceID {
					raidRuns = append(raidRuns, activity)
				}
			}
		}

		if len(raidRuns) > 0 {
			activityRef := CharacterActivityRef{rid, 4}
			playerReport[member.DestinyUserInfo.MembershipID] = NewPlayerReport(member.DestinyUserInfo, activityRef, raidRuns)
		}

		log.Println(playerReport)
	}

	log.Println("Completed fetching member raid activities")

	resp, jsonErr := json.Marshal(playerReport)
	if jsonErr != nil {
		// @todo provide an error response
		log.Println(jsonErr)
	}
	w.Write(resp)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bungo.APIKey = os.Getenv("API_KEY")
	apiRoot = os.Getenv("API_ROOT_PATH")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	// router.HandleFunc("/api/player/{playerID}/raid/stats", playerActivityStats)
	router.HandleFunc("/api/raid/{raidSlug}/leaderboard", clanRaidLeaderboard)
	// router.HandleFunc("/api/raid/leaderboard", clanLeaderboard)

	log.Fatal(http.ListenAndServe(":8888", router))
}
