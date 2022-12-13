package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

type CSV struct {
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
type FINANCE_BankAccount struct {
	Account_id  int           `json:"account_id"`
	App_version int           `json:"app_version"`
	Bank_id     string        `json:"bank_id"`
	Id          sql.NullInt32 `json:"id"`
	Is_paypal   int           `json:"is_paypal"`
	Name        string        `json:"name"`
	Number      string        `json:"number"`
	Ofx_version int           `json:"ofx_version"`
	Public_id   int           `json:"public_id"`
	Username    string        `json:"username"`
	User_id     int           `json:"user_id"`
}

type FINANCE_BankTransactions struct {
	Account_id      sql.NullInt32   `json:"account_id"`
	Amount          float32         `json:"amount"`
	Bank_account_id sql.NullInt32   `json:"bank_account_id"`
	Category        string          `json:"category"`
	Category_id     sql.NullInt32   `json:"category_id"`
	Category_status string          `json:"category_status"`
	Chargeback_id   sql.NullInt32   `json:"chargeback_id"`
	Currency_id     sql.NullInt32   `json:"currency_id"`
	Date            string          `json:"date"`
	Date_reserved   string          `json:"date_reserved"`
	Details         string          `json:"details"`
	Exchange_rate   sql.NullFloat64 `json:"exchange_rate"`
	Expense_id      sql.NullInt32   `json:"expense_id"`
	External_id     string          `json:"external_id"`
	Id              int             `json:"id"`
	Import_code     string          `json:"import_code"`
	Invoices        sql.NullString  `json:"invoices"`
	Invoice_id      sql.NullInt32   `json:"invoice_id"`
	Json_data       string          `json:"json_data"`
	Operator        string          `json:"operator"`
	Payment_id      sql.NullInt32   `json:"payment_id"`
	Private         sql.NullInt32   `json:"private"`
	Problem         string          `json:"problem"`
	Relation_status string          `json:"relation_status"`
	Source          string          `json:"source"`
	Tax_id          sql.NullInt32   `json:"tax_id"`
	Unique_id       sql.NullString  `json:"unique_id"`
	User_id         sql.NullInt32   `json:"user_id"`
}
type IsExistsTransaction struct {
	Transaction_exists bool `json:"transaction_exists"`
}

func GetBankAccount(iban string) (FINANCE_BankAccount, error) {

	db, err := openConnection()
	if err != nil {
		return FINANCE_BankAccount{}, err
	}

	bankAccount := FINANCE_BankAccount{}

	sql := `SELECT 
				account_id, 
				app_version, 
				bank_id,
				id, 
				is_paypal,
				name, 
				number, 
				ofx_version, 
				public_id, 
				Username, 
				user_id 
			FROM bank_accounts where number='` + iban + `' LIMIT 1`

	err = db.QueryRow(sql).Scan(
		&bankAccount.Account_id,
		&bankAccount.App_version,
		&bankAccount.Bank_id,
		&bankAccount.Id,
		&bankAccount.Is_paypal,
		&bankAccount.Name,
		&bankAccount.Number,
		&bankAccount.Ofx_version,
		&bankAccount.Public_id,
		&bankAccount.Username,
		&bankAccount.User_id)

	if err != nil {
		return FINANCE_BankAccount{}, err
	}

	defer db.Close()

	return bankAccount, nil

}

func filter_transaction(transaction Transactions, bankAccount FINANCE_BankAccount) string {

	if transaction.Id != "" {
		return " WHERE external_id='" + transaction.Id + "' LIMIT 1 "
	}

	//valor do account_id definido no código do finance fixo em 1
	account_id := "1"

	sql := " WHERE details='" + transaction.PaymentReference + "' AND "
	sql += " date='" + transaction.ValueDate + "' AND "
	sql += " account_id=" + account_id + " AND "
	sql += " bank_account_id=" + strconv.Itoa(int(bankAccount.Id.Int32))
	sql += " LIMIT 1 "

	return sql

}

func GetBankTransaction(transaction Transactions, bankAccount FINANCE_BankAccount) (FINANCE_BankTransactions, bool, error) {

	//external_id é o mesmo que o transaction_id
	var sql string

	db, err := openConnection()
	if err != nil {
		fmt.Println("erro banco:", err)
		return FINANCE_BankTransactions{}, false, err
	}

	fields := `SELECT 
			account_id, 
			amount, 
			bank_account_id, 
			category, 
			category_id, 
			category_status, 
			chargeback_id, 
			currency_id, 
			date, 
			date_reserved, 
			details, 
			exchange_rate, 
			expense_id, 
			external_id, 
			id, 
			import_code, 
			invoices, 
			invoice_id, 
			json_data, 
			operator, 
			payment_id, 
			private, 
			problem, 
			relation_status, 
			source, 
			tax_id, 
			unique_id, 
			user_id 
		FROM bank_transactions `

	sql = fields + filter_transaction(transaction, bankAccount)

	bankTransaction := FINANCE_BankTransactions{}

	err = db.QueryRow(sql).Scan(
		&bankTransaction.Account_id,
		&bankTransaction.Amount,
		&bankTransaction.Bank_account_id,
		&bankTransaction.Category,
		&bankTransaction.Category_id,
		&bankTransaction.Category_status,
		&bankTransaction.Chargeback_id,
		&bankTransaction.Currency_id,
		&bankTransaction.Date,
		&bankTransaction.Date_reserved,
		&bankTransaction.Details,
		&bankTransaction.Exchange_rate,
		&bankTransaction.Expense_id,
		&bankTransaction.External_id,
		&bankTransaction.Id,
		&bankTransaction.Import_code,
		&bankTransaction.Invoices,
		&bankTransaction.Invoice_id,
		&bankTransaction.Json_data,
		&bankTransaction.Operator,
		&bankTransaction.Payment_id,
		&bankTransaction.Private,
		&bankTransaction.Problem,
		&bankTransaction.Relation_status,
		&bankTransaction.Source,
		&bankTransaction.Tax_id,
		&bankTransaction.Unique_id,
		&bankTransaction.User_id)

	if err != nil {
		return FINANCE_BankTransactions{}, false, err
	}

	db.Close()

	return bankTransaction, true, nil

}
