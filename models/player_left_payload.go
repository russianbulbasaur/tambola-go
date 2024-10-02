package models

import (
	"encoding/json"
	"log"
)

type playerLeftPayload struct {
	Player *Player `json:"player"`
}

type PlayerLeftPayload interface {
	GetJson() map[string]interface{}
	GetPlayer() *Player
}

func NewPlayerLeftPayload(player *Player) PlayerLeftPayload {
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

func (u *playerLeftPayload) GetPlayer() *Player {
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
