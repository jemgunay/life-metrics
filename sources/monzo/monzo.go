package monzo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jemgunay/life-metrics/sources"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// Monzo TODO
type Monzo struct {
	accessToken  string
	refreshToken string
	expiresIn    int64
}

// New TODO
func New() *Monzo {
	return &Monzo{}
}

// Name TODO
func (m *Monzo) Name() string {
	return "monzo"
}

type account struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type accountsResult struct {
	Accounts []account `json:"accounts"`
}

// Collect TODO
func (m *Monzo) Collect(start, end time.Time) ([]sources.Result, error) {
	if m.accessToken == "" {
		return nil, errors.New("access token not set - oauth setup required")
	}

	account, err := m.getAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %s", err)
	}

	// get transactions for account
	transactions, err := m.getTransactions(account.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %s", err)
	}

	//log.Printf("%+v", transactions)
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
	req.Header.Add("Authorization", "Bearer "+m.accessToken)

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

type transaction struct {
	AccountBalance int       `json:"account_balance"`
	Amount         int       `json:"amount"`
	Created        time.Time `json:"created"`
	Currency       string    `json:"currency"`
	Description    string    `json:"description"`
	ID             string    `json:"id"`
	Merchant       string    `json:"merchant"`
	Notes    string    `json:"notes"`
	IsLoad   bool      `json:"is_load"`
	Settled  time.Time `json:"settled"`
	Category string    `json:"category"`
}

type transactionsResult struct {
	Transactions []transaction `json:"transactions"`
}

func (m *Monzo) getTransactions(accountID string) (transactionsResult, error) {
	var transactions transactionsResult

	req, err := http.NewRequest(http.MethodGet, "https://api.monzo.com/transactions?account_id="+accountID, nil)
	if err != nil {
		return transactions, fmt.Errorf("failed to create accounts request: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+m.accessToken)

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
