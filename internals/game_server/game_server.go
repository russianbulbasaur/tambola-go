package game_server

import (
	"cmd/tambola/internals/repositories"
	"cmd/tambola/models"
	"cmd/tambola/utils"
	"context"
	"fmt"
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
	syncer      GameSyncer
}

type GameServer interface {
	StartGameServer()
	AddPlayer(*Player)
	RemovePlayer(*Player)
	BroadcastMessage([]byte)
	Log(string)
	GetGameId() string
}

func NewGameServer(host *Player, servicePipe chan<- string, gameRepo repositories.GameRepository) GameServer {
	gameID := fmt.Sprintf("TMB%d", gameRepo.CreateGame(host.User))
	fmt.Printf("Making new game server with game id %s", gameID)
	ctx, cancel := context.WithCancel(context.Background())
	childCtx := context.WithValue(ctx, "game_id", gameID)
	logger := utils.NewTambolaLogger(childCtx)
	gameState := NewGameState(host, logger)
	syncer := NewGameSyncer(gameID, gameRepo, gameState)
	return &gameServer{
		id:          gameID,
		join:        make(chan *Player),
		leave:       make(chan *Player),
		broadcast:   make(chan []byte),
		state:       gameState,
		servicePipe: servicePipe,
		gameCtx:     childCtx,
		cancel:      cancel,
		gameLogger:  logger,
		syncer:      syncer,
	}
}

func (gs *gameServer) GetGameId() string {
	return gs.id
}

func (gs *gameServer) Log(text string) {
	gs.gameLogger.LogChannel <- text
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
	go gs.gameLogger.StartLogging(gs.gameCtx)
	go gs.syncer.Sync(gs.gameCtx)
	for {
		select {
		case user := <-gs.join:
			gs.registerPlayer(user)
		case user := <-gs.leave:
			gs.unregisterPlayer(user)
		case message := <-gs.broadcast:
			gs.broadcastMessage(message)
		case <-gs.gameCtx.Done():
			gs.Log(fmt.Sprintf("Stopping game server %s", gs.id))
			gs.servicePipe <- gs.id
			return
		}
	}
}

func (gs *gameServer) registerPlayer(player *Player) {
	serverPlayer := &models.User{
		Id:   -1,
		Name: "Server",
	}
	userJoinedPayload := models.NewPlayerJoinedPayload(player.User)
	message := models.NewMessage(-1, models.PlayerJoinedEvent, serverPlayer, userJoinedPayload)

	//server register
	gs.state.AddPlayer(player)

	gs.broadcastMessage(message.EncodeToJson())
	gs.sendGameStateToJoinee(player)
}

func (gs *gameServer) sendGameStateToJoinee(player *Player) {
	players := make([]*models.User, 0)
	for memberPlayer := range gs.state.GetPlayers() {
		players = append(players, memberPlayer.User)
	}
	playersAlreadyInLobbyPayload := models.NewPlayersAlreadyInLobbyPayload(players, gs.id)
	serverPlayer := &models.User{
		Id:   -1,
		Name: "Server",
	}
	message := models.NewMessage(rand.Int64(),
		models.PlayersInLobbyEvent,
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

	serverPlayer := &models.User{
		Id:   -1,
		Name: "Server",
	}
	userLeftPayload := models.NewPlayerLeftPayload(player.User)
	message := models.NewMessage(-1, models.PlayerLeftEvent, serverPlayer, userLeftPayload)
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
	gs.Log("Initiating server kill")
	gs.Log(fmt.Sprintf("Goroutines : %d", runtime.NumGoroutine()))

	//Warning : kills all the goroutines for the game server (logger + syncer)
	gs.cancel()
}
