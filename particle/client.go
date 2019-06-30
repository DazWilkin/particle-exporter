package particle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	client *http.Client
	token  string
}

func NewClient(token string) Client {
	return Client{
		client: &http.Client{},
		token:  token,
	}
}
func (c *Client) get(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func (c *Client) GetDevices() (devices Devices, err error) {
	body, err := c.get(urlDevices)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
func (c *Client) GetDiagnostics(device string) (DiagnosticsResponse, error) {
	body, err := c.get(fmt.Sprintf("%s/%s", urlDiagnostics, device))
	if err != nil {
		return DiagnosticsResponse{}, err
	}
	dr := DiagnosticsResponse{}
	dr.Diagnostics = []Diagnostic{}
	err = json.Unmarshal(body, &dr)
	if err != nil {
		return DiagnosticsResponse{}, err
	}
	return dr, nil
}
func (c *Client) GetIntegrations() (Integrations, error) {
	body, err := c.get(urlIntegrations)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))
	ii := Integrations{}
	json.Unmarshal(body, &ii)
	log.Println(len(ii))
	return ii, nil
}
func (c *Client) GetIntegration(id string) (Integration, error) {
	body, err := c.get(fmt.Sprintf("%s/%s", urlIntegrations, id))
	if err != nil {
		return Integration{}, err
	}
	log.Println(string(body))

	ir := IntegrationResponse{}
	ir.Integration.Logs = []Log{}
	ir.Integration.Errors = []IError{}
	ir.Integration.Counters = []Counter{}

	json.Unmarshal(body, &ir)
	return ir.Integration, nil
}
