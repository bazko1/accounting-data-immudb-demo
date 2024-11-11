package client

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

func (c ImmuDBClient) ListCollectionsName(ledger string) ([]string, error) {
	resp, err := c.DoGetRequest(fmt.Sprintf("/ledger/%s/collections",
		ledger),
		nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed GET ledger %q collections", ledger)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read body")
	}

	if err := CheckResponse(resp); err != nil {
		return nil, errors.Wrap(err, "GET list collections failed")
	}
	collections := struct{ Collections []struct{ Name string } }{}
	if err := json.Unmarshal(data, &collections); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal collections")
	}

	out := make([]string, len(collections.Collections))
	for _, c := range collections.Collections {
		out = append(out, c.Name)
	}

	return out, nil
}

// Create document response
type CreateDocumentSuccessResponse struct {
	TransactionID string `json:"transactionId,omitempty"`
	DocumentID    string `json:"documentId,omitempty"`
}
