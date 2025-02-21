package main

import (
	"encoding/json"
	"log"
	"strconv"
)

func msgHandler(c *Client, msg map[string]string) {
	if msg["ACTION"] == "AUTH" {
		if c.isAuthed {
			c.writeError("Already authenticated")
			log.Println("Duplicated authentication, skipping")
			return
		}

		if msg["INSTANCEID"] == "" {
			c.writeError("Missing field: INSTANCEID")
			log.Println("Authentication request with missing 'instance ID' rejected")
			return
		}
		if msg["USERID"] == "" {
			c.writeError("Missing field: USERID")
			log.Println("Authentication request with missing 'user ID' rejected")
			return
		}
		if msg["USERNAME"] == "" {
			c.writeError("Missing field: USERNAME")
			log.Println("Authentication request with missing 'username' rejected")
			return
		}
		if msg["ICONURL"] == "" {
			c.writeError("Missing field: ICONURL")
			log.Println("Authentication request with missing 'icon URL' rejected")
			return
		}

		var instance = server.instances[msg["INSTANCEID"]]

		if instance != nil && instance.clients[msg["USERID"]] != nil {
			c.writeError("ID is already authenticated with different client")
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

		c.writeOk()

		c.broadcastToInstance(map[string]string{"ACTION": "JOIN", "USERID": c.id, "USERNAME": c.name, "ICONURL": c.iconUrl})
		for _, client := range c.instance.clients {
			if !client.isAuthed || client == c {
				continue
			}
			c.writeJson(map[string]string{"ACTION": "JOIN", "USERID": client.id, "USERNAME": client.name, "ICONURL": client.iconUrl, "SEAT": strconv.Itoa(client.seat)})
		}

		log.Println("Client authenticated")
		return
	}

	if !c.isAuthed {
		c.writeError("Forbidden: Not authenticated")
		log.Println("Unauthenticated SWITCH request rejected")
		return
	}

	switch msg["ACTION"] {
	case "SIT":
		if c.table != nil {
			c.writeError("Already seated")
			log.Println("Seated SIT request rejected")
			return
		}

		err := c.instance.table.seatPlayer(c)
		if err != nil {
			c.writeError("Table is full")
			log.Println("Full table SIT request rejected")
			return
		}
		c.writeOk()
		c.broadcastToInstance(map[string]string{"ACTION": "SIT", "USERID": c.id, "SEAT": strconv.Itoa(c.seat)})
		// TODO: If game was already running, show client their cards

	case "UNSIT":
		if c.table == nil {
			c.writeError("Not seated")
			log.Println("Unseated UNSIT request rejected")
			return
		}

		c.table.unseatPlayer(c)
		c.writeOk()
		c.broadcastToInstance(map[string]string{"ACTION": "UNSIT", "USERID": c.id})

	case "SWITCH":
		if c.table == nil {
			c.writeError("Client is not seated")
			log.Println("Unseated SWITCH request rejected")
			return
		}

		// TODO: implement xd

		c.writeOk()
		c.broadcastToInstance(map[string]string{"ACTION": "SWITCH", "USERID": c.id, "SEAT": msg["SEAT"]})
		log.Println("Client switched seats")

	default:
		c.writeError("Unknown or missing action")
		log.Println("Unknown or missing action skipped")
		return
	}
}

func toJson(msg map[string]string) []byte {
	r, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return []byte("{\"ACTION\":\"ERROR\",\"MESSAGE\":\"Server error in JSON marshalling\"}")
	}
	return r
}
