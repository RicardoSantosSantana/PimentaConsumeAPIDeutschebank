package main

import (
	"fmt"
)

func main() {

	token, err := GetToken()

	if err != nil {
		fmt.Println(err)
	}
	printToken(token)
	fmt.Println("**********************************")
	fmt.Println(token.Access_token)

}
