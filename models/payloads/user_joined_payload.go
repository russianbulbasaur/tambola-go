package payloads

import "cmd/tambola/models"

type UserJoinedPayload struct {
	User *models.User `json:"user"`
}
