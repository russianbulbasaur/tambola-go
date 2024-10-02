package models

import (
	"encoding/json"
	"log"
)

type userLeftPayload struct {
	User *User `json:"user"`
}

type UserLeftPayload interface {
	GetJson() map[string]interface{}
	GetPlayer() *User
}

func NewUserLeftPayload(user *User) UserLeftPayload {
	return &userLeftPayload{user}
}

func ParseUserLeftPayload(encoded string) UserLeftPayload {
	var payload userLeftPayload
	err := json.Unmarshal([]byte(encoded), &payload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &payload
}

func (u *userLeftPayload) GetPlayer() *User {
	return u.User
}

func (u *userLeftPayload) GetJson() map[string]interface{} {
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
