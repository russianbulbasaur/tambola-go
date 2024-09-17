package services

import (
	"cmd/tambola/models"
	"github.com/gorilla/websocket"
	"math/rand/v2"
	"sync"
)

type gameService struct {
	games       map[int32]*models.GameServer
	activeGames int64
	mutex       sync.Mutex
}

type GameService interface {
	CreateGame(int64, string, *websocket.Conn)
	JoinGame(int32, int64, string, *websocket.Conn)
	DeleteGame(int32)
}

func NewGameService() GameService {
	return &gameService{
		games:       make(map[int32]*models.GameServer),
		activeGames: 0,
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
	go gameServer.StartGameServer(gs)
	user.GameServer = gameServer
	gs.games[gameId] = gameServer
	gameServer.Join <- user
}

func (gs *gameService) DeleteGame(gameId int32) {
	gs.mutex.Lock()
	if _, exists := gs.games[gameId]; exists {
		gs.activeGames--
		delete(gs.games, gameId)
	}
	gs.mutex.Unlock()
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
