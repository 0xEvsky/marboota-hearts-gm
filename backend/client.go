package main

import (
	"errors"

	"github.com/gorilla/websocket"
)

type ClientState int

const (
	ClientUnavailable ClientState = iota
	ClientIdle
	ClientSpectating
	ClientWaiting
	ClientPlaying
)

type Client struct {
	conn       *websocket.Conn
	instance   *Instance
	id         string
	name       string
	iconUrl    string
	state      ClientState
	table      *Table
	seat       int
	isTurn     bool
	write      func(msg []byte) error
	writeJson  func(msg map[string]string) error
	writeError func(msg string) error
	writeOk    func() error
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:     conn,
		instance: nil,
		id:       "",
		name:     "",
		iconUrl:  "",
		state:    ClientUnavailable,
		table:    nil,
		seat:     0,
		isTurn:   false,
		write: func(msg []byte) error {
			return errors.New("Client uninitialized")
		},
		writeJson: func(msg map[string]string) error {
			return errors.New("Client uninitialized")
		},
		writeError: func(msg string) error {
			return errors.New("Client uninitialized")
		},
		writeOk: func() error {
			return errors.New("Client uninitialized")
		},
	}
}

func (c *Client) broadcastToInstance(msg map[string]string) error {
	if c.instance == nil {
		return errors.New("Client not authenticated")
	}

	for _, client := range c.instance.clients {
		if client == c {
			continue
		}
		client.writeJson(msg)
	}

	return nil
}
