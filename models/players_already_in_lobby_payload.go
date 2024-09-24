package models

import (
	"encoding/json"
	"log"
)

type playersAlreadyInLobbyPayload struct {
	Players []*User `json:"players"`
	GameId  int32   `json:"game_id"`
}

type PlayersAlreadyInLobbyPayload interface {
	GetJson() []byte
}

func NewPlayersAlreadyInLobbyPayload(users []*User,
	gameId int32) PlayersAlreadyInLobbyPayload {
	return &playersAlreadyInLobbyPayload{users, gameId}
}

func ParsePlayersAlreadyInLobbyPayload(encoded string) PlayersAlreadyInLobbyPayload {
	var payload playersAlreadyInLobbyPayload
	err := json.Unmarshal([]byte(encoded), &payload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &payload
}

func (u *playersAlreadyInLobbyPayload) GetJson() []byte {
	encoded, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		return nil
	}
	return encoded
}
