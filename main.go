package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type PocketRequest struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	State       string `json:"state"`
}

type PocketResponse struct {
	List map[string]interface{} `json:"list"`
}

const POCKET_GET_URL = "https://getpocket.com/v3/get"
const POCKET_GET_CONTENTTYPE = "application/json; charset=UTF-8"
const SLEEP_TIME = 12 * time.Hour
const DATABASE_PATH = "./data/pocket.db"

func main() {
	consumerKey := os.Getenv("POCKET_CONSUMER_KEY")
	accessToken := os.Getenv("POCKET_ACCESS_TOKEN")

	if consumerKey == "" || accessToken == "" {
		log.Println("Please set POCKET_CONSUMER_KEY and POCKET_ACCESS_TOKEN environment variables")
		return
	}

	go startHTTPServer()

	for {
		fetchAndSave(consumerKey, accessToken)
		log.Printf("Sleeping for %v hours", SLEEP_TIME.Hours())
		time.Sleep(SLEEP_TIME)
	}
}

func startHTTPServer() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/unread", getUnreadCountHandler)

	log.Println("Starting HTTP server on port 8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Printf("Error starting HTTP server: %v", err)
	}
}

func getUnreadCountHandler(w http.ResponseWriter, r *http.Request) {
	unreadCounts, err := getUnreadCountsForPastYearFromDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching unread count: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(unreadCounts)
}

func getUnreadCountsForPastYearFromDB() ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	var unreadCounts []map[string]interface{}
	rows, err := db.Query("select date, count from unread_count where date >= date('now', '-1 year') order by date")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		unreadCounts = append(unreadCounts, map[string]interface{}{
			"date":         date,
			"unread_count": count,
		})
	}

	return unreadCounts, nil
}

func fetchAndSave(consumerKey, accessToken string) {
	unreadCount, err := getUnreadCount(consumerKey, accessToken)
	if err != nil {
		log.Printf("Error getting unread count: %v\n", err)
		return
	}

	err = saveToDatabase(unreadCount)
	if err != nil {
		log.Printf("Error saving to database: %v\n", err)
		return
	}

	log.Printf("Unread count: %d\n", unreadCount)
}

func getUnreadCount(consumerKey, accessToken string) (int, error) {
	requestBody := PocketRequest{
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
		State:       "unread",
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error encoding request body: %v\n", err)
		return 0, fmt.Errorf("error encoding request body: %v", err)
	}

	resp, err := http.Post(POCKET_GET_URL, POCKET_GET_CONTENTTYPE, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Printf("Error fetching data from Pocket: %v\n", err)
		return 0, fmt.Errorf("error fetching data from Pocket: %v", err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Error fetching data from Pocket: %v\n", resp.Status)
		return 0, fmt.Errorf("error fetching data from Pocket: %v", resp.Status)
	}
	defer resp.Body.Close()

	var pocketResponse PocketResponse
	if err := json.NewDecoder(resp.Body).Decode(&pocketResponse); err != nil {
		fmt.Printf("Error decoding Pocket response: %v\n", err)
		return 0, fmt.Errorf("error decoding Pocket response: %v", err)
	}

	return len(pocketResponse.List), nil
}

func saveToDatabase(unreadCount int) error {
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`create table if not exists unread_count (
        date date primary key,
        count integer
    )`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	date := time.Now().Format("2006-01-02")
	_, err = db.Exec("insert or replace into unread_count (date, count) values (?, ?)", date, unreadCount)
	if err != nil {
		return fmt.Errorf("error inserting into table: %v", err)
	}

	return nil
}
