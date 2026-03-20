package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/qiudao/rlstudy/pkg/env"
)

type Client struct {
	base   string
	client *http.Client
}

// New creates a client connecting via TCP to the given base URL.
func New(baseURL string) *Client {
	return &Client{base: baseURL, client: &http.Client{}}
}

// NewUnix creates a client connecting via Unix domain socket.
func NewUnix(socketPath string) *Client {
	return &Client{
		base: "http://unix",
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
		},
	}
}

func (c *Client) Info() (env.InfoResponse, error) {
	var info env.InfoResponse
	resp, err := c.client.Get(c.base + "/info")
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}

func (c *Client) Reset() error {
	resp, err := c.client.Post(c.base+"/reset", "application/json", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("reset: status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) Step(action int) (env.StepResponse, error) {
	var result env.StepResponse
	body, _ := json.Marshal(env.StepRequest{Action: action})
	resp, err := c.client.Post(c.base+"/step", "application/json", bytes.NewReader(body))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return result, fmt.Errorf("step: status %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
