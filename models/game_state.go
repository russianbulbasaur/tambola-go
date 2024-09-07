package models

import (
	"fmt"
	"log"
)

const Playing = "playing"
const Waiting = "waiting"
const Closed = "closed"

type GameState struct {
	Players       map[*User]bool `json:"players"`
	Host          *User          `json:"host"`
	State         *GameState     `json:"state"`
	Status        string
	alerts        []string
	claimed       []string
	numbersCalled []int64
}

func NewGameState(host *User) *GameState {
	return &GameState{
		Host:    host,
		Players: make(map[*User]bool),
		Status:  "waiting",
	}
}

func (gameState *GameState) addNumber(number int64) {
	log.Println(fmt.Sprintf("Called %d number", number))
	gameState.numbersCalled = append(gameState.numbersCalled, number)
}

func (gameState *GameState) addAlert(alert string) {
	gameState.alerts = append(gameState.alerts, alert)
}

func (gameState *GameState) addPlayer(player *User) {
	log.Println(fmt.Sprintf("User %s joined", player.Name))
	if player.IsHost {
		return
	}
}

func (gameState *GameState) removePlayer(player *User) {
	log.Println(fmt.Sprintf("User %s left", player.Name))
}

func (gameState *GameState) updateGameStatus(status string) {
	log.Println(fmt.Sprintf("Updating game state to %s", status))
	gameState.Status = status
}

func (gameState *GameState) updateGameState(data []byte) bool {
	message := Decode(data)
	switch message.Event {
	case UserJoinedEvent:
		gameState.addPlayer(message.UserJoinedPayload.User)
		break
	case UserLeftEvent:
		gameState.removePlayer(message.UserLeftPayload.User)
		break
	case NumberCalledEvent:
		gameState.addNumber(message.NumberPayload.Number)
		break
	case UpdateGameStatusEvent:
		gameState.updateGameStatus(message.GameStatusPayload.Status)
		break
	}
	return message.Sender.Id == -1 || !message.Sender.IsHost
}
