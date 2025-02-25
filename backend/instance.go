package main

type Instance struct {
	id      string
	clients map[string]*Client // key is userid
	table   Table
}

func joinInstance(c *Client, id string) *Instance {
	if instance := server.instances[id]; instance != nil {
		instance.clients[c.id] = c
		return instance
	}

	return newInstance(c, id)
}

func newInstance(c *Client, id string) *Instance {
	var newInstance = Instance{
		id:      id,
		clients: map[string]*Client{c.id: c},
		table:   newTable(),
	}

	// TODO: Investigate mutex
	server.instances[id] = &newInstance

	return &newInstance
}
