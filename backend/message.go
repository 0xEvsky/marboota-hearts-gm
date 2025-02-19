package main

import (
	"encoding/json"
	"log"
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

		var instance = server.getInstanceById(msg["INSTANCEID"].(string))

		if instance != nil && instance.getClientById(msg["USERID"].(string)) != nil {
			c.writeError("ID is already authenticated with different client")
			log.Println("Failed authentication, ID is already authenticated with different client")
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
			c.writeJson(map[string]string{"ACTION": "JOIN", "USERID": client.id, "USERNAME": client.name, "ICONURL": client.iconUrl})
		}

		log.Println("Client authenticated")
	default:

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
