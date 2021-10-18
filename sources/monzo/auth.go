package monzo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	clientID     = "oauth2client_0000A6qzWDoTLNP5UsXcP3"
	clientSecret = "mnzconf.1r2wL7WSwlexh3ApOPMQUHYurwzgjVQWfAULv9cMWD4pzd3nfJfgT6pN+gVH8+Fc17Qr1cihxna7EBgASvIivQ=="
)

type oauthCallbackResp struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
}

// StartOauth redirects the user to Monzo to authenticate the account.
func (m *Monzo) StartOauth(w http.ResponseWriter, r *http.Request) {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", "http://localhost:8080/api/auth/monzo/callback")
	q.Set("response_type", "code")
	authURL := "https://auth.monzo.com?" + q.Encode()
	fmt.Println(authURL)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// CompleteOauth
func (m *Monzo) CompleteOauth(w http.ResponseWriter, r *http.Request) {
	// read temporary access code from query
	code := r.URL.Query().Get("code")

	// use the temporary auth code to get an access token
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("redirect_uri", "http://localhost:8080/api/auth/monzo/callback")
	form.Set("code", code)

	req, err := http.NewRequest(http.MethodPost, "https://api.monzo.com/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("failed to create token request: %s", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("failed to perform token request: %s", err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read token response body: %s", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("non-200 status for token request: %s, body: %s", resp.Status, b)
		return
	}

	var authCallback oauthCallbackResp
	if err := json.Unmarshal(b, &authCallback); err != nil {
		log.Printf("failed to JSON decode token response body: %s, %s", err, b)
		return
	}

	m.accessToken = authCallback.AccessToken
	m.refreshToken = authCallback.RefreshToken
	m.expiresIn = authCallback.ExpiresIn

	http.RedirectHandler("/sources", http.StatusTemporaryRedirect)
}
