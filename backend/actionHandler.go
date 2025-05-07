package main

import (
	"errors"
	"log"
	"slices"
	"strconv"
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
	var instance = server.instances[instanceId]
	server.mu.Unlock()

	if instance != nil && instance.clients[userId] != nil {
		return errors.New("ID is already authenticated with different client")
	}

	c.id = userId
	c.name = userName
	c.iconUrl = iconUrl
	c.state = ClientIdle
	c.isAuthed = true

	if instance != nil {
		c.instance = joinInstance(c, instanceId)
		log.Println("Client " + c.id + " joined existing instance")
	} else {
		c.instance = newInstance(c, instanceId)
		log.Println("New instance created (" + instanceId + ")")
	}

	c.writeOk()

	c.broadcastToMates(map[string]string{"ACTION": "JOIN", "USERID": c.id, "USERNAME": c.name, "ICONURL": c.iconUrl})
	for _, client := range c.instance.clients {
		if !client.isAuthed || client == c {
			continue
		}
		// Join catch-up
		c.writeJson(map[string]string{"ACTION": "JOIN", "USERID": client.id, "USERNAME": client.name, "ICONURL": client.iconUrl})
		// Seat catch-up
		if client.state == ClientSeated {
			c.writeJson(map[string]string{"ACTION": "SIT", "USERID": c.id, "SEAT": strconv.Itoa(client.player.seat)})
			// Ready catchup
			if client.player.state == PlayerReady {
				c.writeJson(map[string]string{"ACTION": "READY", "USERID": c.id})
			}
		}
		// TODO: Trump catchup
		// TODO: Play catchup
	}

	return nil
}

func seatClient(c *Client, seatStr string) error {
	var seat, err = strconv.Atoi(seatStr)
	if err != nil || seat < 1 || seat > 4 {
		return errors.New("invalid or missing seat (1-4)")
	}

	if c.state == ClientSeated && seat == c.player.seat {
		return errors.New("already seated")
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
		c.instance.table.startTrump()
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
	if c.instance.table.state != TableTrumping || c.instance.table.trump.isDone {
		return errors.New("table not in trump state")
	}

	// Check if it's player's turn
	if !p.isTurn {
		return errors.New("not player turn")
	}

	// Check if player is in caller list
	if !slices.Contains(c.instance.table.trump.callers, p) {
		return errors.New("player is not a valid caller")
	}

	// Check if player is passing
	if scoreStr == "PASS" {
		c.instance.table.trump.callers = slices.DeleteFunc(c.instance.table.trump.callers, func(caller *Player) bool {
			return caller == p
		})
	} else {
		var score, err = strconv.Atoi(scoreStr)
		// Check if score is between 7 and 13
		if err != nil || score < 7 || score > 13 {
			return errors.New("invalid or missing score (7-13/PASS)")
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

	if len(c.instance.table.trump.callers) <= 1 || c.instance.table.trump.highestCall >= 13 {
		// Toggle trump as done
		c.instance.table.trump.isDone = true
		// Ask for trump suit
		c.instance.table.trump.highestCaller.client.writeJson(map[string]string{"ACTION": "YOURTRUMPSUIT", "SCORE": strconv.Itoa(c.instance.table.trump.highestCall)})
	} else {
		// Advance turn
		p.isTurn = false
		for {
			c.instance.table.turn += 1
			c.instance.table.turn %= 4
			if slices.Contains(c.instance.table.trump.callers, c.instance.table.players[c.instance.table.turn]) {
				break
			}
		}
		c.instance.table.players[c.instance.table.turn].isTurn = true
		c.instance.table.players[c.instance.table.turn].client.writeJson(map[string]string{"ACTION": "YOURTRUMPCALL", "MINSCORE": strconv.Itoa(c.instance.table.trump.highestCall + 1)})
	}

	return nil
}

func endTrump(c *Client, suit string) error {
	// Check if client is trumping
	if c.state != ClientSeated || c.player.state != PlayerTrumping {
		return errors.New("client not playing")
	}

	// Check if table is trumping
	if c.instance.table.state != TableTrumping {
		return errors.New("table not in trump state")
	}

	// Check if trump is done
	if !c.instance.table.trump.isDone {
		return errors.New("table not waiting for trump suit yet")
	}

	// Check if request is from highest caller
	if c != c.instance.table.trump.highestCaller.client {
		return errors.New("player not highest caller")
	}

	// Check if suit is valid
	var suits = []string{"SPADES", "HEARTS", "CLUBS", "DIAMONDS"}
	if !slices.Contains(suits, suit) {
		return errors.New("invalid suit (SPADES, HEARTS, CLUBS, DIAMONDS)")
	}

	// Announce trump end with suit and score
	c.instance.Broadcast(map[string]string{"ACTION": "TRUMPEND", "SUIT": suit, "SCORE": strconv.Itoa(c.instance.table.trump.highestCall)})

	// Start play
	c.instance.table.startPlay()
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

	// Check if card is playable
	var playables, _ = c.player.getPlayableCards()

	if !slices.Contains(playables, card) {
		return errors.New("card not playable")
	}

	// Add cards to played cards
	c.instance.table.play.cards = append(c.instance.table.play.cards, card)

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

	// TODO: Check if play is complete

	c.instance.table.play.round += 1
	return nil
}
