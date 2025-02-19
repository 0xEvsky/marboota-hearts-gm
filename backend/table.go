package main

import (
	"errors"
)

type TableState int

const (
	TableWaiting TableState = iota
	TableTrumping
	TablePlaying
)

type Table struct {
	players    [4]*Client
	spectators []*Client
	state      TableState
	turn       int
}

func NewTable() Table {
	return Table{
		players:    [4]*Client{},
		spectators: []*Client{},
		state:      TableWaiting,
		turn:       0,
	}
}

func (t *Table) seatPlayer(c *Client) error {
	for i, p := range t.players {
		if p == nil {
			t.players[i] = c
			c.table = t
			c.seat = i + 1
			c.state = ClientWaiting
			if t.state != TableWaiting {
				c.state = ClientPlaying
				c.isTurn = t.turn == c.seat
			}
			return nil
		}
	}
	return errors.New("Table is full")
}

func (t *Table) unseatPlayer(p *Client) {
	t.players[p.seat] = nil
	p.table = nil
	p.seat = 0
	p.state = ClientIdle
	p.isTurn = false
}
