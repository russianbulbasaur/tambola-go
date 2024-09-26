package models

import (
	"cmd/tambola/utils"
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
	gameLogger  *utils.TambolaLogger
	gameCtx     context.Context
	cancel      context.CancelFunc
}

type GameServer interface {
	StartGameServer()
	AddPlayer(*User)
	RemovePlayer(*User)
	BroadcastMessage([]byte)
	Log(string)
}

func NewGameServer(gameID int32, host *User, servicePipe chan<- int32) GameServer {
	log.Println(fmt.Sprintf("Making new game server with game id %d", gameID))
	ctx, cancel := context.WithCancel(context.Background())
	childCtx := context.WithValue(ctx, "game_id", gameID)
	logger := utils.NewTambolaLogger(childCtx)
	return &gameServer{
		id:          gameID,
		join:        make(chan *User),
		leave:       make(chan *User),
		broadcast:   make(chan []byte),
		state:       NewGameState(host, logger),
		servicePipe: servicePipe,
		gameCtx:     childCtx,
		cancel:      cancel,
		gameLogger:  logger,
	}
}

func (gs *gameServer) Log(text string) {
	gs.gameLogger.Log(text)
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
			gs.unregisterPlayer(user)
		case message := <-gs.broadcast:
			gs.broadcastMessage(message)
		case <-gs.gameCtx.Done():
			gs.gameLogger.Log(fmt.Sprintf("Stopping game server %d", gs.id))
			gs.servicePipe <- gs.id
			gs.gameLogger.Log(fmt.Sprintf("Stopped game server %d", gs.id))
			return
		}
	}
}

func (gs *gameServer) registerPlayer(user *User) {
	serverUser := &User{
		Id:   -1,
		Name: "Server",
	}
	userJoinedPayload := NewUserJoinedPayload(user)
	message := NewMessage(-1, UserJoinedEvent, serverUser, userJoinedPayload)

	//server register
	gs.state.AddPlayer(user)

	gs.broadcastMessage(message.EncodeToJson())
	gs.sendGameStateToJoinee(user)
}

func (gs *gameServer) sendGameStateToJoinee(player *User) {
	players := make([]*User, 0)
	for memberPlayer := range gs.state.GetPlayers() {
		players = append(players, memberPlayer)
	}
	playersAlreadyInLobbyPayload := NewPlayersAlreadyInLobbyPayload(players, gs.id)
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

	//server register
	gs.state.RemovePlayer(player)

	serverUser := &User{
		Id:   -1,
		Name: "Server",
	}
	userLeftPayload := NewUserLeftPayload(player)
	message := NewMessage(-1, UserLeftEvent, serverUser, userLeftPayload)
	gs.broadcastMessage(message.EncodeToJson())
}

func (gs *gameServer) broadcastMessage(data []byte) {
	gs.state.UpdateGameState(data)
	playersInGame := gs.state.GetPlayers()
	for player := range playersInGame {
		player.Lock.Lock()
		player.Send <- data
		player.Lock.Unlock()
	}
}

func (gs *gameServer) killServer() {
	gs.gameLogger.Log("Initiating server kill")
	gs.gameLogger.Log(fmt.Sprintf("Goroutines : %d", runtime.NumGoroutine()))
	gs.cancel()
}
