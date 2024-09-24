package services

import (
	"cmd/tambola/models"
	"github.com/gorilla/websocket"
	"log"
	"math/rand/v2"
	"runtime"
	"sync"
)

type gameService struct {
	games       map[int32]models.GameServer
	activeGames int64
	mutex       sync.Mutex
	servicePipe chan int32
}

type GameService interface {
	CreateGame(int64, string, *websocket.Conn)
	JoinGame(int32, int64, string, *websocket.Conn)
}

func NewGameService() GameService {
	gameMap := make(map[int32]models.GameServer)
	servicePipe := make(chan int32)
	service := &gameService{
		games:       gameMap,
		activeGames: 0,
		servicePipe: servicePipe,
	}
	go deleteGameService(service)
	return service
}

func deleteGameService(gs *gameService) {
	for {
		select {
		case gameId := <-gs.servicePipe:
			log.Printf("Deleting game %d", gameId)
			gs.mutex.Lock()
			if _, exists := gs.games[gameId]; exists {
				gs.activeGames--
				delete(gs.games, gameId)
			}
			gs.mutex.Unlock()
			log.Printf("Deleted game %d successfully", gameId)
			log.Println("Stats : ", runtime.NumGoroutine())
		}
	}
}

func (gs *gameService) CreateGame(userId int64, name string, conn *websocket.Conn) {
	host := &models.User{
		Name:   name,
		Id:     userId,
		Send:   make(chan []byte, 500),
		Conn:   conn,
		IsHost: true,
	}
	gameId := generateGameCode()
	gameServer := models.NewGameServer(gameId, host, gs.servicePipe)

	go gameServer.StartGameServer()
	gameServer.AddPlayer(host)

	//Add game to list of games
	host.GameServer = gameServer
	gs.games[gameId] = gameServer
}

func generateGameCode() int32 {
	return rand.Int32()
}

func (gs *gameService) JoinGame(code int32, userId int64, name string, conn *websocket.Conn) {
	gameServer := gs.games[code]
	if gameServer == nil {
		conn.Close()
		return
	}
	player := &models.User{Id: userId,
		Name: name, GameServer: gameServer,
		Conn: conn, Send: make(chan []byte, 500), IsHost: false}
	gameServer.AddPlayer(player)
}
