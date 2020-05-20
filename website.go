package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

	db.AutoMigrate(&Website{}, &WebsiteHealthStatusHistory{})
}

// WebsiteHealthStatusHistory to set website-status-history
type WebsiteHealthStatusHistory struct {
	gorm.Model
	WebsiteID            uint
	WebsiteCheckDateTime time.Time `json: "websiteCheckDateTime"`
	IsSuccess            bool      `json: "isSuccess"`
}

// Website content
type Website struct {
	gorm.Model
	URL                string                       `json: "URL"`
	Method             string                       `json: "method"`
	Body               []byte                       `json: "body"`
	Header             []byte                       `json: "header"`
	ExpectedStatusCode int                          `json: "expectedStatusCode"`
	CheckInterval      int                          `json: "checkInterval"`
	HealthStatus       []WebsiteHealthStatusHistory `json: "healthStatus" gorm:"foreignkey:WebsiteRefer"`
}

type regWebReqBody struct {
	Websites []Website
}

// Handle Registration of Website
func registerWebsite(w http.ResponseWriter, r *http.Request) {
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

	for _, website := range reqBody.Websites {

		var web Website
		if db.Where("URL = ?", website.URL).First(&web).RecordNotFound() {
			db.Create(&website)
			_ = setCron(website)
		}
	}

	fmt.Fprintf(w, "Website(s) Successfully Created.")
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

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonBytes, _ := json.Marshal(websites)
	w.Write(jsonBytes)
}

func getWebsite(w http.ResponseWriter, r *http.Request) {
	// Get Website Details
	db, err = gorm.Open("sqlite3", "website.db")
	if err != nil {
		panic("Could not connect to database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	websiteID := vars["id"]
	var statusHistory []WebsiteHealthStatusHistory
	db.Where("website_id=?", websiteID).Find(&statusHistory)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonBytes, _ := json.Marshal(statusHistory)
	w.Write(jsonBytes)
}

func checkLink(website Website) {

	db, err = gorm.Open("sqlite3", "website.db")
	if err != nil {
		panic("Could not connect to database")
	}
	defer db.Close()
	fmt.Println("GETTING CALLED!!!!!")
	switch website.Method {
	case "GET":
		res, _ := http.Get(website.URL)

		healthStatus := WebsiteHealthStatusHistory{
			WebsiteID:            website.ID,
			WebsiteCheckDateTime: time.Now().UTC(),
			IsSuccess:            website.ExpectedStatusCode == res.StatusCode,
		}

		db.Create(&healthStatus)
		fmt.Println(website.HealthStatus)
		return

	default:
		fmt.Println("StatusNotAllowed")
		return
	}

}
