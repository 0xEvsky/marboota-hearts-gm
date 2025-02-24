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
	players [4]*Player
	state   TableState
	turn    int
}

type PlayerState int

const (
	PlayerUnavailable PlayerState = iota
	PlayerWaiting
	PlayerReady
	PlayerTrumping
	PlayerPlaying
)

type Team int

const (
	TeamA Team = iota
	TeamB
)

type Player struct {
	client  *Client
	state   PlayerState
	hand    []string
	seat    int
	team    Team
	score   int
	partner *Player
	isTurn  bool
}

func newTable() Table {
	var players = [4]*Player{}
	for i := range players {
		players[i] = &Player{
			state:   PlayerUnavailable,
			seat:    i + 1,
			team:    Team(i % 2),
			partner: players[(i+2)%4],
			isTurn:  false,
		}
	}

	return Table{
		players: players,
		state:   TableWaiting,
		turn:    0,
	}
}

func (t *Table) seatPlayer(c *Client, s int) error {
	if s < 1 || s > 4 {
		return errors.New("invalid seat")
	}

	var p = t.players[s-1]
	if p.client != nil {
		return errors.New("seat is taken")
	}

	t.unseatPlayer(c)

	p.client = c
	c.player = p
	c.state = ClientSeated
	return nil
}

func (t *Table) unseatPlayer(c *Client) {
	if c.state != ClientSeated {
		return
	}
	c.player.client = nil
	c.state = ClientIdle
}
