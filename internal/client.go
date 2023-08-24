package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const httpClientTimeout = 20 * time.Second

type Client struct {
	httpClient *retryablehttp.Client
	baseURL    string
	accessKey  string
}

func NewClient(baseURL string, accessKey string) Client {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Timeout = httpClientTimeout
	httpClient.Logger = slog.Default()

	return Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		accessKey:  accessKey,
	}
}

func (c *Client) Delete(endpoint string) (*http.Response, error) {
	method := http.MethodDelete
	url := c.buildURL(endpoint)

	slog.Info("request", "method", method, "url", url)

	request, err := retryablehttp.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	return c.do(request)
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	method := http.MethodGet
	url := c.buildURL(endpoint)

	slog.Info("request", "method", method, "url", url)

	request, err := retryablehttp.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	return c.do(request)
}

func (c *Client) Post(endpoint string, data interface{}) (*http.Response, error) {
	method := http.MethodPost
	url := c.buildURL(endpoint)

	slog.Info("request", "method", method, "url", url)
	slog.Debug("POST request payload", "data", data)

	payload, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	request, err := retryablehttp.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	return c.do(request)
}

func (c *Client) buildURL(endpoint string) string {
	return fmt.Sprintf("%s%s", c.baseURL, endpoint)
}

func (c *Client) do(request *retryablehttp.Request) (*http.Response, error) {
	addHeaders(&request.Header, c.accessKey)

	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	status := fmt.Sprintf("%s %s", response.Proto, response.Status)

	if response.StatusCode == http.StatusOK {
		slog.Info(status)
		return response, nil
	} else {
		slog.Error(status)

		res, err := httputil.DumpResponse(response, true)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(string(res))
	}
}

func addHeaders(headers *http.Header, accessKey string) {
	headers.Add("Content-Type", "application/json")
	headers.Add("User-Agent", fmt.Sprintf("msgenctl/%v", Version))
	headers.Add("Ocp-Apim-Subscription-Key", accessKey)
}
