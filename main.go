package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Profile struct {
	Applications []*Application `json:"applications"`
}

type Application struct {
	ApplicationId string `json:"applicationId"`
	Version       string `json:"version"`
}

func main() {
	const url = "https://ad48d74f-9fe5-4ce6-878c-e9f6c70c48bb.mock.pstmn.io/profiles/"

	id1 := &Application{
		ApplicationId: "music_app",
		Version:       "v1.4.10",
	}
	id2 := &Application{
		ApplicationId: "diagnostic_app",
		Version:       "v1.2.6",
	}
	id3 := &Application{
		ApplicationId: "settings_app",
		Version:       "v1.1.5",
	}
	applications := make([]*Application, 3)

	profile := Profile{
		Applications: append(applications, id1, id2, id3),
	}

	file, err := os.Open("./players.csv")
	if err != nil {
		log.Print(err)
		panic(err)
	}
	defer file.Close()

	// read from csv file line by line
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	// call each player to update their application(s)
	for fileScanner.Scan() {
		clientId := strings.Split(fileScanner.Text(), ",")[0]

		// skipping first line (column names)
		if clientId != "mac_addresses" {
			res, err := callPlayerToUpdate(url, clientId, &profile)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res)
		}
	}
}

// initiates the application(s) update on each player
func callPlayerToUpdate(url, clientId string, profile *Profile) (string, error) {
	body, _ := json.Marshal(profile)
	token := "12345"

	log.Printf("Calling player %v...", clientId)
	url += clientId
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
	defer res.Body.Close()

	return res.Status, nil
}
