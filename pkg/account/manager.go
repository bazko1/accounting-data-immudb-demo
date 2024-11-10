package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"accounting-immudb-demo/pkg/client"

	"github.com/pkg/errors"
)

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
	// TODO: Check if given db with ledger/collection is already created
	// and if not create
}

func (am AccountManager) CreateEntry(acc Account) error {
	jsonBytes, err := json.Marshal(acc)
	if err != nil {
		return errors.Wrap(err, "failed to marshal account")
	}

	response, err := am.client.DoPutRequest(fmt.Sprintf("/ledger/%s/collection/%s/documents",
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
	b, _ := io.ReadAll(response.Body)
	fmt.Println(string(b))
	return nil
}

func (am AccountManager) GetEntry(number uint) (Account, error) {
	return Account{}, nil
}
