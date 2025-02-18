package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	intanceId int
	id        string
	name      string
	iconUrl   string
	seat      int
	isTurn    bool
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:      conn,
		intanceId: 0,
		id:        "",
		name:      "",
		iconUrl:   "",
		seat:      0,
		isTurn:    false,
	}
}

type Server struct {
	conns map[*websocket.Conn]*Client
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]*Client),
	}
}

var upgrader = websocket.Upgrader{} // use default options

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: LOG
		log.Println(err)
		return
	}

	// TODO: LOG
	log.Println("new user!")
	var newClient = NewClient(c)

	// TODO: MUTEX
	s.conns[c] = newClient
	// MUTEX

	go s.read(c)
}

func (s *Server) read(ws *websocket.Conn) {
	defer ws.Close()
	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			// TODO: HANDLE & LOG
			delete(s.conns, ws)
			log.Println(err)
			break
		}

		// TODO: HANDLE
		log.Println(string(msg))
		err = ws.WriteMessage(mt, []byte("nice!"))
		if err != nil {
			// TODO: HANDLE & LOG
			log.Println(err)
			break
		}
	}
}

func main() {
	const PORT = "3000"
	const HOST = "localhost"

	var server = NewServer()

	http.HandleFunc("/ws", server.wsHandler)

	// TODO: LOG
	log.Printf("Server is up and listening on port: %s\n", PORT)
	err := http.ListenAndServe(HOST+":"+PORT, nil)
	log.Fatal(err)
}
