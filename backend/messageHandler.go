package main

import (
	"encoding/json"
	"log"
)

func msgHandler(c *Client, rawMsg []byte) {
	var msg map[string]string
	err := json.Unmarshal(rawMsg, &msg)
	if err != nil {
		c.writeError(err.Error())
		log.Println(err)
	}

	c.requestId = msg["REQUESTID"]

	if msg["ACTION"] == "AUTH" {
		err := authClient(c, msg["INSTANCEID"], msg["USERID"], msg["USERNAME"], msg["ICONURL"])
		if err != nil {
			c.writeError(err.Error())
			log.Printf("AUTH request refused: %s\n", err)
			return
		}

		log.Println("AUTH request accepted")
		return
	}

	if !c.isAuthed {
		var err = "not authenticated"
		c.writeError(err)
		log.Printf("Request refused: %s\n", err)
		return
	}

	switch msg["ACTION"] {

	case "SIT":
		err := seatClient(c, msg["SEAT"])
		if err != nil {
			c.writeError(err.Error())
			log.Printf("SIT request refused: %s\n", err)
			return
		}
		c.writeOk()
		log.Println("SIT request accepted")

	case "UNSIT":
		err := unseatClient(c)
		if err != nil {
			c.writeError(err.Error())
			log.Printf("UNSIT request refused: %s\n", err)
			return
		}
		c.writeOk()
		log.Println("UNSIT request accepted")

	case "READY":
		err := setReady(c)
		if err != nil {
			c.writeError(err.Error())
			log.Printf("READY request refused: %s\n", err)
			return
		}
		c.writeOk()
		log.Println("READY request accepted")

	case "UNREADY":
		err := unsetReady(c)
		if err != nil {
			c.writeError(err.Error())
			log.Printf("UNREADY request refused: %s\n", err)
			return
		}
		c.writeOk()
		log.Println("UNREADY request accepted")

	// TODO: case "TRUMPCALL"

	default:
		c.writeError("unknown or missing action")
		log.Println("Unknown or missing action skipped")
		return
	}
}
