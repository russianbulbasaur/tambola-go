package models

import (
	"encoding/json"
	"log"
)

type playerJoinedPayload struct {
	Player *User `json:"player"`
}

type PlayerJoinedPayload interface {
	GetJson() map[string]interface{}
	GetPlayer() *User
}

func NewPlayerJoinedPayload(player *User) PlayerJoinedPayload {
	return &playerJoinedPayload{player}
}

func ParsePlayerJoinedPayload(encoded string) PlayerJoinedPayload {
	var payload playerJoinedPayload
	err := json.Unmarshal([]byte(encoded), &payload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &payload
}

func (u *playerJoinedPayload) GetPlayer() *User {
	return u.Player
}

func (u *playerJoinedPayload) GetJson() map[string]interface{} {
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
