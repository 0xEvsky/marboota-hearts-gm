package main

import "log"

type Instance struct {
	id      string
	clients []*Client
	table   Table
}

func joinInstance(c *Client, id string) *Instance {
	if instance := server.getInstanceById(id); instance != nil {
		instance.clients = append(instance.clients, c)
		log.Println("Joined existing instance")
		return instance
	}

	return NewInstance(c, id)
}

func NewInstance(c *Client, id string) *Instance {
	var newInstance = Instance{
		id:      id,
		clients: []*Client{c},
		table:   NewTable(),
	}
	server.instances = append(server.instances, &newInstance)
	log.Println("New instance created")
	return &newInstance
}

func (i *Instance) getClientById(id string) *Client {
	for _, client := range i.clients {
		if client.id == id {
			return client
		}
	}
	return nil
}
