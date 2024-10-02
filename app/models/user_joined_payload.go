package models

import (
	"encoding/json"
	"log"
)

type userJoinedPayload struct {
	User *User `json:"user"`
}

type UserJoinedPayload interface {
	GetJson() map[string]interface{}
	GetPlayer() *User
}

func NewUserJoinedPayload(user *User) UserJoinedPayload {
	return &userJoinedPayload{user}
}

func ParseUserJoinedPayload(encoded string) UserJoinedPayload {
	var payload userJoinedPayload
	err := json.Unmarshal([]byte(encoded), &payload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &payload
}

func (u *userJoinedPayload) GetPlayer() *User {
	return u.User
}

func (u *userJoinedPayload) GetJson() map[string]interface{} {
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
