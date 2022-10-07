package bunny

import (
	"fmt"
	"net/http"
)

type AddCustomHostnameInput struct {
	Hostname string `json:"Hostname"`
}

func (client *Client) AddCustomHostname(pullZone, hostname string) (err error) {
	err = client.request(requestParams{
		Payload: AddCustomHostnameInput{
			Hostname: hostname,
		},
		Method: http.MethodPost,
		URL:    fmt.Sprintf("/pullzone/%s/addHostname", pullZone),
	}, nil)

	return
}

type RemoveCustomHostnameInput struct {
	Hostname string `json:"Hostname"`
}

func (client *Client) RemoveCustomHostname(pullZone, hostname string) (err error) {
	err = client.request(requestParams{
		Payload: RemoveCustomHostnameInput{
			Hostname: hostname,
		},
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("/pullzone/%s/removeHostname", pullZone),
	}, nil)

	return
}

func (client *Client) LoadFreeCertificate(hostname string) (err error) {
	err = client.request(requestParams{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/pullzone/loadFreeCertificate?hostname=%s", hostname),
	}, nil)

	return
}
