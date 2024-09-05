package services

import (
	"cmd/tambola/models"
	"cmd/tambola/pkg/game_server"
	"math/rand/v2"
	"strconv"
)

type gameService struct {
	games       map[string]*game_server.GameServer
	activeGames int64
}

type GameService interface {
	createGame()
	joinGame(string)
}

func NewGameService() GameService {
	return &gameService{
		games:       make(map[string]*game_server.GameServer),
		activeGames: 0,
	}
}

func (gs *gameService) createGame(host *models.User) {
	gameId := generateGameCode()
	gameServer := game_server.NewGameServer(gameId, host)
	code := generateGameCode()
	gs.games[code] = gameServer
	go gameServer.StartGameServer()
}

func generateGameCode() string {
	return strconv.Itoa(rand.Int())
}

func (gs *gameService) joinGame(code string) {
	gameServer := gs.games[]
	if gameServer == nil {
		return
	}
	gameServer.Join <- user
}
