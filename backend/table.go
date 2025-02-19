package main

type Table struct {
	players    [4]*Client
	spectators []*Client
}

func NewTable() Table {
	return Table{
		players:    [4]*Client{},
		spectators: []*Client{},
	}
}
