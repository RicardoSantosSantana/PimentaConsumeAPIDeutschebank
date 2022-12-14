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
	OriginIban                           string          `json:"originIban"`
	Amount                               float32         `json:"amount"`
	CounterPartyName                     string          `json:"counterPartyName"`
	CounterPartyIban                     string          `json:"counterPartyIban,omitempty"`
	PaymentReference                     string          `json:"paymentReference"`
	BookingDate                          string          `json:"bookingDate"`
	CurrencyCode                         string          `json:"currencyCode"`
	TransactionCode                      string          `json:"transactionCode"`
	ExternalBankTransactionDomainCode    string          `json:"externalBankTransactionDomainCode"`
	ExternalBankTransactionFamilyCode    string          `json:"externalBankTransactionFamilyCode"`
	ExternalBankTransactionSubFamilyCode string          `json:"externalBankTransactionSubFamilyCode"`
	MandateReference                     string          `json:"mandateReference"`
	CreditorId                           string          `json:"creditorId"`
	E2eReference                         string          `json:"e2eReference"`
	PaymentIdentification                string          `json:"paymentIdentification"`
	ValueDate                            string          `json:"valueDate"`
	Id                                   string          `json:"id"`
	Csv                                  json.RawMessage `json:"csv,omitempty"`
}

type ApiBankingTransactions struct {
	TotalItems   int            `json:"totalItems"`
	Offset       int            `json:"offset"`
	Limit        int            `json:"limit"`
	Transactions []Transactions `json:"transactions"`
}

func (bankTransaction ApiBankingTransactions) offset() int {
	return bankTransaction.TotalItems / bankTransaction.Limit
}

func update_field_csv(data string) (string, error) {

	//fill object trasaction with data parameters
	transaction, err := Json_decode(data)

	if err != nil {
		return "", err
	}

	// convert the value to string ( json format )
	transaction_json, err := Json_encode(transaction)
	if err != nil {
		return "", err
	}

	// insert the value string ( json format ) on object transaction
	transaction.Csv = json.RawMessage(transaction_json)

	// convert transaction object into string (json format) with new csv field filled in
	transaction_texto, err := Json_encode(transaction)
	if err != nil {
		return "", err
	}

	return transaction_texto, nil
}

func get_list_transactions(account Account) error {

	bankTransactions, error := get_transactions(account, 0)
	if error != nil {
		fmt.Println("não há transacao")
		panic(error)
	}

	prepare_to_save_transactions(bankTransactions, account)

	offset := bankTransactions.offset()

	if offset > 0 {
		for i := 1; i <= offset; i++ {
			bank_transactions, _ := get_transactions(account, i)
			prepare_to_save_transactions(bank_transactions, account)

		}
	}

	return nil
}

func get_transactions(account Account, offset int) (ApiBankingTransactions, error) {

	settings := Settings()

	params := url.Values{
		"iban":   {account.Iban},
		"limit":  {settings.Api.Limit},
		"offset": {strconv.Itoa(offset)},
	}

	reqUrl := settings.Api.GetTransaction.Uri + "?" + params.Encode()

	transactions := ApiBankingTransactions{}

	req, err := http.NewRequest(settings.Api.GetTransaction.Method, reqUrl, bytes.NewReader([]byte("")))
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

func prepare_to_save_transactions(bankTransactions ApiBankingTransactions, account Account) error {

	// primeira regra se a conta existe
	bankAccount, err := GetBankAccount(account.Iban)
	if err != nil {
		fmt.Println("Bank account not found for: " + account.Iban)
		return err
	}

	transaction := bankTransactions.Transactions

	totalTransactions := len(transaction)

	for i := 0; i < totalTransactions; i++ {

		fmt.Println(i, " de ", totalTransactions)

		external_id := transaction[i].Id

		is_exists_transaction_with_external_id, _ := bank_transaction_search_by_external_id(external_id, bankAccount)

		if is_exists_transaction_with_external_id {
			fmt.Println("Transaction already exists: " + external_id)
			continue
		}

		bank_transaction, is_transaction_exists, _ := bank_transaction_search_by_more_fields(transaction[i], bankAccount)

		if is_transaction_exists {
			bank_transaction.External_id = transaction[i].Id
			bank_transaction.Source = "csv-dbapi"
			transaction_to_json, _ := update_field_csv(string(bank_transaction.Json_data))
			bank_transaction.Json_data = transaction_to_json

			fmt.Println(" ***** SALVAR A TRANSAÇÃO EXISTENTE")
			save_transaction(bank_transaction)
			fmt.Println("Transaction already exists, matched by description, date and bank account: ", transaction[i].ValueDate, " - ", transaction[i].Id)
			continue
		}

		json_data, err_json_encode := Json_encode(transaction[i])

		if err_json_encode != nil {
			return err_json_encode
		}

		new_finance_bank_transation := FINANCE_BankTransactions{
			Json_data:       json_data,
			Date:            transaction[i].ValueDate,
			Details:         transaction[i].PaymentReference,
			Source:          "dbapi",
			External_id:     transaction[i].Id,
			Amount:          transaction[i].Amount,
			Date_reserved:   transaction[i].BookingDate,
			Operator:        transaction[i].CounterPartyName,
			Bank_account_id: bankAccount.Id,
			Category:        "",
			Import_code:     "",
		}

		fmt.Println(" ***** SALVAR NOVA TRANSAÇÃO: ", new_finance_bank_transation.Bank_account_id.Int32)
		save_transaction(new_finance_bank_transation)
		fmt.Println("Transaction created: " + transaction[i].Id)

	}

	return nil
}

func bank_transaction_search_by_more_fields(transaction Transactions, bankAccount FINANCE_BankAccount) (FINANCE_BankTransactions, bool, error) {

	transactions, isTransactionExists, err := GetBankTransaction(transaction, bankAccount)

	if err != nil {
		return FINANCE_BankTransactions{}, false, err
	}

	return transactions, isTransactionExists, nil
}

func bank_transaction_search_by_external_id(external_id string, bankAccount FINANCE_BankAccount) (bool, error) {

	transaction := Transactions{
		PaymentReference: external_id,
	}

	_, isTransactionExists, err := GetBankTransaction(transaction, bankAccount)

	if err != nil {
		return false, nil
	}

	return isTransactionExists, nil
}

func save_transaction(transaction FINANCE_BankTransactions) error {

	fmt.Println(" ")

	db, errConn := openConnection()
	if errConn != nil {
		return errConn
	}

	// //(`id, `deleted_at`, `amount`, `operator`, `date`, `date_reserved`, `details`, `category`, `json_data`, `bank_account_id`, `account_id`, `created_at`, `updated_at`, `import_code`, `category_id`, `private`, `relation_status`, `problem`, `payment_id`, `expense_id`, `chargeback_id`, `exchange_rate`, `unique_id`, `invoice_id`, `tax_id`, `invoices`, `category_status`, `currency_id`, `external_id`, `source`) VALUES ('[value-1]','[value-2]','[value-3]','[value-4]','[value-5]','[value-6]','[value-7]','[value-8]','[value-9]','[value-10]','[value-11]','[value-12]','[value-13]','[value-14]','[value-15]','[value-16]','[value-17]','[value-18]','[value-19]','[value-20]','[value-21]','[value-22]','[value-23]','[value-24]','[value-25]','[value-26]','[value-27]','[value-28]','[value-29]','[value-30]','[value-31]')

	sql := `INSERT INTO bank_transactions (	json_data, date, details, source, external_id, amount, date_reserved, operator, bank_account_id, category, import_code, private ) `

	s := []string{
		"'" + transaction.Json_data + "'",
		strconv.Quote(dateFormat(transaction.Date)),
		strconv.Quote(transaction.Details),
		strconv.Quote(transaction.Source),
		strconv.Quote(transaction.External_id),
		fmt.Sprintf("%g", transaction.Amount),
		strconv.Quote(dateFormat(transaction.Date_reserved)),
		strconv.Quote(transaction.Operator),
		strconv.Quote(strconv.Itoa(int(transaction.Bank_account_id.Int32))),
		strconv.Quote(transaction.Category),
		strconv.Quote(transaction.Import_code),
		strconv.Quote(strconv.Itoa(int(transaction.Private.Int32))),
	}

	values := strings.Join(s, ",")
	//	values := make_values_transactions(transactions)
	fmt.Println(sql + " VALUES (" + values + ")")

	fmt.Println(" ")

	_, errExec := db.Exec(sql + " VALUES (" + values + ")")

	if errExec != nil {
		fmt.Println(errExec)
		return errExec
	}

	defer db.Close()

	return nil
}

func make_values_transactions(transactions FINANCE_BankTransactions) string {
	return ""
	// var v string
	// for i := 0; i <= len(transactions.Transactions)-1; i++ {
	// 	transaction := transactions.Transactions[i]

	// 	s := []string{
	// 		" (" + strconv.Quote(transaction.OriginIban),
	// 		fmt.Sprintf("%g", transaction.Amount),
	// 		strconv.Quote(transaction.CounterPartyName),
	// 		strconv.Quote(transaction.CounterPartyIban),
	// 		strconv.Quote(transaction.PaymentReference),
	// 		strconv.Quote(transaction.BookingDate),
	// 		strconv.Quote(transaction.CurrencyCode),
	// 		strconv.Quote(transaction.TransactionCode),
	// 		strconv.Quote(transaction.ExternalBankTransactionDomainCode),
	// 		strconv.Quote(transaction.ExternalBankTransactionFamilyCode),
	// 		strconv.Quote(transaction.ExternalBankTransactionSubFamilyCode),
	// 		strconv.Quote(transaction.MandateReference),
	// 		strconv.Quote(transaction.CreditorId),
	// 		strconv.Quote(transaction.E2eReference),
	// 		strconv.Quote(transaction.PaymentIdentification),
	// 		strconv.Quote(transaction.ValueDate),
	// 		strconv.Quote(transaction.Id) + "),",
	// 	}

	// 	v += strings.Join(s, ",")

	// }
	// str := v[:len(v)-1]

	// return " values " + str

}
