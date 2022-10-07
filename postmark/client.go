package postmark

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Client struct {
	httpClient      *http.Client
	accountApiToken string
	baseURL         string
}

func NewClient(accountApiToken string) *Client {
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
		accountApiToken: accountApiToken,
		baseURL:         "https://api.postmarkapp.com",
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

	if params.ServerToken != nil {
		req.Header.Add("X-Postmark-Server-Token", *params.ServerToken)
	} else {
		req.Header.Add("X-Postmark-Account-Token", client.accountApiToken)
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if dst == nil || res.StatusCode > 399 {
		var apiRes APIError
		err = json.Unmarshal(body, &apiRes)
		if err != nil {
			return err
		}
		if apiRes.ErrorCode != 0 {
			err = errors.New(apiRes.Message)
			return err
		}
	} else {
		err = json.Unmarshal(body, dst)
	}

	return err
}

type APIError struct {
	ErrorCode int64
	Message   string
}

func (res APIError) Error() string {
	return res.Message
}
