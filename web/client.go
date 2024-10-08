package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const POCKET_GET_URL = "https://getpocket.com/v3/get"
const POCKET_GET_CONTENTTYPE = "application/json; charset=UTF-8"

type PocketRequest struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	State       string `json:"state"`
}

type PocketResponse struct {
	List map[string]interface{} `json:"list"`
}

func GetUnreadCount(consumerKey, accessToken string) (int, error) {
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
