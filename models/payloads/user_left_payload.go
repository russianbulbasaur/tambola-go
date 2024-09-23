package payloads

import (
	"cmd/tambola/models"
	"encoding/json"
	"log"
)

type userLeftPayload struct {
	User *models.User `json:"user"`
}

type UserLeftPayload interface {
	GetJson() []byte
	GetPlayer() *models.User
}

func NewUserLeftPayload(user *models.User) UserLeftPayload {
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

func (u *userLeftPayload) GetPlayer() *models.User {
	return u.User
}

func (u *userLeftPayload) GetJson() []byte {
	encoded, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		return nil
	}
	return encoded
}
