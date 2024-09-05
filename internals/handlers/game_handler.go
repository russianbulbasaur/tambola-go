package handlers

import (
	"cmd/tambola/internals/services"
	"net/http"
)

type gameHandler struct {
	gameService services.GameService
}

type GameHandler interface {
	createGame(http.ResponseWriter, *http.Request)
	joinGame(w http.ResponseWriter, r *http.Request)
}

func NewGameHandler(gs services.GameService) GameHandler {
	return &gameHandler{gs}
}

func (gh *gameHandler) createGame(w http.ResponseWriter, r *http.Request) {

}

func (gh *gameHandler) joinGame(w http.ResponseWriter, r *http.Request) {

}
