package main

import (
	"errors"
	"log"
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

	var instance = server.instances[instanceId]

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
		log.Println("Joined existing instance")
	} else {
		c.instance = newInstance(c, instanceId)
		log.Println("New instance created")
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
		// TODO: Table catchup (game state)
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
		c.instance.table.trumpStart()
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

// TODO: advanceTrump
