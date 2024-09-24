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

type message struct {
	Id          int64  `json:"id"`
	Event       string `json:"event"`
	Sender      *User  `json:"sender"`
	payload     Payload
	PayloadJson interface{} `json:"payload"`
}

type Message interface {
	GetEvent() string
	GetPayloadJson() string
	SetPayload(payload Payload)
	EncodeToJson() []byte
}

func NewMessage(id int64, event string, sender *User, payload Payload) Message {
	return &message{
		Id:      id,
		Event:   event,
		Sender:  sender,
		payload: payload,
	}
}

func (m *message) SetPayload(payload Payload) {
	m.payload = payload
}

func (m *message) GetPayloadJson() string {
	encoded, _ := json2.Marshal(m.PayloadJson)
	return string(encoded)
}

func (m *message) GetEvent() string {
	return m.Event
}

func (m *message) EncodeToJson() []byte {
	m.PayloadJson = m.payload.GetJson()
	encoded, err := json2.Marshal(m)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Encoded Message : %s", encoded)
	return encoded
}

func Decode(data []byte) Message {
	decodedMessage := &message{}
	err := json2.Unmarshal(data, decodedMessage)
	if err != nil {
		log.Println(err)
		return nil
	}
	return decodedMessage
}
