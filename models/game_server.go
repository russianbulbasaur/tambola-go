package models

import (
	"log"
)

type GameServer struct {
	Players   map[*User]bool `json:"players"`
	Host      *User          `json:"host"`
	State     *GameState     `json:"state"`
	Join      chan *User
	Leave     chan *User
	Broadcast chan []byte
}

func NewGameServer(gameID string, host *User) *GameServer {
	log.Println("Making new game server ")
	return &GameServer{
		Players:   make(map[*User]bool),
		Host:      host,
		Join:      make(chan *User),
		Leave:     make(chan *User),
		Broadcast: make(chan []byte),
		State:     &GameState{},
	}
}

func (gs *GameServer) StartGameServer() {
	log.Println("Starting game server")
	for {
		select {
		case user := <-gs.Join:
			gs.registerPlayer(user)
		case user := <-gs.Leave:
			gs.unregisterPlayer(user)
		case message := <-gs.Broadcast:
			gs.broadcast(message)
		}
	}
}

func (gs *GameServer) registerPlayer(user *User) {
	log.Println("Player joined : ", user.Name)
	gs.Players[user] = true
}

func (gs *GameServer) unregisterPlayer(user *User) {
	if gs.Players[user] {
		delete(gs.Players, user)
	}
}

func (gs *GameServer) broadcast(text []byte) {
	for client := range gs.Players {
		client.Send <- text
	}
}
