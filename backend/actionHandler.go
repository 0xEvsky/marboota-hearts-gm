package main

import (
	"errors"
	"slices"
	"strconv"

	"github.com/OmarQurashi868/marboota/backend/clog"
)

func authClient(c *Client, instanceId, userId, userName, iconUrl string) error {
	if c.isAuthed {
		return errors.New("already authenticated")
	}

	if instanceId == "" {
		return errors.New("missing field: INSTANCEID")
	}
	if userId == "" {
		return errors.New("missing field: USERID")
	}
	if userName == "" {
		return errors.New("missing field: USERNAME")
	}
	if iconUrl == "" {
		return errors.New("missing field: ICONURL")
	}

	server.mu.Lock()
	var instance, exists = server.instances[instanceId]
	server.mu.Unlock()

	if exists {
		instance.mu.Lock()
		var _, cExists = instance.clients[userId]
		instance.mu.Unlock()
		if cExists {
			return errors.New("ID is already authenticated with different client")
		}
	}

	c.id = userId
	c.name = userName
	c.iconUrl = iconUrl
	c.state = ClientIdle
	c.isAuthed = true

	if exists {
		c.instance = joinInstance(c, instanceId)
		clog.Debugf("(server) (c:%s) joined existing instance (i:%s)", c.id, instanceId)
	} else {
		c.instance = newInstance(c, instanceId)
		clog.Debugf("(server) (c:%s) new instance created (i:%s)", c.id, instanceId)
	}

	c.writeOk()

	c.broadcastToMates(map[string]string{"ACTION": "JOIN", "USERID": c.id, "USERNAME": c.name, "ICONURL": c.iconUrl})
	c.instance.mu.Lock()
	for _, client := range c.instance.clients {
		if !client.isAuthed || client.id == userId {
			continue
		}
		// Join catch-up
		c.writeJson(map[string]string{"ACTION": "JOIN", "USERID": client.id, "USERNAME": client.name, "ICONURL": client.iconUrl})

		// Seat catch-up
		if client.state == ClientSeated {
			c.writeJson(map[string]string{"ACTION": "SIT", "USERID": client.id, "SEAT": strconv.Itoa(client.player.seat)})
			// Ready catchup
			if client.player.state == PlayerReady {
				c.writeJson(map[string]string{"ACTION": "READY", "USERID": client.id})
			}
		}

		// TODO: Trump catchup

		// TODO: Play catchup
	}
	c.instance.mu.Unlock()

	return nil
}

func seatClient(c *Client, seatStr string) error {
	var seat, err = strconv.Atoi(seatStr)
	if err != nil || seat < 0 || seat > 3 {
		return errors.New("invalid or missing seat (0-3)")
	}

	if c.state == ClientSeated {
		if seat == c.player.seat {
			return errors.New("already seated")
		}
		if c.player.state == PlayerReady {
			return errors.New("can't change seats when ready")
		}
		if c.player.state == PlayerTrumping || c.player.state == PlayerPlaying {
			return errors.New("game already started")
		}
	}

	err = c.instance.table.seatPlayer(c, seat)
	if err != nil {
		return err
	}

	// TODO: If game was already running, show player their hand

	c.broadcastToMates(map[string]string{"ACTION": "SIT", "USERID": c.id, "SEAT": seatStr})
	return nil
}

func unseatClient(c *Client) error {
	if c.instance.table.state != TableWaiting {
		return errors.New("game has started")
	}

	if c.state == ClientIdle {
		return errors.New("already unseated")
	}

	unsetReady(c)
	c.instance.table.unseatPlayer(c)

	c.broadcastToMates(map[string]string{"ACTION": "UNSIT", "USERID": c.id})
	return nil
}

func setReady(c *Client) error {
	if c.instance.table.state != TableWaiting {
		return errors.New("game has started")
	}

	if c.state != ClientSeated {
		return errors.New("not seated")
	}

	if c.player.state == PlayerReady {
		return errors.New("already ready")
	}

	if c.player.state != PlayerWaiting {
		return errors.New("not waiting to be ready")
	}

	c.player.state = PlayerReady

	c.broadcastToMates(map[string]string{"ACTION": "READY", "USERID": c.id})

	// Check if all players are ready
	if c.instance.table.isEveryoneReady() {
		c.instance.table.startGame()
	}

	return nil
}

func unsetReady(c *Client) error {
	if c.state != ClientSeated {
		return errors.New("not seated")
	}

	if c.player.state != PlayerReady {
		return errors.New("not ready")
	}

	c.player.state = PlayerWaiting
	c.broadcastToMates(map[string]string{"ACTION": "UNREADY", "USERID": c.id})
	return nil
}

func advanceTrump(c *Client, scoreStr string) error {
	var p = c.player

	// Check if client is trumping
	if c.state != ClientSeated || c.player.state != PlayerTrumping {
		return errors.New("client not playing")
	}

	// Check if game is in trump state
	if c.instance.table.state != TableTrumping {
		return errors.New("table not in trump state")
	}

	// Check if it's player's turn
	if !p.isTurn {
		return errors.New("not player turn")
	}

	// Check if player is not passing
	if scoreStr != "PASS" {
		var score, err = strconv.Atoi(scoreStr)
		// Check if score is between 7 and 13
		var maxScore = 13
		if c.instance.table.turn == c.instance.table.turnOffset {
			maxScore = 11
		}
		if err != nil || score < 7 || score > maxScore {
			return errors.New("invalid or missing score (7-" + strconv.Itoa(maxScore) + "/PASS)")
		}

		// Check if score is higher than the current highest call
		if score <= c.instance.table.trump.highestCall {
			return errors.New("trump call must be higher than the highest one")
		} else {
			c.instance.table.trump.highestCall = score
			c.instance.table.trump.highestCaller = p
		}
	}

	c.broadcastToMates(map[string]string{"ACTION": "TRUMPCALL", "USERID": c.id, "SCORE": scoreStr})

	if (c.instance.table.turn+1)%4 == c.instance.table.turnOffset || c.instance.table.trump.highestCall >= 13 {
		endTrump(c.instance)
	} else {
		// Advance turn
		p.isTurn = false
		c.instance.table.turn += 1
		c.instance.table.turn %= 4

		c.instance.table.players[c.instance.table.turn].isTurn = true
		c.instance.table.players[c.instance.table.turn].client.writeJson(map[string]string{"ACTION": "YOURTRUMPCALL", "MINSCORE": strconv.Itoa(c.instance.table.trump.highestCall + 1), "MAXSCORE": "13"})
	}

	clog.Debugf("(i:%s) (c:%s) trump advanced, called %s", c.instance.id, c.id, scoreStr)
	return nil
}

func endTrump(i *Instance) error {
	// // Announce trump end with suit and score
	// c.instance.Broadcast(map[string]string{"ACTION": "TRUMPEND", "SUIT": suit, "SCORE": strconv.Itoa(c.instance.table.trump.highestCall)})

	// Start play
	i.table.startPlay()
	clog.Debugf("(i:%s) trump ended (%v)", i.id, i.table.trump.highestCall)
	return nil
}

func advancePlay(c *Client, cardStr string) error {
	// Advance play
	// Check if client is playing
	if c.state != ClientSeated || c.player.state != PlayerPlaying {
		return errors.New("client not playing")
	}

	// Check if table is playing
	if c.instance.table.state != TablePlaying {
		return errors.New("table not playing")
	}

	// Check if it's client's turn
	if c.player != c.instance.table.players[c.instance.table.turn] {
		return errors.New("not client turn")
	}

	// Check if card is valid
	var card, err = getCardByName(cardStr)
	if err != nil {
		return err
	}

	if len(c.instance.table.play.cards) == 0 {
		if c.instance.table.trump.suit == -1 {
			// Check if first play ever is a valid trump
			var playables, _ = c.player.getAvailableTrumps()

			if !slices.Contains(playables, card.suit) {
				return errors.New("invalid suit for trump")
			}

			// Set trump
			c.instance.table.trump.suit = card.suit
		}
	} else {
		// Check if card is playable
		var playables, _ = c.player.getPlayableCards()

		if !slices.Contains(playables, card) {
			return errors.New("card not playable")
		}
	}

	// Add cards to played cards
	c.instance.table.play.cards = append(c.instance.table.play.cards, card)

	// Remove played card from hand
	c.player.hand = slices.DeleteFunc(c.player.hand, func(ec Card) bool {
		return ec.name == card.name
	})

	c.instance.Broadcast(map[string]string{"ACTION": "PLAY", "USERID": c.id, "CARD": card.name})
	clog.Debugf("(i:%s) (c:%s) card played (%s)", c.instance.id, c.id, card.name)

	// Winning card deciding logic
	if len(c.instance.table.play.cards) == 1 {
		c.instance.table.play.curWinCard = card
		c.instance.table.play.curWinPlayer = c.player
	} else {
		if card.suit == c.instance.table.play.curWinCard.suit {
			if card.value > c.instance.table.play.curWinCard.value {
				c.instance.table.play.curWinCard = card
				c.instance.table.play.curWinPlayer = c.player
			}
		} else if card.suit == c.instance.table.trump.suit {
			c.instance.table.play.curWinCard = card
			c.instance.table.play.curWinPlayer = c.player
		}
	}

	// Check if play is complete
	if len(c.instance.table.play.cards) == 4 {
		// Add score & advance round
		c.instance.table.play.curWinPlayer.score += 1
		c.instance.table.playCount += 1

		// Announce play round end
		c.instance.Broadcast(map[string]string{"ACTION": "PLAYEND", "WINNERID": c.instance.table.play.curWinPlayer.client.id})

		// End round if playRound == 13
		if c.instance.table.playCount == 13 {
			endRound(c.instance)
			return nil
		}

		// Set turn to winner
		c.player.isTurn = false
		c.instance.table.turn = c.instance.table.play.curWinPlayer.seat
		c.instance.table.players[c.instance.table.turn].isTurn = true

		// Wipe play
		c.instance.table.play = Play{}

		// Announce new turn
		c.instance.table.players[c.instance.table.turn].client.writeJson(map[string]string{"ACTION": "YOURPLAY", "PLAYABLE": c.instance.table.players[c.instance.table.turn].getHandString()})

		clog.Debugf("(i:%s) play ended, winner (%s)", c.instance.id, c.id)
	} else {
		// Advance turn
		c.player.isTurn = false
		c.instance.table.turn += 1
		c.instance.table.turn %= 4
		c.instance.table.players[c.instance.table.turn].isTurn = true

		var _, nextPlayables = c.instance.table.players[c.instance.table.turn].getPlayableCards()
		c.instance.table.players[c.instance.table.turn].client.writeJson(map[string]string{"ACTION": "YOURPLAY", "PLAYABLE": nextPlayables})

		clog.Debugf("(i:%s) play advanced", c.instance.id)
	}

	return nil
}

func endRound(i *Instance) {
	// Check score of trumpcaller + their partner
	var teamScores = map[Team]int{}
	teamScores[TeamA] = i.table.players[0].score + i.table.players[2].score
	teamScores[TeamB] = i.table.players[1].score + i.table.players[3].score

	var lastGameRound = Round{
		teamAScore: teamScores[TeamA],
		teamBScore: teamScores[TeamB],
	}
	i.table.rounds = append(i.table.rounds, lastGameRound)

	i.Broadcast(map[string]string{"ACTION": "ROUNDEND",
		"TEAMASCORE": strconv.Itoa(teamScores[TeamA]),
		"TEAMBSCORE": strconv.Itoa(teamScores[TeamB])})

	clog.Debugf("(i:%s) round ended (scoreA:%v, scoreB:%v)",
		i.id,
		teamScores[TeamA],
		teamScores[TeamB])

	// Check if one of the teams have 13 score -> endGame (seek)
	if teamScores[TeamA] == 13 {
		endGame(i, TeamA)
		return
	}
	if teamScores[TeamB] == 13 {
		endGame(i, TeamB)
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

	i.Broadcast(map[string]string{"ACTION": "TOTALSCORE",
		"TEAMASCORE": strconv.Itoa(i.table.totalScores[TeamA]),
		"TEAMBSCORE": strconv.Itoa(i.table.totalScores[TeamB])})

	clog.Debugf("(i:%s) total team scores (scoreA:%v, scoreB:%v)",
		i.id,
		i.table.totalScores[TeamA],
		i.table.totalScores[TeamB])

	// Check total scores for winning game
	if i.table.totalScores[TeamA] >= 25 || i.table.totalScores[TeamB] <= -25 {
		endGame(i, TeamA)
		return
	}
	if i.table.totalScores[TeamB] >= 25 || i.table.totalScores[TeamA] <= -25 {
		endGame(i, TeamB)
		return
	}

	// Start new round
	i.table.state = TableWaiting
	i.table.turnOffset += 1
	i.table.turnOffset %= 4

	i.table.startTrump()
}

func endGame(i *Instance, winner Team) {
	// Announce game end
	var winner1 = i.table.players[winner]
	var winner2 = winner1.partner
	i.Broadcast(map[string]string{"ACTION": "GAMEEND", "WINNER1ID": winner1.client.id, "WINNER2ID": winner2.client.id})
	clog.Debugf("(i:%s) game ended (scoreA:%v, scoreB:%v)", i.id, i.table.totalScores[TeamA], i.table.totalScores[TeamB])

	// Save current players
	var curPlayers = [4]*Player{}
	for j := range 4 {
		curPlayers[j] = i.table.players[j]
	}

	// Reset table newTable
	i.table = newTable()
	for j := range 4 {
		i.table.seatPlayer(curPlayers[j].client, j)
	}
}
