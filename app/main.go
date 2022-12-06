package main

import (
	"fmt"
	"time"
)

func main() {

	begin := time.Now()

	fmt.Println("*****************************************************")

	CashAccount, error := response_cash_account()

	if error != nil {
		panic(error)
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
