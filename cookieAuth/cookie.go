package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Session represents a user session
type session struct {
	username string
	expiry   time.Time
}

// Credentials represents the user login credentials
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Dummy users and password map
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var sessions = map[string]session{}

func main() {
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/logout", Logout)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Decode the request body into the credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if username exists and the password is correct
	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new session token
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Add the new session toke in the sessions map
	sessions[sessionToken] = session{
		username: creds.Username,
		expiry:   expiresAt,
	}

	// Store session token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// Get session token from cookie
	c, err := r.Cookie("session_token")

	//If error occurs
	if err != nil {

		//If there is no cookie
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	//Get session token from sessions map
	userSession, exists := sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if session is expired
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Welcome user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", userSession.username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	//Retrieve sessions token from cookie
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// Retrieve session from sessions map
	userSession, exists := sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Store the new session token in the map
	sessions[newSessionToken] = session{
		username: userSession.username,
		expiry:   expiresAt,
	}

	// Delete the old session from the map
	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Retrieve session token from cookie
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	// Delete the session from the sessions map
	delete(sessions, sessionToken)

	// Clear the session token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}
