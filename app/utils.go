package main

import "fmt"

func printToken(token Token) {
	fmt.Println("\nTOKEN TYPE: ", token.Token_type)
	fmt.Println("EXPIRES IN: ", token.Expires_in)
	fmt.Println("WHEN EXPIRES: ", token.When_Expires)

	fmt.Println("\nACESS TOKEN\n", "------------------------------\n", token.Access_token)
	fmt.Println("\nID TOKEN\n", "------------------------------\n", token.Id_token)

	fmt.Println("\nREFRESH TOKEN\n", "------------------------------\n", token.Refresh_token)
	fmt.Println("\nSCOPE\n", "------------------------------\n", token.Scope)
}
