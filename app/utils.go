package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func Json_decode(data string) (Transactions, error) {
	var dat Transactions
	err := json.Unmarshal([]byte(data), &dat)

	return dat, err
}

func Json_encode(data Transactions) (string, error) {
	jsons, err := json.Marshal(data)
	return string(jsons), err
}

func clearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func log(dados any) {
	jsons, _ := json.Marshal(dados)
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println(string(jsons))
	fmt.Println("-----------------------------------------------------------------")

}
