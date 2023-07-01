package graphqlc

import (
	"bytes"
	"net/url"
)

type Interface interface {
	QueryOrMutate() QueryOrMutate
	Subscription() Subscription
}

type Client struct {
	url    *url.URL
	header map[string]string
	req    *GraphRequest
}

func NewClient(rawURL string) (*Client, error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &Client{url: uri}, nil
}

func (c *Client) Header(header map[string]string) *Client {
	c.header = header
	return c
}

func (c *Client) Body(req *GraphRequest) *Client {
	c.req = req
	return c
}

func (c *Client) SetOperationName(name string) *Client {
	c.req.OperationName = name
	return c
}

func (c *Client) QueryOrMutate() QueryOrMutate {
	q := &Query{
		url: c.url,
		req: bytes.NewReader(encode(c.req)),
	}
	switch q.url.Scheme {
	case "http", "https":
	case "ws":
		q.url.Scheme = "http"
	case "wss":
		q.url.Scheme = "https"
	default:
		q.url.Scheme = "https"
	}
	if c.header != nil {
		q.header = c.header
	}
	return q
}

func (c *Client) Subscription() Subscription {
	s := &Subscribe{
		url:     c.url,
		payload: encode(c.req),
		decoder: make(chan Decoder),
	}
	switch s.url.Scheme {
	case "ws", "wss":
	case "http":
		s.url.Scheme = "ws"
	case "https":
		s.url.Scheme = "wss"
	default:
		s.url.Scheme = "wss"
	}
	if c.header != nil {
		s.header = c.header
	}
	return s
}
