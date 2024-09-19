package models

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
)

type gameServer struct {
	id        int32
	join      chan *User
	leave     chan *User
	broadcast chan []byte
	state     GameState
	Lock      sync.Mutex
	gameCtx   context.Context
	cancel    context.CancelFunc
}

type GameServer interface {
	StartGameServer(chan<- int32)
	AddPlayer(*User)
	RemovePlayer(*User)
	BroadcastMessage([]byte)
}

func NewGameServer(gameID int32, host *User) GameServer {
	log.Println(fmt.Sprintf("Making new game server with game id %d", gameID))
	ctx, cancel := context.WithCancel(context.Background())
	return &gameServer{
		id:        gameID,
		join:      make(chan *User),
		leave:     make(chan *User),
		broadcast: make(chan []byte),
		state:     NewGameState(host),
		gameCtx:   ctx,
		cancel:    cancel,
	}
}

func (gs *gameServer) BroadcastMessage(message []byte) {
	gs.Lock.Lock()
	gs.broadcast <- message
	gs.Lock.Unlock()
}

func (gs *gameServer) AddPlayer(player *User) {
	go player.ReadPump(gs.gameCtx)
	go player.WritePump(gs.gameCtx)
	gs.join <- player
}

func (gs *gameServer) RemovePlayer(player *User) {
	gs.leave <- player
}

func (gs *gameServer) StartGameServer(gameServiceDeleteChannel chan<- int32) {
	for {
		select {
		case user := <-gs.join:
			gs.registerPlayer(user)
		case user := <-gs.leave:
			log.Printf("User Leaving : %#v", user)
			gs.unregisterPlayer(user)
		case message := <-gs.broadcast:
			gs.broadcastMessage(message)
		case <-gs.gameCtx.Done():
			log.Printf("Stopping game server %d", gs.id)
			gameServiceDeleteChannel <- gs.id
			log.Printf("Stopped game server %d", gs.id)
			return
		}
	}
}

func (gs *gameServer) registerPlayer(user *User) {
	userJoinedPayload := &UserJoinedPayload{User: user}
	message := Message{
		UserJoinedPayload: userJoinedPayload,
		Id:                -1,
		Event:             UserJoinedEvent,
		Sender: &User{
			Id:   -1,
			Name: "Server",
		},
	}
	encodedMessage, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	gs.broadcastMessage(encodedMessage)
	gs.sendGameStateToJoinee(user)
}

func (gs *gameServer) sendGameStateToJoinee(player *User) {
	var players []*User
	for memberPlayer := range gs.state.GetPlayers() {
		players = append(players, memberPlayer)
	}
	playersAlreadyInLobbyPayload := &PlayersAlreadyInLobbyPayload{Players: players, GameId: gs.id}
	message := &Message{Id: rand.Int64(),
		PlayersAlreadyInLobbyPayload: playersAlreadyInLobbyPayload,
		Sender: &User{
			Id:   -1,
			Name: "Server",
		}, Event: PlayersInLobbyEvent}
	encodedMessage, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	player.Send <- encodedMessage
}

func (gs *gameServer) unregisterPlayer(player *User) {
	if player.IsHost {
		return
	}
	userLeftPayload := &UserLeftPayload{User: player}
	message := Message{
		UserLeftPayload: userLeftPayload,
		Id:              -1,
		Event:           UserLeftEvent,
		Sender: &User{
			Id:   -1,
			Name: "Server",
		},
	}
	encodedMessage, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	gs.broadcastMessage(encodedMessage)
}

func (gs *gameServer) broadcastMessage(data []byte) {
	isForHost := gs.state.UpdateGameState(data)
	if !isForHost {
		gs.state.GetHost().Send <- data
	} else {
		for player := range gs.state.GetPlayers() {
			log.Println(fmt.Sprintf("Sending to player %s", player.Name))
			player.Lock.Lock()
			player.Send <- data
			player.Lock.Unlock()
		}
	}
}
