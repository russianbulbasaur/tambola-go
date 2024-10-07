package services

import (
	"cmd/tambola/internals/game_server"
	"cmd/tambola/internals/repositories"
	"cmd/tambola/models"
	"github.com/gorilla/websocket"
	"log"
	"runtime"
	"sync"
)

type gameService struct {
	games       map[string]game_server.GameServer
	activeGames int64
	mutex       sync.Mutex
	servicePipe chan string
	gameRepo    repositories.GameRepository
}

type GameService interface {
	CreateGame(*models.User, *websocket.Conn)
	JoinGame(string, *models.User, *websocket.Conn)
}

func NewGameService(gameRepo repositories.GameRepository) GameService {
	gameMap := make(map[string]game_server.GameServer)
	servicePipe := make(chan string, 100)
	service := &gameService{
		games:       gameMap,
		activeGames: 0,
		servicePipe: servicePipe,
		gameRepo:    gameRepo,
	}
	go deleteGameService(service)
	return service
}

func deleteGameService(gs *gameService) {
	for {
		select {
		case gameId := <-gs.servicePipe:
			log.Printf("Deleting game %s", gameId)
			gs.mutex.Lock()
			if _, exists := gs.games[gameId]; exists {
				gs.activeGames--
				delete(gs.games, gameId)
			}
			gs.mutex.Unlock()
			log.Printf("Deleted game %s successfully", gameId)
			log.Println("Goroutines : ", runtime.NumGoroutine())
		}
	}
}

func (gs *gameService) CreateGame(user *models.User, conn *websocket.Conn) {
	log.Printf("Create request recieved. GoRoutines : %d", runtime.NumGoroutine())
	user.IsHost = true
	host := &game_server.Player{
		User: user,
		Send: make(chan []byte, 500),
		Conn: conn,
	}
	gameServer := game_server.NewGameServer(host, gs.servicePipe, gs.gameRepo)

	go gameServer.StartGameServer()
	gameServer.AddPlayer(host)

	//Add game to list of games
	host.GameServer = gameServer
	gs.games[gameServer.GetGameId()] = gameServer
}

func (gs *gameService) JoinGame(code string, user *models.User, conn *websocket.Conn) {
	gameServer := gs.games[code]
	if gameServer == nil {
		conn.Close()
		return
	}
	player := &game_server.Player{
		User:       user,
		GameServer: gameServer,
		Conn:       conn, Send: make(chan []byte, 500)}
	gameServer.AddPlayer(player)
}
