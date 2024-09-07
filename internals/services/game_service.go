package services

import (
	"cmd/tambola/models"
	json2 "encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand/v2"
)

type gameService struct {
	games       map[int32]*models.GameServer
	activeGames int64
}

type GameService interface {
	CreateGame(int64, string, *websocket.Conn)
	JoinGame(int32, int64, string, *websocket.Conn)
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
	go gameServer.StartGameServer()
	user.GameServer = gameServer
	gs.games[gameId] = gameServer
	gameServer.Join <- user
	sendGameIdToHost(user, gameId)
}

func sendGameIdToHost(user *models.User, gameId int32) {
	gameIdPayload := &models.GameIdPayload{Id: gameId}
	encodedGameIdPayload, err := json2.Marshal(gameIdPayload)
	if err != nil {
		log.Println(err)
		return
	}
	message := &models.Message{Id: rand.Int64(),
		Sender: &models.User{Id: -1, Name: "Server"},
		Event:  models.GameIdEvent, Payload: string(encodedGameIdPayload)}
	encodedMessage, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	user.Send <- encodedMessage
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
