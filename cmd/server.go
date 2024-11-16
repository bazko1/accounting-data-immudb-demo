package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"accounting-immudb-demo/pkg/account"
	"accounting-immudb-demo/pkg/logger"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AccountResponse struct {
	Accounts []account.Account `json:"accounts"`
	Message  string            `json:"message,omitempty"`
	Status   int               `json:"status,omitempty"`
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
		for i := range 6 {
			_ = manager.CreateAccount(ctx, account.Account{
				Number:  uint(i),
				Name:    fmt.Sprintf("Foo Bar%d", i),
				Iban:    fmt.Sprintf("US12%d", i),
				Address: "Foo street 10",
				Amount:  100 + uint(i),
				Type:    account.TypeSending,
			})
		}
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderAccessControlAllowOrigin},
	}))

	e.Use(middleware.Logger())

	g := e.Group("/api/v1/account")

	g.GET("", func(c echo.Context) error {
		accounts, err := manager.GetAccounts(c.Request().Context())
		if err != nil {
			logger.Error("failed to get accounts",
				zap.Error(err))

			return c.JSON(http.StatusInternalServerError,
				AccountResponse{
					Message:  "Failed to list accounts",
					Status:   http.StatusInternalServerError,
					Accounts: []account.Account{},
				},
			)
		}

		return c.JSON(http.StatusOK, AccountResponse{
			Message:  "OK",
			Status:   http.StatusOK,
			Accounts: accounts,
		},
		)
	})

	g.POST("", func(c echo.Context) error {
		acc := account.Account{}

		decoder := json.NewDecoder(c.Request().Body)

		if err := decoder.Decode(&acc); err != nil {
			logger.Error("Failed to decode account", zap.Error(err))
			return c.JSON(http.StatusInternalServerError,
				AccountResponse{
					Message: "Failed to decode account or missing data.",
					Status:  http.StatusInternalServerError,
				})
		}

		err := manager.CreateAccount(c.Request().Context(), acc)
		if err != nil {
			if errors.Is(err, account.ErrAccountAlreadyExists) {
				return c.JSON(http.StatusInternalServerError,
					AccountResponse{
						Message: err.Error(),
						Status:  http.StatusInternalServerError,
					})
			}
			logger.Error("Failed to create new account",
				zap.Error(err),
				zap.Any("account", acc))
			return c.JSON(http.StatusInternalServerError, AccountResponse{
				Message: "Failed to create new account",
				Status:  http.StatusInternalServerError,
			})

		}

		return c.JSON(http.StatusOK, AccountResponse{
			Message: "Failed to create new account",
			Status:  http.StatusOK,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
