package models

import (
	"encoding/json"
	"log"
)

type alertPayload struct {
	Alert string `json:"alert"`
}
type AlertPayload interface {
	GetJson() []byte
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

func (a *alertPayload) GetJson() []byte {
	encoded, err := json.Marshal(a)
	if err != nil {
		log.Println(err)
		return nil
	}
	return encoded
}
