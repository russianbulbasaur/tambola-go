package models

import (
	"cmd/tambola/models/payloads"
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"runtime"
	"sync"
)

type gameServer struct {
	id          int32
	join        chan *User
	leave       chan *User
	broadcast   chan []byte
	servicePipe chan<- int32
	state       GameState
	Lock        sync.Mutex
	gameCtx     context.Context
	cancel      context.CancelFunc
}

type GameServer interface {
	StartGameServer()
	AddPlayer(*User)
	RemovePlayer(*User)
	BroadcastMessage([]byte)
}

func NewGameServer(gameID int32, host *User, servicePipe chan<- int32) GameServer {
	log.Println(fmt.Sprintf("Making new game server with game id %d", gameID))
	ctx, cancel := context.WithCancel(context.Background())
	return &gameServer{
		id:          gameID,
		join:        make(chan *User),
		leave:       make(chan *User),
		broadcast:   make(chan []byte),
		state:       NewGameState(host),
		servicePipe: servicePipe,
		gameCtx:     ctx,
		cancel:      cancel,
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

func (gs *gameServer) StartGameServer() {
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
			gs.servicePipe <- gs.id
			log.Printf("Stopped game server %d", gs.id)
			return
		}
	}
}

func (gs *gameServer) registerPlayer(user *User) {
	serverUser := &User{
		Id:   -1,
		Name: "Server",
	}
	userJoinedPayload := payloads.NewUserJoinedPayload(user)
	message := NewMessage(-1, UserJoinedEvent, serverUser, userJoinedPayload)
	gs.broadcastMessage(message.EncodeToJson())
	gs.sendGameStateToJoinee(user)
}

func (gs *gameServer) sendGameStateToJoinee(player *User) {
	var players []*User
	for memberPlayer := range gs.state.GetPlayers() {
		players = append(players, memberPlayer)
	}
	playersAlreadyInLobbyPayload := payloads.NewPlayersAlreadyInLobbyPayload(players, gs.id)
	serverUser := &User{
		Id:   -1,
		Name: "Server",
	}
	message := NewMessage(rand.Int64(),
		PlayersInLobbyEvent,
		serverUser, playersAlreadyInLobbyPayload)
	player.Send <- message.EncodeToJson()
}

func (gs *gameServer) unregisterPlayer(player *User) {
	if player.IsHost {
		gs.killServer()
		return
	}
	serverUser := &User{
		Id:   -1,
		Name: "Server",
	}
	userLeftPayload := payloads.NewUserLeftPayload(player)
	message := NewMessage(-1, UserLeftEvent, serverUser, userLeftPayload)
	gs.broadcastMessage(message.EncodeToJson())
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

func (gs *gameServer) killServer() {
	log.Println("Initiating server kill")
	log.Printf("Goroutines : %d", runtime.NumGoroutine())
	gs.cancel()
}
