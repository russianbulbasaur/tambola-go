package services

import (
	"cmd/tambola/models"
	"github.com/gorilla/websocket"
	"math/rand/v2"
	"strconv"
)

type gameService struct {
	games       map[string]*models.GameServer
	activeGames int64
}

type GameService interface {
	CreateGame(string, string, *websocket.Conn)
	JoinGame(string, string, string, *websocket.Conn)
}

func NewGameService() GameService {
	return &gameService{
		games:       make(map[string]*models.GameServer),
		activeGames: 0,
	}
}

func (gs *gameService) CreateGame(userId string, name string, conn *websocket.Conn) {
	user := &models.User{
		Name: name,
		Id:   userId,
		Send: make(chan []byte),
		Conn: conn,
	}
	go user.ReadPump()
	go user.WritePump()
	gameId := generateGameCode()
	gameServer := models.NewGameServer(gameId, user)
	go gameServer.StartGameServer()
	user.Game = gameServer
	gameServer.Join <- user
	code := generateGameCode()
	gameServer.Broadcast <- []byte(code)
	gs.games[code] = gameServer
}

func generateGameCode() string {
	return "2732844840814395882"
	return strconv.Itoa(rand.Int())
}

func (gs *gameService) JoinGame(code string, userId string, name string, conn *websocket.Conn) {
	gameServer := gs.games[code]
	if gameServer == nil {
		return
	}
	user := &models.User{Id: userId,
		Name: name, Game: gameServer,
		Conn: conn, Send: make(chan []byte)}
	go user.ReadPump()
	go user.WritePump()
	gameServer.Join <- user
}
