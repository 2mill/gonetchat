package main

import (
	"fmt"
	"net"
)

type Client struct {
	name string
	conn net.Conn
}

func (self *Client) Message(message string) error {
	msg := []byte(message)
	_, err := self.conn.Write(msg)
	return err
}

func NewClient(name string) Client {
	return Client{
		name,
		nil,
	}
}

func (self *Client) Connect(ip string, port string) error {
	addr := fmt.Sprintf("%s:%s", ip, port)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	self.conn = conn

	return nil
}

func (self *Client) Close() error {
	return self.conn.Close()
}
