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
	id          string
	join        chan *Player
	leave       chan *Player
	broadcast   chan []byte
	servicePipe chan<- string
	state       GameState
	Lock        sync.Mutex
	gameLogger  *utils.TambolaLogger
	gameCtx     context.Context
	cancel      context.CancelFunc
}

type GameServer interface {
	StartGameServer()
	AddPlayer(*Player)
	RemovePlayer(*Player)
	BroadcastMessage([]byte)
	Log(string)
}

func NewGameServer(gameID string, host *Player, servicePipe chan<- string) GameServer {
	log.Println(fmt.Sprintf("Making new game server with game id %s", gameID))
	ctx, cancel := context.WithCancel(context.Background())
	childCtx := context.WithValue(ctx, "game_id", gameID)
	logger := utils.NewTambolaLogger(childCtx)
	return &gameServer{
		id:          gameID,
		join:        make(chan *Player),
		leave:       make(chan *Player),
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

func (gs *gameServer) AddPlayer(player *Player) {
	go player.ReadPump(gs.gameCtx)
	go player.WritePump(gs.gameCtx)
	gs.join <- player
}

func (gs *gameServer) RemovePlayer(player *Player) {
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
			return
		}
	}
}

func (gs *gameServer) registerPlayer(user *Player) {
	serverPlayer := &Player{
		User: &User{
			Id:   -1,
			Name: "Server",
		},
	}
	userJoinedPayload := NewPlayerJoinedPayload(user)
	message := NewMessage(-1, PlayerJoinedEvent, serverPlayer, userJoinedPayload)

	//server register
	gs.state.AddPlayer(user)

	gs.broadcastMessage(message.EncodeToJson())
	gs.sendGameStateToJoinee(user)
}

func (gs *gameServer) sendGameStateToJoinee(player *Player) {
	players := make([]*Player, 0)
	for memberPlayer := range gs.state.GetPlayers() {
		players = append(players, memberPlayer)
	}
	playersAlreadyInLobbyPayload := NewPlayersAlreadyInLobbyPayload(players, gs.id)
	serverPlayer := &Player{
		User: &User{
			Id:   -1,
			Name: "Server",
		},
	}
	message := NewMessage(rand.Int64(),
		PlayersInLobbyEvent,
		serverPlayer, playersAlreadyInLobbyPayload)
	player.Send <- message.EncodeToJson()
}

func (gs *gameServer) unregisterPlayer(player *Player) {
	if player.IsHost {
		gs.killServer()
		return
	}

	//server register
	gs.state.RemovePlayer(player)

	serverPlayer := &Player{
		User: &User{
			Id:   -1,
			Name: "Server",
		},
	}
	userLeftPayload := NewPlayerLeftPayload(player)
	message := NewMessage(-1, PlayerLeftEvent, serverPlayer, userLeftPayload)
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
