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

type authAccessDetails struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
}

// AuthenticateHandler starts the OAuth2 authentication sequence, requesting a temporary access code from Monzo. This
// endpoint also receives callback requests from Monzo with the temporary access code.
func (m *Monzo) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	// first step of oauth - request a temporary access code from monzo
	code := r.URL.Query().Get("code")
	if code == "" {
		q := url.Values{}
		q.Set("client_id", m.currentAuth.ClientID)
		q.Set("redirect_uri", m.serviceRedirectURL)
		q.Set("response_type", "code")
		authURL := "https://auth.monzo.com?" + q.Encode()

		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
		return
	}

	// second step of oauth - monzo sent a temporary access code - request an access token from monzo
	if err := m.fetchAccessToken(code, accessCodeInitial); err != nil {
		log.Printf("failed to fetch access token: %s", err)
	}

	http.Redirect(w, r, m.webAppRedirectURL, http.StatusFound)
}

const (
	accessCodeInitial int = iota
	accessCodeRefresh
)

func (m *Monzo) fetchAccessToken(code string, requestType int) error {
	// use the temporary auth code to get an access token
	form := url.Values{}
	form.Set("client_id", m.currentAuth.ClientID)
	form.Set("client_secret", m.clientSecret)
	if requestType == accessCodeInitial {
		// first access token request
		form.Set("grant_type", "authorization_code")
		form.Set("redirect_uri", m.serviceRedirectURL)
		form.Set("code", code)
	} else {
		// refresh access token request
		form.Set("grant_type", "refresh_token")
		form.Set("refresh_token", code)
	}

	// third step of oauth - exchange the temporary access code for an access token
	req, err := http.NewRequest(http.MethodPost, "https://api.monzo.com/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %s", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform token request: %s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status for token request: %s, body: %s", resp.Status, b)
	}

	var authCallback authAccessDetails
	if err := json.Unmarshal(b, &authCallback); err != nil {
		return fmt.Errorf("failed to JSON decode token response body: %s, %s", err, b)
	}

	// update auth details
	m.authRefreshedChan <- authCallback
	return nil
}
