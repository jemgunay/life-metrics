package monzo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jemgunay/life-metrics/sources"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// Monzo represents the Monzo collection source.
type Monzo struct {
	auth              authAccessDetails
	authRefreshedChan chan authAccessDetails
}

// New TODO
func New() *Monzo {
	m := &Monzo{
		authRefreshedChan: make(chan authAccessDetails, 1),
	}
	refreshTimer := time.NewTimer(time.Hour)
	go func() {
		for {
			select {
			case <-refreshTimer.C:
				m.fetchAccessToken(m.auth.RefreshToken, accessCodeRefresh)

			case m.auth = <-m.authRefreshedChan:
				// reset auth refresh timer
				refreshTimer.Stop()
				timeToRefresh := time.Second * time.Duration(m.auth.ExpiresIn)
				refreshTimer = time.NewTimer(timeToRefresh)
				log.Printf("Monzo authenticated - next authentication in %s", timeToRefresh)
			}
		}
	}()
	return m
}

// Name returns the source name.
func (m *Monzo) Name() string {
	return "monzo"
}

type accountsResult struct {
	Accounts []account `json:"accounts"`
}

type account struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Collect TODO
func (m *Monzo) Collect(start, end time.Time) ([]sources.Result, error) {
	if m.auth.AccessToken == "" {
		return nil, errors.New("access token not set - oauth setup required")
	}

	account, err := m.getAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %s", err)
	}

	// get transactions for account
	transactions, err := m.getTransactions(account.ID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %s", err)
	}

	for _, transaction := range transactions.Transactions {
		if transaction.Category == "eating_out" {
			log.Printf("%s, %s -> %v", transaction.Merchant, transaction.Description, transaction.Amount)
		}
	}

	return nil, nil
}

func (m *Monzo) getAccount() (account, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.monzo.com/accounts", nil)
	if err != nil {
		return account{}, fmt.Errorf("failed to create accounts request: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+m.auth.AccessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return account{}, fmt.Errorf("failed to perform accounts request: %s", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return account{}, fmt.Errorf("failed to read accounts response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return account{}, fmt.Errorf("non-200 status for accounts request: %s, body: %s", resp.Status, b)
	}

	var accounts accountsResult
	if err := json.Unmarshal(b, &accounts); err != nil {
		return account{}, fmt.Errorf("failed to JSON decode accounts response body: %s, %s", err, b)
	}

	if len(accounts.Accounts) == 0 {
		return account{}, errors.New("no accounts found")
	}

	return accounts.Accounts[0], nil
}

type transactionsResult struct {
	Transactions []transaction `json:"transactions"`
}

type transaction struct {
	AccountBalance int    `json:"account_balance"`
	Amount         int    `json:"amount"`
	CreatedTime    string `json:"created"`
	Currency       string `json:"currency"`
	Description    string `json:"description"`
	ID             string `json:"id"`
	Merchant       string `json:"merchant"`
	Notes          string `json:"notes"`
	IsLoad         bool   `json:"is_load"`
	Settled        string `json:"settled"`
	Category       string `json:"category"`
}

func (m *Monzo) getTransactions(accountID string, start, end time.Time) (transactionsResult, error) {
	var transactions transactionsResult

	q := url.Values{}
	q.Set("account_id", accountID)
	q.Set("since", start.Format(time.RFC3339))
	q.Set("before", end.Format(time.RFC3339))

	req, err := http.NewRequest(http.MethodGet, "https://api.monzo.com/transactions?"+q.Encode(), nil)
	if err != nil {
		return transactions, fmt.Errorf("failed to create accounts request: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+m.auth.AccessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return transactions, fmt.Errorf("failed to create transactions request: %s", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return transactions, fmt.Errorf("failed to read transactions response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return transactions, fmt.Errorf("non-200 status for transactions request: %s, body: %s", resp.Status, b)
	}

	if err := json.Unmarshal(b, &transactions); err != nil {
		return transactions, fmt.Errorf("failed to JSON decode transactions response body: %s, %s", err, b)
	}

	return transactions, nil
}

// Shutdown TODO
func (m *Monzo) Shutdown() {
}
