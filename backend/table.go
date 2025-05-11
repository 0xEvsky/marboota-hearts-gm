package main

import (
	"errors"
	"math/rand/v2"
	"slices"

	"github.com/OmarQurashi868/marboota/backend/clog"
)

type TableState int

const (
	TableWaiting TableState = iota
	TableTrumping
	TablePlaying
)

type Table struct {
	instance    *Instance
	players     [4]*Player
	state       TableState
	turn        int
	turnOffset  int
	trump       Trump
	play        Play
	playCount   int
	rounds      []Round
	totalScores map[Team]int
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

type Trump struct {
	highestCall   int
	highestCaller *Player
	isDone        bool
	suit          Suit
}

type Play struct {
	cards        []Card
	curWinCard   Card
	curWinPlayer *Player
}

type Round struct {
	teamAScore int
	teamBScore int
}

func newTable() Table {
	var players = [4]*Player{}
	for i := range players {
		players[i] = &Player{
			seat:    i,
			team:    Team(i % 2),
			partner: players[(i+2)%4],
		}
	}

	return Table{
		players: players,
		state:   TableWaiting,
		trump:   Trump{},
		play:    Play{},
	}
}

func (t *Table) seatPlayer(c *Client, s int) error {
	if s < 0 || s > 3 {
		return errors.New("invalid seat")
	}

	var p = t.players[s]
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
	clog.Debugf("(i:%s) everyone ready", t.instance.id)
	return true
}

func (t *Table) startGame() {
	t.instance.Broadcast(map[string]string{"ACTION": "GAMESTART"})
	clog.Debugf("(i:%s) game started", t.instance.id)
	t.startTrump()
}

func (t *Table) startTrump() {
	t.trump = Trump{}
	t.play = Play{}
	t.playCount = 0
	for _, p := range t.players {
		p.state = PlayerTrumping
		p.score = 0
		p.hand = []Card{}
		p.isTurn = false
	}

	// Reshuffles
	for {
		var deck = newDeck()

		// Shuffle deck
		for i := range deck {
			j := rand.IntN(i + 1)
			deck[i], deck[j] = deck[j], deck[i]
		}

		// Deal hands
		for i := range deck {
			t.players[i/13].hand = append(t.players[i/13].hand, deck[i])
		}

		// Check hand validity
		for _, p := range t.players {
			if p.isHandInvalid() {
				continue
			}
		}
		break
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
		p.client.writeJson(map[string]string{"ACTION": "DEAL", "CARDS": p.getHandString()})
	}

	t.instance.Broadcast(map[string]string{"ACTION": "TRUMPSTART"})
	clog.Debugf("(i:%s) trump started", t.instance.id)
	t.state = TableTrumping
	t.turn = t.turnOffset

	t.players[t.turn].isTurn = true
	t.players[t.turn].client.writeJson(map[string]string{"ACTION": "YOURTRUMPCALL", "MINSCORE": "7", "MAXSCORE": "11"})
}

func (t *Table) startPlay() {
	// StartPlay
	// Change table state to playing
	t.state = TablePlaying

	t.instance.Broadcast(map[string]string{"ACTION": "PLAYSTART"})
	clog.Debugf("(i:%s) play started", t.instance.id)

	for _, p := range t.players {
		p.state = PlayerPlaying
	}

	t.players[t.turn].isTurn = false
	t.turn = t.trump.highestCaller.seat
	t.trump.highestCaller.isTurn = true

	t.playCount = 0

	var _, cardsStr = t.trump.highestCaller.getPlayableCards()

	t.trump.highestCaller.client.writeJson(map[string]string{"ACTION": "YOURPLAY", "PLAYABLE": cardsStr})
}

func (p *Player) getPlayableCards() ([]Card, string) {
	if p.state != PlayerPlaying {
		return []Card{}, ""
	}

	var cards = []Card{}

	for _, c := range p.hand {
		if c.suit == p.client.instance.table.trump.suit {
			cards = append(cards, c)
		}
	}

	if len(cards) == 0 {
		cards = p.hand
	}

	var str = ""

	for i, c := range cards {
		str += c.name
		if i < len(cards)-1 {
			str += ","
		}
	}

	return cards, str
}

func (p *Player) getHandString() string {
	var str = ""

	for i, c := range p.hand {
		str += c.name
		if i < len(p.hand)-1 {
			str += ","
		}
	}

	return str
}

func (p *Player) isHandInvalid() bool {
	var suitCardCounts = map[Suit]int{}
	var faceCardCounts = 0

	for _, c := range p.hand {
		suitCardCounts[c.suit] += 1
		if c.value > 10 {
			faceCardCounts += 1
		}

		if faceCardCounts >= 8 {
			return false
		}
	}

	if faceCardCounts == 0 {
		return false
	}

	for i := range 4 {
		if suitCardCounts[Suit(i)] >= 8 {
			return false
		}
	}

	return true

}

func (p *Player) getAvailableTrumps() ([]string, string) {
	var trumpArr = []string{}
	var trumpsStr = ""
	var suitCounts = map[Suit]int{}
	for _, c := range p.hand {
		suitCounts[c.suit] += 1
	}
	var suits = map[Suit]string{Spades: "SPADES", Hearts: "HEARTS", Clubs: "CLUBS", Diamonds: "DIAMONDS"}
	for k, v := range suits {
		if suitCounts[k]+3 <= p.client.instance.table.trump.highestCall {
			trumpArr = append(trumpArr, v)
			trumpsStr += v
			if k < 4 {
				trumpsStr += ","
			}
		}
	}

	return trumpArr, trumpsStr
}
