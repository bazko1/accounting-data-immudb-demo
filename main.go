package main

import (
	"errors"
	"fmt"
	"os"

	"accounting-immudb-demo/pkg/account"
	"accounting-immudb-demo/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	token := os.Getenv("API_PRIVATE_KEY")

	// there is some sort of protection that does not allows me to create entries in collection
	// that is not named default
	manager := account.NewAccountManager("default", "default", token)

	err := manager.CreateAccountCollection()
	if err != nil {
		logger.Info("main create", zap.String("error", err.Error()))
		return
	}

	err = manager.CreateAccount(account.Account{
		Number:  2,
		Name:    "Foo Bar",
		Iban:    "FOO12",
		Address: "Foo street 10",
		Amount:  100,
		Type:    account.TypeSending,
	})
	if err != nil {
		logger.Info("main", zap.String("error", err.Error()))
		fmt.Println(errors.Is(err, account.ErrAccountAlreadyExists))
		// return
	}
	ac, err := manager.GetAccounts()
	fmt.Println("get accounts", ac, err)
}
