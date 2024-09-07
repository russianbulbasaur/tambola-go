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
const GameIdEvent = "game_id"

type Message struct {
	Id      int64  `json:"id"`
	Event   string `json:"event"`
	Sender  *User  `json:"sender"`
	Payload string `json:"payload"`
}

type GameIdPayload struct {
	Id int32 `json:"game_id"`
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
	Number int64 `json:"number"`
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

func (message *Message) decodePlayerPayload() *UserJoinedPayload {
	log.Println(message.Payload)
	player := &UserJoinedPayload{}
	err := json2.Unmarshal([]byte(message.Payload), player)
	if err != nil {
		log.Println(err)
		return nil
	}
	return player
}

func (message *Message) decodeAlertPayload() *AlertPayload {
	alert := &AlertPayload{}
	err := json2.Unmarshal([]byte(message.Payload), alert)
	if err != nil {
		log.Println(err)
		return nil
	}
	return alert
}

func (message *Message) decodeNumberPayload() *NumberPayload {
	numberPayload := &NumberPayload{}
	err := json2.Unmarshal([]byte(message.Payload), numberPayload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return numberPayload
}

func (message *Message) decodeGameStatusPayload() *GameStatusPayload {
	gameStatusPayload := &GameStatusPayload{}
	err := json2.Unmarshal([]byte(message.Payload), gameStatusPayload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return gameStatusPayload
}
