package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	flagger "github.com/mnbbrown/flagger/pkg"
)

// Client is the flagger client
type Client struct {
	URL     string
	Default bool
}

// List returns a list of all flags
func (c *Client) List() (map[string]map[string]*flagger.Flag, error) {
	resp, err := http.Get(fmt.Sprintf("%s/flags", c.URL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response := make(map[string]map[string]*flagger.Flag)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) get(flag, env string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/flags/%s/%s", c.URL, flag, env))
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()
	v, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return true, err
	}
	asBool, e := strconv.ParseBool(string(v))
	return asBool, e
}

func (c *Client) getURL(flag, env string) string {
	return fmt.Sprintf("%s/flags/%s/%s", c.URL, flag, env)
}

// Set sets a flag
func (c *Client) Set(flag, env, flagType, value string) error {
	f := struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{
		Type:  flagType,
		Value: value,
	}
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	url := c.getURL(flag, env)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Bad request")
	}
	return nil
}

// NewClient creates a new client
func NewClient(url string) *Client {
	return &Client{URL: url, Default: true}
}

// Get returns the flag value as a bool
func (c *Client) Get(flag, env string) bool {
	v, err := c.get(flag, env)
	if err != nil {
		return c.Default
	}
	return v
}
