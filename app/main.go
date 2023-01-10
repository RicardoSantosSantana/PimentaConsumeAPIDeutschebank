package main

import (
	"fmt"
	"time"
	"os"
)

//   - 1 list_accounts[
// 		get_list_transactions [
// 			get_transactions[
// 				save_transactions
// 			]
// 		]
//	]

func main() {
	 
	//https://04c1ffa5ffbf4eebb3f392b6069d38aa@o4504417796882432.ingest.sentry.io/4504418327986176

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
