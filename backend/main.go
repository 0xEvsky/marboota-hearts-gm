package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	conns     map[*websocket.Conn]*Client
	instances []*Instance
}

func NewServer() *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]*Client),
		instances: []*Instance{},
	}
}

var upgrader = websocket.Upgrader{} // use default options

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("New client connected")
	var newClient = NewClient(c)
	newClient.write = func(msg []byte) error {
		return c.WriteMessage(websocket.TextMessage, msg)
	}

	// TODO: MUTEX
	s.conns[c] = newClient
	// MUTEX

	go s.read(c)
}

func (s *Server) read(ws *websocket.Conn) {
	defer ws.Close()
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			delete(s.conns, ws)
			log.Println(err)
			break
		}

		msgHandler(s.conns[ws], msg)
	}
}

var server = NewServer()

func main() {
	const PORT = "3000"
	const HOST = "localhost"

	http.HandleFunc("/ws", server.wsHandler)

	log.Printf("Server is up and listening on port: %s\n", PORT)
	err := http.ListenAndServe(HOST+":"+PORT, nil)
	log.Fatal(err)
}
