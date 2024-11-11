package main

import (
	"os"

	"accounting-immudb-demo/pkg/account"
	"accounting-immudb-demo/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	token := os.Getenv("API_PRIVATE_KEY")

	// there is some sort of protection that does not allows me to create entries in collection
	// that is not named default
	manger := account.NewAccountManager("default", "default", token)

	err := manger.CreateAccountCollection()
	if err != nil {
		logger.Info("main create", zap.String("error", err.Error()))
		return
	}

	err = manger.CreateEntry(account.Account{
		Number:  1,
		Name:    "Foo Bar",
		Iban:    "FOO12",
		Address: "Foo street 10",
		Amount:  100,
		Type:    account.TypeSending,
	})
	if err != nil {
		logger.Info("main", zap.String("error", err.Error()))
		return
	}
}
