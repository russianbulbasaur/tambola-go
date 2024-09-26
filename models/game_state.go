package models

import (
	"cmd/tambola/utils"
	"fmt"
)

const Playing = "playing"
const Waiting = "waiting"
const Closed = "closed"

type gameState struct {
	players       map[*User]bool
	host          *User
	Status        string
	alerts        []string
	claimed       []string
	numbersCalled []int32
	playerCount   int32
	gameLogger    *utils.TambolaLogger
}

type GameState interface {
	GetPlayers() map[*User]bool
	GetHost() *User
	GetStatus() string
	GetPlayerCount() int32
	GetAlerts() []string
	GetClaimed() []string
	GetCalledNumbers() []int32
	AddPlayer(player *User)
	RemovePlayer(player *User)
	UpdateGameState(data []byte) bool
}

func NewGameState(host *User, logger *utils.TambolaLogger) GameState {
	return &gameState{
		host:       host,
		players:    make(map[*User]bool),
		Status:     "waiting",
		gameLogger: logger,
	}
}

func (gs *gameState) GetPlayers() map[*User]bool {
	return gs.players
}

func (gs *gameState) GetHost() *User {
	return gs.host
}

func (gs *gameState) GetStatus() string {
	return gs.Status
}

func (gs *gameState) GetPlayerCount() int32 {
	return gs.playerCount
}

func (gs *gameState) GetAlerts() []string {
	return gs.alerts
}

func (gs *gameState) GetClaimed() []string {
	return gs.claimed
}

func (gs *gameState) GetCalledNumbers() []int32 {
	return gs.numbersCalled
}

func (gs *gameState) addNumber(number int32) {
	gs.numbersCalled = append(gs.numbersCalled, number)
	gs.gameLogger.Log(fmt.Sprintf("Called %d number", number))
}

func (gs *gameState) addAlert(alert string) {
	gs.alerts = append(gs.alerts, alert)
}

func (gs *gameState) AddPlayer(player *User) {
	gs.players[player] = true
	gs.gameLogger.Log(fmt.Sprintf("User %s joined", player.Name))
}

func (gs *gameState) RemovePlayer(player *User) {
	if gs.players[player] {
		delete(gs.players, player)
	}
	gs.gameLogger.Log(fmt.Sprintf("User %s left", player.Name))
}

func (gs *gameState) updateGameStatus(status string) {
	gs.Status = status
	gs.gameLogger.Log(fmt.Sprintf("Updating game state to %s", status))
}

func (gs *gameState) UpdateGameState(data []byte) bool {
	message := Decode(data)
	switch message.GetEvent() {
	case NumberCalledEvent:
		numberCalledPayload := ParseNumberPayload(message.GetPayloadJson())
		gs.addNumber(numberCalledPayload.GetNumber())
	case UpdateGameStatusEvent:
		gameStatusPayload := NewGameStatusPayload(message.GetPayloadJson())
		gs.updateGameStatus(gameStatusPayload.GetGameStatus())
	case AlertEvent:
		alertPayload := NewAlertPayload(message.GetPayloadJson())
		gs.addAlert(alertPayload.GetAlert())
	}
	return false
}
