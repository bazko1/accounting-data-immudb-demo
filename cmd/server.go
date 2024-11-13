package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"accounting-immudb-demo/pkg/account"
	"accounting-immudb-demo/pkg/logger"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	token := os.Getenv("API_PRIVATE_KEY")

	manager := account.NewAccountManager("default", "default", token)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := manager.CreateAccountCollection(ctx)
	if err != nil {
		logger.Error("failed to create collection", zap.String("error", err.Error()))
		panic(err)
	}
	logger.Info("Created account collection")

	e := echo.New()
	g := e.Group("/api/v1/accounts")
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g.GET("/list", func(c echo.Context) error {
		accounts, err := manager.GetAccounts(c.Request().Context())
		if err != nil {
			logger.Error("failed to get accounts",
				zap.Error(err))
			c.Response().Status = http.StatusInternalServerError

			return c.JSON(http.StatusInternalServerError,
				map[string]string{"error": "failed to list accounts"})
		}

		return c.JSON(http.StatusOK, accounts)
	})

	g.POST("/add", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
