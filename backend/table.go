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
	TableModeSelecting
	TableCardsPassing
	TableTrumping
	TablePlaying
)

type Table struct {
	instance         *Instance
	players          [4]*Player
	state            TableState
	turn             int
	turnOffset       int
	trump            Trump
	withTrump        bool
	play             Play
	playCount        int
	rounds           []Round
	totalScores      map[Team]int
	gameMode         GameMode
	roundPassedCards [][]Card
}

type PlayerState int

const (
	PlayerUnavailable PlayerState = iota
	PlayerWaiting
	PlayerReady
	PlayerTrumping
	PlayerPassingCards
	PlayerPassedCards
	PlayerRecievedCards
	PlayerPlaying
)

type Team int

const (
	TeamA Team = iota
	TeamB
)

type Player struct {
	client     *Client
	instance   *Instance
	state      PlayerState
	hand       []Card
	seat       int
	team       Team
	score      int
	partner    int
	isTurn     bool
	lastPrompt map[string]string
}

type Trump struct {
	players       []*Player
	gobool        *Player
	beIstifada    bool
	calls         []string
	highestCall   int
	highestCaller *Player
	suit          Suit
}

type Play struct {
	players      []*Player
	cards        []Card
	curWinCard   Card
	curWinPlayer *Player
}

type Round struct {
	teamAScore int
	teamBScore int
}

func newTable(instance *Instance) Table {
	var players = [4]*Player{}
	for i := range players {
		players[i] = &Player{
			instance: instance,
			seat:     i,
			team:     Team(i % 2),
			partner:  (i + 2) % 4,
		}
	}

	return Table{
		instance: instance,
		players:  players,
		state:    TableWaiting,
		trump: Trump{
			suit: -1,
		},
		play: Play{
			players: []*Player{},
		},
		totalScores: map[Team]int{},
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

func (t *Table) selectMode() {
	t.state = TableModeSelecting
	t.instance.Broadcast(map[string]string{"ACTION": "SELECTMODE"})
	clog.Debugf("(i:%s) selecting mode", t.instance.id)
}

func (t *Table) startGame() {
	t.instance.Broadcast(map[string]string{"ACTION": "GAMESTART"})
	clog.Debugf("(i:%s) game started", t.instance.id)
	t.startShuffle()
}

func (t *Table) startShuffle() {
	t.play = Play{}
	t.playCount = 0
	for _, p := range t.players {
		p.score = 0
		p.hand = []Card{}
		p.isTurn = false
	}

	// Reshuffles
	for {
		var deck = newDeck()

		clog.Printf("(i:%v) shuffling deck", t.instance.id)
		// Shuffle deck
		for i := range deck {
			j := rand.IntN(i + 1)
			deck[i], deck[j] = deck[j], deck[i]
		}

		for _, p := range t.players {
			p.hand = []Card{}
		}

		// Deal hands
		for i := range deck {
			t.players[i/13].hand = append(t.players[i/13].hand, deck[i])
		}

		var reshuffle = false

		// Check hand validity
		for _, p := range t.players {
			if p.isHandInvalid() {
				reshuffle = true
			}
		}

		if reshuffle {
			continue
		}
		break
	}

	// Sort hands
	for _, p := range t.players {
		p.sortHand()
	}

	for _, p := range t.players {
		p.client.writeJson(map[string]string{"ACTION": "DEAL", "CARDS": p.getHandString()})
	}

	// Announce to all
	t.instance.Broadcast(map[string]string{"ACTION": "OTHERDEAL", "COUNT": "13"})
	t.gameMode.onShuffleEnd(t)
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
	var trumpPlayer = t.trump.highestCaller

	if t.trump.beIstifada {
		trumpPlayer = t.trump.gobool
	}

	t.turn = trumpPlayer.seat
	trumpPlayer.isTurn = true

	t.playCount = 0

	var playableCards = ""
	if t.trump.suit == -1 {
		var _, trumps = trumpPlayer.getAvailableTrumps()
		playableCards = trumps
	} else {
		var _, cardsStr = trumpPlayer.getPlayableCards()
		playableCards = cardsStr
	}

	var prompt = map[string]string{"ACTION": "YOURPLAY", "PLAYABLE": playableCards}
	trumpPlayer.lastPrompt = prompt
	trumpPlayer.client.writeJson(prompt)
}

func (p *Player) sortHand() {
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

func (p *Player) recievePassedCards() {
	var seat = p.seat
	var cardsToRecieve = p.instance.table.roundPassedCards[seat]
	p.hand = append(p.hand, cardsToRecieve...)
	p.sortHand()
	p.state = PlayerRecievedCards
	p.client.writeJson(map[string]string{"ACTION": "DEAL", "CARDS": p.getHandString()})
}

func (p *Player) getPlayableCards() ([]Card, string) {
	// if p.state != PlayerPlaying {
	// 	return []Card{}, ""
	// }

	var cards = []Card{}
	if len(p.instance.table.play.cards) == 0 {
		cards = p.hand
	} else {
		for _, c := range p.hand {
			if c.suit == p.instance.table.play.cards[0].suit {
				cards = append(cards, c)
			}
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
			return true
		}
	}

	if faceCardCounts == 0 {
		return true
	}

	for i := range 4 {
		if suitCardCounts[Suit(i)] >= 8 {
			return true
		}
	}

	return false

}

func (p *Player) getAvailableTrumps() ([]Suit, string) {
	var str = ""
	var trumps = []Suit{}
	var suitCounts = map[Suit]int{}
	var trump = p.instance.table.trump
	var callOffset = 3
	if p == trump.gobool && trump.beIstifada {
		callOffset = 2
	}

	for _, c := range p.hand {
		suitCounts[c.suit] += 1
	}
	for i := range 4 {
		if suitCounts[Suit(i)]+callOffset <= p.client.instance.table.trump.highestCall {
			trumps = append(trumps, Suit(i))
		}
	}
	for _, c := range p.hand {
		if slices.Contains(trumps, c.suit) {
			if str != "" {
				str += ","
			}
			str += c.name
		}
	}

	return trumps, str
}
