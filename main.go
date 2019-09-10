package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var apiKey string
var apiRoot string
var clanID = "882490"

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Wolves!")
}

// raids only right now, but could be expanded
func clanLeaderboard(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	// mode 4 == raid activity
	queryModes := []string{"4"}
	maxTop := "25"

	// @todo this endpoint is buggy, we'll have to manually get stats for every member :(
	req, clientErr := http.NewRequest("GET", apiRoot+"/Destiny2/Stats/Leaderboards/Clans/"+clanID+"?modes="+strings.Join(queryModes, ",")+"&maxtop="+maxTop, nil)
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	req.Header.Add("X-API-Key", apiKey)
	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func clanStats(w http.ResponseWriter, r *http.Request) {
	m := getClanMembers()
	client := &http.Client{}
	// just a test value for now. In the real thing we'll loop through the members list
	me, nf := m.Find("Kypothesis")
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
	req, clientErr := http.NewRequest("GET", apiRoot+"/Destiny2/"+memberType+"/Account/"+me.DestinyUserInfo.MembershipID+"/Character/0/Stats?modes="+strings.Join(modes, ",")+"&groups="+strings.Join(statGroups, ","), nil)
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	req.Header.Add("X-API-Key", apiKey)
	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func clanActivities(w http.ResponseWriter, r *http.Request) {
	m := getClanMembers()
	client := &http.Client{}
	// just a test value
	me, nf := m.Find("Kypothesis")
	if nf != nil {
		log.Fatal(nf)
	}
	p := getProfile(me)
	memberType := strconv.Itoa(me.DestinyUserInfo.MembershipType)
	output := []byte{}
	for _, cid := range p.Response.Profile.Data.CharacterIDs {
		req, clientErr := http.NewRequest("GET", apiRoot+"/Destiny2/"+memberType+"/Account/"+me.DestinyUserInfo.MembershipID+"/Character/"+cid+"/Stats/Activities?mode=4", nil)
		if clientErr != nil {
			log.Fatal(clientErr)
		}
		req.Header.Add("X-API-Key", apiKey)
		resp, getErr := client.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}
		defer resp.Body.Close()
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		output = append(output, body...)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// clanMembers gets the clan members and wraps them in a struct
func getClanMembers() Members {
	client := &http.Client{}
	page := "1"
	req, clientErr := http.NewRequest("GET", apiRoot+"/GroupV2/"+clanID+"/Members?currentPage="+page, nil)
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	req.Header.Add("X-API-Key", apiKey)
	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	m := Members{}
	jsonErr := json.Unmarshal(body, &m)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return m
}

func getProfile(m Member) MemberProfile {
	client := &http.Client{}
	memberType := strconv.Itoa(m.DestinyUserInfo.MembershipType)
	req, clientErr := http.NewRequest("GET", apiRoot+"/Destiny2/"+memberType+"/Profile/"+m.DestinyUserInfo.MembershipID+"?components=100", nil)
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	req.Header.Add("X-API-Key", apiKey)
	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	mp := MemberProfile{}
	jsonErr := json.Unmarshal(body, &mp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return mp
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey = os.Getenv("API_KEY")
	apiRoot = os.Getenv("API_ROOT_PATH")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/clan-leaderboard", clanLeaderboard)
	router.HandleFunc("/clan-stats", clanStats)
	router.HandleFunc("/clan-activities", clanActivities)
	log.Fatal(http.ListenAndServe(":8888", router))
}

// -√ get clan members
// -√ get profile for each member
// -√ get characters for each member
// - get raid activity history for each member->character
// - get stats for each activity
