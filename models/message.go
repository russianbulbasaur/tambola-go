package models

import (
	json2 "encoding/json"
	"log"
)

const UserJoinedEvent = "user_joined"
const UserLeftEvent = "user_left"
const AlertEvent = "alert"
const NumberCalledEvent = "number_called"
const UpdateGameStatusEvent = "game_status"
const PlayersInLobbyEvent = "players_already_in_lobby"

type Message struct {
	Id                           int64                         `json:"id"`
	Event                        string                        `json:"event"`
	Sender                       *User                         `json:"sender"`
	PlayersAlreadyInLobbyPayload *PlayersAlreadyInLobbyPayload `json:"players_already_in_lobby_payload,omitempty"`
	UserJoinedPayload            *UserJoinedPayload            `json:"user_joined_payload,omitempty"`
	UserLeftPayload              *UserLeftPayload              `json:"user_left_payload,omitempty"`
	NumberPayload                *NumberPayload                `json:"number_payload,omitempty"`
	GameStatusPayload            *GameStatusPayload            `json:"game_status_payload,omitempty"`
}

type PlayersAlreadyInLobbyPayload struct {
	Players []*User `json:"players"`
	GameId  int32   `json:"game_id"`
}

type UserJoinedPayload struct {
	User *User `json:"user"`
}

type UserLeftPayload struct {
	User *User `json:"user"`
}

type AlertPayload struct {
	Alert string `json:"alert"`
}

type NumberPayload struct {
	Number int32 `json:"number"`
}

type GameStatusPayload struct {
	Status string `json:"status"`
}

func (message *Message) encode() []byte {
	json, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	return json
}

func Decode(data []byte) *Message {
	decodedMessage := &Message{}
	err := json2.Unmarshal(data, decodedMessage)
	if err != nil {
		log.Println(err)
		return nil
	}
	return decodedMessage
}
