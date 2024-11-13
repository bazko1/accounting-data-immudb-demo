package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"accounting-immudb-demo/pkg/account"
	"accounting-immudb-demo/pkg/logger"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MessageRespones struct {
	Message string `json:"message,omitempty"`
}

func main() {
	token := os.Getenv("API_PRIVATE_KEY")

	manager := account.NewAccountManager("default", "default", token)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := manager.CreateAccountCollection(ctx)
	if err != nil {
		logger.Error("failed to create collection",
			zap.String("error", err.Error()))
		panic(err)
	}
	logger.Info("Initialized account manager")

	if os.Getenv("ADD_TEST_DATA") == "true" {
		for i := range 5 {
			_ = manager.CreateAccount(ctx, account.Account{
				Number:  uint(i),
				Name:    "Foo Bar",
				Iban:    "FOO12",
				Address: "Foo street 10",
				Amount:  100,
				Type:    account.TypeSending,
			})
		}
	}

	e := echo.New()
	g := e.Group("/api/v1/account")
	// Middleware
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	g.GET("", func(c echo.Context) error {
		accounts, err := manager.GetAccounts(c.Request().Context())
		if err != nil {
			logger.Error("failed to get accounts",
				zap.Error(err))
			c.Response().Status = http.StatusInternalServerError

			return c.JSON(http.StatusInternalServerError,
				MessageRespones{"Failed to list accounts"})
		}

		return c.JSON(http.StatusOK, accounts)
	})

	g.POST("", func(c echo.Context) error {
		acc := account.Account{}

		decoder := json.NewDecoder(c.Request().Body)

		if err := decoder.Decode(&acc); err != nil {
			logger.Error("Failed to decode account", zap.Error(err))
			return c.JSON(http.StatusInternalServerError,
				MessageRespones{"Failed to decode account or missing data."})

		}

		err := manager.CreateAccount(c.Request().Context(), acc)
		if err != nil {
			if errors.Is(err, account.ErrAccountAlreadyExists) {
				return c.JSON(http.StatusInternalServerError,
					MessageRespones{err.Error()})
			}
			logger.Error("Failed to create new account",
				zap.Error(err),
				zap.Any("account", acc))
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
