package game_server

import "cmd/tambola/models"

type GameServer struct {
	Players   map[*models.User]bool `json:"players"`
	Host      *models.User          `json:"host"`
	State     *models.GameState     `json:"state"`
	Join      chan *models.User
	Leave     chan *models.User
	Broadcast chan string
}

func NewGameServer(gameID string, host *models.User) *GameServer {
	return &GameServer{
		Players:   make(map[*models.User]bool),
		Host:      host,
		Join:      make(chan *models.User),
		Leave:     make(chan *models.User),
		Broadcast: make(chan string),
		State:     &models.GameState{},
	}
}

func (gs *GameServer) StartGameServer() {
	for {
		select {
		case user := <-gs.Join:
			gs.registerClient(user)
		case user := <-gs.Leave:
			gs.unregisterClient(user)
		}
	}
}

func (gs *GameServer) registerClient(user *models.User) {

}

func (gs *GameServer) unregisterClient(user *models.User) {

}
