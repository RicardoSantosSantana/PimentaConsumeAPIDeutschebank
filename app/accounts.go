package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Account struct {
	Iban               string  `json:"iban"`
	CurrencyCode       string  `json:"currencyCode"`
	Bic                string  `json:"bic"`
	AccountType        string  `json:"accountType"`
	CurrentBalance     float32 `json:"currentBalance"`
	ProductDescription string  `json:"productDescription"`
}
type CashAccount struct {
	TotalItems int       `json:"totalItems"`
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
	Accounts   []Account `json:"accounts"`
}

func response_cash_account() (CashAccount, error) {

	cashAccount := CashAccount{}
	settings := Settings()

	req, err := http.NewRequest(settings.Api.GetAccounts.Method, settings.Api.GetAccounts.Uri, bytes.NewReader([]byte("")))
	if err != nil {
		return cashAccount, err
	}
	token, errorToken := GetToken()

	if errorToken != nil {
		return cashAccount, errorToken
	}

	req.Header.Add("Authorization", "Bearer "+token.Access_token)

	req.Close = true

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return cashAccount, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return cashAccount, err
	}

	if strings.Contains(string(body), "error") {
		return cashAccount, errors.New(string(body))
	}

	if err := json.Unmarshal(body, &cashAccount); err != nil {
		return cashAccount, errors.New(string(body))
	}

	return cashAccount, nil
}
