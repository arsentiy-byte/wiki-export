package http

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Client struct {
	client *resty.Client
}

func NewClient(host string, port int, timeout string) *Client {
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}

	baseUrl := fmt.Sprintf("%s", host)

	if port != 0 {
		baseUrl = fmt.Sprintf("%s:%d", baseUrl, port)
	}

	restyClient := resty.New()
	restyClient.SetTimeout(duration)
	restyClient.SetBaseURL(baseUrl)

	return &Client{
		client: restyClient,
	}
}

func (c *Client) GetClient() *resty.Client {
	return c.client
}
