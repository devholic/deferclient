package deferstats

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// GetExpvar captures expvar using ExpvarHost and ExpvarEndpoint parameters
func (c *Client) GetExpvar() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", c.ExpvarHost+c.ExpvarEndpoint, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("Expvars not found")
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
