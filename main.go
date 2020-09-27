package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {

	http.HandleFunc("/login", Login)
	http.HandleFunc("/callback", Callback)

	fmt.Println("server running on port 4004")
	http.ListenAndServe(":4004", nil)
}

var (
	CLIENT_ID     = os.Getenv("CLIENT_ID")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
)

func Login(w http.ResponseWriter, r *http.Request) {

	var scope = "user-read-private user-read-email"
	spotify_url := fmt.Sprintf("%s?response_type=%s&client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		"https://accounts.spotify.com/authorize",
		"code",
		CLIENT_ID,
		scope,
		"http://localhost:4004/callback",
		"qwertyuiopasdfgh",
	)
	http.Redirect(w, r, spotify_url, http.StatusTemporaryRedirect)
}

func Callback(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["code"]
	if !ok {
		fmt.Println("query code missing")
		return
	}
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	code := string(keys[0])

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", "http://localhost:4004/callback")

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(data.Encode())))
	if err != nil {
		fmt.Println("error new request")
		fmt.Println(err.Error())
		return
	}

	tokenStr := fmt.Sprintf("%s:%s", CLIENT_ID, CLIENT_SECRET)
	token := b64.StdEncoding.EncodeToString([]byte(tokenStr))
	basicToken := fmt.Sprintf("basic %s", token)

	req.Header.Set("Authorization", basicToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("err get request")
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {

		var resbody map[string]interface{}

		err = json.NewDecoder(res.Body).Decode(&resbody)
		if err != nil {
			fmt.Println("error decode xml")
			fmt.Println(err.Error())
			return
		}

		fmt.Println(resbody)
	} else {
		fmt.Println(res.StatusCode)
		fmt.Println("cannot authorized")
	}
}
