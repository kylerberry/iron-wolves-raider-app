package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	maxTop := "24"

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
	log.Fatal(http.ListenAndServe(":8888", router))
}

// get clan members
// get raid activity information
// be able to output raid info as csv or json
