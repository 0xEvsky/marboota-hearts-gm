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
		if c.id != "" {
			c.writeError("Already authenticated")
			log.Println("Duplicated authentication, skipping")
			return
		}

		// TODO: check if all required fields are present

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

		c.writeOk()
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
