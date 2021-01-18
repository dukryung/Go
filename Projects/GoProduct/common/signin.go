package common

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http:/localhost:8080/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Pircture      string `json:"picture"`
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state

}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthstate.Value {
		errMsg := fmt.Sprintf("invalid google oauth cookie : %s, state : %s\n", oauthstate.Value, r.FormValue("state"))
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userinfo GoogleUserInfo
	err = json.Unmarshal(data, &userinfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["id"] = userinfo.ID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	res, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Falied to Exchange %s\n", err.Error())
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Falied to Exchange %s\n", err.Error())
	}

	return data, nil
}
