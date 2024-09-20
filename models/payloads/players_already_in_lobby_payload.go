package payloads

import "cmd/tambola/models"

type PlayersAlreadyInLobbyPayload struct {
	Players []*models.User `json:"players"`
	GameId  int32          `json:"game_id"`
}
