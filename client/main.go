package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Client is the flagger client
type Client struct {
	URL     string
	Default bool
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
	return fmt.Sprintf("%s/%s/%s", c.URL, flag, env)
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
	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
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
