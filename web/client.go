package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	POCKET_GET_URL         = "https://getpocket.com/v3/get"
	POCKET_GET_CONTENTTYPE = "application/json; charset=UTF-8"
	MAX_RETRIES            = 5
	TIMEOUT                = 10 * time.Second
)

type PocketRequest struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	State       string `json:"state"`
	DetailType  string `json:"detailType"`
	Sort        string `json:"sort"`
	Count       int    `json:"count"`
	Offset      int    `json:"offset"`
	Total       int    `json:"total"`
}

type PocketResponse struct {
	Status     int                    `json:"status"`
	List       map[string]interface{} `json:"list"`
	Total      int                    `json:"total"`      // this is actually not returned by the API
	MaxActions int                    `json:"maxActions"` // undocumented field
	CacheType  string                 `json:"cacheType"`  // undocumented field
	Error      string                 `json:"error"`      // undocumented field
	Complete   int                    `json:"complete"`   // undocumented field
	Since      int                    `json:"since"`      // undocumented field
}

func GetUnreadCount(consumerKey, accessToken string) (int, error) {
	requestBody := PocketRequest{
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
		State:       "unread",
		DetailType:  "simple",
		Sort:        "newest",
		Count:       30,
		Offset:      0,
		Total:       1,
	}

	// Need to use retryablehttp because API returns HTTP 504 quite often
	client := retryablehttp.NewClient()
	client.RetryMax = MAX_RETRIES
	client.HTTPClient.Timeout = TIMEOUT

	totItems := 0

	for {
		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			log.Printf("Error encoding request body: %v\n", err)
			return 0, fmt.Errorf("error encoding request body: %v", err)
		}

		resp, err := client.Post(POCKET_GET_URL, POCKET_GET_CONTENTTYPE, bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			log.Printf("Error fetching data from Pocket: %v\n", err)
			return 0, fmt.Errorf("error fetching data from Pocket: %v", err)
		}
		if resp.StatusCode != 200 {
			log.Printf("Error fetching data from Pocket: %v\n", resp.Status)
			return 0, fmt.Errorf("error fetching data from Pocket: %v", resp.Status)
		}
		defer resp.Body.Close()

		var pocketResponse PocketResponse
		if err := json.NewDecoder(resp.Body).Decode(&pocketResponse); err != nil {
			log.Printf("Error decoding Pocket response: %v\n", err)
			return 0, fmt.Errorf("error decoding Pocket response: %v", err)
		}

		if len(pocketResponse.List) == 0 {
			break
		}

		log.Printf("[DEBUG] fetched %d new items", len(pocketResponse.List))
		totItems += len(pocketResponse.List)
		requestBody.Offset += len(pocketResponse.List)
	}

	return totItems, nil
}
