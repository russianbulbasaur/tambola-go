package payloads

import (
	"encoding/json"
	"log"
)

type gameStatusPayload struct {
	Status string `json:"status"`
}

type GameStatusPayload interface {
	GetJson() []byte
	GetGameStatus() string
}

func NewGameStatusPayload(gameStatus string) GameStatusPayload {
	return &gameStatusPayload{gameStatus}
}

func ParseGameStatusPayload(encoded string) GameStatusPayload {
	var payload gameStatusPayload
	err := json.Unmarshal([]byte(encoded), &payload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &payload
}

func (g *gameStatusPayload) GetGameStatus() string {
	return g.Status
}

func (g *gameStatusPayload) GetJson() []byte {
	encoded, err := json.Marshal(g)
	if err != nil {
		log.Println(err)
		return nil
	}
	return encoded
}
