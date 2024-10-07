package handlers

import (
	"cmd/tambola/internals/services"
	"cmd/tambola/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type gameHandler struct {
	gameService services.GameService
}

type GameHandler interface {
	CreateGame(http.ResponseWriter, *http.Request)
	JoinGame(w http.ResponseWriter, r *http.Request)
}

func NewGameHandler(gs services.GameService) GameHandler {
	return &gameHandler{gs}
}

func (gh *gameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	user := models.ParseUserFromJson(r.Form.Get("user"))
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	gh.gameService.CreateGame(&user, conn)
}

func (gh *gameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	user := models.ParseUserFromJson(r.Form.Get("user"))
	code := r.URL.Query().Get("code")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	gh.gameService.JoinGame(code, &user, conn)
}
