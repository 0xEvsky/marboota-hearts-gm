package main

import (
	"github.com/OmarQurashi868/marboota/backend/clog"
)

type GameModeID int

const (
	WhistModeID GameModeID = iota
	HeartsModeID
)

type PenaltyConfig struct {
	suit            Suit
	poisonedCard    string
	poisonedPenalty int
}

type GameMode struct {
	id             GameModeID
	penaltyConfig  *PenaltyConfig
	onShuffleEnd   func(t *Table)
	calcRoundScore func(t *Table) int
}

var heartsPenaltyConfig = PenaltyConfig{
	suit:            Hearts,
	poisonedCard:    "S_12",
	poisonedPenalty: 7,
}

var HeartsMode = GameMode{
	id:             HeartsModeID,
	penaltyConfig:  &heartsPenaltyConfig,
	onShuffleEnd:   startCardsPassing,
	calcRoundScore: calculatePlayScore,
}

var WhistMode = GameMode{
	id:            WhistModeID,
	penaltyConfig: nil,
	onShuffleEnd:  startTrump,
}

func startTrump(t *Table) {
	t.trump = Trump{
		players: []*Player{},
		calls:   []string{},
		suit:    -1,
	}
	for _, p := range t.players {
		p.state = PlayerTrumping
	}
	t.withTrump = true

	t.instance.Broadcast(map[string]string{"ACTION": "TRUMPSTART"})
	clog.Debugf("(i:%s) trump started", t.instance.id)

	t.state = TableTrumping
	t.players[t.turn].isTurn = false
	t.turn = t.turnOffset
	t.players[t.turn].isTurn = true

	var prompt = map[string]string{"ACTION": "YOURTRUMPCALL", "MINSCORE": "7", "MAXSCORE": "11"}
	t.players[t.turn].lastPrompt = prompt
	t.players[t.turn].client.writeJson(prompt)
}

func startCardsPassing(t *Table) {
	t.state = TableCardsPassing
	for _, p := range t.players {
		p.state = PlayerPassingCards
	}

	t.roundPassedCards = make([][]Card, 4)
	t.instance.Broadcast(map[string]string{"ACTION": "PASSCARDS"})
}

func startPlayNoTrump(t *Table) {
	// StartPlay
	// Change table state to playing
	t.state = TablePlaying

	t.instance.Broadcast(map[string]string{"ACTION": "PLAYSTART"})
	clog.Debugf("(i:%s) play started", t.instance.id)

	for _, p := range t.players {
		p.state = PlayerPlaying
	}

	t.players[t.turn].isTurn = false
	var firstPlayer = t.players[0]

	t.turn = firstPlayer.seat
	firstPlayer.isTurn = true

	t.playCount = 0

	var _, playableCards = firstPlayer.getPlayableCards()

	var prompt = map[string]string{"ACTION": "YOURPLAY", "PLAYABLE": playableCards}
	firstPlayer.client.writeJson(prompt)
}

func calculatePlayScore(t *Table) int {
	if t.gameMode.penaltyConfig != nil {
		var penaltySuit = t.gameMode.penaltyConfig.suit
		var poisonedCard = t.gameMode.penaltyConfig.poisonedCard
		var poisonedPenalty = t.gameMode.penaltyConfig.poisonedPenalty
		for _, c := range t.play.cards {
			if c.suit == penaltySuit {
				t.play.curWinPlayer.score -= 1
			}

			if c.name == poisonedCard {
				t.play.curWinPlayer.score -= poisonedPenalty
			}
		}
	} else {
		t.play.curWinPlayer.score += 1
	}

	return t.play.curWinPlayer.score
}
