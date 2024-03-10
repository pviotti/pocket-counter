package main

import (
	"log"
	"os"
	"time"

	"github.com/pviotti/pocket-counter/models"
	"github.com/pviotti/pocket-counter/web"
)

const SLEEP_TIME = 12 * time.Hour

func main() {
	consumerKey := os.Getenv("POCKET_CONSUMER_KEY")
	accessToken := os.Getenv("POCKET_ACCESS_TOKEN")

	if consumerKey == "" || accessToken == "" {
		log.Println("Please set POCKET_CONSUMER_KEY and POCKET_ACCESS_TOKEN environment variables")
		return
	}

	go web.StartHTTPServer()

	for {
		fetchAndSave(consumerKey, accessToken)
		log.Printf("Sleeping for %v hours", SLEEP_TIME.Hours())
		time.Sleep(SLEEP_TIME)
	}
}

func fetchAndSave(consumerKey, accessToken string) {
	unreadCount, err := web.GetUnreadCount(consumerKey, accessToken)
	if err != nil {
		log.Printf("Error getting unread count: %v\n", err)
		return
	}

	err = models.SaveToDatabase(unreadCount)
	if err != nil {
		log.Printf("Error saving to database: %v\n", err)
		return
	}

	log.Printf("Unread count: %d\n", unreadCount)
}
