package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

var HTTPConflictResponseErr = errors.New("409 HTTP Conflict response")

func (c ImmuDBClient) ListCollectionsName(ctx context.Context, ledger string) ([]string, error) {
	resp, err := c.DoGetRequest(ctx, fmt.Sprintf("/ledger/%s/collections",
		ledger),
		nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed GET ledger %q collections", ledger)
	}
	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return nil, errors.Wrap(err, "GET list collections failed")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read body")
	}

	collections := struct{ Collections []struct{ Name string } }{}
	if err := json.Unmarshal(data, &collections); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal collections")
	}

	out := make([]string, 0, len(collections.Collections))
	for _, c := range collections.Collections {
		out = append(out, c.Name)
	}

	return out, nil
}

func (c ImmuDBClient) GetCollectionCount(ctx context.Context, ledger, collection string) (int, error) {
	response, err := c.DoPostRequest(ctx, fmt.Sprintf("/ledger/%s/collection/%s/documents/count",
		ledger,
		collection),
		strings.NewReader("{}"))
	if err != nil {
		return 0, errors.Wrap(err, "collection PUT count request failed")
	}
	defer response.Body.Close()

	if err := CheckResponse(response); err != nil {
		return 0, errors.Wrap(err, "PUT requests for count fail")
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, errors.Wrap(err, "failed to read PUT count body")
	}

	countData := struct{ Count int }{}
	if err := json.Unmarshal(data, &countData); err != nil {
		return 0, errors.Wrap(err, "failed to unmarshal collection count")
	}

	return countData.Count, nil
}

// Create document response
type CreateDocumentSuccessResponse struct {
	TransactionID string `json:"transactionId,omitempty"`
	DocumentID    string `json:"documentId,omitempty"`
}
