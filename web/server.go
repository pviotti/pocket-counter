package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pviotti/pocket-counter/models"
)

func StartHTTPServer() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/unread", getUnreadCountHandler)

	log.Println("Starting HTTP server on port 8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Printf("Error starting HTTP server: %v", err)
	}
}

func getUnreadCountHandler(w http.ResponseWriter, r *http.Request) {
	unreadCounts, err := models.GetUnreadCountsForPastYearFromDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching unread count: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(unreadCounts)
}
