package main

import (
	"fmt"
	"os"
	"time"
)

//   - 1 list_accounts[
// 		get_list_transactions [
// 			get_transactions[
// 				save_transactions
// 			]
// 		]
//	]

func main() {

	//list_accounts()
	// var isTransactionExists bool
	// iban := "DE00500700100200000885"
	// bankAccount, err := GetBankAccount(iban)
	// if err != nil {
	// 	fmt.Println("Bank account not found for: " + iban)
	// }

	// fmt.Println(bankAccount)
	//sql.NullInt32{Int32: 10},

	// bankTransaction := BankTransactions{
	// 	Date:            "2022-01-01",
	// 	Details:         "Detalhes",
	// 	Account_id:      sql.NullInt32{Int32: 10},
	// 	Bank_account_id: sql.NullInt32{Int32: 10},
	// }

	// bankTransaction := BankTransactions{
	// 	External_id: "1FUO_RjfGt_d4HfY3gO5BcRhIfR8caz-rEw7h9nPBbjyV3WI7N--",
	// }
	// transaction := Transactions{}

	// bankTransactionFound, isTransactionExists, err := GetBankTransaction(bankTransaction, transaction, bankAccount)

	// if err != nil {
	// 	fmt.Println(err)
	// 	//continue
	// }
	// if isTransactionExists {
	// 	fmt.Println("Transaction already exists: " + transaction.Id)
	// }

	// if isTransactionExists {
	// 	bankTransactionFound.External_id = "external id de teste"
	// 	bankTransactionFound.Source = "csv-dbapi"

	// 	fmt.Println(bankTransactionFound.Bank_account_id.Int32)
	// 	fmt.Println(bankTransactionFound.Source)
	// 	fmt.Println(bankTransactionFound.External_id)
	// }
	// fmt.Println("ola")

	// var listTransaction []Transactions

	// a := Transactions{
	// 	OriginIban: "Original IBAN",
	// 	Amount:     10,
	// }

	// listTransaction = append(listTransaction, a)
	// listTransaction = append(listTransaction, Transactions{
	// 	OriginIban: "Original IBAN 2",
	// 	Amount:     210,
	// })

	// fmt.Println(listTransaction)
	list_accounts()

}

func list_accounts() {

	clearConsole()
	begin := time.Now()

	CashAccount, error := response_cash_account()

	if error != nil {
		fmt.Println("\n  (!) ", error)
		os.Exit(0)
	}

	for i := 0; i < len(CashAccount.Accounts); i++ {

		fmt.Println("---------------------------------------------------")
		fmt.Println("   IBAN: ", CashAccount.Accounts[i].Iban)
		fmt.Println(" ")
		error := get_list_transactions(CashAccount.Accounts[i])

		if error != nil {
			panic(error)
		}

		fmt.Println(" ")

	}

	fmt.Println("Start : ", begin)
	fmt.Println("Finish: ", time.Now())
}
