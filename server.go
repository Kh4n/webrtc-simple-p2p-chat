package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type websocketRWLock struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

type server struct {
	upgrader *websocket.Upgrader
	peers    map[string]*websocketRWLock
}

func newServer() *server {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return &server{
		upgrader: &upgrader,
		peers:    make(map[string]*websocketRWLock),
	}
}

func (s *server) writeJSONToPeer(dat interface{}, peerID string) error {
	if _, ok := s.peers[peerID]; !ok {
		return errors.New("Peer attempting to connect to peerID that does not exist")
	}
	s.peers[peerID].mutex.Lock()
	defer s.peers[peerID].mutex.Unlock()
	err := s.peers[peerID].conn.WriteJSON(dat)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) handleConnection(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Unable to upgrade:", err)
		return
	}
	log.Println("Recieved new connection from:", c.RemoteAddr().String())
	c.SetCloseHandler(s.handleClose)

	defer c.Close()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return
			}
			log.Println("Error reading message:", err)
			return
		}
		switch mt {
		case websocket.TextMessage:
			err = s.handleText(c, msg)
			if err != nil {
				log.Println("Error occurred while reading text message:", err)
			}
		case websocket.BinaryMessage:
			log.Println("Binary message recieved:", msg)
		}
	}
}

func (s *server) handleClose(code int, text string) error {
	log.Println("Connection closed:", code, ":", text)
	return nil
}

func (s *server) handleText(c *websocket.Conn, msg []byte) error {
	var a all
	err := json.Unmarshal(msg, &a)
	if err != nil {
		return err
	}

	log.Println("Type:", a.Type)
	switch a.Type {
	case "forward":
		err = s.handleForward(msg)
	case "register":
		err = s.handleRegister(c, msg)
	case "offer", "answer":
		err = s.handleOfferOrAnswer(msg)
	default:
		e := fmt.Sprint("unknown data recieved: ", a)
		err = errors.New(e)
	}

	return err
}

func (s *server) handleRegister(c *websocket.Conn, msg []byte) error {
	r, err := readRegister(msg)
	if err != nil {
		return err
	}

	if _, ok := s.peers[r.PeerID]; ok {
		return errors.New("Peer attempting to register as peer that already exists")
	}
	s.peers[r.PeerID] = &websocketRWLock{conn: c}
	c.SetCloseHandler(func(code int, text string) error {
		delete(s.peers, r.PeerID)
		log.Println("Removed peer:", r.PeerID)
		return s.handleClose(code, text)
	})
	log.Println("Registered new peer:", r.PeerID)
	return nil
}

func (s *server) handleOfferOrAnswer(msg []byte) error {
	oa, err := readOfferOrAnswer(msg)
	if err != nil {
		return err
	}
	return s.writeJSONToPeer(oa, oa.To)
}

func (s *server) handleForward(msg []byte) error {
	f, err := readForward(msg)
	if err != nil {
		return err
	}
	return s.writeJSONToPeer(f, f.To)
}

func main() {
	s := newServer()
	http.HandleFunc("/", s.handleConnection)
	log.Println("Starting server")
	log.Fatal(http.ListenAndServe("localhost:6503", nil))
}
