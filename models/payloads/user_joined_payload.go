package payloads

import (
	"cmd/tambola/models"
	"encoding/json"
	"log"
)

type userJoinedPayload struct {
	User *models.User `json:"user"`
}

type UserJoinedPayload interface {
	GetJson() []byte
	GetPlayer() *models.User
}

func NewUserJoinedPayload(user *models.User) UserJoinedPayload {
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

func (u *userJoinedPayload) GetPlayer() *models.User {
	return u.User
}

func (u *userJoinedPayload) GetJson() []byte {
	encoded, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		return nil
	}
	return encoded
}
