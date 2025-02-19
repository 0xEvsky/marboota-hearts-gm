package main

import (
	"encoding/json"
	"log"
	"strconv"
)

func msgHandler(c *Client, rawMsg []byte) {
	var msg map[string]any
	err := json.Unmarshal(rawMsg, &msg)
	if err != nil {
		c.writeError("Invalid message: message may not be JSON")
		log.Println(err)
		return
	}

	switch msg["ACTION"] {
	case "AUTH":
		if c.isAuthed {
			c.writeError("Already authenticated")
			log.Println("Duplicated authentication, skipping")
			return
		}

		if msg["INSTANCEID"] == nil {
			c.writeError("Missing field: INSTANCEID")
			log.Println("Authentication request with missing 'instance ID' rejected")
			return
		}
		if msg["USERID"] == nil {
			c.writeError("Missing field: USERID")
			log.Println("Authentication request with missing 'user ID' rejected")
			return
		}
		if msg["USERNAME"] == nil {
			c.writeError("Missing field: USERNAME")
			log.Println("Authentication request with missing 'username' rejected")
			return
		}
		if msg["ICONURL"] == nil {
			c.writeError("Missing field: ICONURL")
			log.Println("Authentication request with missing 'icon URL' rejected")
			return
		}

		var instance = server.instances[msg["INSTANCEID"].(string)]

		if instance != nil && instance.clients[msg["USERID"].(string)] != nil {
			c.writeError("ID is already authenticated with different client")
			log.Println("Failed authentication, user ID is already authenticated with different client")
			return
		}

		if instance != nil {
			c.instance = joinInstance(c, msg["INSTANCEID"].(string))
		} else {
			c.instance = NewInstance(c, msg["INSTANCEID"].(string))
		}
		c.id = msg["USERID"].(string)
		c.name = msg["USERNAME"].(string)
		c.iconUrl = msg["ICONURL"].(string)
		c.state = ClientIdle
		c.isAuthed = true

		c.writeOk()

		c.broadcastToInstance(map[string]string{"ACTION": "JOIN", "USERID": c.id, "USERNAME": c.name, "ICONURL": c.iconUrl})
		for _, client := range c.instance.clients {
			if !client.isAuthed || client == c {
				continue
			}
			c.writeJson(map[string]string{"ACTION": "JOIN", "USERID": client.id, "USERNAME": client.name, "ICONURL": client.iconUrl, "SEAT": strconv.Itoa(client.seat)})
		}

		log.Println("Client authenticated")
	case "SWITCH":
		if !c.isAuthed {
			c.writeError("Forbidden: not authenticated")
			log.Println("Unauthenticated SWITCH request rejected")
			return
		}

		if msg["SEAT"] == nil {
			c.writeError("Missing field: SEAT")
			log.Println("SWITCH request with missing 'seat' rejected")
			return
		}
		var reqSeat, err = strconv.Atoi(msg["SEAT"].(string))
		if err != nil || reqSeat < 0 || reqSeat > 4 {
			c.writeError("Invalid field: SEAT number outside of range (0, 4)")
			log.Println("SWITCH request with invalid 'seat' rejected")
			return
		}

		if c.instance.table.players[reqSeat-1] != nil {
			c.writeError("Invalid request: SEAT is taken")
			log.Println("SWITCH request with taken 'seat' rejected")
			return
		}

		if reqSeat > 0 {
			c.table = &c.instance.table
		} else {
			c.table = nil
		}

		c.seat = reqSeat
		c.writeOk()
		c.broadcastToInstance(map[string]string{"ACTION": "SWITCH", "USERID": c.id, "SEAT": msg["SEAT"].(string)})
		log.Println("Client switched seats")
	default:
		c.writeError("Unknown action")
		log.Println("Unknown action skipped")
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
