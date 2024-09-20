package payloads

import "cmd/tambola/models"

type UserLeftPayload struct {
	User *models.User `json:"user"`
}
