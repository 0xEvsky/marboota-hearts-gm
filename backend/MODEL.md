# Model
The various parts making up the server heirarchy.
> Methods were omitted
> Iota constant groups are basically enums
## Client
```go
const (
	ClientUnavailable ClientState = iota
	ClientIdle
	ClientWaiting
	ClientPlaying
)

type Client struct {
	conn       *websocket.Conn
	isAuthed   bool
	instance   *Instance
	id         string
	name       string
	iconUrl    string
	state      ClientState
	table      *Table
	seat       int
	isTurn     bool
}
```
## Server
```go
type Server struct {
	conns     map[*websocket.Conn]*Client
	instances map[string]*Instance
}
```
## Instance
```go
type Instance struct {
	id      string
	clients map[string]*Client // key is userid
	table   Table
}
```
## Table
```go
const (
	TableWaiting TableState = iota
	TableTrumping
	TablePlaying
)

type Table struct {
	players    [4]*Client
	spectators []*Client
	state      TableState
	turn       int
}
```