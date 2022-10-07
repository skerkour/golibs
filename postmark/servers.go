package postmark

import (
	"fmt"
	"net/http"
)

type Server struct {
	// ID of server
	ID int64 `json:",omitempty"`
	// Name of server
	Name string
	// ApiTokens associated with server.
	ApiTokens []string `json:",omitempty"`
	// ServerLink to your server overview page in Postmark.
	ServerLink string `json:",omitempty"`
	// Color of the server in the rack screen. Purple Blue Turquoise Green Red Yellow Grey
	Color string
	// SmtpApiActivated specifies whether or not SMTP is enabled on this server.
	SmtpApiActivated bool
	// RawEmailEnabled allows raw email to be sent with inbound.
	RawEmailEnabled bool
	// InboundAddress is the inbound email address
	InboundAddress string `json:",omitempty"`
	// InboundHookUrl to POST to every time an inbound event occurs.
	InboundHookUrl string
	// PostFirstOpenOnly - If set to true, only the first open by a particular recipient will initiate the open webhook. Any
	// subsequent opens of the same email by the same recipient will not initiate the webhook.
	PostFirstOpenOnly bool
	// TrackOpens indicates if all emails being sent through this server have open tracking enabled.
	TrackOpens bool
	// InboundDomain is the inbound domain for MX setup
	InboundDomain string
	// InboundHash is the inbound hash of your inbound email address.
	InboundHash string
	// InboundSpamThreshold is the maximum spam score for an inbound message before it's blocked.
	InboundSpamThreshold int64
}

func (client *Client) GetServer(serverID string) (Server, error) {
	res := Server{}
	err := client.request(requestParams{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/servers/%s", serverID),
	}, &res)

	return res, err
}

func (client *Client) EditServer(serverID string, server Server) (Server, error) {
	res := Server{}
	err := client.request(requestParams{
		Method:  http.MethodPut,
		URL:     fmt.Sprintf("/servers/%s", serverID),
		Payload: server,
	}, &res)

	return res, err
}

func (client *Client) CreateServer(server Server) (Server, error) {
	res := Server{}
	err := client.request(requestParams{
		Method:  http.MethodPost,
		URL:     "/servers",
		Payload: server,
	}, &res)

	return res, err
}

func (client *Client) DeleteServer(serverID string) (err error) {
	err = client.request(requestParams{
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("/servers/%s", serverID),
	}, nil)

	return
}
