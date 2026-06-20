package natsclient

import (
	"time"
	"github.com/nats-io/nats.go"
)

type Client struct {
	Conn *nats.Conn
}

func NewClient(url string) (*Client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Client{Conn: nc}, nil
}

func (c *Client) Request(subject string, data []byte, timeout time.Duration) ([]byte, error) {
	msg, err := c.Conn.Request(subject, data, timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}

func (c *Client) RequestMsg(subject string, data []byte, header map[string]string, timeout time.Duration) ([]byte, error) {
	msg := nats.NewMsg(subject)
	msg.Data = data
	for k, v := range header {
		msg.Header.Set(k, v)
	}
	resp, err := c.Conn.RequestMsg(msg, timeout)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) QueueSubscribe(subject, queue string, handler func([]byte) []byte) error {
	_, err := c.Conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		resp := handler(msg.Data)
		if resp != nil {
			msg.Respond(resp)
		}
	})
	return err
}

func (c *Client) Publish(subject string, data []byte) error {
	return c.Conn.Publish(subject, data)
}

func (c *Client) Close() {
	c.Conn.Close()
}
