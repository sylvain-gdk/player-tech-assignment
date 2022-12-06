package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// using struct for the profile so that it is statically typed
type RequestBody struct {
	Profile *Profile
}
type Profile struct {
	Applications []*Application
}

type Application struct {
	ApplicationId string
	Version       string
}

var profile = RequestBody{
	&Profile{
		Applications: []*Application{
			{
				ApplicationId: "music_app",
				Version:       "v1.4.10",
			},
			{
				ApplicationId: "diagnostic_app",
				Version:       "v1.2.6",
			},
			{
				ApplicationId: "settings_app",
				Version:       "v1.1.5",
			},
		},
	},
}

// baseURL would need to be changed for real url
const baseURL string = "https://ad48d74f-9fe5-4ce6-878c-e9f6c70c48bb.mock.pstmn.io/profiles/"

func main() {
	/*
	 * go run command takes the path for the .csv file as the first argument and
	 * the JWT secret as the second argument i.e. go run main.go ./players.csv MyVeryPrivateSecret
	 */
	var (
		filePath     = os.Args[1]
		secretForJWT = os.Args[2]
	)
	// secret is passed from command line argument so it doesn't appear in git repo
	InitUpdatePlayers(baseURL, filePath, secretForJWT)
}

/*
 * creates a very basic JWT with expiration date
 */
func CreateNewToken(secret string, hoursToAdd int) (string, error) {
	currentPlusHours := time.Now().Add(time.Duration(hoursToAdd) * time.Hour)
	exp := time.Unix(currentPlusHours.Unix(), 0)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    "Application Update Tool",
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

/*
 * validates if a token is valid (date not expired, token not malformed)
 */
func IsValidToken(tokenString string) (bool, error) {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenExpired) {
		return false, err
	}
	return true, nil
}

/*
 * validates if the base url is in an acceptable format
 */
func IsValidURL(URI string) error {
	_, err := url.ParseRequestURI(URI)
	if err != nil {
		return err
	}
	return nil
}

/*
 * validates if a client id (mac address) is in the right format
 */
func IsValidClientId(id string) bool {
	pattern := "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})|([0-9a-fA-F]{4}\\.[0-9a-fA-F]{4}\\.[0-9a-fA-F]{4})$"
	result, _ := regexp.MatchString(pattern, id)
	return result
}

/*
 * starts the update process
 */
func InitUpdatePlayers(baseURL, filePath, secretForJWT string) {
	// execution should stop immediately if baseUrl is malformated
	err := IsValidURL(baseURL)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// execution should stop immediately if token fails to create
	token, err := CreateNewToken(secretForJWT, -6)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// open and read from .csv file line by line
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	// call each player to update their application(s)
	for fileScanner.Scan() {
		// using only the first column in .csv (mac_addresses)
		clientId := strings.Split(fileScanner.Text(), ",")[0]

		// skipping first line if column name
		if clientId != "mac_addresses" {
			// verify each clientId before sending request
			// (should log maybe in a file but not exit on error)
			if !IsValidClientId(clientId) {
				log.Printf("clientId: %v is not valid", clientId)
			} else {
				// calling player (should log error in a file but not exit on error)
				statusCode, err := CallPlayerToUpdate(baseURL, clientId, token)
				if err != nil {
					log.Println(err)
				}
				log.Printf("calling %v: %v", clientId, statusCode)
			}
		}
	}
}

/*
 * initiates the application(s) update process on each player
 */
func CallPlayerToUpdate(baseURL, clientId, token string) (int, error) {
	req, err := CreateRequest(baseURL, clientId, token)
	if err != nil {
		return 0, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	return res.StatusCode, nil
}

/*
 * creates a new request for a specific player, based on it's client id
 */
func CreateRequest(baseURL, clientId, token string) (*http.Request, error) {
	// creates url for a specific player
	url := baseURL + clientId

	// marshalling the profile into []byte before sending as a buffer
	body, _ := json.Marshal(profile)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	// setting headers on the request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-client-id", clientId)
	req.Header.Add("x-authentication-token", token)

	return req, nil
}
