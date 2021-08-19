package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const httpClientTimeout = 20 * time.Second

type Client struct {
	httpClient http.Client
	baseURL    string
	accessKey  string
}

func NewClient(baseURL string, accessKey string) Client {
	return Client{
		httpClient: http.Client{
			Timeout: httpClientTimeout,
		},
		baseURL:   baseURL,
		accessKey: accessKey,
	}
}

func (c *Client) Delete(endpoint string) (*http.Response, error) {
	method := http.MethodDelete
	url := c.buildURL(endpoint)

	zap.S().Infof("%s %s", method, url)

	request, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	return c.do(request)
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	method := http.MethodGet
	url := c.buildURL(endpoint)

	zap.S().Infof("%s %s", method, url)

	request, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	return c.do(request)
}

func (c *Client) Post(endpoint string, data interface{}) (*http.Response, error) {
	method := http.MethodPost
	url := c.buildURL(endpoint)

	zap.S().Infof("%s %s", method, url)
	zap.S().Debugw("POST request payload", "data", data)

	payload, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	return c.do(request)
}

func (c *Client) buildURL(endpoint string) string {
	return fmt.Sprintf("%s%s", c.baseURL, endpoint)
}

func (c *Client) do(request *http.Request) (*http.Response, error) {
	addHeaders(&request.Header, c.accessKey)

	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	status := fmt.Sprintf("%s %s", response.Proto, response.Status)

	if response.StatusCode == http.StatusOK {
		zap.S().Info(status)
		return response, nil
	} else {
		zap.S().Error(status)
		return nil, fmt.Errorf(response.Status)
	}
}

func addHeaders(headers *http.Header, accessKey string) {
	headers.Add("Content-Type", "application/json")
	headers.Add("User-Agent", fmt.Sprintf("msgenctl/%v", Version))
	headers.Add("Ocp-Apim-Subscription-Key", accessKey)
}
