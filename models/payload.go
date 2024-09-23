package models

type Payload interface {
	GetJson() []byte
}
