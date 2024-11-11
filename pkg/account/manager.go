package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"accounting-immudb-demo/pkg/client"
	"accounting-immudb-demo/pkg/logger"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// account data manager based on immudb
type AccountManager struct {
	Ledger     string
	Collection string
	token      string
	client     client.ImmuDBClient
	logger     *zap.Logger
}

func NewAccountManager(ledger, collection, token string) AccountManager {
	return AccountManager{
		Ledger:     ledger,
		Collection: collection,
		token:      token,
		client: client.NewImmuDBClient(
			client.ImmuDBClientConfig{
				Token: token,
			}),
	}
}

// CreateAccountCollection creates new account collection if it does not exist
func (am AccountManager) CreateAccountCollection() error {
	names, err := am.client.ListCollectionsName(am.Ledger)
	if err != nil {
		return errors.Wrap(err, "failed to list collections")
	}

	for _, name := range names {
		if name == am.Collection {
			logger.Debug("CreateAccountCollection",
				zap.String("collection", am.Collection),
				zap.String("status", "found"))
			return nil
		}
	}

	// TODO: In the future add force flag that deletes / recreates collection

	jsonBytes, err := json.Marshal(map[string]any{
		"idFieldName": "number",
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal account")
	}
	response, err := am.client.DoPutRequest(fmt.Sprintf("/ledger/%s/collection/%s",
		am.Ledger,
		am.Collection),
		bytes.NewBuffer(jsonBytes))
	if err != nil {
		return errors.Wrapf(err, "failed to create collection for accounts sending PUT request")
	}
	defer response.Body.Close()
	if err := client.CheckResponse(response); err != nil {
		return errors.Wrap(err, "failed to create accounts collection")
	}

	logger.Debug("CreateAccountCollection",
		zap.String("collection", am.Collection),
		zap.String("ledger", am.Ledger),
		zap.String("status", "created"))

	return nil
}

func (am AccountManager) CreateEntry(acc Account) error {
	jsonBytes, err := json.Marshal(acc)
	if err != nil {
		return errors.Wrap(err, "failed to marshal account")
	}

	response, err := am.client.DoPutRequest(fmt.Sprintf("/ledger/%s/collection/%s/document",
		am.Ledger,
		am.Collection),
		bytes.NewBuffer(jsonBytes))
	if err != nil {
		return errors.Wrapf(err, "failed to create entry sending PUT request")
	}
	defer response.Body.Close()
	if err := client.CheckResponse(response); err != nil {
		return errors.Wrap(err, "failed to create new entry")
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Debug("entry created", zap.String("respBody", string(b)))
	}

	return nil
}

func (am AccountManager) GetEntry(number uint) (Account, error) {
	return Account{}, nil
}
