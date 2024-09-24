package models

type Payload interface {
	GetJson() map[string]interface{}
}
