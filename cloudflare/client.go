package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiToken   string
	baseURL    string
}

func NewClient(apiToken string) *Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Client{
		httpClient: &http.Client{
			Transport: transport,
		},
		apiToken: apiToken,
		baseURL:  "https://api.cloudflare.com",
	}
}

type requestParams struct {
	Method      string
	URL         string
	Payload     interface{}
	ServerToken *string
}

func (client *Client) request(params requestParams, dst interface{}) error {
	url := client.baseURL + params.URL
	var apiRes ApiResponse

	req, err := http.NewRequest(params.Method, url, nil)
	if err != nil {
		return err
	}

	if params.Payload != nil {
		payloadData, err := json.Marshal(params.Payload)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(payloadData))
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+client.apiToken)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&apiRes)
	if err != nil {
		err = fmt.Errorf("cloudflare: decoding JSON body: %w", err)
		return err
	}

	if len(apiRes.Errors) != 0 {
		err = fmt.Errorf("cloudflare: %s", apiRes.Errors[0].Error())
		return err
	}

	if dst != nil {
		err = json.Unmarshal(apiRes.Result, dst)
		if err != nil {
			err = fmt.Errorf("cloudflare: decoding JSON result: %w", err)
			return err
		}
	}

	return nil
}

type ApiResponse struct {
	Result  json.RawMessage `json:"result"`
	Success bool            `json:"success"`
	Errors  []ApiError      `json:"errors"`
}

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (res ApiError) Error() string {
	return res.Message
}
