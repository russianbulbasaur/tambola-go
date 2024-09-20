package models

import (
	"cmd/tambola/models/payloads"
	"fmt"
	"log"
)

const Playing = "playing"
const Waiting = "waiting"
const Closed = "closed"

type gameState struct {
	Players       map[*User]bool `json:"players"`
	Host          *User          `json:"host"`
	Status        string
	alerts        []string
	claimed       []string
	numbersCalled []int32
	playerCount   int32
}

type GameState interface {
	GetPlayers() map[*User]bool
	GetHost() *User
	GetStatus() string
	GetPlayerCount() int32
	GetAlerts() []string
	GetClaimed() []string
	GetCalledNumbers() []int32
	UpdateGameState([]byte) bool
}

func NewGameState(host *User) GameState {
	return &gameState{
		Host:    host,
		Players: make(map[*User]bool),
		Status:  "waiting",
	}
}

func (gs *gameState) GetPlayers() map[*User]bool {
	return gs.Players
}

func (gs *gameState) GetHost() *User {
	return gs.Host
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
	log.Println(fmt.Sprintf("Called %d number", number))
	gs.numbersCalled = append(gs.numbersCalled, number)
}

func (gs *gameState) addAlert(alert string) {
	gs.alerts = append(gs.alerts, alert)
}

func (gs *gameState) addPlayer(player *User) {
	log.Println(fmt.Sprintf("User %s joined", player.Name))
	if player.IsHost {
		return
	}
}

func (gs *gameState) removePlayer(player *User) {
	log.Println(fmt.Sprintf("User %s left", player.Name))
}

func (gs *gameState) updateGameStatus(status string) {
	log.Println(fmt.Sprintf("Updating game state to %s", status))
	gs.Status = status
}

func (gs *gameState) UpdateGameState(data []byte) bool {
	println(string(data))
	message := Decode(data)
	switch message.GetEvent() {
	case UserJoinedEvent:
		gs.addPlayer(message.UserJoinedPayload.User)
		break
	case UserLeftEvent:
		gs.removePlayer(message.UserLeftPayload.User)
		break
	case NumberCalledEvent:
		payload := payloads.ParseNumberPayload(message.GetJsonPayload())
		gs.addNumber(payload.GetNumber())
		break
	case UpdateGameStatusEvent:
		gs.updateGameStatus(message.GameStatusPayload.Status)
		break
	}
	return message.Sender.Id == -1 || !message.Sender.IsHost
}
