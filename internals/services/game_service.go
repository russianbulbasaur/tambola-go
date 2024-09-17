package services

import (
	"cmd/tambola/models"
	"github.com/gorilla/websocket"
	"math/rand/v2"
	"sync"
)

type gameService struct {
	games             map[int32]*models.GameServer
	activeGames       int64
	mutex             sync.Mutex
	deleteGameChannel chan int32
}

type GameService interface {
	CreateGame(int64, string, *websocket.Conn)
	JoinGame(int32, int64, string, *websocket.Conn)
}

func NewGameService() GameService {
	gameMap := make(map[int32]*models.GameServer)
	deleteChannel := make(chan int32)
	service := &gameService{
		games:             gameMap,
		activeGames:       0,
		deleteGameChannel: deleteChannel,
	}
	go deleteGame(service)
	return service
}

func deleteGame(gs *gameService) {
	for {
		select {
		case gameId := <-gs.deleteGameChannel:
			gs.mutex.Lock()
			if _, exists := gs.games[gameId]; exists {
				gs.activeGames--
				delete(gs.games, gameId)
			}
			gs.mutex.Unlock()
		}
	}
}

func (gs *gameService) CreateGame(userId int64, name string, conn *websocket.Conn) {
	user := &models.User{
		Name:   name,
		Id:     userId,
		Send:   make(chan []byte, 500),
		Conn:   conn,
		IsHost: true,
	}
	go user.ReadPump()
	go user.WritePump()
	gameId := generateGameCode()
	gameServer := models.NewGameServer(gameId, user)
	go gameServer.StartGameServer(gs.deleteGameChannel)
	user.GameServer = gameServer
	gs.games[gameId] = gameServer
	gameServer.Join <- user
}

func generateGameCode() int32 {
	return rand.Int32()
}

func (gs *gameService) JoinGame(code int32, userId int64, name string, conn *websocket.Conn) {
	gameServer := gs.games[code]
	if gameServer == nil || !(gameServer.State.Status == models.Waiting) {
		conn.Close()
		return
	}
	user := &models.User{Id: userId,
		Name: name, GameServer: gameServer,
		Conn: conn, Send: make(chan []byte, 500), IsHost: false}
	go user.ReadPump()
	go user.WritePump()
	gameServer.Join <- user
}
