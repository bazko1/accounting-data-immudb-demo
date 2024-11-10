package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultTimeout    = time.Minute * 5
	defaultApiBaseURL = "https://vault.immudb.io/ics/api/v1/"
)

type ImmuDBClient struct {
	baseURL string
	token   string
}
type ImmuDBClientConfig struct {
	BaseUrl string
	Token   string
}

func NewImmuDBClient(config ImmuDBClientConfig) ImmuDBClient {
	baseUrl := config.BaseUrl
	if baseUrl == "" {
		baseUrl = defaultApiBaseURL
	}
	return ImmuDBClient{
		baseURL: baseUrl,
		token:   config.Token,
	}
}

func (c ImmuDBClient) DoPutRequest(endpoint string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + "ledger/default/collection/default/document"
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new request with context")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-Key", c.token)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sent PUT request")
	}
	return resp, nil
}

func CheckResponse(resp *http.Response) error {
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Got %d status code instead of 200. \nResponse Details: %v",
			resp.StatusCode,
			string(b))
	}
	return nil
}

// Create document response
type CreateDocumentSuccessResponse struct {
	TransactionID string `json:"transactionId,omitempty"`
	DocumentID    string `json:"documentId,omitempty"`
}
