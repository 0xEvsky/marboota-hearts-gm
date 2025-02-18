package main

import (
	"encoding/json"
	"log"
)

func msgHandler(c *Client, rawMsg []byte) {
	var msg map[string]any
	err := json.Unmarshal(rawMsg, &msg)
	if err != nil {
		c.write([]byte(err.Error()))
		log.Println(err)
		return
	}

	switch msg["ACTION"] {
	case "CONNECT":
		// TODO: check if already authenticated
		// TODO: check if all required fields are present
		var newInstance = joinInstance(c, msg["INSTANCEID"].(string))
		c.instance = newInstance
		c.id = msg["USERID"].(string)
		c.name = msg["USERNAME"].(string)
		c.iconUrl = msg["ICONURL"].(string)
		c.state = ClientIdle

		r, err := json.Marshal(map[string]string{"ACTION": "OK"})
		if err != nil {
			c.write([]byte(err.Error()))
			log.Println(err)
			return
		}
		c.write(r)
		log.Println("Client authenticated")
	default:

	}
}
