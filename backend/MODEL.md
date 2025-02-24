# Model
The various parts making up the server heirarchy.
> Iota constant groups are basically enums
> Methods were omitted
## Client
```go
const (
	ClientUnavailable ClientState = iota
	ClientIdle
	ClientSeated
)

type Client struct {
	mu       sync.Mutex
	conn     *websocket.Conn
	isAuthed bool
	instance *Instance
	id       string
	name     string
	iconUrl  string
	state    ClientState
	player   *Player
}
```
## Player
```go
type PlayerState int

const (
	PlayerUnavailable PlayerState = iota
	PlayerWaiting
	PlayerReady
	PlayerTrumping
	PlayerPlaying
)

type Team int

const (
	TeamA Team = iota
	TeamB
)

type Player struct {
	client  *Client
	state   PlayerState
	hand    []string
	seat    int
	team    Team
	score   int
	partner *Player
	isTurn  bool
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
type TableState int

const (
	TableWaiting TableState = iota
	TableTrumping
	TablePlaying
)

type Table struct {
	players [4]*Player
	state   TableState
	turn    int
}
```