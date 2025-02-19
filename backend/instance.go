package main

import "log"

type Instance struct {
	id      string
	clients map[string]*Client // key is userid
	table   Table
}

func joinInstance(c *Client, id string) *Instance {
	if instance := server.instances[id]; instance != nil {
		instance.clients[c.id] = c
		log.Println("Joined existing instance")
		return instance
	}

	return NewInstance(c, id)
}

func NewInstance(c *Client, id string) *Instance {
	var newInstance = Instance{
		id:      id,
		clients: map[string]*Client{id: c},
		table:   NewTable(),
	}
	server.instances[id] = &newInstance
	log.Println("New instance created")
	return &newInstance
}
