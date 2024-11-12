package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	address string
	conn    net.Conn
}

func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

func (c *Client) Start() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	c.conn = conn
	defer c.conn.Close()

	go c.readMessages()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Connected to the server. Type your messages:")
	for scanner.Scan() {
		message := scanner.Text()
		_, err := c.conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Failed to send message:", err)
			return err
		}
	}

	return nil
}

func (c *Client) readMessages() {
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closed by server:", err)
			return
		}
		fmt.Println("Message from server:", string(buf[:n]))
	}
}
