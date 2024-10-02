package models

import (
	"encoding/json"
	"log"
)

type numberPayload struct {
	Number int32 `json:"number"`
}

type NumberPayload interface {
	GetJson() map[string]interface{}
	GetNumber() int32
}

func NewNumberPayload(number int32) NumberPayload {
	return &numberPayload{number}
}

func ParseNumberPayload(encoded string) NumberPayload {
	var payload numberPayload
	err := json.Unmarshal([]byte(encoded), &payload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &payload
}

func (n *numberPayload) GetNumber() int32 {
	return n.Number
}

func (n *numberPayload) GetJson() map[string]interface{} {
	encoded, err := json.Marshal(n)
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
