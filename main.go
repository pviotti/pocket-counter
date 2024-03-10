package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type PocketRequest struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	State       string `json:"state"`
}

type PocketResponse struct {
	List map[string]interface{} `json:"list"`
}

func main() {
	consumerKey := os.Getenv("POCKET_CONSUMER_KEY")
	accessToken := os.Getenv("POCKET_ACCESS_TOKEN")

	if consumerKey == "" || accessToken == "" {
		fmt.Println("Please set POCKET_CONSUMER_KEY and POCKET_ACCESS_TOKEN environment variables")
		return
	}

	requestBody := PocketRequest{
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
		State:       "unread",
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error encoding request body: %v\n", err)
		return
	}

	url := "https://getpocket.com/v3/get"

	resp, err := http.Post(url, "application/json; charset=UTF-8", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Printf("Error fetching data from Pocket: %v\n", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Error fetching data from Pocket: %v\n", resp.Status)
		return
	}
	defer resp.Body.Close()
	// bodyBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Printf("Error reading response body: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Response body: %v\n", string(bodyBytes))

	var pocketResponse PocketResponse
	if err := json.NewDecoder(resp.Body).Decode(&pocketResponse); err != nil {
		fmt.Printf("Error decoding Pocket response: %v\n", err)
		return
	}

	fmt.Printf("Number of unread articles in Pocket: %d\n", len(pocketResponse.List))
}
