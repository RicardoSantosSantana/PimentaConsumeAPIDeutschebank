package main

import (
	"fmt"
	"net/url"
)

func main() {

	// token, err := GetToken()

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// printToken(token)
	// fmt.Println("______________________________________________")
	// fmt.Println(token.Access_token)

	CashAccount, error := response_cash_account()

	if error != nil {
		panic(error)
	}
	settings := Settings()

	for i := 0; i < len(CashAccount.Accounts); i++ {

		params := url.Values{
			"iban":   {CashAccount.Accounts[i].Iban},
			"limit":  {"200"},
			"offset": {"0"},
		}

		reqUrl := settings.GetTransaction.Uri + "?" + params.Encode()
		fmt.Println(reqUrl)

		transactions, error := get_transactions(CashAccount.Accounts[i], ApiEndPoint{
			Method: "GET",
			Uri:    reqUrl,
		})

		if error != nil {
			panic(error)
		}
		fmt.Println(transactions)
	}

}
