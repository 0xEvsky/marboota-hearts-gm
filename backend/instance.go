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

	return NewInstance(id, c)
}

func NewInstance(id string, c *Client) *Instance {
	var newInstance = Instance{
		id:      id,
		clients: []*Client{c},
		table:   NewTable(),
	}
	server.instances = append(server.instances, &newInstance)
	log.Println("New instance created, id: " + id)
	return &newInstance
}

func isInstanceExists(id string) *Instance {
	for _, instance := range server.instances {
		if instance.id == id {
			return instance
		}
	}
	return nil
}
