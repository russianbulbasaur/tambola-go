package models

import (
	"encoding/json"
	"log"
)

type alertPayload struct {
	Alert string `json:"alert"`
}
type AlertPayload interface {
	GetJson() map[string]interface{}
	GetAlert() string
}

func NewAlertPayload(alert string) AlertPayload {
	return &alertPayload{alert}
}

func ParseAlertPayload(encoded string) AlertPayload {
	var payload alertPayload
	err := json.Unmarshal([]byte(encoded), &payload)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &payload
}

func (a *alertPayload) GetAlert() string {
	return a.Alert
}

func (a *alertPayload) GetJson() map[string]interface{} {
	encoded, err := json.Marshal(a)
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
