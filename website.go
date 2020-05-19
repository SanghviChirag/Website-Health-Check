package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Set DB Connection
var db *gorm.DB
var err error

// InitialMigration to migrate model
func InitialMigration() {

	db, err = gorm.Open("sqlite3", "website.db")
	if err != nil {
		panic("Failed to connect!")
	}
	defer db.Close()

	db.AutoMigrate(&Website{})
}

// WebsiteHealthStatusHistory to set website-status-history
type WebsiteHealthStatusHistory struct {
	websiteCheckDateTime time.Time
	isSuccess            bool
}

// Website content
type Website struct {
	gorm.Model
	URL                string
	method             string
	body               []byte
	header             []byte
	expectedStatusCode int
	checkInterval      int
	healthStatus       []WebsiteHealthStatusHistory
}

type regWebReqBody struct {
	links []Website
}

// Handle Registration of Website
func registerWebsite(w http.ResponseWriter, r *http.Request) {
	// Option #1: Keep checking using channels
	// Option #2: Set crons
	db, err = gorm.Open("sqlite3", "website.db")
	if err != nil {
		panic("Could not connect to database")
	}
	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var reqBody regWebReqBody
	err := decoder.Decode(&reqBody)
	if err != nil {
		panic(err)
	}
	for _, website := range reqBody.links {
		db.Create(website)
	}

	fmt.Fprintf(w, "New User Successfully Created.")
}

// Get all websites data
func getAllWebsiteInfo(w http.ResponseWriter, r *http.Request) {
	// Get Website Details
	db, err = gorm.Open("sqlite3", "website.db")
	if err != nil {
		panic("Could not connect to database")
	}
	defer db.Close()

	var websites []Website
	db.Find(&websites)
	json.NewEncoder(w).Encode(websites)
}

func getWebsite(w http.ResponseWriter, r *http.Request) {
	// Get Website Details
	db, err = gorm.Open("sqlite3", "website.db")
	if err != nil {
		panic("Could not connect to database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	url := vars["url"]
	var websites Website
	db.Where("URL=?", url).Find(&websites)
	json.NewEncoder(w).Encode(websites)
}

// func main() {
// 	links := []Website{
// 		Website{
// 			URL:                "http://google.com",
// 			method:             "GET",
// 			expectedStatusCode: 200,
// 			checkInterval:      50,
// 		},
// 		Website{
// 			URL:                "http://amazon.com",
// 			method:             "GET",
// 			expectedStatusCode: 200,
// 			checkInterval:      50,
// 		},
// 	}

// 	c := make(chan Website)

// 	for _, website := range links {
// 		go checkLink(website, c)
// 	}

// 	fmt.Println(<-c)
// 	fmt.Println(<-c)
// }

func checkLink(website Website) {

	db, err = gorm.Open("sqlite3", "website.db")
	if err != nil {
		panic("Could not connect to database")
	}
	defer db.Close()

	switch website.method {
	case "GET":
		res, _ := http.Get(website.URL)

		healthStatus := WebsiteHealthStatusHistory{
			websiteCheckDateTime: time.Now().UTC(),
			isSuccess:            website.expectedStatusCode == res.StatusCode,
		}
		website.healthStatus = append(website.healthStatus, healthStatus)
		db.Model(&website).Updates(website)
		return

	default:
		fmt.Println("StatusNotAllowed")
		return
	}

}
