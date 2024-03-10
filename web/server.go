package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/pviotti/pocket-counter/models"
)

func StartHTTPServer() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/unread", getUnreadCountHandler)

	log.Println("Starting HTTP server on port 8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Printf("Error starting HTTP server: %v", err)
	}
}

func serveIndex(w http.ResponseWriter, _ *http.Request) {
	unreadURL := os.Getenv("UNREAD_API_URL")
	if unreadURL == "" {
		http.Error(w, "UNREAD_API_URL environment variable not set", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		UnreadURL string
	}{
		UnreadURL: unreadURL,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
