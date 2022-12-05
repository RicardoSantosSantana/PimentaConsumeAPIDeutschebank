package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Transactions struct {
	OriginIban                           string  `json:"originIban"`
	Amount                               float32 `json:"amount"`
	CounterPartyName                     string  `json:"counterPartyName"`
	PaymentReference                     string  `json:"paymentReference"`
	BookingDate                          string  `json:"bookingDate"`
	CurrencyCode                         string  `json:"currencyCode"`
	TransactionCode                      string  `json:"transactionCode"`
	ExternalBankTransactionDomainCode    string  `json:"externalBankTransactionDomainCode"`
	ExternalBankTransactionFamilyCode    string  `json:"externalBankTransactionFamilyCode"`
	ExternalBankTransactionSubFamilyCode string  `json:"externalBankTransactionSubFamilyCode"`
	MandateReference                     string  `json:"mandateReference"`
	CreditorId                           string  `json:"creditorId"`
	E2eReference                         string  `json:"e2eReference"`
	PaymentIdentification                string  `json:"paymentIdentification"`
	ValueDate                            string  `json:"valueDate"`
	Id                                   string  `json:"id"`
}

type BankingTransactions struct {
	TotalItems   int            `json:"totalItems"`
	Offset       int            `json:"offset"`
	Limit        int            `json:"limit"`
	Transactions []Transactions `json:"transactions"`
}

//https://simulator-api.db.com:443/gw/dbapi/banking/transactions/v2?iban=DE00500700100200000867&bookingDateFrom=2022-05-14&bookingDateTo=2022-05-14

func (cashAccount CashAccount) offset() int {
	fmt.Println("cash Account: ", cashAccount.TotalItems, "Limit: ", cashAccount.Limit)
	return cashAccount.TotalItems / cashAccount.Limit
}

func get_transactions(account Account, endPoint ApiEndPoint) (BankingTransactions, error) {

	transactions := BankingTransactions{}

	req, err := http.NewRequest(endPoint.Method, endPoint.Uri, bytes.NewReader([]byte("")))
	if err != nil {
		panic(err)
	}

	token, errorToken := GetToken()

	if errorToken != nil {
		panic(errorToken)
	}

	req.Header.Add("Authorization", "Bearer "+token.Access_token)

	req.Close = true

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return transactions, err
	}

	if strings.Contains(string(body), "error") {
		return transactions, errors.New(string(body))
	}

	if err := json.Unmarshal(body, &transactions); err != nil {
		return transactions, errors.New(string(body))
	}

	return transactions, nil
}
