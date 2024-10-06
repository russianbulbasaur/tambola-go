package models

import (
	"encoding/json"
	"log"
)

type playerLeftPayload struct {
	Player *User `json:"player"`
}

type PlayerLeftPayload interface {
	GetJson() map[string]interface{}
	GetPlayer() *User
}

func NewPlayerLeftPayload(player *User) PlayerLeftPayload {
	return &playerLeftPayload{player}
}

func ParsePlayerLeftPayload(encoded string) PlayerLeftPayload {
	var payload playerLeftPayload
	err := json.Unmarshal([]byte(encoded), &payload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &payload
}

func (u *playerLeftPayload) GetPlayer() *User {
	return u.Player
}

func (u *playerLeftPayload) GetJson() map[string]interface{} {
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
