package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"accounting-immudb-demo/pkg/logger"

	"github.com/pkg/errors"
	"go.uber.org/zap"
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

func (c ImmuDBClient) DoPutRequest(endpoint string,
	body io.Reader,
) (*http.Response, error) {
	return c.DoRequest(http.MethodPut, endpoint, body)
}

func (c ImmuDBClient) DoPostRequest(endpoint string,
	body io.Reader,
) (*http.Response, error) {
	return c.DoRequest(http.MethodPost, endpoint, body)
}

func (c ImmuDBClient) DoGetRequest(endpoint string,
	body io.Reader,
) (*http.Response, error) {
	return c.DoRequest(http.MethodGet, endpoint, body)
}

func (c ImmuDBClient) DoRequest(method string,
	endpoint string,
	body io.Reader) (*http.Response,
	error,
) {
	path, err := url.JoinPath(c.baseURL, endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create url path %q", path)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	logger.Debug("DoRequest", zap.String("method", method), zap.String("url", path))
	req, err := http.NewRequestWithContext(ctx,
		method,
		path,
		body)
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
	switch resp.StatusCode {
	case 200:
	case 409:
		return HTTPConflictResponseErr
	default:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("CheckResponse failed to read response body",
				zap.Error(err))
		}

		return fmt.Errorf("Got %d status code instead of 200.\nResponse body: %v",
			resp.StatusCode,
			string(b))
	}
	return nil
}
