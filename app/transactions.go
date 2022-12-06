package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Transactions struct {
	OriginIban                           string  `json:"originIban"`
	Amount                               float32 `json:"amount"`
	CounterPartyName                     string  `json:"counterPartyName"`
	CounterPartyIban                     string  `json:"counterPartyIban"`
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

func (bankTransaction BankingTransactions) offset() int {
	return bankTransaction.TotalItems / bankTransaction.Limit
}

func get_list_transactions(account Account) error {

	transactions, error := get_transactions(account, 0)
	if error != nil {
		panic(error)
	}

	error_save_offset := save_transactions(transactions)

	if error_save_offset != nil {
		panic(error_save_offset)
		//return error_transaction_loop
	}

	offset := transactions.offset()

	if offset > 0 {
		for i := 1; i <= offset; i++ {
			banking_transactions, error_transaction_loop := get_transactions(account, i)
			error_save := save_transactions(banking_transactions)

			if error_transaction_loop != nil {
				panic(error_transaction_loop)
				//return error_transaction_loop
			}

			if error_save != nil {
				panic(error_save)
				//return error_save
			}
		}
	}

	return nil
}

func get_transactions(account Account, offset int) (BankingTransactions, error) {

	settings := Settings()

	params := url.Values{
		"iban":   {account.Iban},
		"limit":  {settings.LimitToGetTransactions},
		"offset": {strconv.Itoa(offset)},
	}

	reqUrl := settings.GetTransaction.Uri + "?" + params.Encode()

	transactions := BankingTransactions{}

	req, err := http.NewRequest(settings.GetTransaction.Method, reqUrl, bytes.NewReader([]byte("")))
	if err != nil {
		return transactions, err
	}

	token, errorToken := GetToken()

	if errorToken != nil {
		return transactions, errorToken
	}

	req.Header.Add("Authorization", "Bearer "+token.Access_token)

	req.Close = true

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return transactions, err
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

	fmt.Println("   Total Itens: ", transactions.TotalItems)
	fmt.Println("   URI: ", reqUrl)
	fmt.Println(" ")

	return transactions, nil
}

func save_transactions(transactions BankingTransactions) error {

	fmt.Println(" ")

	db, errConn := openConnection()
	if errConn != nil {
		return errConn
	}
	sql := `insert into transactions (
		originiban, 
		amount, 
		counterpartyname, 
		counterpartyiban, 
		paymentreference, 
		bookingdate, 
		currencycode, 
		transactioncode, 
		externalbanktransactiondomaincode, 
		externalbanktransactionfamilycode, 
		externalbanktransactionsubfamilycode, 
		mandatereference, 
		creditorid, 
		e2ereference, 
		paymentidentification, 
		valuedate, 
		id)`

	values := make_values_transactions(transactions)

	_, errExec := db.Exec(sql + values)

	if errExec != nil {
		return errExec
	}

	defer db.Close()

	return nil
}

func make_values_transactions(transactions BankingTransactions) string {

	var v string
	for i := 0; i <= len(transactions.Transactions)-1; i++ {
		transaction := transactions.Transactions[i]

		s := []string{
			" (" + strconv.Quote(transaction.OriginIban),
			fmt.Sprintf("%g", transaction.Amount),
			strconv.Quote(transaction.CounterPartyName),
			strconv.Quote(transaction.CounterPartyIban),
			strconv.Quote(transaction.PaymentReference),
			strconv.Quote(transaction.BookingDate),
			strconv.Quote(transaction.CurrencyCode),
			strconv.Quote(transaction.TransactionCode),
			strconv.Quote(transaction.ExternalBankTransactionDomainCode),
			strconv.Quote(transaction.ExternalBankTransactionFamilyCode),
			strconv.Quote(transaction.ExternalBankTransactionSubFamilyCode),
			strconv.Quote(transaction.MandateReference),
			strconv.Quote(transaction.CreditorId),
			strconv.Quote(transaction.E2eReference),
			strconv.Quote(transaction.PaymentIdentification),
			strconv.Quote(transaction.ValueDate),
			strconv.Quote(transaction.Id) + "),",
		}

		v += strings.Join(s, ",")

	}
	str := v[:len(v)-1]

	return " values " + str

}
