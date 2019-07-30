package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

// ConnectionEventType not sure I need this...
type ConnectionEventType string

const (
	CONNECTION_EVENT_TYPE_NEW_CONNECTION           ConnectionEventType = "new_connection"
	CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED    ConnectionEventType = "connection_terminated"
	CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR ConnectionEventType = "general_error"
)

// Client holds info about connection
type Client struct {
	ID                string /* client is responsible of generating a ulid for each request, it will be sent in the response from the server so that client will know what request generated this response */
	DeviceID          string /* a unique id generated from the client itself */
	conn              net.Conn
	onConnectionEvent func(c *Client, eventType ConnectionEventType, e error) /* function for handling new connections */
	onDataEvent       func(c *Client, data []byte)                            /* function for handling new date events */
}

// NewClient creates a new instance of Client
func NewClient(conn net.Conn, onConnectionEvent func(c *Client, eventType ConnectionEventType, e error), onDataEvent func(c *Client, data []byte)) *Client {
	return &Client{
		conn:              conn,
		onConnectionEvent: onConnectionEvent,
		onDataEvent:       onDataEvent,
	}
}

// Read client data from channel
func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)

		switch err {
		case io.EOF:
			// connection terminated
			c.conn.Close()
			c.onConnectionEvent(c, CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED, err)
			return
		case nil:
			// new data available
			c.onDataEvent(c, buf[:n])
		default:
			log.Fatalf("Receive data failed:%s", err)
			c.conn.Close()
			c.onConnectionEvent(c, CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR, err)
			return
		}
	}
}

// Send text message to client
func (c *Client) Send(message []byte) error {
	_, err := c.conn.Write(message)
	return err
}

// SendBytes sends bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

// Conn returns a connection
func (c *Client) Conn() net.Conn {
	return c.conn
}

// Close closes a connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil

}
