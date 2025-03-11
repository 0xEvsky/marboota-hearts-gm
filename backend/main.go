package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	mu        sync.Mutex
	conns     map[*websocket.Conn]*Client
	instances map[string]*Instance // key is instanceid
}

func newServer() *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]*Client),
		instances: make(map[string]*Instance),
	}
}

var upgrader = websocket.Upgrader{} // use default options

func (s *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: handle auth
	w.Write([]byte("W.I.P"))
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("New client connected")
	var newClient = newClient(c)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns[c] = newClient

	go s.read(c)
}

func (s *Server) read(ws *websocket.Conn) {
	defer ws.Close()
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			var c = s.conns[ws]
			if c.isAuthed {
				c.broadcastToMates(map[string]string{"ACTION": "LEAVE", "USERID": s.conns[ws].id})
				if c.state == ClientSeated {
					c.instance.table.unseatPlayer(c)
				}
				delete(c.instance.clients, s.conns[ws].id)

				if len(c.instance.clients) == 0 {
					delete(s.instances, s.conns[ws].instance.id)
					log.Printf("Deleted empty instance")
				}
			}

			delete(s.conns, ws)
			log.Printf("Connection closed: %s\n", err)
			break
		}

		msgHandler(s.conns[ws], msg)
	}
}

var server = newServer()

func main() {
	const PORT = "3000"
	//const HOST = "localhost"

	http.HandleFunc("/auth", server.authHandler)
	http.HandleFunc("/ws", server.wsHandler)

	log.Printf("Server is up and listening on port: %s\n", PORT)
	err := http.ListenAndServe(":"+PORT, nil)
	log.Fatal(err)
}
