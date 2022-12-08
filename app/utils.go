package main

import (
	"fmt"
	"os"
	"os/exec"
)

func printToken(token Token) {
	fmt.Println("\nTOKEN TYPE: ", token.Token_type)
	fmt.Println("EXPIRES IN: ", token.Expires_in)
	fmt.Println("WHEN EXPIRES: ", token.When_Expires)

	fmt.Println("\nACESS TOKEN\n", "------------------------------\n", token.Access_token)
	fmt.Println("\nID TOKEN\n", "------------------------------\n", token.Id_token)

	fmt.Println("\nREFRESH TOKEN\n", "------------------------------\n", token.Refresh_token)
	fmt.Println("\nSCOPE\n", "------------------------------\n", token.Scope)
}

func printConfig() {
	// cfg, error := OpenConfig()
	// if error != nil {
	// 	panic(error)
	// }
	// fmt.Println("DATABASE: ")
	// fmt.Println("  hostname: ", cfg.Database.Hostname)
	// fmt.Println("  name: ", cfg.Database.Name)
	// fmt.Println("  password: ", cfg.Database.Password)
	// fmt.Println("  port: ", cfg.Database.Port)
	// fmt.Println("  username: ", cfg.Database.Username)

	// fmt.Println("API: ")
	// fmt.Println("  limit:", cfg.Api.Limit)
	// fmt.Println("  Auth:")
	// fmt.Println("    client_id: ", cfg.Api.Auth.Client_id)
	// fmt.Println("    client_secret: ", cfg.Api.Auth.Client_secret)
	// fmt.Println("    code: ", cfg.Api.Auth.Code)
	// fmt.Println("    redirect_url", cfg.Api.Auth.Redirect_uri)
	// fmt.Println("  Get_token: ")
	// fmt.Println("    method: ", cfg.Api.Get_token.Method)
	// fmt.Println("    uri: ", cfg.Api.Get_token.Uri)
	// fmt.Println("  Get_transaction: ")
	// fmt.Println("    method: ", cfg.Api.Get_transaction.Method)
	// fmt.Println("    uri: ", cfg.Api.Get_transaction.Uri)
	// fmt.Println("  Get_accounts: ")
	// fmt.Println("    method: ", cfg.Api.Get_accounts.Method)
	// fmt.Println("    uri: ", cfg.Api.Get_accounts.Uri)
}
func clearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
