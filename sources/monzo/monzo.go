package monzo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/jemgunay/life-metrics/config"
	"github.com/jemgunay/life-metrics/sources"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// Monzo represents the Monzo collection source.
type Monzo struct {
	currentAuth       authAccessDetails
	clientSecret      string
	authRefreshedChan chan authAccessDetails
	collectionChan    chan sources.Period
}

// New initialises the Monzo source and manages auth token refreshing.
func New(conf config.Monzo, exporter sources.Exporter) *Monzo {
	m := &Monzo{
		authRefreshedChan: make(chan authAccessDetails, 1),
		collectionChan:    make(chan sources.Period),
		currentAuth: authAccessDetails{
			ClientID: conf.ClientID,
		},
		clientSecret: conf.ClientSecret,
	}

	// start polling for oauth initial, oauth refresh and collection requests
	refreshTimer := time.NewTimer(time.Hour)
	go func() {
		for {
			select {
			case <-refreshTimer.C:
				log.Print("starting Monzo authentication refresh")
				m.fetchAccessToken(m.currentAuth.RefreshToken, accessCodeRefresh)

			case m.currentAuth = <-m.authRefreshedChan:
				// reset auth refresh timer
				refreshTimer.Stop()
				timeToRefresh := (time.Second * time.Duration(m.currentAuth.ExpiresIn)) - time.Minute*5
				refreshTimer = time.NewTimer(timeToRefresh)
				log.Printf("Monzo authenticated - next authentication in %s", timeToRefresh)

			case period := <-m.collectionChan:
				results, err := m.performCollection(period)
				if err != nil {
					log.Printf("failed to perform collection for Monzo: %s", err)
					continue
				}

				// write collected source data to influx
				if err := exporter.Write(m.Name(), results...); err != nil {
					log.Printf("writing data to influx failed for Monzo: %s", err)
				}
			}
		}
	}()
	return m
}

// Name returns the source name.
func (m *Monzo) Name() string {
	return "monzo"
}

// Collect enqueues a Monzo collection request.
func (m *Monzo) Collect(period sources.Period) {
	select {
	case m.collectionChan <- period:
	default:
		log.Print("collection failed for Monzo as collection queue is full")
	}
}

func (m *Monzo) performCollection(period sources.Period) ([]sources.Result, error) {
	if m.currentAuth.AccessToken == "" {
		return nil, errors.New("access token not set - oauth setup required")
	}

	account, err := m.getAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %s", err)
	}

	// get transactions for account
	transactions, err := m.getTransactions(account.ID, period.Start, period.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions list: %s", err)
	}

	var results []sources.Result
	for _, transaction := range transactions.Transactions {
		// process eating out/take away data
		if transaction.Category == "eating_out" {
			createdTime, err := time.Parse(time.RFC3339, transaction.CreatedTime)
			if err != nil {
				log.Printf("failed to parse Monzo created time for %s: %s", transaction.CreatedTime, err)
				continue
			}

			// convert price from pence (1234) to pounds (12.34)
			amountInPounds := math.Abs(float64(transaction.Amount)) / 100

			results = append(results, sources.Result{
				Time: createdTime,
				Tags: map[string]string{
					"category":        transaction.Category,
					"restaurant_name": transaction.Merchant.Name,
					"restaurant_city": transaction.Merchant.Address.City,
				},
				Fields: map[string]interface{}{
					"price":                amountInPounds,
					"currency":             transaction.Currency,
					"restaurant_latitude":  transaction.Merchant.Address.Latitude,
					"restaurant_longitude": transaction.Merchant.Address.Longitude,
				},
			})
		}
	}

	return results, nil
}

type accountsResult struct {
	Accounts []account `json:"accounts"`
}

type account struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (m *Monzo) getAccount() (account, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.monzo.com/accounts", nil)
	if err != nil {
		return account{}, fmt.Errorf("failed to create accounts request: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+m.currentAuth.AccessToken)

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
	AccountBalance int      `json:"account_balance"`
	Amount         int      `json:"amount"`
	CreatedTime    string   `json:"created"`
	Currency       string   `json:"currency"`
	Description    string   `json:"description"`
	ID             string   `json:"id"`
	Merchant       merchant `json:"merchant"`
	Notes          string   `json:"notes"`
	IsLoad         bool     `json:"is_load"`
	Settled        string   `json:"settled"`
	Category       string   `json:"category"`
}

type merchant struct {
	Address struct {
		Address   string  `json:"address"`
		City      string  `json:"city"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Postcode  string  `json:"postcode"`
		Region    string  `json:"region"`
	} `json:"address"`
	Created  string `json:"created"`
	GroupID  string `json:"group_id"`
	ID       string `json:"id"`
	Logo     string `json:"logo"`
	Emoji    string `json:"emoji"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

func (m *Monzo) getTransactions(accountID string, start, end time.Time) (transactionsResult, error) {
	var transactions transactionsResult

	q := url.Values{}
	q.Set("account_id", accountID)
	q.Set("since", start.Format(time.RFC3339))
	q.Set("before", end.Format(time.RFC3339))
	// enrich transaction with merchant data
	q.Set("expand[]", "merchant")
	fmt.Println(q.Encode())

	req, err := http.NewRequest(http.MethodGet, "https://api.monzo.com/transactions?"+q.Encode(), nil)
	if err != nil {
		return transactions, fmt.Errorf("failed to create accounts request: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+m.currentAuth.AccessToken)

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

// Shutdown does nothing for this module.
func (m *Monzo) Shutdown() {
}
