package main

import (
	"errors"
	"math/rand/v2"
	"slices"
)

type TableState int

const (
	TableWaiting TableState = iota
	TableTrumping
	TablePlaying
)

type Table struct {
	instance   *Instance
	players    [4]*Player
	state      TableState
	turn       int
	turnOffset int
	trump      Trump
}

type Trump struct {
	highestCall   int
	highestCaller *Player
	callers       []*Player
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
	hand    []Card
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
			seat:    i + 1,
			team:    Team(i % 2),
			partner: players[(i+2)%4],
		}
	}

	return Table{
		players: players,
		state:   TableWaiting,
		trump: Trump{
			callers: players[:],
		},
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

	// Change depending on game state
	p.state = PlayerWaiting
	if t.state == TableTrumping {
		p.state = PlayerTrumping
	}
	if t.state == TablePlaying {
		p.state = PlayerPlaying
	}

	return nil
}

func (t *Table) unseatPlayer(c *Client) {
	if c.state != ClientSeated {
		return
	}

	var p = c.player

	p.client = nil
	p.state = PlayerUnavailable
	c.player = nil

	c.state = ClientIdle
}

func (t *Table) isEveryoneReady() bool {
	for _, p := range t.players {
		if p.state != PlayerReady {
			return false
		}
	}
	return true
}

func (t *Table) trumpStart() {
	var deck = newDeck()

	// shuffle deck
	for i := range deck {
		j := rand.IntN(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	// deal hands
	for i := range deck {
		t.players[i/13].hand = append(t.players[i/13].hand, deck[i])
	}

	// Sort hands
	for _, p := range t.players {
		slices.SortFunc(p.hand, func(i Card, j Card) int {
			if i.suit < j.suit {
				return -1
			}

			if i.suit == j.suit {
				if i.value > j.value {
					return -1
				} else {
					return 1
				}
			}

			return 1
		})
	}

	for _, p := range t.players {
		var cardstring = ""
		for i, c := range p.hand {
			cardstring += c.name
			if i < len(p.hand)-1 {
				cardstring += ","
			}
		}

		p.client.writeJson(map[string]string{"ACTION": "DEAL", "CARDS": cardstring})
	}

	t.instance.Broadcast(map[string]string{"ACTION": "TRUMPSTART"})
	t.state = TableTrumping
	t.turn = t.turnOffset

	for _, p := range t.players {
		p.state = PlayerTrumping
	}
	t.players[t.turn].isTurn = true
	t.players[t.turn].client.writeJson(map[string]string{"ACTION": "YOURTRUMPCALL", "MINSCORE": "7"})
}
