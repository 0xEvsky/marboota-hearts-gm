package main

type Table struct {
	players [4]Client
}

func NewTable() Table {
	return Table{}
}
