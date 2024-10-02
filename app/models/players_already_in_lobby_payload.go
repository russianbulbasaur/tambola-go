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
	GetJson() map[string]interface{}
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

func (u *playersAlreadyInLobbyPayload) GetJson() map[string]interface{} {
	encoded, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		return nil
	}
	var rawJson map[string]interface{}
	err = json.Unmarshal(encoded, &rawJson)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rawJson
}
