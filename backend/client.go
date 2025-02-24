package main

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientState int

const (
	ClientUnavailable ClientState = iota
	ClientIdle
	ClientSeated
)

type Client struct {
	mu       sync.Mutex
	conn     *websocket.Conn
	isAuthed bool
	instance *Instance
	id       string
	name     string
	iconUrl  string
	state    ClientState
	player   *Player
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) writeJson(msg map[string]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(msg)
}

func (c *Client) writeError(requestId string, msg string) error {
	return c.writeJson(map[string]string{"ACTION": "ERROR", "REQUESTID": requestId, "MESSAGE": msg})
}

func (c *Client) writeOk(requestId string) error {
	return c.writeJson(map[string]string{"ACTION": "OK", "REQUESTID": requestId})
}

func (c *Client) broadcastToInstance(msg map[string]string) error {
	if c.instance == nil {
		return errors.New("Client not authenticated")
	}

	for _, client := range c.instance.clients {
		if !client.isAuthed || client == c {
			continue
		}
		client.writeJson(msg)
	}

	return nil
}
