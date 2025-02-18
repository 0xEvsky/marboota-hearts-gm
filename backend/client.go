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
	conn     *websocket.Conn
	instance *Instance
	id       string
	name     string
	iconUrl  string
	state    ClientState
	table    *Table
	seat     int
	isTurn   bool
	write    func(msg []byte) error
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
	}
}
