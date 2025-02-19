package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	conns     map[*websocket.Conn]*Client
	instances map[string]*Instance
}

func NewServer() *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]*Client),
		instances: make(map[string]*Instance),
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
	newClient.writeJson = func(msg map[string]string) error {
		return c.WriteMessage(websocket.TextMessage, toJson(msg))
	}
	newClient.writeError = func(msg string) error {
		return c.WriteMessage(websocket.TextMessage, toJson(map[string]string{"ACTION": "ERROR", "MESSAGE": msg}))
	}
	newClient.writeOk = func() error {
		return c.WriteMessage(websocket.TextMessage, toJson(map[string]string{"ACTION": "OK"}))
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
			if s.conns[ws].isAuthed {
				s.conns[ws].broadcastToInstance(map[string]string{"ACTION": "LEAVE", "USERID": s.conns[ws].id})
			}
			// TODO: Delete client
			delete(s.conns[ws].instance.clients, s.conns[ws].id)
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

	// TODO: auth route
	http.HandleFunc("/ws", server.wsHandler)

	log.Printf("Server is up and listening on port: %s\n", PORT)
	err := http.ListenAndServe(HOST+":"+PORT, nil)
	log.Fatal(err)
}
