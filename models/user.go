package models

import (
	"encoding/json"
	"log"
)

type User struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Token  string `json:"token"`
	IsHost bool
}

func (user *User) GetName() string {
	return user.Name
}

func (user *User) EncodeToJson() []byte {
	encoded, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}
	return encoded
}

func ParseUserFromJson(userString string) User {
	var user User
	err := json.Unmarshal([]byte(userString), &user)
	if err != nil {
		log.Fatalln(err)
	}
	return user
}
