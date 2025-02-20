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
	ClientWaiting
	ClientPlaying
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
	table    *Table
	seat     int
	isTurn   bool
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:     conn,
		isAuthed: false,
		instance: nil,
		id:       "",
		name:     "",
		iconUrl:  "",
		state:    ClientUnavailable,
		table:    nil,
		seat:     0,
		isTurn:   false,
	}
}

func (c *Client) write(msg []byte) error {
	// TODO: Investigate mutex
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *Client) writeJson(msg map[string]string) error {
	return c.write(toJson(msg))
}

func (c *Client) writeError(msg string) error {
	return c.write(toJson(map[string]string{"ACTION": "ERROR", "MESSAGE": msg}))
}

func (c *Client) writeOk() error {
	return c.write(toJson(map[string]string{"ACTION": "OK"}))
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
