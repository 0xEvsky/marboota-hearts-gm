package main

import "log"

type Instance struct {
	id      string
	clients []*Client
	table   Table
}

func joinInstance(c *Client, id string) *Instance {
	if instance := isInstanceExists(id); instance != nil {
		instance.clients = append(instance.clients, c)
		log.Println("Joined existing instance, id: " + id)
		return instance
	}

	var newInstance = NewInstance(id, c)
	log.Println("New instance created, id: " + id)
	return newInstance
}

func NewInstance(id string, c *Client) *Instance {
	return &Instance{
		id:      id,
		clients: []*Client{c},
		table:   NewTable(),
	}
}

func isInstanceExists(id string) *Instance {
	for _, instance := range server.instances {
		if instance.id == id {
			return instance
		}
	}
	return nil
}
