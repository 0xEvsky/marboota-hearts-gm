package main

import (
	"log"
	"strconv"
)

func msgHandler(c *Client, msg map[string]string) {
	if msg["ACTION"] == "AUTH" {
		if c.isAuthed {
			c.writeError(msg["REQUESTID"], "already authenticated")
			log.Println("Duplicated authentication, skipping")
			return
		}

		if msg["INSTANCEID"] == "" {
			c.writeError(msg["REQUESTID"], "missing field: INSTANCEID")
			log.Println("Authentication request with missing 'instance ID' refused")
			return
		}
		if msg["USERID"] == "" {
			c.writeError(msg["REQUESTID"], "missing field: USERID")
			log.Println("Authentication request with missing 'user ID' refused")
			return
		}
		if msg["USERNAME"] == "" {
			c.writeError(msg["REQUESTID"], "missing field: USERNAME")
			log.Println("Authentication request with missing 'username' refused")
			return
		}
		if msg["ICONURL"] == "" {
			c.writeError(msg["REQUESTID"], "missing field: ICONURL")
			log.Println("Authentication request with missing 'icon URL' refused")
			return
		}

		var instance = server.instances[msg["INSTANCEID"]]

		if instance != nil && instance.clients[msg["USERID"]] != nil {
			c.writeError(msg["REQUESTID"], "ID is already authenticated with different client")
			log.Println("Failed authentication, user ID is already authenticated with different client")
			return
		}

		c.id = msg["USERID"]
		c.name = msg["USERNAME"]
		c.iconUrl = msg["ICONURL"]
		c.state = ClientIdle
		c.isAuthed = true

		if instance != nil {
			c.instance = joinInstance(c, msg["INSTANCEID"])
		} else {
			c.instance = newInstance(c, msg["INSTANCEID"])
		}

		c.writeOk(msg["REQUESTID"])

		c.broadcastToInstance(map[string]string{"ACTION": "JOIN", "USERID": c.id, "USERNAME": c.name, "ICONURL": c.iconUrl})
		for _, client := range c.instance.clients {
			if !client.isAuthed || client == c {
				continue
			}
			// Catch-up
			c.writeJson(map[string]string{"ACTION": "JOIN", "USERID": client.id, "USERNAME": client.name, "ICONURL": client.iconUrl})
			if client.state == ClientSeated {
				// Seat catch-up
				c.writeJson(map[string]string{"ACTION": "SIT", "USERID": c.id, "SEAT": strconv.Itoa(client.player.seat)})
			}
		}

		log.Println("Client authenticated")
		return
	}

	if !c.isAuthed {
		c.writeError(msg["REQUESTID"], "not authenticated")
		log.Println("Unauthenticated request refused")
		return
	}

	switch msg["ACTION"] {
	case "SIT":
		var seat, err = strconv.Atoi(msg["SEAT"])
		if err != nil || seat < 0 || seat > 4 {
			c.writeError(msg["REQUESTID"], "invalid or missing seat (1-4)")
			log.Println("SIT request with invalid or missing seat refused")
			return
		}

		if seat == 0 {
			if c.state == ClientIdle {
				c.writeOk(msg["REQUESTID"])
				log.Println("Unseated SIT request (seat=0) skipped")
				return
			}

			c.instance.table.unseatPlayer(c)
		}

		if seat > 0 {
			if c.state == ClientSeated && seat == c.player.seat {
				c.writeError(msg["REQUESTID"], "already seated")
				log.Println("Seated SIT request refused")
				return
			}

			err = c.instance.table.seatPlayer(c, seat)
			if err != nil {
				c.writeError(msg["REQUESTID"], err.Error())
				log.Println("SIT request with taken seat refused")
				return
			}
		}

		c.writeOk(msg["REQUESTID"])
		c.broadcastToInstance(map[string]string{"ACTION": "SIT", "USERID": c.id, "SEAT": msg["SEAT"]})
		log.Println("SIT request accepted")
		// TODO: If game was already running, show player their hand

	default:
		c.writeError(msg["REQUESTID"], "unknown or missing action")
		log.Println("Unknown or missing action skipped")
		return
	}
}
