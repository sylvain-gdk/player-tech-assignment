package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// using struct for the profile so that it is statically typed
type Profile struct {
	Applications []*Application `json:"applications"`
}

type Application struct {
	ApplicationId string `json:"applicationId"`
	Version       string `json:"version"`
}

func main() {
	var filePath = os.Args[1]
	// "./players.csv"
	const url = "https://ad48d74f-9fe5-4ce6-878c-e9f6c70c48bb.mock.pstmn.io/profiles/"

	applications := []*Application{
		{"music_app", "v1.4.10"},
		{"diagnostic_app", "v1.2.6"},
		{"settings_app", "v1.1.5"},
	}
	profile := Profile{
		Applications: applications,
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Print(err)
		panic(err)
	}
	defer file.Close()

	// read from .csv file line by line
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	// call each player to update their application(s)
	for fileScanner.Scan() {
		// using only the first column (mac_addresses)
		clientId := strings.Split(fileScanner.Text(), ",")[0]

		// skipping first line if its column name
		if clientId != "mac_addresses" {
			res, err := callPlayerToUpdate(url, clientId, &profile)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(res)
		}
	}
}

// initiates the application(s) update on each player
func callPlayerToUpdate(url, clientId string, profile *Profile) (string, error) {
	token := "12345" // TODO: use JWT
	url += clientId

	log.Printf("Calling player: %v...", clientId)
	// marshalling the profile into []byte
	body, _ := json.Marshal(profile)
	// using http.NewRequest so I can add headers to the request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	// setting headers to the request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-client-id", clientId)
	req.Header.Add("x-authentication-token", token)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	// returning status for testing
	return res.Status, nil
}
