package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kylerberry/iron-wolves-raider-app/server/member"
	"github.com/kylerberry/iron-wolves-raider-app/server/storage"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
)

var raidIDs = map[string]int{
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

type spaHandler struct {
	staticPath string
	indexPath  string
}

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

// member.clanMembers gets the clan members and wraps them in a struct
func getClanMembers() member.ClanMembers {
	// all members fit on 1 page atm
	page := "1"
	body := bungo.Get(apiRoot + "/GroupV2/" + clanID + "/Members?currentPage=" + page)
	m := member.ClanMembers{}
	jsonErr := json.Unmarshal(body, &m)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return m
}

func getActivities(m member.Member) CharactersActivities {
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

func getProfile(m member.Member) MemberProfile {
	memberType := strconv.Itoa(m.DestinyUserInfo.MembershipType)
	body := bungo.Get(apiRoot + "/Destiny2/" + memberType + "/Profile/" + m.DestinyUserInfo.MembershipID + "?components=100")
	mp := MemberProfile{}
	jsonErr := json.Unmarshal(body, &mp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return mp
}

// gets the clan leaderboard for an exact raid activity
func clanRaidLeaderboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clan := getClanMembers()
	rid := raidIDs[vars["raidSlug"]]
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
	}

	log.Println("Completed fetching member raid activities")

	resp, jsonErr := json.Marshal(playerReport)
	if jsonErr != nil {
		// @todo provide an error response
		log.Println(jsonErr)
	}
	w.Write(resp)
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func cronPopulateMembers() {
	log.Println("BEGIN cron to fetch clan members.")
	memberStore := storage.NewJSONStore("members")
	resp := []byte{}
	for i := 0; i < 2; i++ {
		page := strconv.Itoa(i)
		body := bungo.Get(apiRoot + "/GroupV2/" + clanID + "/Members?currentPage=" + page)
		resp = append(resp, body...)
	}
	memberStore.Write(resp)
	log.Println("END cron to fetch clan members.")
}

func main() {
	// establish Bungie env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bungo.APIKey = os.Getenv("API_KEY")
	apiRoot = os.Getenv("API_ROOT_PATH")

	// establish cron to populate storage
	c := cron.New()
	// On the minute, every 15 minutes
	c.AddFunc("1-14/15", cronPopulateMembers)
	// On the 5th minute, every 15 minutes
	// c.AddFunc("5-14/15", cronPopulateProfiles)
	// On the 10th minute, every 15 minutes
	// c.AddFunc("10-14/15", cronPopulateActivities)
	c.Start()

	router := mux.NewRouter()

	// api routes
	router.HandleFunc("/api/raid/{raidSlug}/leaderboard", clanRaidLeaderboard)

	// static root
	spa := spaHandler{staticPath: "build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
