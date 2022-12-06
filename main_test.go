package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a = Args{
	BaseURL:      "",
	FilePath:     "./fixtures/clients.csv",
	SecretForJWT: "secret",
	Token:        "",
}

/*
 * creates a fake server to mock the responses when calling API
 */
func mockServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientId := r.Header.Get("x-client-id")
			token := r.Header.Get("x-authentication-token")
			tokenValid, _ := IsValidToken(token)

			if clientId == "a1:bb:cc:dd:ee:ff" && tokenValid {
				w.WriteHeader(http.StatusOK)
			} else if !tokenValid || !IsValidClientId(clientId) {
				w.WriteHeader(http.StatusUnauthorized)
			} else if clientId != "a1:bb:cc:dd:ee:ff" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))
}

func TestIsValidTokenTrue(t *testing.T) {
	token, _ := CreateNewToken("somesecret", 6)
	got, _ := IsValidToken(token)
	want := true

	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestIsValidTokenFalse(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ.9eyJpc3MiOiJBcHBsaWNhdGlvbiBVcGRhdGUgVG9vbCIsImV4cCI6MTY3MDI5ODMxMH0.COX6oFXHTy7z2q08puv0X2jBtttn7lhqKe766P9C52U"
	_, got := IsValidToken(token)
	want := fmt.Errorf("invalid token")

	if got.Error() != want.Error() {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestIsValidTokenFalseExpired(t *testing.T) {
	token, _ := CreateNewToken("somesecret", -6)
	_, got := IsValidToken(token)
	want := fmt.Errorf("invalid token")

	if got.Error() != want.Error() {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestIsValidClientIdTrue(t *testing.T) {
	got := IsValidClientId("a1:bb:cc:dd:ee:ff")
	want := true

	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestIsValidClientIdFalse(t *testing.T) {
	got := IsValidClientId("a1:bb:cc:dd:ee:f")
	want := false

	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestCallPlayerToUpdateStatusOk(t *testing.T) {
	mock := mockServer()
	defer mock.Close()

	clientId := "a1:bb:cc:dd:ee:ff"
	token, _ := CreateNewToken("somesecret", 6)
	a.Token = token
	a.BaseURL = fmt.Sprintf(mock.URL + "/profiles/")

	got, _ := a.CallPlayerToUpdate(clientId)
	want := http.StatusOK

	if got != want {
		t.Errorf("got StatusCode %v, want StatusCode %v", got, want)
	}
}

func TestCallPlayerToUpdateStatusClientNotfound(t *testing.T) {
	mock := mockServer()
	defer mock.Close()

	clientId := "a2:bb:cc:dd:ee:ff"
	token, _ := CreateNewToken("somesecret", 6)
	a.Token = token
	a.BaseURL = fmt.Sprintf(mock.URL + "/profiles/")

	got, _ := a.CallPlayerToUpdate(clientId)
	want := http.StatusNotFound

	if got != want {
		t.Errorf("got statusCode %v, want statusCode %v", got, want)
	}
}

func TestCallPlayerToUpdateStatusUnauthorizedToken(t *testing.T) {
	mock := mockServer()
	defer mock.Close()

	clientId := "a1:bb:cc:dd:ee:ff"
	// creates an expired token
	token, _ := CreateNewToken("secret", -6)
	a.Token = token
	a.BaseURL = fmt.Sprintf(mock.URL + "/profiles/")

	got, _ := a.CallPlayerToUpdate(clientId)
	want := http.StatusUnauthorized

	if got != want {
		t.Errorf("got statusCode %v, want statusCode %v", got, want)
	}
}
