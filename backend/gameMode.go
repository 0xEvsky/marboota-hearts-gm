package main

import (
	"strconv"

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
	calcPlayScore  func(t *Table) int
	modeRoundScore RoundScore
	onRoundEnd     func(i *Instance)
}

type RoundScore interface {
	addScore(i *Instance)
	getScores(i *Instance) map[string]int
}

type TeamRound struct {
	teamAScore int
	teamBScore int
}

type FFARound struct {
	scores map[int]int // seat -> score
}

func (r *TeamRound) addScore(i *Instance) {
	r.teamAScore = i.table.players[0].score + i.table.players[2].score
	r.teamBScore = i.table.players[1].score + i.table.players[3].score
}

func (r *TeamRound) getScores(i *Instance) map[string]int {
	return map[string]int{
		"TEAMASCORE": r.teamAScore,
		"TEAMBSCORE": r.teamBScore,
	}
}

func (r *FFARound) addScore(i *Instance) {
	if r.scores == nil {
		r.scores = make(map[int]int)
	}

	for _, p := range i.table.players {
		r.scores[p.seat] = p.score
	}
}
func (r *FFARound) getScores(i *Instance) map[string]int {
	return map[string]int{
		"0": r.scores[0],
		"1": r.scores[1],
		"2": r.scores[2],
		"3": r.scores[3],
	}
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
	calcPlayScore:  calculatePlayScore,
	modeRoundScore: &FFARound{},
	onRoundEnd:     onFFARoundEnd,
}

var WhistMode = GameMode{
	id:             WhistModeID,
	penaltyConfig:  nil,
	onShuffleEnd:   startTrump,
	calcPlayScore:  calculatePlayScore,
	modeRoundScore: &TeamRound{},
	onRoundEnd:     onTeamRoundEnd,
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

func onTeamRoundEnd(i *Instance) {
	var teamScores = map[Team]int{}
	teamScores[TeamA] = i.table.gameMode.modeRoundScore.getScores(i)["TEAMASCORE"]
	teamScores[TeamB] = i.table.gameMode.modeRoundScore.getScores(i)["TEAMBSCORE"]

	i.Broadcast(map[string]string{
		"ACTION":     "TEAMROUNDEND",
		"TEAMASCORE": strconv.Itoa(teamScores[TeamA]),
		"TEAMBSCORE": strconv.Itoa(teamScores[TeamB])},
	)

	clog.Debugf("(i:%s) round ended (scoreA:%v, scoreB:%v)",
		i.id,
		teamScores[TeamA],
		teamScores[TeamB])

	// Check if one of the teams have 13 score -> endGame (seek)
	if teamScores[TeamA] == 13 {
		teamGameEnd(i, TeamA)
		return
	}
	if teamScores[TeamB] == 13 {
		teamGameEnd(i, TeamB)
		return
	}

	// Calculate scores for round
	var trumpCallerTeam = i.table.trump.highestCaller.team
	if teamScores[trumpCallerTeam] >= i.table.trump.highestCall {
		i.table.totalScores[trumpCallerTeam] += teamScores[trumpCallerTeam]
	} else {
		var otherTeam = (trumpCallerTeam + 1) % 2
		i.table.totalScores[otherTeam] += teamScores[otherTeam]
		i.table.totalScores[trumpCallerTeam] -= i.table.trump.highestCall
	}

	i.Broadcast(map[string]string{"ACTION": "TEAMTOTALSCORE",
		"TEAMASCORE": strconv.Itoa(i.table.totalScores[TeamA]),
		"TEAMBSCORE": strconv.Itoa(i.table.totalScores[TeamB])})

	clog.Debugf("(i:%s) total team scores (scoreA:%v, scoreB:%v)",
		i.id,
		i.table.totalScores[TeamA],
		i.table.totalScores[TeamB])

	// Check total scores for winning game
	if i.table.totalScores[TeamA] >= 25 || i.table.totalScores[TeamB] <= -25 {
		teamGameEnd(i, TeamA)
		return
	}
	if i.table.totalScores[TeamB] >= 25 || i.table.totalScores[TeamA] <= -25 {
		teamGameEnd(i, TeamB)
		return
	}

	// Start new round
	i.table.state = TableWaiting
	i.table.turnOffset += 1
	i.table.turnOffset %= 4

	i.table.startShuffle()
}

func teamGameEnd(i *Instance, winner Team) {
	// Announce game end
	var winner1 = i.table.players[winner]
	var winner2 = i.table.players[winner1.partner]
	i.Broadcast(map[string]string{"ACTION": "TEAMGAMEEND", "WINNER1ID": winner1.client.id, "WINNER2ID": winner2.client.id})
	clog.Debugf("(i:%s) game ended (scoreA:%v, scoreB:%v)", i.id, i.table.totalScores[TeamA], i.table.totalScores[TeamB])

	// Save current players
	var curPlayers = [4]*Player{}
	for j := range 4 {
		curPlayers[j] = i.table.players[j]
	}

	// Reset table newTable
	i.table = newTable(i)
	for j := range 4 {
		i.table.seatPlayer(curPlayers[j].client, j)
	}
}

func onFFARoundEnd(i *Instance) {
	var lastRoundScore = map[int]int{
		0: i.table.rounds[len(i.table.rounds)-1].score.getScores(i)["0"],
		1: i.table.rounds[len(i.table.rounds)-1].score.getScores(i)["1"],
		2: i.table.rounds[len(i.table.rounds)-1].score.getScores(i)["2"],
		3: i.table.rounds[len(i.table.rounds)-1].score.getScores(i)["3"],
	}

	i.Broadcast(map[string]string{
		"ACTION": "FFAROUNDEND",
		"0":      strconv.Itoa(lastRoundScore[0]),
		"1":      strconv.Itoa(lastRoundScore[1]),
		"2":      strconv.Itoa(lastRoundScore[2]),
		"3":      strconv.Itoa(lastRoundScore[3]),
	})

	clog.Debugf("(i:%s) round ended (player0:%v, player1:%v, player2:%v, player3:%v)",
		i.id,
		lastRoundScore[0],
		lastRoundScore[1],
		lastRoundScore[2],
		lastRoundScore[3],
	)

	var gameScores = i.table.gameMode.modeRoundScore.getScores(i)

	i.Broadcast(map[string]string{"ACTION": "FFATOTALSCORE",
		"0": strconv.Itoa(gameScores["0"]),
		"1": strconv.Itoa(gameScores["1"]),
		"2": strconv.Itoa(gameScores["2"]),
		"3": strconv.Itoa(gameScores["3"]),
	})

	clog.Debugf("(i:%s) round ended (player0:%v, player1:%v, player2:%v, player3:%v)",
		i.id,
		gameScores["0"],
		gameScores["1"],
		gameScores["2"],
		gameScores["3"],
	)

	for _, val := range gameScores {
		if val < -35 {
			endFFAgame(i, gameScores)
			return
		}
	}

	i.table.roundPassedCards = make([][]Card, 4)
	i.table.startShuffle()
}

func endFFAgame(i *Instance, score map[string]int) {
	var highest_score = -100
	var winner_seat = 0
	for key, val := range score {
		if val > highest_score {
			winner_seat, _ = strconv.Atoi(key)
		}
	}

	var winner = i.table.players[winner_seat]
	i.Broadcast(map[string]string{
		"ACTION":   "FFAGAMEEND",
		"WINNERID": winner.client.id,
	})
}
