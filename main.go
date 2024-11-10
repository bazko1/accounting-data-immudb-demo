package main

import (
	"fmt"
	"os"

	"accounting-immudb-demo/pkg/account"
)

func main() {
	token := os.Getenv("API_PRIVATE_KEY")
	manger := account.NewAccountManager("default", "default", token)
	err := manger.CreateEntry(account.Account{
		Number:  1,
		Name:    "Foo Bar",
		Iban:    "FOO12",
		Address: "Foo street 10",
		Amount:  100,
		Type:    account.TypeSending,
	})
	fmt.Println(err)
}
