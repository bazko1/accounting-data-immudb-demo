package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"accounting-immudb-demo/pkg/client"
	"accounting-immudb-demo/pkg/logger"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var ErrAccountAlreadyExists = errors.New("account with given number already exists")

// account data manager based on immudb
type AccountManager struct {
	Ledger     string
	Collection string
	token      string
	client     client.ImmuDBClient
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
func (am AccountManager) CreateAccountCollection(ctx context.Context) error {
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
		"idFieldName": "_id",
		// TODO: I might want to create reflect based generator for that
		"fields": []map[string]any{
			{"name": "number", "type": "INTEGER"},
			{"name": "name", "type": "STRING"},
			{"name": "iban", "type": "STRING"},
			{"name": "address", "type": "STRING"},
			{"name": "amount", "type": "INTEGER"},
			{"name": "type", "type": "STRING"},
		},

		"indexes": []map[string]any{
			{"fields": []string{"number"}, "isUnique": true},
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal account")
	}
	response, err := am.client.DoPutRequest(ctx, fmt.Sprintf("/ledger/%s/collection/%s",
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

func (am AccountManager) CreateAccount(ctx context.Context, acc Account) error {
	jsonBytes, err := json.Marshal(acc)
	if err != nil {
		return errors.Wrap(err, "failed to marshal account")
	}

	response, err := am.client.DoPutRequest(ctx, fmt.Sprintf("/ledger/%s/collection/%s/document",
		am.Ledger,
		am.Collection),
		bytes.NewBuffer(jsonBytes))
	if err != nil {
		return errors.Wrapf(err, "failed to create entry sending PUT request")
	}
	defer response.Body.Close()
	if err := client.CheckResponse(response); err != nil {
		if errors.Is(err, client.HTTPConflictResponseErr) {
			mainErr := err
			b, err := io.ReadAll(response.Body)
			if err != nil {
				logger.Error("failed to read conflict err body", zap.Error(err))
				return errors.Wrap(mainErr, "failed to create new entry")
			}

			errData := struct {
				Error  string
				Status string
			}{}
			err = json.Unmarshal(b, &errData)
			if err != nil {
				logger.Error("failed to unmarshal json error", zap.Error(err))
				return errors.Wrap(mainErr, "failed to create new entry")
			}

			if errData.Status == "Conflict" &&
				errData.Error == "unable to create document, error document already exists" {
				return fmt.Errorf("can't create document with number %d: %w", acc.Number, ErrAccountAlreadyExists)
			}
		}

		return errors.Wrap(err, "failed to create new entry")
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Debug("entry created", zap.String("respBody", string(b)))
	}

	return nil
}

func (am AccountManager) GetAccounts(ctx context.Context) ([]Account, error) {
	count, err := am.client.GetCollectionCount(am.Ledger, am.Collection)
	if err != nil {
		return nil, errors.Wrap(err, "get accounts failed to list collection")
	}

	logger.Debug("GetAccount",
		zap.Int("collectionCount", count),
		zap.String("collection", am.Collection))

	if count == 0 {
		return []Account{}, nil
	}

	// TODO: If count is higher than 100 we need to loop over pages
	// FIXME: count > 100 page looping and appending
	if count > 100 {
		count = 100
	}

	searchSchema := struct {
		Page    int `json:"page"`
		PerPage int `json:"perPage"`
	}{Page: 1, PerPage: 100}

	jsonBytes, err := json.Marshal(searchSchema)
	if err != nil {
		return nil, errors.Wrap(err, "GetAccounts failed to marshal search schema")
	}

	response, err := am.client.DoPostRequest(ctx, fmt.Sprintf("/ledger/%s/collection/%s/documents/search",
		am.Ledger,
		am.Collection),
		bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list collection for accounts sending POST request")
	}

	if err := client.CheckResponse(response); err != nil {
		return nil, errors.Wrapf(err, "failed to list collection for accounts sending POST request")
	}

	type accountDocument struct {
		Document Account `json:"document"`
	}
	type SearchAccountResponse struct {
		Revisions []accountDocument `json:"revisions"`
	}
	searchResponse := SearchAccountResponse{Revisions: make([]accountDocument, 0, count)}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to list collection for accounts while reading body")
	}

	err = json.Unmarshal(data, &searchResponse)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to list collection for accounts while unmarshaling body")
	}
	outAccounts := make([]Account, 0, count)
	for _, rev := range searchResponse.Revisions {
		outAccounts = append(outAccounts, rev.Document)
	}

	return outAccounts, nil
}
