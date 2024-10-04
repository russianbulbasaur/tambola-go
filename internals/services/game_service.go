package services

import (
	"cmd/tambola/models"
	"github.com/gorilla/websocket"
	"log"
	"math/rand/v2"
	"runtime"
	"strconv"
	"sync"
)

type gameService struct {
	games       map[string]models.GameServer
	activeGames int64
	mutex       sync.Mutex
	servicePipe chan string
}

type GameService interface {
	CreateGame(*models.User, *websocket.Conn)
	JoinGame(string, *models.User, *websocket.Conn)
}

func NewGameService() GameService {
	gameMap := make(map[string]models.GameServer)
	servicePipe := make(chan string)
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
			log.Println("Goroutines : ", runtime.NumGoroutine())
		}
	}
}

func (gs *gameService) CreateGame(user *models.User, conn *websocket.Conn) {
	host := &models.Player{
		User: user,
		Send: make(chan []byte, 500),
		Conn: conn,
	}
	gameId := generateGameCode()
	gameServer := models.NewGameServer(gameId, host, gs.servicePipe)

	go gameServer.StartGameServer()
	gameServer.AddPlayer(host)

	//Add game to list of games
	host.GameServer = gameServer
	gs.games[gameId] = gameServer
}

func generateGameCode() string {
	return strconv.Itoa(rand.Int())
}

func (gs *gameService) JoinGame(code string, user *models.User, conn *websocket.Conn) {
	gameServer := gs.games[code]
	if gameServer == nil {
		conn.Close()
		return
	}
	player := &models.Player{
		User:       user,
		GameServer: gameServer,
		Conn:       conn, Send: make(chan []byte, 500)}
	gameServer.AddPlayer(player)
}
