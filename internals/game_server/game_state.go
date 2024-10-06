package game_server

import (
	"cmd/tambola/models"
	"cmd/tambola/utils"
	"encoding/json"
	"fmt"
	"log"
)

const Playing = "playing"
const Waiting = "waiting"
const Closed = "closed"

type gameState struct {
	players       map[*Player]bool
	host          *Player
	Status        string   `json:"status"`
	Alerts        []string `json:"alerts"`
	Claimed       []string `json:"claimed"`
	NumbersCalled []int32  `json:"numbers"`
	playerCount   int32
	logger        *utils.TambolaLogger
}

type GameState interface {
	GetPlayers() map[*Player]bool
	GetHost() *Player
	GetStatus() string
	GetPlayerCount() int32
	GetAlerts() []string
	GetClaimed() []string
	GetCalledNumbers() []int32
	AddPlayer(player *Player)
	RemovePlayer(player *Player)
	UpdateGameState(data []byte) bool
	GetStateJson() string
}

func NewGameState(host *Player, logger *utils.TambolaLogger) GameState {
	return &gameState{
		host:    host,
		players: make(map[*Player]bool),
		Status:  "waiting",
		logger:  logger,
	}
}

func (gs *gameState) GetPlayers() map[*Player]bool {
	return gs.players
}

func (gs *gameState) GetHost() *Player {
	return gs.host
}

func (gs *gameState) GetStatus() string {
	return gs.Status
}

func (gs *gameState) GetPlayerCount() int32 {
	return gs.playerCount
}

func (gs *gameState) GetAlerts() []string {
	return gs.Alerts
}

func (gs *gameState) GetClaimed() []string {
	return gs.Claimed
}

func (gs *gameState) GetCalledNumbers() []int32 {
	return gs.NumbersCalled
}

func (gs *gameState) addNumber(number int32) {
	gs.NumbersCalled = append(gs.NumbersCalled, number)
	gs.logger.LogChannel <- fmt.Sprintf("Called %d number", number)
}

func (gs *gameState) addAlert(alert string) {
	gs.Alerts = append(gs.Alerts, alert)
}

func (gs *gameState) AddPlayer(player *Player) {
	gs.players[player] = true
	gs.logger.LogChannel <- fmt.Sprintf("Player %s joined", player.GetName())
}

func (gs *gameState) RemovePlayer(player *Player) {
	if gs.players[player] {
		delete(gs.players, player)
	}
	gs.logger.LogChannel <- fmt.Sprintf("Player %s left", player.GetName())
}

func (gs *gameState) updateGameStatus(status string) {
	gs.Status = status
	gs.logger.LogChannel <- fmt.Sprintf("Updating game state to %s", status)
}

func (gs *gameState) UpdateGameState(data []byte) bool {
	message := models.Decode(data)
	switch message.GetEvent() {
	case models.NumberCalledEvent:
		numberCalledPayload := models.ParseNumberPayload(message.GetPayloadJson())
		gs.addNumber(numberCalledPayload.GetNumber())
	case models.UpdateGameStatusEvent:
		gameStatusPayload := models.ParseGameStatusPayload(message.GetPayloadJson())
		gs.updateGameStatus(gameStatusPayload.GetGameStatus())
	case models.AlertEvent:
		alertPayload := models.NewAlertPayload(message.GetPayloadJson())
		gs.addAlert(alertPayload.GetAlert())
	}
	return false
}

func (gs *gameState) GetStateJson() string {
	encoded, err := json.Marshal(gs)
	if err != nil {
		log.Fatalln(err)
	}
	return string(encoded)
}
